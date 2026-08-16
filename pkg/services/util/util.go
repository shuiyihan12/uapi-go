// Package util provides the SOAP client implementation for the Travelport Util
// (utility) service (util_v55_0). Its UtilServicePort interface corresponds
// one-to-one with the *PortTypes in the WSDL, covering all utility operations
// such as tax calculation, currency conversion, MCO, credit card authorization,
// and reference data. All methods return the WSDL-generated strongly typed
// responses, so callers no longer need to hand-write operation strings.
package util

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	utilxsd "github.com/shuiyihan12/uapi-go/pkg/generated/util"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// UtilServicePort mirrors the *PortType operations of util_v55_0 one-to-one, making it easy to
// swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated util package
// (script generated; do not edit by hand).
type UtilServicePort interface {
	// AgencyServiceFeeCreate corresponds to AgencyCreateServiceFeePortType and creates an agency service fee.
	AgencyServiceFeeCreate(ctx context.Context, req *utilxsd.AgencyServiceFeeCreateReq) (*utilxsd.AgencyServiceFeeCreateRsp, error)
	// BrandedFareAdmin corresponds to BrandedFareAdminPortType and manages branded Fares.
	BrandedFareAdmin(ctx context.Context, req *utilxsd.BrandedFareAdminReq) (*utilxsd.BrandedFareAdminRsp, error)
	// BrandedFareSearch corresponds to BrandedFareSearchPortType and searches branded Fares.
	BrandedFareSearch(ctx context.Context, req *utilxsd.BrandedFareSearchReq) (*utilxsd.BrandedFareSearchRsp, error)
	// CalculateTax corresponds to CalculateTaxPortType and calculates taxes.
	CalculateTax(ctx context.Context, req *utilxsd.CalculateTaxReq) (*utilxsd.CalculateTaxRsp, error)
	// ContentProviderRetrieve corresponds to ContentProviderRetrievePortType and retrieves content provider data.
	ContentProviderRetrieve(ctx context.Context, req *utilxsd.ContentProviderRetrieveReq) (*utilxsd.ContentProviderRetrieveRsp, error)
	// CreateAgencyFeeMco corresponds to McoCreateAgencyFeePortType and creates an agency fee MCO.
	CreateAgencyFeeMco(ctx context.Context, req *utilxsd.CreateAgencyFeeMcoReq) (*utilxsd.CreateAgencyFeeMcoRsp, error)
	// CreateAirlineFeeMco corresponds to the airline fee variant of McoCreateAgencyFeePortType and creates an airline fee MCO.
	CreateAirlineFeeMco(ctx context.Context, req *utilxsd.CreateAirlineFeeMcoReq) (*utilxsd.CreateAirlineFeeMcoRsp, error)
	// CreditCardAuth corresponds to UtilCreditCardAuthPortType and performs credit card authorization.
	CreditCardAuth(ctx context.Context, req *utilxsd.CreditCardAuthReq) (*utilxsd.CreditCardAuthRsp, error)
	// CurrencyConversion corresponds to UtilCurrencyConversionPortType and performs currency conversion.
	CurrencyConversion(ctx context.Context, req *utilxsd.CurrencyConversionReq) (*utilxsd.CurrencyConversionRsp, error)
	// FindEmployeesOnFlight corresponds to FindEmployeesOnFlightServicePortType and finds employees on a flight.
	FindEmployeesOnFlight(ctx context.Context, req *utilxsd.FindEmployeesOnFlightReq) (*utilxsd.FindEmployeesOnFlightRsp, error)
	// MCOCreate corresponds to MCOCreatePortType and creates an MCO.
	MCOCreate(ctx context.Context, req *utilxsd.MCOCreateReq) (*utilxsd.MCOCreateRsp, error)
	// MCOExchange corresponds to MCOExchangePortType and exchanges an MCO.
	MCOExchange(ctx context.Context, req *utilxsd.MCOExchangeReq) (*utilxsd.MCOExchangeRsp, error)
	// MCOIssue corresponds to MCOIssuePortType and issues an MCO.
	MCOIssue(ctx context.Context, req *utilxsd.MCOIssueReq) (*utilxsd.MCOIssueRsp, error)
	// MCORetrieve corresponds to MCORetrievePortType and retrieves an MCO.
	MCORetrieve(ctx context.Context, req *utilxsd.MCORetrieveReq) (*utilxsd.MCORetrieveRsp, error)
	// McoSearch corresponds to McoSearchPortType and searches MCOs.
	McoSearch(ctx context.Context, req *utilxsd.McoSearchReq) (*utilxsd.McoSearchRsp, error)
	// McoVoid corresponds to McoVoidPortType and voids an MCO.
	McoVoid(ctx context.Context, req *utilxsd.McoVoidReq) (*utilxsd.McoVoidRsp, error)
	// MctCount corresponds to MctCountPortType and counts minimum connect times.
	MctCount(ctx context.Context, req *utilxsd.MctCountReq) (*utilxsd.MctCountRsp, error)
	// MctLookup corresponds to MctLookupPortType and looks up minimum connect times.
	MctLookup(ctx context.Context, req *utilxsd.MctLookupReq) (*utilxsd.MctLookupRsp, error)
	// MirReportRetrieve corresponds to ReportRetrievePortType and retrieves MIR reports.
	MirReportRetrieve(ctx context.Context, req *utilxsd.MirReportRetrieveReq) (*utilxsd.MirReportRetrieveRsp, error)
	// ReferenceDataRetrieve corresponds to ReferenceDataRetrievePortType and retrieves reference master data.
	ReferenceDataRetrieve(ctx context.Context, req *utilxsd.ReferenceDataRetrieveReq) (*utilxsd.ReferenceDataRetrieveRsp, error)
	// ReferenceDataSearch corresponds to ReferenceDataLookupPortType and searches reference master data.
	ReferenceDataSearch(ctx context.Context, req *utilxsd.ReferenceDataSearchReq) (*utilxsd.ReferenceDataSearchRsp, error)
	// ReferenceDataUpdate corresponds to ReferenceDataUpdatePortType and updates reference master data.
	ReferenceDataUpdate(ctx context.Context, req *utilxsd.ReferenceDataUpdateReq) (*utilxsd.ReferenceDataUpdateRsp, error)
	// UpsellAdmin corresponds to UpsellAdminPortType and manages upsell rules.
	UpsellAdmin(ctx context.Context, req *utilxsd.UpsellAdminReq) (*utilxsd.UpsellAdminRsp, error)
	// UpsellSearch corresponds to UpsellAdminSearchPortType and searches upsell offers.
	UpsellSearch(ctx context.Context, req *utilxsd.UpsellSearchReq) (*utilxsd.UpsellSearchRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// UtilService is the SOAP implementation of UtilServicePort.
type UtilService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *UtilService must satisfy the UtilServicePort interface.
var _ UtilServicePort = (*UtilService)(nil)

// infraInjectable describes request types whose trace identifier can be
// injected by the service side before sending.
type infraInjectable interface {
	InjectInfrastructure(traceID string)
}

// prepareReq ensures the call carries a trace_id and injects one as a fallback
// when the request body's TraceId is empty. Authorization/business fields such
// as the billing point of sale and TargetBranch are no longer injected by code;
// the caller (API user) must provide them explicitly in the request body.
// Authentication (Authorization) is supplied by the caller as an HTTP request
// header and passed through the context to the SOAP call (startup-time
// environment variables are no longer used). InjectInfrastructure only fills
// the request body's TraceId with the trace_id when it is empty; a business
// value already set is not overwritten.
func prepareReq[T infraInjectable](ctx context.Context, req T) context.Context {
	_, traceID := trace.Ensure(ctx)
	req.InjectInfrastructure(traceID)
	return ctx
}

// NewUtilService builds a Util service client from the given SOAP configuration and logger.
func NewUtilService(config client.SOAPConfig, logger logging.Logger) (*UtilService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "util-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create util service client: %w", err)
	}

	return &UtilService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// AgencyServiceFeeCreate issues the AgencyCreateServiceFee SOAP call and returns the response.
func (s *UtilService) AgencyServiceFeeCreate(ctx context.Context, req *utilxsd.AgencyServiceFeeCreateReq) (*utilxsd.AgencyServiceFeeCreateRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.AgencyServiceFeeCreateRsp](s.client, ctx, "AgencyServiceFeeCreate", req)
}

