// Package manager creates and uniformly manages the Travelport SOAP service
// clients (currently Hotel and Universal among others), providing service
// discovery, configuration management, health checks and lifecycle teardown.
package manager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	systemxsd "github.com/shuiyihan12/uapi-go/pkg/generated/system"
	"github.com/shuiyihan12/uapi-go/pkg/services/air"
	"github.com/shuiyihan12/uapi-go/pkg/services/gdsQueue"
	"github.com/shuiyihan12/uapi-go/pkg/services/hotel"
	"github.com/shuiyihan12/uapi-go/pkg/services/rail"
	"github.com/shuiyihan12/uapi-go/pkg/services/sharedBooking"
	"github.com/shuiyihan12/uapi-go/pkg/services/sharedUprofile"
	"github.com/shuiyihan12/uapi-go/pkg/services/system"
	"github.com/shuiyihan12/uapi-go/pkg/services/terminal"
	"github.com/shuiyihan12/uapi-go/pkg/services/universal"
	"github.com/shuiyihan12/uapi-go/pkg/services/uprofile"
	"github.com/shuiyihan12/uapi-go/pkg/services/util"
	"github.com/shuiyihan12/uapi-go/pkg/services/vehicle"
	"go.uber.org/zap"
)

// serviceFactory constructs a client instance for a specific service. Held
// by the factory registry so that "adding a service" collapses into one
// registration entry instead of another duplicate getter (see the generic
// Get[T]).
type serviceFactory func(config client.SOAPConfig, logger logging.Logger) (any, error)

// ServiceManager uniformly manages the lifecycle of per-domain SOAP services
// (creation, caching, shutdown and stats).
//
// Service instances are cached as map[string]any and handed out through the
// type-safe generic Get[T](key). Construction logic lives in the factories
// registry, so the manager's type dependencies on concrete service packages
// appear only inside factory closures — adding a domain means appending one
// registration, not another typed getter (open/closed principle).
type ServiceManager struct {
	config    ServiceConfig
	logger    logging.Logger
	services  map[string]any
	factories map[string]serviceFactory
	mu        sync.RWMutex
}

// ServiceConfig describes the connection, timeout and logging configuration
// for ServiceManager and the services it manages.
//
// Authorization and Region are request-level configuration: the caller sends
// them in every HTTP request header and pkg/requestctx forwards them to
// pkg/client, so no Authorization field is kept here.
type ServiceConfig struct {
	// Base configuration.
	BaseEndpoint string // default endpoint prefix (UAPI_ENDPOINT or the apac production environment); used when no region is specified
	Environment  string // "test" or "production"

	// Timeout settings (all at the HTTP transport level, constraining ONE
	// SOAP call, in time.Duration units).
	//
	// RequestTimeout: total timeout of a single SOAP call (connect, send,
	// read, transfer), from UAPI_REQUEST_TIMEOUT. Maps to
	// http.Client.Timeout and is the hard cap per request.
	// ConnectionTimeout: max time to establish the TCP/TLS connection to the
	// GDS, from UAPI_CONNECTION_TIMEOUT. Maps to
	// http.Transport.DialContext.Timeout.
	// ReadTimeout: after connecting, max wait for the GDS to start returning
	// response headers, from UAPI_READ_TIMEOUT. Maps to
	// http.Transport.ResponseHeaderTimeout (times out when the server never
	// starts responding).
	//
	// Note: these three only bound the network time of ONE SOAP call; they
	// are unrelated to business-level cumulative pagination time (e.g.
	// HotelDetails' NextResultReference paging). This project does not
	// retry — any single failed call returns immediately.
	RequestTimeout    time.Duration
	ConnectionTimeout time.Duration
	ReadTimeout       time.Duration

	// Outbound keep-alive pool sizing: warm idle connections to the GDS,
	// avoiding a TCP+TLS re-handshake per request. Per-host is capped by the
	// global value. Overridable via UAPI_MAX_IDLE_CONNS /
	// UAPI_MAX_IDLE_CONNS_PER_HOST.
	MaxIdleConns        int
	MaxIdleConnsPerHost int

	// Logging configuration.
	LogLevel      string
	IsDevelopment bool

	// SkipTLSVerify controls whether upstream TLS certificate verification is
	// skipped. Securely defaults to false (always verify); only enabled
	// explicitly via UAPI_SKIP_TLS_VERIFY=1 for private environments with
	// self-signed certificates, never implicitly tied to Environment.
	SkipTLSVerify bool
}

