// Package gdsQueue provides the SOAP client implementation for the Travelport GdsQueue
// service (the gdsqueue-mapped namespace).
// Its GdsQueueServicePort interface corresponds one-to-one with the WSDL *PortType operations.
// Request/response types come directly from the generator-produced gdsqueue package
// (script generated; do not edit by hand).
// This package keeps no hand-written request models; all infrastructure fields
// are injected uniformly by prepareReq before sending.
package gdsQueue

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	gdsqueuexsd "github.com/shuiyihan12/uapi-go/pkg/generated/gdsqueue"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// GdsQueueServicePort mirrors the *PortType operations of the gdsqueue-mapped namespace one-to-one,
// making it easy to swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated gdsqueue package.
type GdsQueueServicePort interface {
	// GdsEnterQueue corresponds to the GdsEnterQueueReq operation of the GdsQueue service.
	GdsEnterQueue(ctx context.Context, req *gdsqueuexsd.GdsEnterQueueReq) (*gdsqueuexsd.GdsEnterQueueRsp, error)
	// GdsExitQueue corresponds to the GdsExitQueueReq operation of the GdsQueue service.
	GdsExitQueue(ctx context.Context, req *gdsqueuexsd.GdsExitQueueReq) (*gdsqueuexsd.GdsExitQueueRsp, error)
	// GdsNextOnQueue corresponds to the GdsNextOnQueueReq operation of the GdsQueue service.
	GdsNextOnQueue(ctx context.Context, req *gdsqueuexsd.GdsNextOnQueueReq) (*gdsqueuexsd.GdsNextOnQueueRsp, error)
	// GdsQueueAgentList corresponds to the GdsQueueAgentListReq operation of the GdsQueue service.
	GdsQueueAgentList(ctx context.Context, req *gdsqueuexsd.GdsQueueAgentListReq) (*gdsqueuexsd.GdsQueueAgentListRsp, error)
	// GdsQueueCount corresponds to the GdsQueueCountReq operation of the GdsQueue service.
	GdsQueueCount(ctx context.Context, req *gdsqueuexsd.GdsQueueCountReq) (*gdsqueuexsd.GdsQueueCountRsp, error)
	// GdsQueueList corresponds to the GdsQueueListReq operation of the GdsQueue service.
	GdsQueueList(ctx context.Context, req *gdsqueuexsd.GdsQueueListReq) (*gdsqueuexsd.GdsQueueListRsp, error)
	// GdsQueuePlace corresponds to the GdsQueuePlaceReq operation of the GdsQueue service.
	GdsQueuePlace(ctx context.Context, req *gdsqueuexsd.GdsQueuePlaceReq) (*gdsqueuexsd.GdsQueuePlaceRsp, error)
	// GdsQueueRemove corresponds to the GdsQueueRemoveReq operation of the GdsQueue service.
	GdsQueueRemove(ctx context.Context, req *gdsqueuexsd.GdsQueueRemoveReq) (*gdsqueuexsd.GdsQueueRemoveRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// GdsQueueService is the SOAP implementation of GdsQueueServicePort.
type GdsQueueService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *GdsQueueService must satisfy the GdsQueueServicePort interface.
var _ GdsQueueServicePort = (*GdsQueueService)(nil)

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

// NewGdsQueueService builds a GdsQueue service client from the given SOAP configuration and logger.
func NewGdsQueueService(config client.SOAPConfig, logger logging.Logger) (*GdsQueueService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "gdsQueue-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create gdsQueue service client: %w", err)
	}

	return &GdsQueueService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// GdsEnterQueue issues the GdsEnterQueueReq SOAP call and returns the strongly typed response.
func (s *GdsQueueService) GdsEnterQueue(ctx context.Context, req *gdsqueuexsd.GdsEnterQueueReq) (*gdsqueuexsd.GdsEnterQueueRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[gdsqueuexsd.GdsEnterQueueRsp](s.client, ctx, "GdsEnterQueue", req)
}

// GdsExitQueue issues the GdsExitQueueReq SOAP call and returns the strongly typed response.
func (s *GdsQueueService) GdsExitQueue(ctx context.Context, req *gdsqueuexsd.GdsExitQueueReq) (*gdsqueuexsd.GdsExitQueueRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[gdsqueuexsd.GdsExitQueueRsp](s.client, ctx, "GdsExitQueue", req)
}

// GdsNextOnQueue issues the GdsNextOnQueueReq SOAP call and returns the strongly typed response.
func (s *GdsQueueService) GdsNextOnQueue(ctx context.Context, req *gdsqueuexsd.GdsNextOnQueueReq) (*gdsqueuexsd.GdsNextOnQueueRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[gdsqueuexsd.GdsNextOnQueueRsp](s.client, ctx, "GdsNextOnQueue", req)
}

// GdsQueueAgentList issues the GdsQueueAgentListReq SOAP call and returns the strongly typed response.
func (s *GdsQueueService) GdsQueueAgentList(ctx context.Context, req *gdsqueuexsd.GdsQueueAgentListReq) (*gdsqueuexsd.GdsQueueAgentListRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[gdsqueuexsd.GdsQueueAgentListRsp](s.client, ctx, "GdsQueueAgentList", req)
}

// GdsQueueCount issues the GdsQueueCountReq SOAP call and returns the strongly typed response.
func (s *GdsQueueService) GdsQueueCount(ctx context.Context, req *gdsqueuexsd.GdsQueueCountReq) (*gdsqueuexsd.GdsQueueCountRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[gdsqueuexsd.GdsQueueCountRsp](s.client, ctx, "GdsQueueCount", req)
}

// GdsQueueList issues the GdsQueueListReq SOAP call and returns the strongly typed response.
func (s *GdsQueueService) GdsQueueList(ctx context.Context, req *gdsqueuexsd.GdsQueueListReq) (*gdsqueuexsd.GdsQueueListRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[gdsqueuexsd.GdsQueueListRsp](s.client, ctx, "GdsQueueList", req)
}

// GdsQueuePlace issues the GdsQueuePlaceReq SOAP call and returns the strongly typed response.
func (s *GdsQueueService) GdsQueuePlace(ctx context.Context, req *gdsqueuexsd.GdsQueuePlaceReq) (*gdsqueuexsd.GdsQueuePlaceRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[gdsqueuexsd.GdsQueuePlaceRsp](s.client, ctx, "GdsQueuePlace", req)
}

// GdsQueueRemove issues the GdsQueueRemoveReq SOAP call and returns the strongly typed response.
func (s *GdsQueueService) GdsQueueRemove(ctx context.Context, req *gdsqueuexsd.GdsQueueRemoveReq) (*gdsqueuexsd.GdsQueueRemoveRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[gdsqueuexsd.GdsQueueRemoveRsp](s.client, ctx, "GdsQueueRemove", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType
// that performs a single SOAP call and decodes it into the strongly typed
// response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *GdsQueueService) Close() error {
	return s.client.Close()
}
