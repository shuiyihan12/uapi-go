// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the Rail service (air): methods map 1:1 to the RailServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	airxsd "github.com/shuiyihan12/uapi-go/pkg/generated/air"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	railsvc "github.com/shuiyihan12/uapi-go/pkg/services/rail"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// RailFacade orchestrates the Travelport Rail service use cases; methods map 1:1 to the service PortType.
type RailFacade struct {
	manager *manager.ServiceManager
}

// NewRailFacade creates the Rail use-case layer.
func NewRailFacade(serviceManager *manager.ServiceManager) *RailFacade {
	return &RailFacade{manager: serviceManager}
}

// getService lazily retrieves the Rail service client, handling nil manager and lookup failures uniformly.
func (f *RailFacade) getService() (*railsvc.RailService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*railsvc.RailService](f.manager, "rail")
	if err != nil {
		return nil, fmt.Errorf("failed to get rail service: %w", err)
	}
	return svc, nil
}

// RailAvailabilitySearch corresponds to RailServicePort.RailAvailabilitySearch。
func (f *RailFacade) RailAvailabilitySearch(ctx context.Context, req *airxsd.RailAvailabilitySearchReq) (*airxsd.RailAvailabilitySearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.RailAvailabilitySearch(ctx, req)
}

// RailExchange corresponds to RailServicePort.RailExchange。
func (f *RailFacade) RailExchange(ctx context.Context, req *airxsd.RailExchangeReq) (*airxsd.RailExchangeRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.RailExchange(ctx, req)
}

// RailExchangeQuote corresponds to RailServicePort.RailExchangeQuote。
func (f *RailFacade) RailExchangeQuote(ctx context.Context, req *airxsd.RailExchangeQuoteReq) (*airxsd.RailExchangeQuoteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.RailExchangeQuote(ctx, req)
}

// RailRefund corresponds to RailServicePort.RailRefund。
func (f *RailFacade) RailRefund(ctx context.Context, req *airxsd.RailRefundReq) (*airxsd.RailRefundRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.RailRefund(ctx, req)
}

// RailRefundQuote corresponds to RailServicePort.RailRefundQuote。
func (f *RailFacade) RailRefundQuote(ctx context.Context, req *airxsd.RailRefundQuoteReq) (*airxsd.RailRefundQuoteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.RailRefundQuote(ctx, req)
}

// RailSeatMap corresponds to RailServicePort.RailSeatMap。
func (f *RailFacade) RailSeatMap(ctx context.Context, req *airxsd.RailSeatMapReq) (*airxsd.RailSeatMapRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.RailSeatMap(ctx, req)
}