// BrandedFareAdmin issues the BrandedFareAdmin SOAP call and returns the response.
func (s *UtilService) BrandedFareAdmin(ctx context.Context, req *utilxsd.BrandedFareAdminReq) (*utilxsd.BrandedFareAdminRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.BrandedFareAdminRsp](s.client, ctx, "BrandedFareAdmin", req)
}

// BrandedFareSearch issues the BrandedFareSearch SOAP call and returns the response.
func (s *UtilService) BrandedFareSearch(ctx context.Context, req *utilxsd.BrandedFareSearchReq) (*utilxsd.BrandedFareSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.BrandedFareSearchRsp](s.client, ctx, "BrandedFareSearch", req)
}

// CalculateTax issues the CalculateTax SOAP call and returns the response.
func (s *UtilService) CalculateTax(ctx context.Context, req *utilxsd.CalculateTaxReq) (*utilxsd.CalculateTaxRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.CalculateTaxRsp](s.client, ctx, "CalculateTax", req)
}

// ContentProviderRetrieve issues the ContentProviderRetrieve SOAP call and returns the response.
func (s *UtilService) ContentProviderRetrieve(ctx context.Context, req *utilxsd.ContentProviderRetrieveReq) (*utilxsd.ContentProviderRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.ContentProviderRetrieveRsp](s.client, ctx, "ContentProviderRetrieve", req)
}

