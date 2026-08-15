// Package manager_test contains external tests for pkg/manager, verifying
// service manager creation, health checks and caching behavior.
package manager_test

import (
	"context"
	"os"
	"testing"
	"time"

	hotel "github.com/shuiyihan12/uapi-go/pkg/services/hotel"
	universal "github.com/shuiyihan12/uapi-go/pkg/services/universal"

	"github.com/shuiyihan12/uapi-go/pkg/manager"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
)

// TestServiceManager verifies that the service manager creates the
// hotel/universal services correctly, passes health checks, returns service
// stats and caches services.
func TestServiceManager(t *testing.T) {
	config := manager.DefaultServiceConfig()
	config.Environment = "test"

	serviceManager, err := manager.NewServiceManager(config)
	if err != nil {
		t.Fatalf("Failed to create service manager: %v", err)
	}
	defer serviceManager.Close()

	t.Run("HotelService", func(t *testing.T) {
		hotelService, err := manager.Get[*hotel.HotelService](serviceManager, "hotel")
		if err != nil {
			t.Errorf("Failed to get hotel service: %v", err)
		}
		if hotelService == nil {
			t.Error("Hotel service is nil")
		}
	})

	t.Run("HealthCheck", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// HealthCheck always performs a real upstream SystemPing with auth
		// forwarded via context (same chain as business requests). Without
		// injected credentials the upstream must reject (401 / network
		// unreachable), so an error with the health: prefix is expected.
		if err := serviceManager.HealthCheck(ctx); err == nil {
			t.Error("Expected HealthCheck without credentials to fail against upstream")
		}

		// Credentialed success path: only exercised when both
		// UAPI_INTEGRATION=1 and UAPI_TEST_AUTHORIZATION are set (requires a
		// reachable Travelport environment).
		if os.Getenv("UAPI_INTEGRATION") != "" {
			auth := os.Getenv("UAPI_TEST_AUTHORIZATION")
			if auth == "" {
				t.Fatal("UAPI_INTEGRATION is set but UAPI_TEST_AUTHORIZATION is missing")
			}
			actx := requestctx.WithAuthorization(ctx, auth)
			if err := serviceManager.HealthCheck(actx); err != nil {
				t.Errorf("Upstream health check with credentials failed: %v", err)
			}
		} else {
			t.Log("set UAPI_INTEGRATION=1 and UAPI_TEST_AUTHORIZATION to exercise credentialed HealthCheck")
		}
	})

	t.Run("ServiceStats", func(t *testing.T) {
		manager.Get[*hotel.HotelService](serviceManager, "hotel")
		manager.Get[*universal.UniversalService](serviceManager, "universal")

		stats := serviceManager.GetServiceStats()
		if stats == nil {
			t.Error("Service stats is nil")
		}

		totalServices, ok := stats["total_services"]
		if !ok {
			t.Error("Missing total_services in stats")
		}

		if totalServices.(int) < 2 {
			t.Errorf("Expected at least 2 services, got %v", totalServices)
		}
	})

	t.Run("ServiceCaching", func(t *testing.T) {
		hotelService1, err := manager.Get[*hotel.HotelService](serviceManager, "hotel")
		if err != nil {
			t.Fatalf("Failed to get hotel service first time: %v", err)
		}

		hotelService2, err := manager.Get[*hotel.HotelService](serviceManager, "hotel")
		if err != nil {
			t.Fatalf("Failed to get hotel service second time: %v", err)
		}

		if hotelService1 != hotelService2 {
			t.Error("Services are not cached properly")
		}
	})

	t.Run("UniversalService", func(t *testing.T) {
		universalService, err := manager.Get[*universal.UniversalService](serviceManager, "universal")
		if err != nil {
			t.Errorf("Failed to get universal service: %v", err)
		}
		if universalService == nil {
			t.Error("Universal service is nil")
		}
	})
}

// TestDefaultServiceConfig verifies that DefaultServiceConfig returns
// sensible default values.
func TestDefaultServiceConfig(t *testing.T) {
	config := manager.DefaultServiceConfig()

	if config.Environment != "test" {
		t.Errorf("Expected default environment to be 'test', got %s", config.Environment)
	}

	if config.RequestTimeout != 60*time.Second {
		t.Errorf("Expected default timeout to be 60s, got %v", config.RequestTimeout)
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected default log level to be 'info', got %s", config.LogLevel)
	}

	if !config.IsDevelopment {
		t.Error("Expected default development mode to be true")
	}

	if config.BaseEndpoint == "" {
		t.Error("Expected default base endpoint")
	}

	// Security regression: TLS certificate verification must be on by
	// default (SkipTLSVerify=false), never implicitly tied to Environment;
	// skipping is only possible via explicit configuration
	// (UAPI_SKIP_TLS_VERIFY=1).
	if config.SkipTLSVerify {
		t.Error("Expected default SkipTLSVerify to be false (TLS verification must be on by default)")
	}
}

// BenchmarkServiceManagerCreation benchmarks service manager creation and
// teardown.
func BenchmarkServiceManagerCreation(b *testing.B) {
	config := manager.DefaultServiceConfig()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		serviceManager, err := manager.NewServiceManager(config)
		if err != nil {
			b.Fatalf("Failed to create service manager: %v", err)
		}
		serviceManager.Close()
	}
}

// BenchmarkServiceRetrieval benchmarks hotel service retrieval from the
// service manager.
func BenchmarkServiceRetrieval(b *testing.B) {
	config := manager.DefaultServiceConfig()

	serviceManager, err := manager.NewServiceManager(config)
	if err != nil {
		b.Fatalf("Failed to create service manager: %v", err)
	}
	defer serviceManager.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := manager.Get[*hotel.HotelService](serviceManager, "hotel")
		if err != nil {
			b.Fatalf("Failed to get hotel service: %v", err)
		}
	}
}
