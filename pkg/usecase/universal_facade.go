// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the Universal service (air): methods map 1:1 to the UniversalServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	airxsd "github.com/shuiyihan12/uapi-go/pkg/generated/air"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	universalsvc "github.com/shuiyihan12/uapi-go/pkg/services/universal"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// UniversalFacade orchestrates the Travelport Universal service use cases; methods map 1:1 to the service PortType.niversalServicePort。
type UniversalFacade struct {
	manager *manager.ServiceManager
}

// NewUniversalFacade creates the Universal use-case layer.
func NewUniversalFacade(serviceManager *manager.ServiceManager) *UniversalFacade {
	return &UniversalFacade{manager: serviceManager}
}

// getService lazily retrieves the Universal service client, handling nil manager and lookup failures uniformly.
func (f *UniversalFacade) getService() (*universalsvc.UniversalService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*universalsvc.UniversalService](f.manager, "universal")
	if err != nil {
		return nil, fmt.Errorf("failed to get universal service: %w", err)
	}
	return svc, nil
}

// AckScheduleChange corresponds to UniversalServicePort.AckScheduleChange。
func (f *UniversalFacade) AckScheduleChange(ctx context.Context, req *airxsd.AckScheduleChangeReq) (*airxsd.AckScheduleChangeRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AckScheduleChange(ctx, req)
}

// AirCancel corresponds to UniversalServicePort.AirCancel。
func (f *UniversalFacade) AirCancel(ctx context.Context, req *airxsd.AirCancelReq) (*airxsd.AirCancelRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirCancel(ctx, req)
}

// AirCreateReservation corresponds to UniversalServicePort.AirCreateReservation。
func (f *UniversalFacade) AirCreateReservation(ctx context.Context, req *airxsd.AirCreateReservationReq) (*airxsd.AirCreateReservationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirCreateReservation(ctx, req)
}

// AirMerchandisingFulfillment corresponds to UniversalServicePort.AirMerchandisingFulfillment。
func (f *UniversalFacade) AirMerchandisingFulfillment(ctx context.Context, req *airxsd.AirMerchandisingFulfillmentReq) (*airxsd.AirMerchandisingFulfillmentRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirMerchandisingFulfillment(ctx, req)
}

// HotelCancel corresponds to UniversalServicePort.HotelCancel。
func (f *UniversalFacade) HotelCancel(ctx context.Context, req *airxsd.HotelCancelReq) (*airxsd.HotelCancelRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.HotelCancel(ctx, req)
}

// HotelCreateReservation corresponds to UniversalServicePort.HotelCreateReservation。
func (f *UniversalFacade) HotelCreateReservation(ctx context.Context, req *airxsd.HotelCreateReservationReq) (*airxsd.HotelCreateReservationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.HotelCreateReservation(ctx, req)
}

// PassiveCancel corresponds to UniversalServicePort.PassiveCancel。
func (f *UniversalFacade) PassiveCancel(ctx context.Context, req *airxsd.PassiveCancelReq) (*airxsd.PassiveCancelRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.PassiveCancel(ctx, req)
}

// PassiveCreateReservation corresponds to UniversalServicePort.PassiveCreateReservation。
func (f *UniversalFacade) PassiveCreateReservation(ctx context.Context, req *airxsd.PassiveCreateReservationReq) (*airxsd.PassiveCreateReservationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.PassiveCreateReservation(ctx, req)
}

// ProviderReservationDisplayDetails corresponds to UniversalServicePort.ProviderReservationDisplayDetails。
func (f *UniversalFacade) ProviderReservationDisplayDetails(ctx context.Context, req *airxsd.ProviderReservationDisplayDetailsReq) (*airxsd.ProviderReservationDisplayDetailsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProviderReservationDisplayDetails(ctx, req)
}

// ProviderReservationDivide corresponds to UniversalServicePort.ProviderReservationDivide。
func (f *UniversalFacade) ProviderReservationDivide(ctx context.Context, req *airxsd.ProviderReservationDivideReq) (*airxsd.ProviderReservationDivideRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProviderReservationDivide(ctx, req)
}