// CreateAgencyFeeMco issues the CreateAgencyFeeMco SOAP call and returns the response.
func (s *UtilService) CreateAgencyFeeMco(ctx context.Context, req *utilxsd.CreateAgencyFeeMcoReq) (*utilxsd.CreateAgencyFeeMcoRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.CreateAgencyFeeMcoRsp](s.client, ctx, "CreateAgencyFeeMco", req)
}

// CreateAirlineFeeMco issues the CreateAirlineFeeMco SOAP call and returns the response.
func (s *UtilService) CreateAirlineFeeMco(ctx context.Context, req *utilxsd.CreateAirlineFeeMcoReq) (*utilxsd.CreateAirlineFeeMcoRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.CreateAirlineFeeMcoRsp](s.client, ctx, "CreateAirlineFeeMco", req)
}

// CreditCardAuth issues the CreditCardAuth SOAP call and returns the response.
func (s *UtilService) CreditCardAuth(ctx context.Context, req *utilxsd.CreditCardAuthReq) (*utilxsd.CreditCardAuthRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.CreditCardAuthRsp](s.client, ctx, "CreditCardAuth", req)
}

// CurrencyConversion issues the CurrencyConversion SOAP call and returns the response.
func (s *UtilService) CurrencyConversion(ctx context.Context, req *utilxsd.CurrencyConversionReq) (*utilxsd.CurrencyConversionRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.CurrencyConversionRsp](s.client, ctx, "CurrencyConversion", req)
}

// FindEmployeesOnFlight issues the FindEmployeesOnFlight SOAP call and returns the response.
func (s *UtilService) FindEmployeesOnFlight(ctx context.Context, req *utilxsd.FindEmployeesOnFlightReq) (*utilxsd.FindEmployeesOnFlightRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.FindEmployeesOnFlightRsp](s.client, ctx, "FindEmployeesOnFlight", req)
}

// MCOCreate issues the MCOCreate SOAP call and returns the response.
func (s *UtilService) MCOCreate(ctx context.Context, req *utilxsd.MCOCreateReq) (*utilxsd.MCOCreateRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.MCOCreateRsp](s.client, ctx, "MCOCreate", req)
}

