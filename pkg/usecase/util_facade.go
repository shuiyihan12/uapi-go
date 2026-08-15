// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the Util service (util_v55_0): methods map 1:1 to the UtilServicePort,
// each retrieves the service client, injects the trace context and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	utilxsd "github.com/shuiyihan12/uapi-go/pkg/generated/util"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	utilsvc "github.com/shuiyihan12/uapi-go/pkg/services/util"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// UtilFacade orchestrates the Travelport Util service use cases; methods map 1:1 to the service PortType.
type UtilFacade struct {
	manager *manager.ServiceManager
}

// NewUtilFacade creates the Util use-case layer.
func NewUtilFacade(serviceManager *manager.ServiceManager) *UtilFacade {
	return &UtilFacade{manager: serviceManager}
}

// getService lazily retrieves the Util service client, handling nil manager and lookup failures uniformly.
func (f *UtilFacade) getService() (*utilsvc.UtilService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*utilsvc.UtilService](f.manager, "util")
	if err != nil {
		return nil, fmt.Errorf("failed to get util service: %w", err)
	}
	return svc, nil
}

// AgencyServiceFeeCreate creates an agency service fee (corresponds to AgencyCreateServiceFeePortType).
func (f *UtilFacade) AgencyServiceFeeCreate(ctx context.Context, req *utilxsd.AgencyServiceFeeCreateReq) (*utilxsd.AgencyServiceFeeCreateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AgencyServiceFeeCreate(ctx, req)
}

// BrandedFareAdmin manages branded fares (corresponds to BrandedFareAdminPortType).
func (f *UtilFacade) BrandedFareAdmin(ctx context.Context, req *utilxsd.BrandedFareAdminReq) (*utilxsd.BrandedFareAdminRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BrandedFareAdmin(ctx, req)
}

// BrandedFareSearch searches for branded fares (corresponds to BrandedFareSearchPortType).
func (f *UtilFacade) BrandedFareSearch(ctx context.Context, req *utilxsd.BrandedFareSearchReq) (*utilxsd.BrandedFareSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.BrandedFareSearch(ctx, req)
}

// CalculateTax calculates taxes (corresponds to CalculateTaxPortType).
func (f *UtilFacade) CalculateTax(ctx context.Context, req *utilxsd.CalculateTaxReq) (*utilxsd.CalculateTaxRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.CalculateTax(ctx, req)
}

// ContentProviderRetrieve retrieves content provider data (corresponds to ContentProviderRetrievePortType).
func (f *UtilFacade) ContentProviderRetrieve(ctx context.Context, req *utilxsd.ContentProviderRetrieveReq) (*utilxsd.ContentProviderRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ContentProviderRetrieve(ctx, req)
}

// CreateAgencyFeeMco creates an agency-fee MCO (corresponds to McoCreateAgencyFeePortType).
func (f *UtilFacade) CreateAgencyFeeMco(ctx context.Context, req *utilxsd.CreateAgencyFeeMcoReq) (*utilxsd.CreateAgencyFeeMcoRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.CreateAgencyFeeMco(ctx, req)
}

// CreateAirlineFeeMco creates an airline-fee MCO (corresponds to McoCreateAgencyFeePortType, airline-fee variant).
func (f *UtilFacade) CreateAirlineFeeMco(ctx context.Context, req *utilxsd.CreateAirlineFeeMcoReq) (*utilxsd.CreateAirlineFeeMcoRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.CreateAirlineFeeMco(ctx, req)
}

// CreditCardAuth authorizes a credit card (corresponds to UtilCreditCardAuthPortType).
func (f *UtilFacade) CreditCardAuth(ctx context.Context, req *utilxsd.CreditCardAuthReq) (*utilxsd.CreditCardAuthRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.CreditCardAuth(ctx, req)
}

// CurrencyConversion converts currency (corresponds to UtilCurrencyConversionPortType).
func (f *UtilFacade) CurrencyConversion(ctx context.Context, req *utilxsd.CurrencyConversionReq) (*utilxsd.CurrencyConversionRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.CurrencyConversion(ctx, req)
}

