// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the Vehicle service (vehicle): methods map 1:1 to the VehicleServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	vehiclexsd "github.com/shuiyihan12/uapi-go/pkg/generated/vehicle"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	vehiclesvc "github.com/shuiyihan12/uapi-go/pkg/services/vehicle"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// VehicleFacade orchestrates the Travelport Vehicle service use cases; methods map 1:1 to the service PortType.
type VehicleFacade struct {
	manager *manager.ServiceManager
}

// NewVehicleFacade creates the Vehicle use-case layer.
func NewVehicleFacade(serviceManager *manager.ServiceManager) *VehicleFacade {
	return &VehicleFacade{manager: serviceManager}
}

// getService lazily retrieves the Vehicle service client, handling nil manager and lookup failures uniformly.
func (f *VehicleFacade) getService() (*vehiclesvc.VehicleService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*vehiclesvc.VehicleService](f.manager, "vehicle")
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle service: %w", err)
	}
	return svc, nil
}

// VehicleKeyword corresponds to VehicleServicePort.VehicleKeyword。
func (f *VehicleFacade) VehicleKeyword(ctx context.Context, req *vehiclexsd.VehicleKeywordReq) (*vehiclexsd.VehicleKeywordRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleKeyword(ctx, req)
}

// VehicleLocation corresponds to VehicleServicePort.VehicleLocation。
func (f *VehicleFacade) VehicleLocation(ctx context.Context, req *vehiclexsd.VehicleLocationReq) (*vehiclexsd.VehicleLocationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleLocation(ctx, req)
}

// VehicleLocationDetail corresponds to VehicleServicePort.VehicleLocationDetail。
func (f *VehicleFacade) VehicleLocationDetail(ctx context.Context, req *vehiclexsd.VehicleLocationDetailReq) (*vehiclexsd.VehicleLocationDetailRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleLocationDetail(ctx, req)
}

// VehicleMediaLinks corresponds to VehicleServicePort.VehicleMediaLinks。
func (f *VehicleFacade) VehicleMediaLinks(ctx context.Context, req *vehiclexsd.VehicleMediaLinksReq) (*vehiclexsd.VehicleMediaLinksRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleMediaLinks(ctx, req)
}

// VehicleRetrieve corresponds to VehicleServicePort.VehicleRetrieve。
func (f *VehicleFacade) VehicleRetrieve(ctx context.Context, req *vehiclexsd.VehicleRetrieveReq) (*vehiclexsd.VehicleRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleRetrieve(ctx, req)
}

// VehicleRules corresponds to VehicleServicePort.VehicleRules。
func (f *VehicleFacade) VehicleRules(ctx context.Context, req *vehiclexsd.VehicleRulesReq) (*vehiclexsd.VehicleRulesRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleRules(ctx, req)
}

// VehicleSearchAvailability corresponds to VehicleServicePort.VehicleSearchAvailability。
func (f *VehicleFacade) VehicleSearchAvailability(ctx context.Context, req *vehiclexsd.VehicleSearchAvailabilityReq) (*vehiclexsd.VehicleSearchAvailabilityRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleSearchAvailability(ctx, req)
}

// VehicleUpsellSearchAvailability corresponds to VehicleServicePort.VehicleUpsellSearchAvailability。
func (f *VehicleFacade) VehicleUpsellSearchAvailability(ctx context.Context, req *vehiclexsd.VehicleUpsellSearchAvailabilityReq) (*vehiclexsd.VehicleUpsellSearchAvailabilityRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.VehicleUpsellSearchAvailability(ctx, req)
}
