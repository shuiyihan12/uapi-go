// Package air provides the SOAP client implementation for the Travelport Air
// service (the air-mapped namespace). Its AirServicePort interface corresponds
// one-to-one with the WSDL *PortType operations. Request/response types come
// directly from the generator-produced air package (script generated; do not
// edit by hand). This package keeps no hand-written request models; all
// infrastructure fields are injected uniformly by prepareReq before sending.
package air

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	airxsd "github.com/shuiyihan12/uapi-go/pkg/generated/air"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// AirServicePort mirrors the *PortType operations of the air-mapped namespace
// one-to-one, making it easy to swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated air package.
type AirServicePort interface {
	// AirExchange corresponds to the AirExchangeReq operation of the Air service.
	AirExchange(ctx context.Context, req *airxsd.AirExchangeReq) (*airxsd.AirExchangeRsp, error)
	// AirExchangeEligibility corresponds to the AirExchangeEligibilityReq operation of the Air service.
	AirExchangeEligibility(ctx context.Context, req *airxsd.AirExchangeEligibilityReq) (*airxsd.AirExchangeEligibilityRsp, error)
	// AirExchangeMultiQuote corresponds to the AirExchangeMultiQuoteReq operation of the Air service.
	AirExchangeMultiQuote(ctx context.Context, req *airxsd.AirExchangeMultiQuoteReq) (*airxsd.AirExchangeMultiQuoteRsp, error)
	// AirExchangeQuote corresponds to the AirExchangeQuoteReq operation of the Air service.
	AirExchangeQuote(ctx context.Context, req *airxsd.AirExchangeQuoteReq) (*airxsd.AirExchangeQuoteRsp, error)
	// AirExchangeTicketing corresponds to the AirExchangeTicketingReq operation of the Air service.
	AirExchangeTicketing(ctx context.Context, req *airxsd.AirExchangeTicketingReq) (*airxsd.AirExchangeTicketingRsp, error)
	// AirFareDisplay corresponds to the AirFareDisplayReq operation of the Air service.
	AirFareDisplay(ctx context.Context, req *airxsd.AirFareDisplayReq) (*airxsd.AirFareDisplayRsp, error)
	// AirFareRules corresponds to the AirFareRulesReq operation of the Air service.
	AirFareRules(ctx context.Context, req *airxsd.AirFareRulesReq) (*airxsd.AirFareRulesRsp, error)
	// AirMerchandisingDetails corresponds to the AirMerchandisingDetailsReq operation of the Air service.
	AirMerchandisingDetails(ctx context.Context, req *airxsd.AirMerchandisingDetailsReq) (*airxsd.AirMerchandisingDetailsRsp, error)
	// AirMerchandisingOfferAvailability corresponds to the AirMerchandisingOfferAvailabilityReq operation of the Air service.
	AirMerchandisingOfferAvailability(ctx context.Context, req *airxsd.AirMerchandisingOfferAvailabilityReq) (*airxsd.AirMerchandisingOfferAvailabilityRsp, error)
	// AirPrePay corresponds to the AirPrePayReq operation of the Air service.
	AirPrePay(ctx context.Context, req *airxsd.AirPrePayReq) (*airxsd.AirPrePayRsp, error)
	// AirPrice corresponds to the AirPriceReq operation of the Air service.
	AirPrice(ctx context.Context, req *airxsd.AirPriceReq) (*airxsd.AirPriceRsp, error)
	// AirRefund corresponds to the AirRefundReq operation of the Air service.
	AirRefund(ctx context.Context, req *airxsd.AirRefundReq) (*airxsd.AirRefundRsp, error)
	// AirRefundQuote corresponds to the AirRefundQuoteReq operation of the Air service.
	AirRefundQuote(ctx context.Context, req *airxsd.AirRefundQuoteReq) (*airxsd.AirRefundQuoteRsp, error)
	// AirReprice corresponds to the AirRepriceReq operation of the Air service.
	AirReprice(ctx context.Context, req *airxsd.AirRepriceReq) (*airxsd.AirRepriceRsp, error)
	// AirRetrieveDocument corresponds to the AirRetrieveDocumentReq operation of the Air service.
	AirRetrieveDocument(ctx context.Context, req *airxsd.AirRetrieveDocumentReq) (*airxsd.AirRetrieveDocumentRsp, error)
	// AirTicketing corresponds to the AirTicketingReq operation of the Air service.
	AirTicketing(ctx context.Context, req *airxsd.AirTicketingReq) (*airxsd.AirTicketingRsp, error)
	// AirUpsellSearch corresponds to the AirUpsellSearchReq operation of the Air service.
	AirUpsellSearch(ctx context.Context, req *airxsd.AirUpsellSearchReq) (*airxsd.AirUpsellSearchRsp, error)
	// AirVoidDocument corresponds to the AirVoidDocumentReq operation of the Air service.
	AirVoidDocument(ctx context.Context, req *airxsd.AirVoidDocumentReq) (*airxsd.AirVoidDocumentRsp, error)
	// AvailabilitySearch corresponds to the AvailabilitySearchReq operation of the Air service.
	AvailabilitySearch(ctx context.Context, req *airxsd.AvailabilitySearchReq) (*airxsd.AvailabilitySearchRsp, error)
	// EMDIssuance corresponds to the EMDIssuanceReq operation of the Air service.
	EMDIssuance(ctx context.Context, req *airxsd.EMDIssuanceReq) (*airxsd.EMDIssuanceRsp, error)
	// EMDRetrieve corresponds to the EMDRetrieveReq operation of the Air service.
	EMDRetrieve(ctx context.Context, req *airxsd.EMDRetrieveReq) (*airxsd.EMDRetrieveRsp, error)
	// FlightDetails corresponds to the FlightDetailsReq operation of the Air service.
	FlightDetails(ctx context.Context, req *airxsd.FlightDetailsReq) (*airxsd.FlightDetailsRsp, error)
	// FlightInformation corresponds to the FlightInformationReq operation of the Air service.
	FlightInformation(ctx context.Context, req *airxsd.FlightInformationReq) (*airxsd.FlightInformationRsp, error)
	// FlightTimeTable corresponds to the FlightTimeTableReq operation of the Air service.
	FlightTimeTable(ctx context.Context, req *airxsd.FlightTimeTableReq) (*airxsd.FlightTimeTableRsp, error)
	// LowFareSearch corresponds to the LowFareSearchReq operation of the Air service.
	LowFareSearch(ctx context.Context, req *airxsd.LowFareSearchReq) (*airxsd.LowFareSearchRsp, error)
	// ScheduleSearch corresponds to the ScheduleSearchReq operation of the Air service.
	ScheduleSearch(ctx context.Context, req *airxsd.ScheduleSearchReq) (*airxsd.ScheduleSearchRsp, error)
	// SeatMap corresponds to the SeatMapReq operation of the Air service.
	SeatMap(ctx context.Context, req *airxsd.SeatMapReq) (*airxsd.SeatMapRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// AirService is the SOAP implementation of AirServicePort.
type AirService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *AirService must satisfy the AirServicePort interface.
var _ AirServicePort = (*AirService)(nil)

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

// NewAirService builds an Air service client from the given SOAP configuration and logger.
func NewAirService(config client.SOAPConfig, logger logging.Logger) (*AirService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "air-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create air service client: %w", err)
	}

	return &AirService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// AirExchange issues the AirExchangeReq SOAP call and returns the strongly typed response.
func (s *AirService) AirExchange(ctx context.Context, req *airxsd.AirExchangeReq) (*airxsd.AirExchangeRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirExchangeRsp](s.client, ctx, "AirExchange", req)
}

