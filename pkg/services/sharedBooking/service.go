// Package sharedBooking provides the SOAP client implementation for the Travelport SharedBooking
// service (the sharedbooking-mapped namespace).
// Its SharedBookingServicePort interface corresponds one-to-one with the WSDL *PortType operations.
// Request/response types come directly from the generator-produced sharedbooking package
// (script generated; do not edit by hand).
// This package keeps no hand-written request models; all infrastructure fields
// are injected uniformly by prepareReq before sending.
package sharedBooking

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	sharedbookingxsd "github.com/shuiyihan12/uapi-go/pkg/generated/sharedbooking"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// SharedBookingServicePort mirrors the *PortType operations of the sharedbooking-mapped namespace one-to-one,
// making it easy to swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated sharedbooking package.
type SharedBookingServicePort interface {
	// BookingAirExchange corresponds to the BookingAirExchangeReq operation of the SharedBooking service.
	BookingAirExchange(ctx context.Context, req *sharedbookingxsd.BookingAirExchangeReq) (*sharedbookingxsd.BookingAirExchangeRsp, error)
	// BookingAirExchangeQuote corresponds to the BookingAirExchangeQuoteReq operation of the SharedBooking service.
	BookingAirExchangeQuote(ctx context.Context, req *sharedbookingxsd.BookingAirExchangeQuoteReq) (*sharedbookingxsd.BookingAirExchangeQuoteRsp, error)
	// BookingAirPnrElement corresponds to the BookingAirPnrElementReq operation of the SharedBooking service.
	BookingAirPnrElement(ctx context.Context, req *sharedbookingxsd.BookingAirPnrElementReq) (*struct{}, error)
	// BookingAirSegment corresponds to the BookingAirSegmentReq operation of the SharedBooking service.
	BookingAirSegment(ctx context.Context, req *sharedbookingxsd.BookingAirSegmentReq) (*sharedbookingxsd.BookingAirSegmentRsp, error)
	// BookingDisplay corresponds to the BookingDisplayReq operation of the SharedBooking service.
	BookingDisplay(ctx context.Context, req *sharedbookingxsd.BookingDisplayReq) (*struct{}, error)
	// BookingEnd corresponds to the BookingEndReq operation of the SharedBooking service.
	BookingEnd(ctx context.Context, req *sharedbookingxsd.BookingEndReq) (*struct{}, error)
	// BookingHotelPnrElement corresponds to the BookingHotelPnrElementReq operation of the SharedBooking service.
	BookingHotelPnrElement(ctx context.Context, req *sharedbookingxsd.BookingHotelPnrElementReq) (*struct{}, error)
	// BookingHotelSegment corresponds to the BookingHotelSegmentReq operation of the SharedBooking service.
	BookingHotelSegment(ctx context.Context, req *sharedbookingxsd.BookingHotelSegmentReq) (*struct{}, error)
	// BookingPnrElement corresponds to the BookingPnrElementReq operation of the SharedBooking service.
	BookingPnrElement(ctx context.Context, req *sharedbookingxsd.BookingPnrElementReq) (*struct{}, error)
	// BookingPricing corresponds to the BookingPricingReq operation of the SharedBooking service.
	BookingPricing(ctx context.Context, req *sharedbookingxsd.BookingPricingReq) (*struct{}, error)
	// BookingRetrieveDocument corresponds to the BookingRetrieveDocumentReq operation of the SharedBooking service.
	BookingRetrieveDocument(ctx context.Context, req *sharedbookingxsd.BookingRetrieveDocumentReq) (*sharedbookingxsd.BookingRetrieveDocumentRsp, error)
	// BookingSeatAssignment corresponds to the BookingSeatAssignmentReq operation of the SharedBooking service.
	BookingSeatAssignment(ctx context.Context, req *sharedbookingxsd.BookingSeatAssignmentReq) (*struct{}, error)
	// BookingStart corresponds to the BookingStartReq operation of the SharedBooking service.
	BookingStart(ctx context.Context, req *sharedbookingxsd.BookingStartReq) (*struct{}, error)
	// BookingTerminal corresponds to the BookingTerminalReq operation of the SharedBooking service.
	BookingTerminal(ctx context.Context, req *sharedbookingxsd.BookingTerminalReq) (*sharedbookingxsd.BookingTerminalRsp, error)
	// BookingTraveler corresponds to the BookingTravelerReq operation of the SharedBooking service.
	BookingTraveler(ctx context.Context, req *sharedbookingxsd.BookingTravelerReq) (*struct{}, error)
	// BookingVehiclePnrElement corresponds to BookingVehiclePnrElementPortType (new
	// in v55) and adds/updates/deletes vehicle elements on the shared booking PNR.
	BookingVehiclePnrElement(ctx context.Context, req *sharedbookingxsd.BookingVehiclePnrElementReq) (*sharedbookingxsd.BookingVehiclePnrElementRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// SharedBookingService is the SOAP implementation of SharedBookingServicePort.
type SharedBookingService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *SharedBookingService must satisfy the SharedBookingServicePort interface.
var _ SharedBookingServicePort = (*SharedBookingService)(nil)

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

// NewSharedBookingService builds a SharedBooking service client from the given SOAP configuration and logger.
func NewSharedBookingService(config client.SOAPConfig, logger logging.Logger) (*SharedBookingService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "sharedBooking-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create sharedBooking service client: %w", err)
	}

	return &SharedBookingService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// BookingAirExchange issues the BookingAirExchangeReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingAirExchange(ctx context.Context, req *sharedbookingxsd.BookingAirExchangeReq) (*sharedbookingxsd.BookingAirExchangeRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[sharedbookingxsd.BookingAirExchangeRsp](s.client, ctx, "BookingAirExchange", req)
}

