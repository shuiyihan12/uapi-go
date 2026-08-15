// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the Air service (air): methods map 1:1 to the AirServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	airxsd "github.com/shuiyihan12/uapi-go/pkg/generated/air"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	airsvc "github.com/shuiyihan12/uapi-go/pkg/services/air"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// AirFacade orchestrates the Travelport Air service use cases; methods map 1:1 to the service PortType.
type AirFacade struct {
	manager *manager.ServiceManager
}

// NewAirFacade creates the Air use-case layer.
func NewAirFacade(serviceManager *manager.ServiceManager) *AirFacade {
	return &AirFacade{manager: serviceManager}
}

// getService lazily retrieves the Air service client, handling nil manager and lookup failures uniformly.
func (f *AirFacade) getService() (*airsvc.AirService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*airsvc.AirService](f.manager, "air")
	if err != nil {
		return nil, fmt.Errorf("failed to get air service: %w", err)
	}
	return svc, nil
}

// AirExchange corresponds to AirServicePort.AirExchange。
func (f *AirFacade) AirExchange(ctx context.Context, req *airxsd.AirExchangeReq) (*airxsd.AirExchangeRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirExchange(ctx, req)
}

// AirExchangeEligibility corresponds to AirServicePort.AirExchangeEligibility。
func (f *AirFacade) AirExchangeEligibility(ctx context.Context, req *airxsd.AirExchangeEligibilityReq) (*airxsd.AirExchangeEligibilityRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirExchangeEligibility(ctx, req)
}

// AirExchangeMultiQuote corresponds to AirServicePort.AirExchangeMultiQuote。
func (f *AirFacade) AirExchangeMultiQuote(ctx context.Context, req *airxsd.AirExchangeMultiQuoteReq) (*airxsd.AirExchangeMultiQuoteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirExchangeMultiQuote(ctx, req)
}

// AirExchangeQuote corresponds to AirServicePort.AirExchangeQuote。
func (f *AirFacade) AirExchangeQuote(ctx context.Context, req *airxsd.AirExchangeQuoteReq) (*airxsd.AirExchangeQuoteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirExchangeQuote(ctx, req)
}

// AirExchangeTicketing corresponds to AirServicePort.AirExchangeTicketing。
func (f *AirFacade) AirExchangeTicketing(ctx context.Context, req *airxsd.AirExchangeTicketingReq) (*airxsd.AirExchangeTicketingRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirExchangeTicketing(ctx, req)
}

// AirFareDisplay corresponds to AirServicePort.AirFareDisplay。
func (f *AirFacade) AirFareDisplay(ctx context.Context, req *airxsd.AirFareDisplayReq) (*airxsd.AirFareDisplayRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirFareDisplay(ctx, req)
}

// AirFareRules corresponds to AirServicePort.AirFareRules。
func (f *AirFacade) AirFareRules(ctx context.Context, req *airxsd.AirFareRulesReq) (*airxsd.AirFareRulesRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirFareRules(ctx, req)
}

// AirMerchandisingDetails corresponds to AirServicePort.AirMerchandisingDetails。
func (f *AirFacade) AirMerchandisingDetails(ctx context.Context, req *airxsd.AirMerchandisingDetailsReq) (*airxsd.AirMerchandisingDetailsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirMerchandisingDetails(ctx, req)
}

// AirMerchandisingOfferAvailability corresponds to AirServicePort.AirMerchandisingOfferAvailability。
func (f *AirFacade) AirMerchandisingOfferAvailability(ctx context.Context, req *airxsd.AirMerchandisingOfferAvailabilityReq) (*airxsd.AirMerchandisingOfferAvailabilityRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirMerchandisingOfferAvailability(ctx, req)
}

// AirPrePay corresponds to AirServicePort.AirPrePay。
func (f *AirFacade) AirPrePay(ctx context.Context, req *airxsd.AirPrePayReq) (*airxsd.AirPrePayRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirPrePay(ctx, req)
}

// AirPrice corresponds to AirServicePort.AirPrice。
func (f *AirFacade) AirPrice(ctx context.Context, req *airxsd.AirPriceReq) (*airxsd.AirPriceRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirPrice(ctx, req)
}

