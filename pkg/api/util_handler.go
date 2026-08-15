// Package api provides the HTTP interface layer for the Util service, plus unified error
// response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// UtilHandler handles the HTTP interface for the Travelport Util service.
type UtilHandler struct {
	facade *usecase.UtilFacade
}

// NewUtilHandler creates the Util handler.
func NewUtilHandler(facade *usecase.UtilFacade) *UtilHandler {
	return &UtilHandler{facade: facade}
}

// RegisterRoutes registers all endpoints of the Util service.
func (h *UtilHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/util/agency-service-fee-create", f.AgencyServiceFeeCreate)
	registerPortHandler(mux, apiBasePath+"/util/branded-fare-admin", f.BrandedFareAdmin)
	registerPortHandler(mux, apiBasePath+"/util/branded-fare-search", f.BrandedFareSearch)
	registerPortHandler(mux, apiBasePath+"/util/calculate-tax", f.CalculateTax)
	registerPortHandler(mux, apiBasePath+"/util/content-provider-retrieve", f.ContentProviderRetrieve)
	registerPortHandler(mux, apiBasePath+"/util/create-agency-fee-mco", f.CreateAgencyFeeMco)
	registerPortHandler(mux, apiBasePath+"/util/create-airline-fee-mco", f.CreateAirlineFeeMco)
	registerPortHandler(mux, apiBasePath+"/util/credit-card-auth", f.CreditCardAuth)
	registerPortHandler(mux, apiBasePath+"/util/currency-conversion", f.CurrencyConversion)
	registerPortHandler(mux, apiBasePath+"/util/find-employees-on-flight", f.FindEmployeesOnFlight)
	registerPortHandler(mux, apiBasePath+"/util/mco-create", f.MCOCreate)
	registerPortHandler(mux, apiBasePath+"/util/mco-exchange", f.MCOExchange)
	registerPortHandler(mux, apiBasePath+"/util/mco-issue", f.MCOIssue)
	registerPortHandler(mux, apiBasePath+"/util/mco-retrieve", f.MCORetrieve)
	registerPortHandler(mux, apiBasePath+"/util/mco-search", f.McoSearch)
	registerPortHandler(mux, apiBasePath+"/util/mco-void", f.McoVoid)
	registerPortHandler(mux, apiBasePath+"/util/mct-count", f.MctCount)
	registerPortHandler(mux, apiBasePath+"/util/mct-lookup", f.MctLookup)
	registerPortHandler(mux, apiBasePath+"/util/mir-report-retrieve", f.MirReportRetrieve)
	registerPortHandler(mux, apiBasePath+"/util/reference-data-retrieve", f.ReferenceDataRetrieve)
	registerPortHandler(mux, apiBasePath+"/util/reference-data-search", f.ReferenceDataSearch)
	registerPortHandler(mux, apiBasePath+"/util/reference-data-update", f.ReferenceDataUpdate)
	registerPortHandler(mux, apiBasePath+"/util/upsell-admin", f.UpsellAdmin)
	registerPortHandler(mux, apiBasePath+"/util/upsell-search", f.UpsellSearch)
}
