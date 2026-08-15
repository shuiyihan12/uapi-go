// Package rail provides the SOAP client implementation for the Travelport Rail
// service (the air-mapped namespace). Its RailServicePort interface corresponds
// one-to-one with the WSDL *PortType operations. Request/response types come
// directly from the generator-produced air package (script generated; do not
// edit by hand). This package keeps no hand-written request models; all
// infrastructure fields are injected uniformly by prepareReq before sending.
package rail

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/internal/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	airxsd "github.com/shuiyihan12/uapi-go/pkg/generated/air"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// RailServicePort mirrors the *PortType operations of the air-mapped namespace
// one-to-one, making it easy to swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated air package.
type RailServicePort interface {
	// RailAvailabilitySearch corresponds to the RailAvailabilitySearchReq operation of the Rail service.
	RailAvailabilitySearch(ctx context.Context, req *airxsd.RailAvailabilitySearchReq) (*airxsd.RailAvailabilitySearchRsp, error)
	// RailExchange corresponds to the RailExchangeReq operation of the Rail service.
	RailExchange(ctx context.Context, req *airxsd.RailExchangeReq) (*airxsd.RailExchangeRsp, error)
	// RailExchangeQuote corresponds to the RailExchangeQuoteReq operation of the Rail service.
	RailExchangeQuote(ctx context.Context, req *airxsd.RailExchangeQuoteReq) (*airxsd.RailExchangeQuoteRsp, error)
	// RailRefund corresponds to the RailRefundReq operation of the Rail service.
	RailRefund(ctx context.Context, req *airxsd.RailRefundReq) (*airxsd.RailRefundRsp, error)
	// RailRefundQuote corresponds to the RailRefundQuoteReq operation of the Rail service.
	RailRefundQuote(ctx context.Context, req *airxsd.RailRefundQuoteReq) (*airxsd.RailRefundQuoteRsp, error)
	// RailSeatMap corresponds to the RailSeatMapReq operation of the Rail service.
	RailSeatMap(ctx context.Context, req *airxsd.RailSeatMapReq) (*airxsd.RailSeatMapRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// RailService is the SOAP implementation of RailServicePort.
type RailService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *RailService must satisfy the RailServicePort interface.
var _ RailServicePort = (*RailService)(nil)

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

// NewRailService builds a Rail service client from the given SOAP configuration and logger.
func NewRailService(config client.SOAPConfig, logger logging.Logger) (*RailService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "rail-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create rail service client: %w", err)
	}

	return &RailService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// RailAvailabilitySearch issues the RailAvailabilitySearchReq SOAP call and returns the strongly typed response.
func (s *RailService) RailAvailabilitySearch(ctx context.Context, req *airxsd.RailAvailabilitySearchReq) (*airxsd.RailAvailabilitySearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.RailAvailabilitySearchRsp](s.client, ctx, "RailAvailabilitySearch", req)
}

// RailExchange issues the RailExchangeReq SOAP call and returns the strongly typed response.
func (s *RailService) RailExchange(ctx context.Context, req *airxsd.RailExchangeReq) (*airxsd.RailExchangeRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.RailExchangeRsp](s.client, ctx, "RailExchange", req)
}

// RailExchangeQuote issues the RailExchangeQuoteReq SOAP call and returns the strongly typed response.
func (s *RailService) RailExchangeQuote(ctx context.Context, req *airxsd.RailExchangeQuoteReq) (*airxsd.RailExchangeQuoteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.RailExchangeQuoteRsp](s.client, ctx, "RailExchangeQuote", req)
}

// RailRefund issues the RailRefundReq SOAP call and returns the strongly typed response.
func (s *RailService) RailRefund(ctx context.Context, req *airxsd.RailRefundReq) (*airxsd.RailRefundRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.RailRefundRsp](s.client, ctx, "RailRefund", req)
}

// RailRefundQuote issues the RailRefundQuoteReq SOAP call and returns the strongly typed response.
func (s *RailService) RailRefundQuote(ctx context.Context, req *airxsd.RailRefundQuoteReq) (*airxsd.RailRefundQuoteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.RailRefundQuoteRsp](s.client, ctx, "RailRefundQuote", req)
}

// RailSeatMap issues the RailSeatMapReq SOAP call and returns the strongly typed response.
func (s *RailService) RailSeatMap(ctx context.Context, req *airxsd.RailSeatMapReq) (*airxsd.RailSeatMapRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.RailSeatMapRsp](s.client, ctx, "RailSeatMap", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType that
// performs a single SOAP call and decodes it into the strongly typed response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *RailService) Close() error {
	return s.client.Close()
}