// RailCreateReservation corresponds to UniversalServicePort.RailCreateReservation。
func (f *UniversalFacade) RailCreateReservation(ctx context.Context, req *airxsd.RailCreateReservationReq) (*airxsd.RailCreateReservationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.RailCreateReservation(ctx, req)
}

// SavedTripCreate corresponds to UniversalServicePort.SavedTripCreate。
func (f *UniversalFacade) SavedTripCreate(ctx context.Context, req *airxsd.SavedTripCreateReq) (*airxsd.SavedTripCreateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.SavedTripCreate(ctx, req)
}

// SavedTripDelete corresponds to UniversalServicePort.SavedTripDelete。
func (f *UniversalFacade) SavedTripDelete(ctx context.Context, req *airxsd.SavedTripDeleteReq) (*airxsd.SavedTripDeleteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.SavedTripDelete(ctx, req)
}

// SavedTripModify corresponds to UniversalServicePort.SavedTripModify。
func (f *UniversalFacade) SavedTripModify(ctx context.Context, req *airxsd.SavedTripModifyReq) (*airxsd.SavedTripModifyRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.SavedTripModify(ctx, req)
}

// SavedTripRetrieve corresponds to UniversalServicePort.SavedTripRetrieve。
func (f *UniversalFacade) SavedTripRetrieve(ctx context.Context, req *airxsd.SavedTripRetrieveReq) (*airxsd.SavedTripRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.SavedTripRetrieve(ctx, req)
}

// SavedTripSearch corresponds to UniversalServicePort.SavedTripSearch。
func (f *UniversalFacade) SavedTripSearch(ctx context.Context, req *airxsd.SavedTripSearchReq) (*airxsd.SavedTripSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.SavedTripSearch(ctx, req)
}

// UniversalRecordCancel corresponds to UniversalServicePort.UniversalRecordCancel。
func (f *UniversalFacade) UniversalRecordCancel(ctx context.Context, req *airxsd.UniversalRecordCancelReq) (*airxsd.UniversalRecordCancelRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UniversalRecordCancel(ctx, req)
}

// UniversalRecordImport corresponds to UniversalServicePort.UniversalRecordImport。
func (f *UniversalFacade) UniversalRecordImport(ctx context.Context, req *airxsd.UniversalRecordImportReq) (*airxsd.UniversalRecordImportRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UniversalRecordImport(ctx, req)
}

// UniversalRecordModify corresponds to UniversalServicePort.UniversalRecordModify。
func (f *UniversalFacade) UniversalRecordModify(ctx context.Context, req *airxsd.UniversalRecordModifyReq) (*airxsd.UniversalRecordModifyRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UniversalRecordModify(ctx, req)
}

// UniversalRecordRetrieve corresponds to UniversalServicePort.UniversalRecordRetrieve。
func (f *UniversalFacade) UniversalRecordRetrieve(ctx context.Context, req *airxsd.UniversalRecordRetrieveReq) (*airxsd.UniversalRecordRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UniversalRecordRetrieve(ctx, req)
}

// UniversalRecordSearch corresponds to UniversalServicePort.UniversalRecordSearch。
func (f *UniversalFacade) UniversalRecordSearch(ctx context.Context, req *airxsd.UniversalRecordSearchReq) (*airxsd.UniversalRecordSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UniversalRecordSearch(ctx, req)
}

// VehicleCancel corresponds to UniversalServicePort.VehicleCancel。
func (f *UniversalFacade) VehicleCancel(ctx context.Context, req *airxsd.VehicleCancelReq) (*airxsd.VehicleCancelRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleCancel(ctx, req)
}

// VehicleCreateReservation corresponds to UniversalServicePort.VehicleCreateReservation。
func (f *UniversalFacade) VehicleCreateReservation(ctx context.Context, req *airxsd.VehicleCreateReservationReq) (*airxsd.VehicleCreateReservationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleCreateReservation(ctx, req)
}