// AirExchangeEligibility issues the AirExchangeEligibilityReq SOAP call and returns the strongly typed response.
func (s *AirService) AirExchangeEligibility(ctx context.Context, req *airxsd.AirExchangeEligibilityReq) (*airxsd.AirExchangeEligibilityRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirExchangeEligibilityRsp](s.client, ctx, "AirExchangeEligibility", req)
}

// AirExchangeMultiQuote issues the AirExchangeMultiQuoteReq SOAP call and returns the strongly typed response.
func (s *AirService) AirExchangeMultiQuote(ctx context.Context, req *airxsd.AirExchangeMultiQuoteReq) (*airxsd.AirExchangeMultiQuoteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirExchangeMultiQuoteRsp](s.client, ctx, "AirExchangeMultiQuote", req)
}

// AirExchangeQuote issues the AirExchangeQuoteReq SOAP call and returns the strongly typed response.
func (s *AirService) AirExchangeQuote(ctx context.Context, req *airxsd.AirExchangeQuoteReq) (*airxsd.AirExchangeQuoteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirExchangeQuoteRsp](s.client, ctx, "AirExchangeQuote", req)
}

// AirExchangeTicketing issues the AirExchangeTicketingReq SOAP call and returns the strongly typed response.
func (s *AirService) AirExchangeTicketing(ctx context.Context, req *airxsd.AirExchangeTicketingReq) (*airxsd.AirExchangeTicketingRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirExchangeTicketingRsp](s.client, ctx, "AirExchangeTicketing", req)
}