// DefaultServiceConfig returns the default service configuration for Hotel
// and Universal.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		Environment:         "test",
		BaseEndpoint:        "https://apac.universal-api.travelport.com/B2BGateway/connect/uAPI",
		RequestTimeout:      60 * time.Second,
		ConnectionTimeout:   30 * time.Second,
		ReadTimeout:         60 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		LogLevel:            "info",
		IsDevelopment:       true,
	}
}

// NewServiceManager builds a ServiceManager from the given configuration,
// initializing the logger and resolving the endpoint.
func NewServiceManager(config ServiceConfig) (*ServiceManager, error) {
	logger, err := logging.NewLogger(config.LogLevel, config.IsDevelopment)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	// In production without an explicit endpoint, fall back to the apac
	// production endpoint as the default base.
	if config.Environment == "production" && config.BaseEndpoint == "" {
		config.BaseEndpoint = "https://apac.universal-api.travelport.com/B2BGateway/connect/uAPI"
	}

	manager := &ServiceManager{
		config:    config,
		logger:    logger,
		services:  make(map[string]any),
		factories: buildServiceFactories(),
	}

	logger.Info("Service manager created successfully",
		zap.String("environment", config.Environment))

	return manager, nil
}

// buildServiceFactories returns the service key -> constructor registry.
// Adding a domain service means appending one line here — no new typed
// getters (open/closed principle).
func buildServiceFactories() map[string]serviceFactory {
	return map[string]serviceFactory{
		"air":      func(c client.SOAPConfig, l logging.Logger) (any, error) { return air.NewAirService(c, l) },
		"rail":     func(c client.SOAPConfig, l logging.Logger) (any, error) { return rail.NewRailService(c, l) },
		"vehicle":  func(c client.SOAPConfig, l logging.Logger) (any, error) { return vehicle.NewVehicleService(c, l) },
		"gdsQueue": func(c client.SOAPConfig, l logging.Logger) (any, error) { return gdsQueue.NewGdsQueueService(c, l) },
		"sharedBooking": func(c client.SOAPConfig, l logging.Logger) (any, error) {
			return sharedBooking.NewSharedBookingService(c, l)
		},
		"uprofile": func(c client.SOAPConfig, l logging.Logger) (any, error) { return uprofile.NewUprofileService(c, l) },
		"sharedUprofile": func(c client.SOAPConfig, l logging.Logger) (any, error) {
			return sharedUprofile.NewSharedUprofileService(c, l)
		},
		"terminal":  func(c client.SOAPConfig, l logging.Logger) (any, error) { return terminal.NewTerminalService(c, l) },
		"universal": func(c client.SOAPConfig, l logging.Logger) (any, error) { return universal.NewUniversalService(c, l) },
		"system":    func(c client.SOAPConfig, l logging.Logger) (any, error) { return system.NewSystemService(c, l) },
		"util":      func(c client.SOAPConfig, l logging.Logger) (any, error) { return util.NewUtilService(c, l) },
		"hotel":     func(c client.SOAPConfig, l logging.Logger) (any, error) { return hotel.NewHotelService(c, l) },
	}
}

