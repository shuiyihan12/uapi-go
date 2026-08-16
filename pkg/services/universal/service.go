// Package universal provides the SOAP client implementation for the Travelport Universal
// service (the air-mapped namespace).
// Its UniversalServicePort interface corresponds one-to-one with the WSDL *PortType operations.
// Request/response types come directly from the generator-produced air package
// (script generated; do not edit by hand).
// This package keeps no hand-written request models; all infrastructure fields
// are injected uniformly by prepareReq before sending.
package universal

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	airxsd "github.com/shuiyihan12/uapi-go/pkg/generated/air"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// UniversalServicePort mirrors the *PortType operations of the air-mapped namespace one-to-one,
// making it easy to swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated air package.
type UniversalServicePort interface {
	// AckScheduleChange corresponds to the AckScheduleChangeReq operation of the Universal service.
	AckScheduleChange(ctx context.Context, req *airxsd.AckScheduleChangeReq) (*airxsd.AckScheduleChangeRsp, error)
	// AirCancel corresponds to the AirCancelReq operation of the Universal service.
	AirCancel(ctx context.Context, req *airxsd.AirCancelReq) (*airxsd.AirCancelRsp, error)
	// AirCreateReservation corresponds to the AirCreateReservationReq operation of the Universal service.
	AirCreateReservation(ctx context.Context, req *airxsd.AirCreateReservationReq) (*airxsd.AirCreateReservationRsp, error)
	// AirMerchandisingFulfillment corresponds to the AirMerchandisingFulfillmentReq operation of the Universal service.
	AirMerchandisingFulfillment(ctx context.Context, req *airxsd.AirMerchandisingFulfillmentReq) (*airxsd.AirMerchandisingFulfillmentRsp, error)
	// HotelCancel corresponds to the HotelCancelReq operation of the Universal service.
	HotelCancel(ctx context.Context, req *airxsd.HotelCancelReq) (*airxsd.HotelCancelRsp, error)
	// HotelCreateReservation corresponds to the HotelCreateReservationReq operation of the Universal service.
	HotelCreateReservation(ctx context.Context, req *airxsd.HotelCreateReservationReq) (*airxsd.HotelCreateReservationRsp, error)
	// PassiveCancel corresponds to the PassiveCancelReq operation of the Universal service.
	PassiveCancel(ctx context.Context, req *airxsd.PassiveCancelReq) (*airxsd.PassiveCancelRsp, error)
	// PassiveCreateReservation corresponds to the PassiveCreateReservationReq operation of the Universal service.
	PassiveCreateReservation(ctx context.Context, req *airxsd.PassiveCreateReservationReq) (*airxsd.PassiveCreateReservationRsp, error)
	// ProviderReservationDisplayDetails corresponds to the ProviderReservationDisplayDetailsReq operation of the Universal service.
	ProviderReservationDisplayDetails(ctx context.Context, req *airxsd.ProviderReservationDisplayDetailsReq) (*airxsd.ProviderReservationDisplayDetailsRsp, error)
	// ProviderReservationDivide corresponds to the ProviderReservationDivideReq operation of the Universal service.
	ProviderReservationDivide(ctx context.Context, req *airxsd.ProviderReservationDivideReq) (*airxsd.ProviderReservationDivideRsp, error)
	// RailCreateReservation corresponds to the RailCreateReservationReq operation of the Universal service.
	RailCreateReservation(ctx context.Context, req *airxsd.RailCreateReservationReq) (*airxsd.RailCreateReservationRsp, error)
	// SavedTripCreate corresponds to the SavedTripCreateReq operation of the Universal service.
	SavedTripCreate(ctx context.Context, req *airxsd.SavedTripCreateReq) (*airxsd.SavedTripCreateRsp, error)
	// SavedTripDelete corresponds to the SavedTripDeleteReq operation of the Universal service.
	SavedTripDelete(ctx context.Context, req *airxsd.SavedTripDeleteReq) (*airxsd.SavedTripDeleteRsp, error)
	// SavedTripModify corresponds to the SavedTripModifyReq operation of the Universal service.
	SavedTripModify(ctx context.Context, req *airxsd.SavedTripModifyReq) (*airxsd.SavedTripModifyRsp, error)
	// SavedTripRetrieve corresponds to the SavedTripRetrieveReq operation of the Universal service.
	SavedTripRetrieve(ctx context.Context, req *airxsd.SavedTripRetrieveReq) (*airxsd.SavedTripRetrieveRsp, error)
	// SavedTripSearch corresponds to the SavedTripSearchReq operation of the Universal service.
	SavedTripSearch(ctx context.Context, req *airxsd.SavedTripSearchReq) (*airxsd.SavedTripSearchRsp, error)
	// UniversalRecordCancel corresponds to the UniversalRecordCancelReq operation of the Universal service.
	UniversalRecordCancel(ctx context.Context, req *airxsd.UniversalRecordCancelReq) (*airxsd.UniversalRecordCancelRsp, error)
	// UniversalRecordImport corresponds to the UniversalRecordImportReq operation of the Universal service.
	UniversalRecordImport(ctx context.Context, req *airxsd.UniversalRecordImportReq) (*airxsd.UniversalRecordImportRsp, error)
	// UniversalRecordModify corresponds to the UniversalRecordModifyReq operation of the Universal service.
	UniversalRecordModify(ctx context.Context, req *airxsd.UniversalRecordModifyReq) (*airxsd.UniversalRecordModifyRsp, error)
	// UniversalRecordRetrieve corresponds to the UniversalRecordRetrieveReq operation of the Universal service.
	UniversalRecordRetrieve(ctx context.Context, req *airxsd.UniversalRecordRetrieveReq) (*airxsd.UniversalRecordRetrieveRsp, error)
	// UniversalRecordSearch corresponds to the UniversalRecordSearchReq operation of the Universal service.
	UniversalRecordSearch(ctx context.Context, req *airxsd.UniversalRecordSearchReq) (*airxsd.UniversalRecordSearchRsp, error)
	// VehicleCancel corresponds to the VehicleCancelReq operation of the Universal service.
	VehicleCancel(ctx context.Context, req *airxsd.VehicleCancelReq) (*airxsd.VehicleCancelRsp, error)
	// VehicleCreateReservation corresponds to the VehicleCreateReservationReq operation of the Universal service.
	VehicleCreateReservation(ctx context.Context, req *airxsd.VehicleCreateReservationReq) (*airxsd.VehicleCreateReservationRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// UniversalService is the SOAP implementation of UniversalServicePort.
type UniversalService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *UniversalService must satisfy the UniversalServicePort interface.
var _ UniversalServicePort = (*UniversalService)(nil)

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

// NewUniversalService builds a Universal service client from the given SOAP configuration and logger.
func NewUniversalService(config client.SOAPConfig, logger logging.Logger) (*UniversalService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "universal-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create universal service client: %w", err)
	}

	return &UniversalService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// AckScheduleChange issues the AckScheduleChangeReq SOAP call and returns the strongly typed response.
func (s *UniversalService) AckScheduleChange(ctx context.Context, req *airxsd.AckScheduleChangeReq) (*airxsd.AckScheduleChangeRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AckScheduleChangeRsp](s.client, ctx, "AckScheduleChange", req)
}