// AirFareDisplay issues the AirFareDisplayReq SOAP call and returns the strongly typed response.
func (s *AirService) AirFareDisplay(ctx context.Context, req *airxsd.AirFareDisplayReq) (*airxsd.AirFareDisplayRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirFareDisplayRsp](s.client, ctx, "AirFareDisplay", req)
}

// AirFareRules issues the AirFareRulesReq SOAP call and returns the strongly typed response.
func (s *AirService) AirFareRules(ctx context.Context, req *airxsd.AirFareRulesReq) (*airxsd.AirFareRulesRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirFareRulesRsp](s.client, ctx, "AirFareRules", req)
}

// AirMerchandisingDetails issues the AirMerchandisingDetailsReq SOAP call and returns the strongly typed response.
func (s *AirService) AirMerchandisingDetails(ctx context.Context, req *airxsd.AirMerchandisingDetailsReq) (*airxsd.AirMerchandisingDetailsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirMerchandisingDetailsRsp](s.client, ctx, "AirMerchandisingDetails", req)
}

// AirMerchandisingOfferAvailability issues the AirMerchandisingOfferAvailabilityReq SOAP call and returns the strongly typed response.
func (s *AirService) AirMerchandisingOfferAvailability(ctx context.Context, req *airxsd.AirMerchandisingOfferAvailabilityReq) (*airxsd.AirMerchandisingOfferAvailabilityRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirMerchandisingOfferAvailabilityRsp](s.client, ctx, "AirMerchandisingOfferAvailability", req)
}

// AirPrePay issues the AirPrePayReq SOAP call and returns the strongly typed response.
func (s *AirService) AirPrePay(ctx context.Context, req *airxsd.AirPrePayReq) (*airxsd.AirPrePayRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirPrePayRsp](s.client, ctx, "AirPrePay", req)
}

// AirPrice issues the AirPriceReq SOAP call and returns the strongly typed response.
func (s *AirService) AirPrice(ctx context.Context, req *airxsd.AirPriceReq) (*airxsd.AirPriceRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirPriceRsp](s.client, ctx, "AirPrice", req)
}

// AirRefund issues the AirRefundReq SOAP call and returns the strongly typed response.
func (s *AirService) AirRefund(ctx context.Context, req *airxsd.AirRefundReq) (*airxsd.AirRefundRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirRefundRsp](s.client, ctx, "AirRefund", req)
}

// AirRefundQuote issues the AirRefundQuoteReq SOAP call and returns the strongly typed response.
func (s *AirService) AirRefundQuote(ctx context.Context, req *airxsd.AirRefundQuoteReq) (*airxsd.AirRefundQuoteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirRefundQuoteRsp](s.client, ctx, "AirRefundQuote", req)
}

// AirReprice issues the AirRepriceReq SOAP call and returns the strongly typed response.
func (s *AirService) AirReprice(ctx context.Context, req *airxsd.AirRepriceReq) (*airxsd.AirRepriceRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirRepriceRsp](s.client, ctx, "AirReprice", req)
}

// AirRetrieveDocument issues the AirRetrieveDocumentReq SOAP call and returns the strongly typed response.
func (s *AirService) AirRetrieveDocument(ctx context.Context, req *airxsd.AirRetrieveDocumentReq) (*airxsd.AirRetrieveDocumentRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirRetrieveDocumentRsp](s.client, ctx, "AirRetrieveDocument", req)
}

