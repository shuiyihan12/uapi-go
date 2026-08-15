// Package vehicle provides the SOAP client implementation for the Travelport
// Vehicle service (the vehicle-mapped namespace). Its VehicleServicePort
// interface corresponds one-to-one with the WSDL *PortType operations.
// Request/response types come directly from the generator-produced vehicle
// package (script generated; do not edit by hand). This package keeps no
// hand-written request models; all infrastructure fields are injected uniformly
// by prepareReq before sending.
package vehicle

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/internal/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	vehiclexsd "github.com/shuiyihan12/uapi-go/pkg/generated/vehicle"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// VehicleServicePort mirrors the *PortType operations of the vehicle-mapped
// namespace one-to-one, making it easy to swap implementations and plug in test
// stubs. Request/response types come directly from the WSDL-generated vehicle
// package.
type VehicleServicePort interface {
	// VehicleKeyword corresponds to the VehicleKeywordReq operation of the Vehicle service.
	VehicleKeyword(ctx context.Context, req *vehiclexsd.VehicleKeywordReq) (*vehiclexsd.VehicleKeywordRsp, error)
	// VehicleLocation corresponds to the VehicleLocationReq operation of the Vehicle service.
	VehicleLocation(ctx context.Context, req *vehiclexsd.VehicleLocationReq) (*vehiclexsd.VehicleLocationRsp, error)
	// VehicleLocationDetail corresponds to the VehicleLocationDetailReq operation of the Vehicle service.
	VehicleLocationDetail(ctx context.Context, req *vehiclexsd.VehicleLocationDetailReq) (*vehiclexsd.VehicleLocationDetailRsp, error)
	// VehicleMediaLinks corresponds to the VehicleMediaLinksReq operation of the Vehicle service.
	VehicleMediaLinks(ctx context.Context, req *vehiclexsd.VehicleMediaLinksReq) (*vehiclexsd.VehicleMediaLinksRsp, error)
	// VehicleRetrieve corresponds to the VehicleRetrieveReq operation of the Vehicle service.
	VehicleRetrieve(ctx context.Context, req *vehiclexsd.VehicleRetrieveReq) (*vehiclexsd.VehicleRetrieveRsp, error)
	// VehicleRules corresponds to the VehicleRulesReq operation of the Vehicle service.
	VehicleRules(ctx context.Context, req *vehiclexsd.VehicleRulesReq) (*vehiclexsd.VehicleRulesRsp, error)
	// VehicleSearchAvailability corresponds to the VehicleSearchAvailabilityReq operation of the Vehicle service.
	VehicleSearchAvailability(ctx context.Context, req *vehiclexsd.VehicleSearchAvailabilityReq) (*vehiclexsd.VehicleSearchAvailabilityRsp, error)
	// VehicleUpsellSearchAvailability corresponds to the VehicleUpsellSearchAvailabilityReq operation of the Vehicle service.
	VehicleUpsellSearchAvailability(ctx context.Context, req *vehiclexsd.VehicleUpsellSearchAvailabilityReq) (*vehiclexsd.VehicleUpsellSearchAvailabilityRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// VehicleService is the SOAP implementation of VehicleServicePort.
type VehicleService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *VehicleService must satisfy the VehicleServicePort interface.
var _ VehicleServicePort = (*VehicleService)(nil)

// prepareReq ensures the call carries a trace_id and injects one as a fallback
// when the request body's TraceId is empty. Authorization/business fields such
// as the billing point of sale and TargetBranch are no longer injected by code;
// the caller (API user) must provide them explicitly in the request body.
// Authentication (Authorization) is supplied by the caller as an HTTP request
// header and passed through the context to the SOAP call (startup-time
// environment variables are no longer used). Injection follows a "fallback"
// strategy: InjectInfrastructure only fills the request body's TraceId with the
// trace_id when it is empty; a TraceId business value already set by the caller
// is not overwritten. Requests without an InjectInfrastructure implementation
// are skipped (the call is unaffected).
func prepareReq(ctx context.Context, req any) context.Context {
	ctx, traceID := trace.Ensure(ctx)
	if inj, ok := req.(interface{ InjectInfrastructure(traceID string) }); ok {
		inj.InjectInfrastructure(traceID)
	}
	return ctx
}

// NewVehicleService builds a Vehicle service client from the given SOAP configuration and logger.
func NewVehicleService(config client.SOAPConfig, logger logging.Logger) (*VehicleService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "vehicle-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle service client: %w", err)
	}

	return &VehicleService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// VehicleKeyword issues the VehicleKeywordReq SOAP call and returns the strongly typed response.
func (s *VehicleService) VehicleKeyword(ctx context.Context, req *vehiclexsd.VehicleKeywordReq) (*vehiclexsd.VehicleKeywordRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[vehiclexsd.VehicleKeywordRsp](s.client, ctx, "VehicleKeyword", req)
}

// VehicleLocation issues the VehicleLocationReq SOAP call and returns the strongly typed response.
func (s *VehicleService) VehicleLocation(ctx context.Context, req *vehiclexsd.VehicleLocationReq) (*vehiclexsd.VehicleLocationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[vehiclexsd.VehicleLocationRsp](s.client, ctx, "VehicleLocation", req)
}

// VehicleLocationDetail issues the VehicleLocationDetailReq SOAP call and returns the strongly typed response.
func (s *VehicleService) VehicleLocationDetail(ctx context.Context, req *vehiclexsd.VehicleLocationDetailReq) (*vehiclexsd.VehicleLocationDetailRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[vehiclexsd.VehicleLocationDetailRsp](s.client, ctx, "VehicleLocationDetail", req)
}

// VehicleMediaLinks issues the VehicleMediaLinksReq SOAP call and returns the strongly typed response.
func (s *VehicleService) VehicleMediaLinks(ctx context.Context, req *vehiclexsd.VehicleMediaLinksReq) (*vehiclexsd.VehicleMediaLinksRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[vehiclexsd.VehicleMediaLinksRsp](s.client, ctx, "VehicleMediaLinks", req)
}

// VehicleRetrieve issues the VehicleRetrieveReq SOAP call and returns the strongly typed response.
func (s *VehicleService) VehicleRetrieve(ctx context.Context, req *vehiclexsd.VehicleRetrieveReq) (*vehiclexsd.VehicleRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[vehiclexsd.VehicleRetrieveRsp](s.client, ctx, "VehicleRetrieve", req)
}

// VehicleRules issues the VehicleRulesReq SOAP call and returns the strongly typed response.
func (s *VehicleService) VehicleRules(ctx context.Context, req *vehiclexsd.VehicleRulesReq) (*vehiclexsd.VehicleRulesRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[vehiclexsd.VehicleRulesRsp](s.client, ctx, "VehicleRules", req)
}

// VehicleSearchAvailability issues the VehicleSearchAvailabilityReq SOAP call and returns the strongly typed response.
func (s *VehicleService) VehicleSearchAvailability(ctx context.Context, req *vehiclexsd.VehicleSearchAvailabilityReq) (*vehiclexsd.VehicleSearchAvailabilityRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[vehiclexsd.VehicleSearchAvailabilityRsp](s.client, ctx, "VehicleSearchAvailability", req)
}

// VehicleUpsellSearchAvailability issues the VehicleUpsellSearchAvailabilityReq SOAP call and returns the strongly typed response.
func (s *VehicleService) VehicleUpsellSearchAvailability(ctx context.Context, req *vehiclexsd.VehicleUpsellSearchAvailabilityReq) (*vehiclexsd.VehicleUpsellSearchAvailabilityRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[vehiclexsd.VehicleUpsellSearchAvailabilityRsp](s.client, ctx, "VehicleUpsellSearchAvailability", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType that
// performs a single SOAP call and decodes it into the strongly typed response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *VehicleService) Close() error {
	return s.client.Close()
}
