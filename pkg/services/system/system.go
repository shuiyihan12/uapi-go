// Package system provides the SOAP client implementation for the Travelport
// System service (system_v32_0). Its SystemServicePort interface corresponds
// one-to-one with the *PortTypes in the WSDL:
//   - SystemPingPortType          -> Ping
//   - SystemInfoPortType          -> Info
//   - SystemTimePortType          -> Time
//   - ExternalCacheAccessPortType -> ExternalCacheAccess
//
// Request/response types come directly from the generator-produced system
// package (script generated; do not edit by hand): the generator now unifies
// the camelCase JSON contract, emits unqualified attributes per
// attributeFormDefault=unqualified, and generates InjectInfrastructure for
// BaseCoreReq to inject the billing point of sale and the trace identifier.
// Therefore this package keeps no hand-written request models; all
// infrastructure fields are injected uniformly by prepareReq before sending.
package system

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	systemxsd "github.com/shuiyihan12/uapi-go/pkg/generated/system"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// SystemServicePort mirrors the *PortType operations of system_v32_0 one-to-one, making it easy to
// swap implementations and plug in test stubs.
// Request/response types come directly from the generator-produced system package.
type SystemServicePort interface {
	// Ping corresponds to SystemPingPortType and probes service connectivity.
	Ping(ctx context.Context, req *systemxsd.PingReq) (*systemxsd.PingRsp, error)
	// Info corresponds to SystemInfoPortType and retrieves gateway/system information.
	Info(ctx context.Context, req *systemxsd.SystemInfoReq) (*systemxsd.SystemInfoRsp, error)
	// Time corresponds to SystemTimePortType and retrieves the GDS current time.
	Time(ctx context.Context, req *systemxsd.TimeReq) (*systemxsd.TimeRsp, error)
	// ExternalCacheAccess corresponds to ExternalCacheAccessPortType and accesses
	// the external cache (e.g. clearing cache entries).
	ExternalCacheAccess(ctx context.Context, req *systemxsd.ExternalCacheAccessReq) (*systemxsd.ExternalCacheAccessRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// SystemService is the SOAP implementation of SystemServicePort.
type SystemService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *SystemService must satisfy the SystemServicePort interface.
var _ SystemServicePort = (*SystemService)(nil)

// prepareReq injects the call's trace identifier (TraceId) into the request.
// Authorization/business fields such as the billing point of sale and
// TargetBranch are no longer injected by code; the caller (API user) must
// provide them explicitly in the request body. Only UAPI_AUTHORIZATION is
// passed via the HTTP header.
func prepareReq[T infraInjectable](ctx context.Context, req T) context.Context {
	_, traceID := trace.Ensure(ctx)
	req.InjectInfrastructure(traceID)
	return ctx
}

// infraInjectable describes request types whose trace identifier can be
// injected by the service side before sending.
type infraInjectable interface {
	InjectInfrastructure(traceID string)
}

// NewSystemService builds a System service client from the given SOAP configuration and logger.
func NewSystemService(config client.SOAPConfig, logger logging.Logger) (*SystemService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "system-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create system service client: %w", err)
	}

	return &SystemService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// Ping issues the SystemPing SOAP call and returns the system probe response.
func (s *SystemService) Ping(ctx context.Context, req *systemxsd.PingReq) (*systemxsd.PingRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[systemxsd.PingRsp](s.client, ctx, "SystemPing", req)
}

// Info issues the SystemInfo SOAP call and returns the system info response.
func (s *SystemService) Info(ctx context.Context, req *systemxsd.SystemInfoReq) (*systemxsd.SystemInfoRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[systemxsd.SystemInfoRsp](s.client, ctx, "SystemInfo", req)
}

// Time issues the SystemTime SOAP call and returns the GDS current time response.
func (s *SystemService) Time(ctx context.Context, req *systemxsd.TimeReq) (*systemxsd.TimeRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[systemxsd.TimeRsp](s.client, ctx, "SystemTime", req)
}

// ExternalCacheAccess issues the ExternalCacheAccess SOAP call and returns the cache access response.
func (s *SystemService) ExternalCacheAccess(ctx context.Context, req *systemxsd.ExternalCacheAccessReq) (*systemxsd.ExternalCacheAccessRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[systemxsd.ExternalCacheAccessRsp](s.client, ctx, "ExternalCacheAccess", req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *SystemService) Close() error {
	return s.client.Close()
}