// AirTicketing issues the AirTicketingReq SOAP call and returns the strongly typed response.
func (s *AirService) AirTicketing(ctx context.Context, req *airxsd.AirTicketingReq) (*airxsd.AirTicketingRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirTicketingRsp](s.client, ctx, "AirTicketing", req)
}

// AirUpsellSearch issues the AirUpsellSearchReq SOAP call and returns the strongly typed response.
func (s *AirService) AirUpsellSearch(ctx context.Context, req *airxsd.AirUpsellSearchReq) (*airxsd.AirUpsellSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirUpsellSearchRsp](s.client, ctx, "AirUpsellSearch", req)
}

// AirVoidDocument issues the AirVoidDocumentReq SOAP call and returns the strongly typed response.
func (s *AirService) AirVoidDocument(ctx context.Context, req *airxsd.AirVoidDocumentReq) (*airxsd.AirVoidDocumentRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirVoidDocumentRsp](s.client, ctx, "AirVoidDocument", req)
}

// AvailabilitySearch issues the AvailabilitySearchReq SOAP call and returns the strongly typed response.
func (s *AirService) AvailabilitySearch(ctx context.Context, req *airxsd.AvailabilitySearchReq) (*airxsd.AvailabilitySearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AvailabilitySearchRsp](s.client, ctx, "AvailabilitySearch", req)
}

// EMDIssuance issues the EMDIssuanceReq SOAP call and returns the strongly typed response.
func (s *AirService) EMDIssuance(ctx context.Context, req *airxsd.EMDIssuanceReq) (*airxsd.EMDIssuanceRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.EMDIssuanceRsp](s.client, ctx, "EMDIssuance", req)
}

// EMDRetrieve issues the EMDRetrieveReq SOAP call and returns the strongly typed response.
func (s *AirService) EMDRetrieve(ctx context.Context, req *airxsd.EMDRetrieveReq) (*airxsd.EMDRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.EMDRetrieveRsp](s.client, ctx, "EMDRetrieve", req)
}

// FlightDetails issues the FlightDetailsReq SOAP call and returns the strongly typed response.
func (s *AirService) FlightDetails(ctx context.Context, req *airxsd.FlightDetailsReq) (*airxsd.FlightDetailsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.FlightDetailsRsp](s.client, ctx, "FlightDetails", req)
}

// FlightInformation issues the FlightInformationReq SOAP call and returns the strongly typed response.
func (s *AirService) FlightInformation(ctx context.Context, req *airxsd.FlightInformationReq) (*airxsd.FlightInformationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.FlightInformationRsp](s.client, ctx, "FlightInformation", req)
}

// FlightTimeTable issues the FlightTimeTableReq SOAP call and returns the strongly typed response.
func (s *AirService) FlightTimeTable(ctx context.Context, req *airxsd.FlightTimeTableReq) (*airxsd.FlightTimeTableRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.FlightTimeTableRsp](s.client, ctx, "FlightTimeTable", req)
}

// LowFareSearch issues the LowFareSearchReq SOAP call and returns the strongly typed response.
func (s *AirService) LowFareSearch(ctx context.Context, req *airxsd.LowFareSearchReq) (*airxsd.LowFareSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.LowFareSearchRsp](s.client, ctx, "LowFareSearch", req)
}

// ScheduleSearch issues the ScheduleSearchReq SOAP call and returns the strongly typed response.
func (s *AirService) ScheduleSearch(ctx context.Context, req *airxsd.ScheduleSearchReq) (*airxsd.ScheduleSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.ScheduleSearchRsp](s.client, ctx, "ScheduleSearch", req)
}

// SeatMap issues the SeatMapReq SOAP call and returns the strongly typed response.
func (s *AirService) SeatMap(ctx context.Context, req *airxsd.SeatMapReq) (*airxsd.SeatMapRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.SeatMapRsp](s.client, ctx, "SeatMap", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType that
// performs a single SOAP call and decodes it into the strongly typed response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *AirService) Close() error {
	return s.client.Close()
}
