// Package main is the UAPI service daemon: it serves the business API
// (/api/<domain>/<op>) and the operations endpoints (/health, /ready,
// /version, /stats, /metrics) on a single HTTP port, with graceful shutdown
// on SIGINT/SIGTERM.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shuiyihan12/uapi-go/internal/metrics"
	"github.com/shuiyihan12/uapi-go/pkg/api"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// version is injected at build time via -ldflags "-X main.version=..."
// (see Dockerfile / release scripts); it stays "dev" under plain go run.
var version = "dev"

// main starts the UAPI service daemon: it parses configuration, creates the
// service manager, wires the operations endpoints and business routes on a
// single port, and waits for SIGINT/SIGTERM for graceful shutdown.
//
// Startup configuration has only two knobs (both overridable via environment
// variables): -env (test/production, affects logging shape) and -port
// (default 8080). Operations endpoints such as health checks and metrics
// share the same port as the business API, simplifying deployment (one port,
// one TLS certificate, one load-balancer configuration).
func main() {
	var (
		environment = flag.String("env", envString("UAPI_ENV", "test"), "Environment (test/production)")
		port        = flag.String("port", envString("PORT", "8080"), "Server port (business + ops endpoints)")
	)
	flag.Parse()

	// Authorization and Region are request-level configuration: the caller
	// sends them in every HTTP request header and pkg/requestctx forwards
	// them to pkg/client, so no UAPI_AUTHORIZATION is required at startup.
	// UAPI_ENDPOINT serves as the default endpoint prefix (used when no
	// region is specified); it defaults to the apac production environment.

	// Build the service configuration.
	config := manager.DefaultServiceConfig()
	config.BaseEndpoint = envString("UAPI_ENDPOINT", config.BaseEndpoint)
	config.Environment = *environment
	config.LogLevel = "info"
	config.IsDevelopment = *environment != "production"
	config.ConnectionTimeout = envDurationMS("UAPI_CONNECTION_TIMEOUT", config.ConnectionTimeout) // connect timeout (ms): max time to establish the TCP/TLS connection
	config.ReadTimeout = envDurationMS("UAPI_READ_TIMEOUT", config.ReadTimeout)                   // response-header timeout (ms): max wait for the GDS to start responding after connecting
	config.RequestTimeout = envDurationMS("UAPI_REQUEST_TIMEOUT", config.RequestTimeout)          // per-request total timeout (ms): hard cap over connect+send+read+transfer
	config.MaxIdleConns = envInt("UAPI_MAX_IDLE_CONNS", config.MaxIdleConns)                      // warm keep-alive pool across all upstream hosts
	config.MaxIdleConnsPerHost = envInt("UAPI_MAX_IDLE_CONNS_PER_HOST", config.MaxIdleConnsPerHost)
	// Instrument SOAP calls with the daemon's Prometheus collector; SDK
	// consumers default to the no-op implementation.
	config.Metrics = metrics.GetMetrics() // warm keep-alive pool per GDS host (capped by the global pool)
	// TLS certificate verification is on by default; only private
	// environments (self-signed certificates) explicitly skip it via
	// UAPI_SKIP_TLS_VERIFY=1.
	config.SkipTLSVerify = envBool("UAPI_SKIP_TLS_VERIFY")

	// Create the service manager.
	serviceManager, err := manager.NewServiceManager(config)
	if err != nil {
		log.Fatalf("Failed to create service manager: %v", err)
	}

	hotelFacade := usecase.NewHotelFacade(serviceManager)
	systemFacade := usecase.NewSystemFacade(serviceManager)
	utilFacade := usecase.NewUtilFacade(serviceManager)
	airFacade := usecase.NewAirFacade(serviceManager)
	railFacade := usecase.NewRailFacade(serviceManager)
	vehicleFacade := usecase.NewVehicleFacade(serviceManager)
	gdsQueueFacade := usecase.NewGdsQueueFacade(serviceManager)
	sharedBookingFacade := usecase.NewSharedBookingFacade(serviceManager)
	uprofileFacade := usecase.NewUprofileFacade(serviceManager)
	sharedUprofileFacade := usecase.NewSharedUprofileFacade(serviceManager)
	terminalFacade := usecase.NewTerminalFacade(serviceManager)
	universalFacade := usecase.NewUniversalFacade(serviceManager)

	// Signal-driven graceful shutdown context: ctx cancels on
	// SIGINT/SIGTERM, the main loop exits and triggers Shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ---- Single HTTP port: operations endpoints + business API ----
	// Operations endpoints (/health, /ready, /stats, /metrics) and business
	// APIs (/api/*) share the same mux without path conflicts; deployments
	// expose a single port and can split by path at the ingress / Service
	// layer when isolation is needed.
	mux := http.NewServeMux()

	// /health issues a real upstream SystemPing. Auth follows this project's
	// "pass everything through" principle: the probing caller (monitoring
	// system / k8s probe) sends Authorization in the request header and the
	// gateway forwards it verbatim to UAPI. Without credentials the upstream
	// rejects and /health returns 503 — truthfully reflecting "unable to
	// verify upstream reachability".
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		hctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		hctx = requestctx.WithAuthorization(hctx, r.Header.Get(requestctx.HeaderAuthorization))
		hctx = requestctx.WithRegion(hctx, r.Header.Get(requestctx.HeaderRegion))

		if err := serviceManager.HealthCheck(hctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Health check failed: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ready")
	})
	// /version reports the build version injected at compile time via
	// -ldflags "-X main.version=..." (defaults to "dev" under plain go run).
	// It is a cheap, upstream-independent probe so ops/monitoring can tell
	// which binary is running without touching /health (which depends on the
	// UAPI upstream). Keep it dependency-free and fast.
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"service":"uapi-go","version":%q}`, version)
	})
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := serviceManager.GetServiceStats()
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintf(w, `{
			"total_services": %v,
			"active_services": %v,
			"environment": "%v"
		}`,
			stats["total_services"],
			len(stats["active_services"].([]string)),
			stats["config"].(map[string]interface{})["environment"])
	})
	mux.Handle("/metrics", promhttp.Handler())

	api.NewHotelHandler(hotelFacade, universalFacade).RegisterRoutes(mux)
	api.NewSystemHandler(systemFacade).RegisterRoutes(mux)
	api.NewUtilHandler(utilFacade).RegisterRoutes(mux)
	api.NewAirHandler(airFacade, universalFacade).RegisterRoutes(mux)
	api.NewRailHandler(railFacade, universalFacade).RegisterRoutes(mux)
	api.NewVehicleHandler(vehicleFacade, universalFacade).RegisterRoutes(mux)
	api.NewGdsQueueHandler(gdsQueueFacade).RegisterRoutes(mux)
	api.NewSharedBookingHandler(sharedBookingFacade).RegisterRoutes(mux)
	api.NewUprofileHandler(uprofileFacade).RegisterRoutes(mux)
	api.NewSharedUprofileHandler(sharedUprofileFacade).RegisterRoutes(mux)
	api.NewTerminalHandler(terminalFacade).RegisterRoutes(mux)
	api.NewUniversalHandler(universalFacade).RegisterRoutes(mux)

	// passive has no WSDL/portType of its own; its booking/cancel
	// capabilities are exposed through UniversalRecord.
	api.NewPassiveHandler(universalFacade).RegisterRoutes(mux)

	srv := &http.Server{Addr: ":" + *port, Handler: mux}

	go func() {
		log.Printf("API server starting on port %s", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server failed: %v", err)
		}
	}()

	log.Println("UAPI service daemon started successfully")
	log.Printf("Version:           %s", version)
	log.Printf("Version endpoint:  GET  http://localhost:%s/version", *port)
	log.Printf("Health:            GET  http://localhost:%s/health", *port)
	log.Printf("Readiness:         GET  http://localhost:%s/ready", *port)
	log.Printf("Stats:             GET  http://localhost:%s/stats", *port)
	log.Printf("Prometheus:        GET  http://localhost:%s/metrics", *port)
	log.Printf("Business APIs:     POST http://localhost:%s/api/{domain}/{op}", *port)

	// Block until SIGINT/SIGTERM. Health checks are pull-based (initiated by
	// the probing caller); there is no periodic in-process probing —
	// background tasks have no request-level credentials to forward, and
	// periodic probing would only create noise.
	<-ctx.Done()

	// ---- Graceful shutdown: stop accepting new requests, drain in-flight
	// ones, then release service connections ----
	log.Println("Shutting down, draining in-flight requests...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("API server shutdown error: %v", err)
	}
	if err := serviceManager.Close(); err != nil {
		log.Printf("Service manager close error: %v", err)
	}
	log.Println("Shutdown complete")
}

// envString reads environment variable key, trims surrounding whitespace and
// returns it; falls back when unset or empty.
func envString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

// envDurationMS reads environment variable key, parses it as milliseconds
// into a time.Duration; falls back when unset, empty or unparsable.
func envDurationMS(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return time.Duration(parsed) * time.Millisecond
}

// envInt reads integer environment variable key; falls back when unset,
// empty or unparsable.
func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// envBool reads boolean environment variable key; returns true only for the
// values 1 / true (case-insensitive).
func envBool(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true":
		return true
	default:
		return false
	}
}