// FindEmployeesOnFlight queries employees on a flight (corresponds to FindEmployeesOnFlightServicePortType).
func (f *UtilFacade) FindEmployeesOnFlight(ctx context.Context, req *utilxsd.FindEmployeesOnFlightReq) (*utilxsd.FindEmployeesOnFlightRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.FindEmployeesOnFlight(ctx, req)
}

// MCOCreate creates an MCO (corresponds to MCOCreatePortType).
func (f *UtilFacade) MCOCreate(ctx context.Context, req *utilxsd.MCOCreateReq) (*utilxsd.MCOCreateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.MCOCreate(ctx, req)
}

// MCOExchange exchanges an MCO (corresponds to MCOExchangePortType).
func (f *UtilFacade) MCOExchange(ctx context.Context, req *utilxsd.MCOExchangeReq) (*utilxsd.MCOExchangeRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.MCOExchange(ctx, req)
}

// MCOIssue issues an MCO (corresponds to MCOIssuePortType).
func (f *UtilFacade) MCOIssue(ctx context.Context, req *utilxsd.MCOIssueReq) (*utilxsd.MCOIssueRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.MCOIssue(ctx, req)
}

// MCORetrieve retrieves an MCO (corresponds to MCORetrievePortType).
func (f *UtilFacade) MCORetrieve(ctx context.Context, req *utilxsd.MCORetrieveReq) (*utilxsd.MCORetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.MCORetrieve(ctx, req)
}

// McoSearch searches for MCOs (corresponds to McoSearchPortType).
func (f *UtilFacade) McoSearch(ctx context.Context, req *utilxsd.McoSearchReq) (*utilxsd.McoSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.McoSearch(ctx, req)
}

// McoVoid voids an MCO (corresponds to McoVoidPortType).
func (f *UtilFacade) McoVoid(ctx context.Context, req *utilxsd.McoVoidReq) (*utilxsd.McoVoidRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.McoVoid(ctx, req)
}

// MctCount counts minimum connecting times (corresponds to MctCountPortType).
func (f *UtilFacade) MctCount(ctx context.Context, req *utilxsd.MctCountReq) (*utilxsd.MctCountRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.MctCount(ctx, req)
}

// MctLookup looks up minimum connecting times (corresponds to MctLookupPortType).
func (f *UtilFacade) MctLookup(ctx context.Context, req *utilxsd.MctLookupReq) (*utilxsd.MctLookupRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.MctLookup(ctx, req)
}

// MirReportRetrieve retrieves an MIR report (corresponds to ReportRetrievePortType).
func (f *UtilFacade) MirReportRetrieve(ctx context.Context, req *utilxsd.MirReportRetrieveReq) (*utilxsd.MirReportRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.MirReportRetrieve(ctx, req)
}

// ReferenceDataRetrieve retrieves reference data (corresponds to ReferenceDataRetrievePortType).
func (f *UtilFacade) ReferenceDataRetrieve(ctx context.Context, req *utilxsd.ReferenceDataRetrieveReq) (*utilxsd.ReferenceDataRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ReferenceDataRetrieve(ctx, req)
}

// ReferenceDataSearch queries reference data (corresponds to ReferenceDataLookupPortType).
func (f *UtilFacade) ReferenceDataSearch(ctx context.Context, req *utilxsd.ReferenceDataSearchReq) (*utilxsd.ReferenceDataSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ReferenceDataSearch(ctx, req)
}

// ReferenceDataUpdate updates reference data (corresponds to ReferenceDataUpdatePortType).
func (f *UtilFacade) ReferenceDataUpdate(ctx context.Context, req *utilxsd.ReferenceDataUpdateReq) (*utilxsd.ReferenceDataUpdateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ReferenceDataUpdate(ctx, req)
}

// UpsellAdmin manages upsell rules (corresponds to UpsellAdminPortType).
func (f *UtilFacade) UpsellAdmin(ctx context.Context, req *utilxsd.UpsellAdminReq) (*utilxsd.UpsellAdminRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UpsellAdmin(ctx, req)
}

// UpsellSearch searches for upsell offers (corresponds to UpsellAdminSearchPortType).
func (f *UtilFacade) UpsellSearch(ctx context.Context, req *utilxsd.UpsellSearchReq) (*utilxsd.UpsellSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UpsellSearch(ctx, req)
}
