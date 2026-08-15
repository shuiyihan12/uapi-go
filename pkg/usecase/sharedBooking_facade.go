// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the SharedBooking service (sharedbooking): methods map 1:1 to the SharedBookingServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	sharedbookingxsd "github.com/shuiyihan12/uapi-go/pkg/generated/sharedbooking"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	sharedBookingsvc "github.com/shuiyihan12/uapi-go/pkg/services/sharedBooking"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// SharedBookingFacade orchestrates the Travelport SharedBooking service use cases; methods map 1:1 to the service PortType.ngServicePort。
type SharedBookingFacade struct {
	manager *manager.ServiceManager
}

// NewSharedBookingFacade creates the SharedBooking use-case layer.
func NewSharedBookingFacade(serviceManager *manager.ServiceManager) *SharedBookingFacade {
	return &SharedBookingFacade{manager: serviceManager}
}

// getService lazily retrieves the SharedBooking service client, handling nil manager and lookup failures uniformly.
func (f *SharedBookingFacade) getService() (*sharedBookingsvc.SharedBookingService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*sharedBookingsvc.SharedBookingService](f.manager, "sharedBooking")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharedBooking service: %w", err)
	}
	return svc, nil
}

// BookingAirExchange corresponds to SharedBookingServicePort.BookingAirExchange。
func (f *SharedBookingFacade) BookingAirExchange(ctx context.Context, req *sharedbookingxsd.BookingAirExchangeReq) (*sharedbookingxsd.BookingAirExchangeRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingAirExchange(ctx, req)
}

// BookingAirExchangeQuote corresponds to SharedBookingServicePort.BookingAirExchangeQuote。
func (f *SharedBookingFacade) BookingAirExchangeQuote(ctx context.Context, req *sharedbookingxsd.BookingAirExchangeQuoteReq) (*sharedbookingxsd.BookingAirExchangeQuoteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingAirExchangeQuote(ctx, req)
}

// BookingAirPnrElement corresponds to SharedBookingServicePort.BookingAirPnrElement。
func (f *SharedBookingFacade) BookingAirPnrElement(ctx context.Context, req *sharedbookingxsd.BookingAirPnrElementReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingAirPnrElement(ctx, req)
}

// BookingAirSegment corresponds to SharedBookingServicePort.BookingAirSegment。
func (f *SharedBookingFacade) BookingAirSegment(ctx context.Context, req *sharedbookingxsd.BookingAirSegmentReq) (*sharedbookingxsd.BookingAirSegmentRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingAirSegment(ctx, req)
}

// BookingDisplay corresponds to SharedBookingServicePort.BookingDisplay。
func (f *SharedBookingFacade) BookingDisplay(ctx context.Context, req *sharedbookingxsd.BookingDisplayReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingDisplay(ctx, req)
}

// BookingEnd corresponds to SharedBookingServicePort.BookingEnd。
func (f *SharedBookingFacade) BookingEnd(ctx context.Context, req *sharedbookingxsd.BookingEndReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingEnd(ctx, req)
}

// BookingHotelPnrElement corresponds to SharedBookingServicePort.BookingHotelPnrElement。
func (f *SharedBookingFacade) BookingHotelPnrElement(ctx context.Context, req *sharedbookingxsd.BookingHotelPnrElementReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingHotelPnrElement(ctx, req)
}

// BookingHotelSegment corresponds to SharedBookingServicePort.BookingHotelSegment。
func (f *SharedBookingFacade) BookingHotelSegment(ctx context.Context, req *sharedbookingxsd.BookingHotelSegmentReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingHotelSegment(ctx, req)
}

// BookingPnrElement corresponds to SharedBookingServicePort.BookingPnrElement。
func (f *SharedBookingFacade) BookingPnrElement(ctx context.Context, req *sharedbookingxsd.BookingPnrElementReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingPnrElement(ctx, req)
}

// BookingPricing corresponds to SharedBookingServicePort.BookingPricing。
func (f *SharedBookingFacade) BookingPricing(ctx context.Context, req *sharedbookingxsd.BookingPricingReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingPricing(ctx, req)
}

// BookingRetrieveDocument corresponds to SharedBookingServicePort.BookingRetrieveDocument。
func (f *SharedBookingFacade) BookingRetrieveDocument(ctx context.Context, req *sharedbookingxsd.BookingRetrieveDocumentReq) (*sharedbookingxsd.BookingRetrieveDocumentRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingRetrieveDocument(ctx, req)
}

// BookingSeatAssignment corresponds to SharedBookingServicePort.BookingSeatAssignment。
func (f *SharedBookingFacade) BookingSeatAssignment(ctx context.Context, req *sharedbookingxsd.BookingSeatAssignmentReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingSeatAssignment(ctx, req)
}

// BookingStart corresponds to SharedBookingServicePort.BookingStart。
func (f *SharedBookingFacade) BookingStart(ctx context.Context, req *sharedbookingxsd.BookingStartReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingStart(ctx, req)
}

// BookingTerminal corresponds to SharedBookingServicePort.BookingTerminal。
func (f *SharedBookingFacade) BookingTerminal(ctx context.Context, req *sharedbookingxsd.BookingTerminalReq) (*sharedbookingxsd.BookingTerminalRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingTerminal(ctx, req)
}

// BookingTraveler corresponds to SharedBookingServicePort.BookingTraveler。
func (f *SharedBookingFacade) BookingTraveler(ctx context.Context, req *sharedbookingxsd.BookingTravelerReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingTraveler(ctx, req)
}

// BookingVehiclePnrElement corresponds to SharedBookingServicePort.BookingVehiclePnrElement (v55 addition),
// adds / updates / deletes vehicle elements on a shared booking PNR.
func (f *SharedBookingFacade) BookingVehiclePnrElement(ctx context.Context, req *sharedbookingxsd.BookingVehiclePnrElementReq) (*sharedbookingxsd.BookingVehiclePnrElementRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BookingVehiclePnrElement(ctx, req)
}