// AirRefund corresponds to AirServicePort.AirRefund。
func (f *AirFacade) AirRefund(ctx context.Context, req *airxsd.AirRefundReq) (*airxsd.AirRefundRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirRefund(ctx, req)
}

// AirRefundQuote corresponds to AirServicePort.AirRefundQuote。
func (f *AirFacade) AirRefundQuote(ctx context.Context, req *airxsd.AirRefundQuoteReq) (*airxsd.AirRefundQuoteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirRefundQuote(ctx, req)
}

// AirReprice corresponds to AirServicePort.AirReprice。
func (f *AirFacade) AirReprice(ctx context.Context, req *airxsd.AirRepriceReq) (*airxsd.AirRepriceRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirReprice(ctx, req)
}

// AirRetrieveDocument corresponds to AirServicePort.AirRetrieveDocument。
func (f *AirFacade) AirRetrieveDocument(ctx context.Context, req *airxsd.AirRetrieveDocumentReq) (*airxsd.AirRetrieveDocumentRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirRetrieveDocument(ctx, req)
}

// AirTicketing corresponds to AirServicePort.AirTicketing。
func (f *AirFacade) AirTicketing(ctx context.Context, req *airxsd.AirTicketingReq) (*airxsd.AirTicketingRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirTicketing(ctx, req)
}

// AirUpsellSearch corresponds to AirServicePort.AirUpsellSearch。
func (f *AirFacade) AirUpsellSearch(ctx context.Context, req *airxsd.AirUpsellSearchReq) (*airxsd.AirUpsellSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirUpsellSearch(ctx, req)
}

// AirVoidDocument corresponds to AirServicePort.AirVoidDocument。
func (f *AirFacade) AirVoidDocument(ctx context.Context, req *airxsd.AirVoidDocumentReq) (*airxsd.AirVoidDocumentRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirVoidDocument(ctx, req)
}

// AvailabilitySearch corresponds to AirServicePort.AvailabilitySearch。
func (f *AirFacade) AvailabilitySearch(ctx context.Context, req *airxsd.AvailabilitySearchReq) (*airxsd.AvailabilitySearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AvailabilitySearch(ctx, req)
}

// EMDIssuance corresponds to AirServicePort.EMDIssuance。
func (f *AirFacade) EMDIssuance(ctx context.Context, req *airxsd.EMDIssuanceReq) (*airxsd.EMDIssuanceRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.EMDIssuance(ctx, req)
}

// EMDRetrieve corresponds to AirServicePort.EMDRetrieve。
func (f *AirFacade) EMDRetrieve(ctx context.Context, req *airxsd.EMDRetrieveReq) (*airxsd.EMDRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.EMDRetrieve(ctx, req)
}

// FlightDetails corresponds to AirServicePort.FlightDetails。
func (f *AirFacade) FlightDetails(ctx context.Context, req *airxsd.FlightDetailsReq) (*airxsd.FlightDetailsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.FlightDetails(ctx, req)
}

// FlightInformation corresponds to AirServicePort.FlightInformation。
func (f *AirFacade) FlightInformation(ctx context.Context, req *airxsd.FlightInformationReq) (*airxsd.FlightInformationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.FlightInformation(ctx, req)
}

// FlightTimeTable corresponds to AirServicePort.FlightTimeTable。
func (f *AirFacade) FlightTimeTable(ctx context.Context, req *airxsd.FlightTimeTableReq) (*airxsd.FlightTimeTableRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.FlightTimeTable(ctx, req)
}

// LowFareSearch corresponds to AirServicePort.LowFareSearch。
func (f *AirFacade) LowFareSearch(ctx context.Context, req *airxsd.LowFareSearchReq) (*airxsd.LowFareSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.LowFareSearch(ctx, req)
}

// ScheduleSearch corresponds to AirServicePort.ScheduleSearch。
func (f *AirFacade) ScheduleSearch(ctx context.Context, req *airxsd.ScheduleSearchReq) (*airxsd.ScheduleSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ScheduleSearch(ctx, req)
}

// SeatMap corresponds to AirServicePort.SeatMap。
func (f *AirFacade) SeatMap(ctx context.Context, req *airxsd.SeatMapReq) (*airxsd.SeatMapRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.SeatMap(ctx, req)
}