// BookingAirExchangeQuote issues the BookingAirExchangeQuoteReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingAirExchangeQuote(ctx context.Context, req *sharedbookingxsd.BookingAirExchangeQuoteReq) (*sharedbookingxsd.BookingAirExchangeQuoteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[sharedbookingxsd.BookingAirExchangeQuoteRsp](s.client, ctx, "BookingAirExchangeQuote", req)
}

// BookingAirPnrElement issues the BookingAirPnrElementReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingAirPnrElement(ctx context.Context, req *sharedbookingxsd.BookingAirPnrElementReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingAirPnrElement", req)
}

// BookingAirSegment issues the BookingAirSegmentReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingAirSegment(ctx context.Context, req *sharedbookingxsd.BookingAirSegmentReq) (*sharedbookingxsd.BookingAirSegmentRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[sharedbookingxsd.BookingAirSegmentRsp](s.client, ctx, "BookingAirSegment", req)
}

// BookingDisplay issues the BookingDisplayReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingDisplay(ctx context.Context, req *sharedbookingxsd.BookingDisplayReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingDisplay", req)
}

// BookingEnd issues the BookingEndReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingEnd(ctx context.Context, req *sharedbookingxsd.BookingEndReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingEnd", req)
}

// BookingHotelPnrElement issues the BookingHotelPnrElementReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingHotelPnrElement(ctx context.Context, req *sharedbookingxsd.BookingHotelPnrElementReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingHotelPnrElement", req)
}

// BookingHotelSegment issues the BookingHotelSegmentReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingHotelSegment(ctx context.Context, req *sharedbookingxsd.BookingHotelSegmentReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingHotelSegment", req)
}

// BookingPnrElement issues the BookingPnrElementReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingPnrElement(ctx context.Context, req *sharedbookingxsd.BookingPnrElementReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingPnrElement", req)
}

// BookingPricing issues the BookingPricingReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingPricing(ctx context.Context, req *sharedbookingxsd.BookingPricingReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingPricing", req)
}

// BookingRetrieveDocument issues the BookingRetrieveDocumentReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingRetrieveDocument(ctx context.Context, req *sharedbookingxsd.BookingRetrieveDocumentReq) (*sharedbookingxsd.BookingRetrieveDocumentRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[sharedbookingxsd.BookingRetrieveDocumentRsp](s.client, ctx, "BookingRetrieveDocument", req)
}

// BookingSeatAssignment issues the BookingSeatAssignmentReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingSeatAssignment(ctx context.Context, req *sharedbookingxsd.BookingSeatAssignmentReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingSeatAssignment", req)
}

// BookingStart issues the BookingStartReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingStart(ctx context.Context, req *sharedbookingxsd.BookingStartReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingStart", req)
}

// BookingTerminal issues the BookingTerminalReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingTerminal(ctx context.Context, req *sharedbookingxsd.BookingTerminalReq) (*sharedbookingxsd.BookingTerminalRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[sharedbookingxsd.BookingTerminalRsp](s.client, ctx, "BookingTerminal", req)
}

// BookingTraveler issues the BookingTravelerReq SOAP call and returns the strongly typed response.
func (s *SharedBookingService) BookingTraveler(ctx context.Context, req *sharedbookingxsd.BookingTravelerReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "BookingTraveler", req)
}

// BookingVehiclePnrElement issues the BookingVehiclePnrElementReq SOAP call (a
// PortType new in v55) and returns the strongly typed response carrying
// VehicleRateChangedInfo.
func (s *SharedBookingService) BookingVehiclePnrElement(ctx context.Context, req *sharedbookingxsd.BookingVehiclePnrElementReq) (*sharedbookingxsd.BookingVehiclePnrElementRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[sharedbookingxsd.BookingVehiclePnrElementRsp](s.client, ctx, "BookingVehiclePnrElement", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType
// that performs a single SOAP call and decodes it into the strongly typed
// response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *SharedBookingService) Close() error {
	return s.client.Close()
}