// Get lazily creates and caches the service client for the given key in a
// type-safe way (a generic free function; Go disallows generic methods, so
// it is called as manager.Get[T](m, key)).
//
// The caller declares the expected service type as the type parameter and
// gets compile-time checking, eliminating the type erasure and DRY issues of
// the former map[string]interface{} plus 12 duplicate getters:
//
//	svc, err := manager.Get[*hotel.HotelService](serviceManager, "hotel")
//
// The first access constructs via the factory registry and caches; later
// accesses return the cached instance (double-checked locking).
func Get[T any](m *ServiceManager, key string) (T, error) {
	var zero T
	if svc, ok := m.lookup(key); ok {
		if typed, ok := svc.(T); ok {
			return typed, nil
		}
		return zero, fmt.Errorf("service %q cached as %T, incompatible with requested %T", key, svc, zero)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if svc, ok := m.services[key]; ok {
		if typed, ok := svc.(T); ok {
			return typed, nil
		}
		return zero, fmt.Errorf("service %q cached as %T, incompatible with requested %T", key, svc, zero)
	}

	factory, ok := m.factories[key]
	if !ok {
		return zero, fmt.Errorf("no service factory registered for key %q", key)
	}
	cfg, err := m.createSOAPConfig(key)
	if err != nil {
		return zero, err
	}
	svc, err := factory(cfg, m.logger)
	if err != nil {
		return zero, fmt.Errorf("failed to create %s service: %w", key, err)
	}
	typed, ok := svc.(T)
	if !ok {
		return zero, fmt.Errorf("factory for %q produced %T, incompatible with requested %T", key, svc, zero)
	}
	m.services[key] = svc
	m.logger.Info("service created and cached", zap.String("service", key))
	return typed, nil
}

// lookup finds a cached service under the read lock, returning (svc, true)
// on hit.
func (m *ServiceManager) lookup(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	svc, ok := m.services[key]
	return svc, ok
}

// createSOAPConfig builds the SOAP call configuration for a target service
// from the manager configuration and service name.
//
// Auth and region are not fixed at construction time: Authorization is
// forwarded per request via context from the caller's header, and Region
// resolves the endpoint at runtime, so only the default base and the UAPI
// service name (appended to the endpoint) are set here. TLS certificate
// verification is on by default and skipped only when the configuration
// explicitly sets SkipTLSVerify (independent of Environment).
func (m *ServiceManager) createSOAPConfig(serviceName string) (client.SOAPConfig, error) {
	suffix, err := serviceSuffix(serviceName)
	if err != nil {
		return client.SOAPConfig{}, err
	}
	return client.SOAPConfig{
		BaseEndpoint:        m.config.BaseEndpoint,
		ServiceName:         suffix,
		Timeout:             m.config.RequestTimeout,
		ConnectionTimeout:   m.config.ConnectionTimeout,
		ReadTimeout:         m.config.ReadTimeout,
		MaxIdleConns:        m.config.MaxIdleConns,
		MaxIdleConnsPerHost: m.config.MaxIdleConnsPerHost,
		SkipTLSVerify:       m.config.SkipTLSVerify,
	}, nil
}

// HealthCheck performs a real SystemPing against the System service and
// considers the system healthy only when the GDS is reachable and answers
// successfully; network unreachability, auth failure or upstream errors all
// return an error, and the caller decides whether to answer 503 or raise an
// alert.
//
// Auth follows this project's "pass everything through" principle: this
// package holds no credentials; the Authorization and Region used by the
// probe are injected into the context by the caller (the prober of /health),
// over the same pass-through chain as business requests.
func (m *ServiceManager) HealthCheck(ctx context.Context) error {
	sysSvc, err := Get[*system.SystemService](m, "system")
	if err != nil {
		return fmt.Errorf("health: cannot resolve system service: %w", err)
	}
	if _, err := sysSvc.Ping(ctx, &systemxsd.PingReq{}); err != nil {
		return fmt.Errorf("health: upstream system ping failed: %w", err)
	}
	return nil
}

// Config returns a snapshot of the manager's current configuration for the
// use-case layer.
func (m *ServiceManager) Config() ServiceConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// Close closes every created service and releases the underlying
// connections.
func (m *ServiceManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for serviceName, service := range m.services {
		if closer, ok := service.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close %s service: %v", serviceName, err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred while closing services: %v", errs)
	}

	m.logger.Info("All services closed successfully")
	return nil
}

// GetServiceStats returns service runtime stats: managed service count,
// active service names and a summary of key configuration.
func (m *ServiceManager) GetServiceStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["total_services"] = len(m.services)
	stats["active_services"] = make([]string, 0, len(m.services))

	for serviceName := range m.services {
		stats["active_services"] = append(stats["active_services"].([]string), serviceName)
	}

	stats["config"] = map[string]interface{}{
		"environment":     m.config.Environment,
		"base_endpoint":   m.config.BaseEndpoint,
		"request_timeout": m.config.RequestTimeout.String(),
	}

	return stats
}

// serviceSuffix maps an internal service key to the Travelport UAPI SOAP
// service name (appended to the endpoint). Production only; sharedUprofile
// keeps the uAPI path (its separate /uProfile/ endpoint is ignored). Unknown
// keys return an error instead of silently falling back to HotelService, so
// a mis-routed endpoint stays easy to diagnose.
func serviceSuffix(key string) (string, error) {
	switch key {
	case "air":
		return "AirService", nil
	case "rail":
		return "RailService", nil
	case "universal":
		return "UniversalRecordService", nil
	case "cruise":
		return "CruiseService", nil
	case "vehicle":
		return "VehicleService", nil
	case "gdsQueue":
		return "GdsQueueService", nil
	case "passive":
		return "PassiveService", nil
	case "sharedBooking":
		return "BookingService", nil
	case "uprofile":
		return "UProfileService", nil
	case "sharedUprofile":
		return "SharedUprofileService", nil
	case "terminal":
		return "TerminalService", nil
	case "sessionContext":
		return "SessionService", nil
	case "system":
		return "SystemService", nil
	case "util":
		return "UtilService", nil
	case "hotel":
		return "HotelService", nil
	default:
		return "", fmt.Errorf("unknown service key %q: no UAPI service name mapping", key)
	}
}