// AirCancel issues the AirCancelReq SOAP call and returns the strongly typed response.
func (s *UniversalService) AirCancel(ctx context.Context, req *airxsd.AirCancelReq) (*airxsd.AirCancelRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirCancelRsp](s.client, ctx, "AirCancel", req)
}

// AirCreateReservation issues the AirCreateReservationReq SOAP call and returns the strongly typed response.
func (s *UniversalService) AirCreateReservation(ctx context.Context, req *airxsd.AirCreateReservationReq) (*airxsd.AirCreateReservationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirCreateReservationRsp](s.client, ctx, "AirCreateReservation", req)
}

// AirMerchandisingFulfillment issues the AirMerchandisingFulfillmentReq SOAP call and returns the strongly typed response.
func (s *UniversalService) AirMerchandisingFulfillment(ctx context.Context, req *airxsd.AirMerchandisingFulfillmentReq) (*airxsd.AirMerchandisingFulfillmentRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.AirMerchandisingFulfillmentRsp](s.client, ctx, "AirMerchandisingFulfillment", req)
}

// HotelCancel issues the HotelCancelReq SOAP call and returns the strongly typed response.
func (s *UniversalService) HotelCancel(ctx context.Context, req *airxsd.HotelCancelReq) (*airxsd.HotelCancelRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.HotelCancelRsp](s.client, ctx, "HotelCancel", req)
}

// HotelCreateReservation issues the HotelCreateReservationReq SOAP call and returns the strongly typed response.
func (s *UniversalService) HotelCreateReservation(ctx context.Context, req *airxsd.HotelCreateReservationReq) (*airxsd.HotelCreateReservationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.HotelCreateReservationRsp](s.client, ctx, "HotelCreateReservation", req)
}

// PassiveCancel issues the PassiveCancelReq SOAP call and returns the strongly typed response.
func (s *UniversalService) PassiveCancel(ctx context.Context, req *airxsd.PassiveCancelReq) (*airxsd.PassiveCancelRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.PassiveCancelRsp](s.client, ctx, "PassiveCancel", req)
}

// PassiveCreateReservation issues the PassiveCreateReservationReq SOAP call and returns the strongly typed response.
func (s *UniversalService) PassiveCreateReservation(ctx context.Context, req *airxsd.PassiveCreateReservationReq) (*airxsd.PassiveCreateReservationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.PassiveCreateReservationRsp](s.client, ctx, "PassiveCreateReservation", req)
}

// ProviderReservationDisplayDetails issues the ProviderReservationDisplayDetailsReq SOAP call and returns the strongly typed response.
func (s *UniversalService) ProviderReservationDisplayDetails(ctx context.Context, req *airxsd.ProviderReservationDisplayDetailsReq) (*airxsd.ProviderReservationDisplayDetailsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.ProviderReservationDisplayDetailsRsp](s.client, ctx, "ProviderReservationDisplayDetails", req)
}

