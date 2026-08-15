// Package api provides the HTTP interface layer for the Air service, plus unified error
// response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// AirHandler handles the HTTP interface for the Travelport Air service.
type AirHandler struct {
	facade    *usecase.AirFacade
	universal *usecase.UniversalFacade
}

// NewAirHandler creates the Air handler. universal proxies the cross-product unified
// booking/cancellation capabilities.
func NewAirHandler(facade *usecase.AirFacade, universal *usecase.UniversalFacade) *AirHandler {
	return &AirHandler{facade: facade, universal: universal}
}

// RegisterRoutes registers all endpoints of the Air service.
func (h *AirHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/air/air-exchange", f.AirExchange)
	registerPortHandler(mux, apiBasePath+"/air/air-exchange-eligibility", f.AirExchangeEligibility)
	registerPortHandler(mux, apiBasePath+"/air/air-exchange-multi-quote", f.AirExchangeMultiQuote)
	registerPortHandler(mux, apiBasePath+"/air/air-exchange-quote", f.AirExchangeQuote)
	registerPortHandler(mux, apiBasePath+"/air/air-exchange-ticketing", f.AirExchangeTicketing)
	registerPortHandler(mux, apiBasePath+"/air/air-fare-display", f.AirFareDisplay)
	registerPortHandler(mux, apiBasePath+"/air/air-fare-rules", f.AirFareRules)
	registerPortHandler(mux, apiBasePath+"/air/air-merchandising-details", f.AirMerchandisingDetails)
	registerPortHandler(mux, apiBasePath+"/air/air-merchandising-offer-availability", f.AirMerchandisingOfferAvailability)
	registerPortHandler(mux, apiBasePath+"/air/air-pre-pay", f.AirPrePay)
	registerPortHandler(mux, apiBasePath+"/air/air-price", f.AirPrice)
	registerPortHandler(mux, apiBasePath+"/air/air-refund", f.AirRefund)
	registerPortHandler(mux, apiBasePath+"/air/air-refund-quote", f.AirRefundQuote)
	registerPortHandler(mux, apiBasePath+"/air/air-reprice", f.AirReprice)
	registerPortHandler(mux, apiBasePath+"/air/air-retrieve-document", f.AirRetrieveDocument)
	registerPortHandler(mux, apiBasePath+"/air/air-ticketing", f.AirTicketing)
	registerPortHandler(mux, apiBasePath+"/air/air-upsell-search", f.AirUpsellSearch)
	registerPortHandler(mux, apiBasePath+"/air/air-void-document", f.AirVoidDocument)
	registerPortHandler(mux, apiBasePath+"/air/availability-search", f.AvailabilitySearch)
	registerPortHandler(mux, apiBasePath+"/air/emd-issuance", f.EMDIssuance)
	registerPortHandler(mux, apiBasePath+"/air/emd-retrieve", f.EMDRetrieve)
	registerPortHandler(mux, apiBasePath+"/air/flight-details", f.FlightDetails)
	registerPortHandler(mux, apiBasePath+"/air/flight-information", f.FlightInformation)
	registerPortHandler(mux, apiBasePath+"/air/flight-time-table", f.FlightTimeTable)
	registerPortHandler(mux, apiBasePath+"/air/low-fare-search", f.LowFareSearch)
	registerPortHandler(mux, apiBasePath+"/air/schedule-search", f.ScheduleSearch)
	registerPortHandler(mux, apiBasePath+"/air/seat-map", f.SeatMap)

	// Productized aliases: booking/cancellation is delegated to the UniversalRecord engine
	// (cross-product unified ownership, no semantic loss, no duplicated code).
	registerPortHandler(mux, apiBasePath+"/air/book", h.universal.AirCreateReservation)
	registerPortHandler(mux, apiBasePath+"/air/cancel", h.universal.AirCancel)
}