// MCOExchange issues the MCOExchange SOAP call and returns the response.
func (s *UtilService) MCOExchange(ctx context.Context, req *utilxsd.MCOExchangeReq) (*utilxsd.MCOExchangeRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.MCOExchangeRsp](s.client, ctx, "MCOExchange", req)
}

// MCOIssue issues the MCOIssue SOAP call and returns the response.
func (s *UtilService) MCOIssue(ctx context.Context, req *utilxsd.MCOIssueReq) (*utilxsd.MCOIssueRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.MCOIssueRsp](s.client, ctx, "MCOIssue", req)
}

// MCORetrieve issues the MCORetrieve SOAP call and returns the response.
func (s *UtilService) MCORetrieve(ctx context.Context, req *utilxsd.MCORetrieveReq) (*utilxsd.MCORetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.MCORetrieveRsp](s.client, ctx, "MCORetrieve", req)
}

// McoSearch issues the McoSearch SOAP call and returns the response.
func (s *UtilService) McoSearch(ctx context.Context, req *utilxsd.McoSearchReq) (*utilxsd.McoSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.McoSearchRsp](s.client, ctx, "McoSearch", req)
}

// McoVoid issues the McoVoid SOAP call and returns the response.
func (s *UtilService) McoVoid(ctx context.Context, req *utilxsd.McoVoidReq) (*utilxsd.McoVoidRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.McoVoidRsp](s.client, ctx, "McoVoid", req)
}

// MctCount issues the MctCount SOAP call and returns the response.
func (s *UtilService) MctCount(ctx context.Context, req *utilxsd.MctCountReq) (*utilxsd.MctCountRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.MctCountRsp](s.client, ctx, "MctCount", req)
}

// MctLookup issues the MctLookup SOAP call and returns the response.
func (s *UtilService) MctLookup(ctx context.Context, req *utilxsd.MctLookupReq) (*utilxsd.MctLookupRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.MctLookupRsp](s.client, ctx, "MctLookup", req)
}

// MirReportRetrieve issues the MirReportRetrieve SOAP call and returns the response.
func (s *UtilService) MirReportRetrieve(ctx context.Context, req *utilxsd.MirReportRetrieveReq) (*utilxsd.MirReportRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.MirReportRetrieveRsp](s.client, ctx, "MirReportRetrieve", req)
}

// ReferenceDataRetrieve issues the ReferenceDataRetrieve SOAP call and returns the response.
func (s *UtilService) ReferenceDataRetrieve(ctx context.Context, req *utilxsd.ReferenceDataRetrieveReq) (*utilxsd.ReferenceDataRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.ReferenceDataRetrieveRsp](s.client, ctx, "ReferenceDataRetrieve", req)
}

// ReferenceDataSearch issues the ReferenceDataSearch SOAP call and returns the response.
func (s *UtilService) ReferenceDataSearch(ctx context.Context, req *utilxsd.ReferenceDataSearchReq) (*utilxsd.ReferenceDataSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.ReferenceDataSearchRsp](s.client, ctx, "ReferenceDataSearch", req)
}

// ReferenceDataUpdate issues the ReferenceDataUpdate SOAP call and returns the response.
func (s *UtilService) ReferenceDataUpdate(ctx context.Context, req *utilxsd.ReferenceDataUpdateReq) (*utilxsd.ReferenceDataUpdateRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.ReferenceDataUpdateRsp](s.client, ctx, "ReferenceDataUpdate", req)
}

// UpsellAdmin issues the UpsellAdmin SOAP call and returns the response.
func (s *UtilService) UpsellAdmin(ctx context.Context, req *utilxsd.UpsellAdminReq) (*utilxsd.UpsellAdminRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.UpsellAdminRsp](s.client, ctx, "UpsellAdmin", req)
}

// UpsellSearch issues the UpsellSearch SOAP call and returns the response.
func (s *UtilService) UpsellSearch(ctx context.Context, req *utilxsd.UpsellSearchReq) (*utilxsd.UpsellSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return client.CallPortType[utilxsd.UpsellSearchRsp](s.client, ctx, "UpsellSearch", req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *UtilService) Close() error {
	return s.client.Close()
}