// ProviderReservationDivide issues the ProviderReservationDivideReq SOAP call and returns the strongly typed response.
func (s *UniversalService) ProviderReservationDivide(ctx context.Context, req *airxsd.ProviderReservationDivideReq) (*airxsd.ProviderReservationDivideRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.ProviderReservationDivideRsp](s.client, ctx, "ProviderReservationDivide", req)
}

// RailCreateReservation issues the RailCreateReservationReq SOAP call and returns the strongly typed response.
func (s *UniversalService) RailCreateReservation(ctx context.Context, req *airxsd.RailCreateReservationReq) (*airxsd.RailCreateReservationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.RailCreateReservationRsp](s.client, ctx, "RailCreateReservation", req)
}

// SavedTripCreate issues the SavedTripCreateReq SOAP call and returns the strongly typed response.
func (s *UniversalService) SavedTripCreate(ctx context.Context, req *airxsd.SavedTripCreateReq) (*airxsd.SavedTripCreateRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.SavedTripCreateRsp](s.client, ctx, "SavedTripCreate", req)
}

// SavedTripDelete issues the SavedTripDeleteReq SOAP call and returns the strongly typed response.
func (s *UniversalService) SavedTripDelete(ctx context.Context, req *airxsd.SavedTripDeleteReq) (*airxsd.SavedTripDeleteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.SavedTripDeleteRsp](s.client, ctx, "SavedTripDelete", req)
}

// SavedTripModify issues the SavedTripModifyReq SOAP call and returns the strongly typed response.
func (s *UniversalService) SavedTripModify(ctx context.Context, req *airxsd.SavedTripModifyReq) (*airxsd.SavedTripModifyRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.SavedTripModifyRsp](s.client, ctx, "SavedTripModify", req)
}

// SavedTripRetrieve issues the SavedTripRetrieveReq SOAP call and returns the strongly typed response.
func (s *UniversalService) SavedTripRetrieve(ctx context.Context, req *airxsd.SavedTripRetrieveReq) (*airxsd.SavedTripRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.SavedTripRetrieveRsp](s.client, ctx, "SavedTripRetrieve", req)
}

// SavedTripSearch issues the SavedTripSearchReq SOAP call and returns the strongly typed response.
func (s *UniversalService) SavedTripSearch(ctx context.Context, req *airxsd.SavedTripSearchReq) (*airxsd.SavedTripSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.SavedTripSearchRsp](s.client, ctx, "SavedTripSearch", req)
}

// UniversalRecordCancel issues the UniversalRecordCancelReq SOAP call and returns the strongly typed response.
func (s *UniversalService) UniversalRecordCancel(ctx context.Context, req *airxsd.UniversalRecordCancelReq) (*airxsd.UniversalRecordCancelRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.UniversalRecordCancelRsp](s.client, ctx, "UniversalRecordCancel", req)
}

// UniversalRecordImport issues the UniversalRecordImportReq SOAP call and returns the strongly typed response.
func (s *UniversalService) UniversalRecordImport(ctx context.Context, req *airxsd.UniversalRecordImportReq) (*airxsd.UniversalRecordImportRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.UniversalRecordImportRsp](s.client, ctx, "UniversalRecordImport", req)
}

// UniversalRecordModify issues the UniversalRecordModifyReq SOAP call and returns the strongly typed response.
func (s *UniversalService) UniversalRecordModify(ctx context.Context, req *airxsd.UniversalRecordModifyReq) (*airxsd.UniversalRecordModifyRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.UniversalRecordModifyRsp](s.client, ctx, "UniversalRecordModify", req)
}

// UniversalRecordRetrieve issues the UniversalRecordRetrieveReq SOAP call and returns the strongly typed response.
func (s *UniversalService) UniversalRecordRetrieve(ctx context.Context, req *airxsd.UniversalRecordRetrieveReq) (*airxsd.UniversalRecordRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.UniversalRecordRetrieveRsp](s.client, ctx, "UniversalRecordRetrieve", req)
}

// UniversalRecordSearch issues the UniversalRecordSearchReq SOAP call and returns the strongly typed response.
func (s *UniversalService) UniversalRecordSearch(ctx context.Context, req *airxsd.UniversalRecordSearchReq) (*airxsd.UniversalRecordSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.UniversalRecordSearchRsp](s.client, ctx, "UniversalRecordSearch", req)
}

// VehicleCancel issues the VehicleCancelReq SOAP call and returns the strongly typed response.
func (s *UniversalService) VehicleCancel(ctx context.Context, req *airxsd.VehicleCancelReq) (*airxsd.VehicleCancelRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.VehicleCancelRsp](s.client, ctx, "VehicleCancel", req)
}

// VehicleCreateReservation issues the VehicleCreateReservationReq SOAP call and returns the strongly typed response.
func (s *UniversalService) VehicleCreateReservation(ctx context.Context, req *airxsd.VehicleCreateReservationReq) (*airxsd.VehicleCreateReservationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[airxsd.VehicleCreateReservationRsp](s.client, ctx, "VehicleCreateReservation", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType
// that performs a single SOAP call and decodes it into the strongly typed
// response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *UniversalService) Close() error {
	return s.client.Close()
}
