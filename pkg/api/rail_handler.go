// Package api provides the HTTP interface layer for the Rail service, plus unified error
// response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// RailHandler handles the HTTP interface for the Travelport Rail service.
type RailHandler struct {
	facade    *usecase.RailFacade
	universal *usecase.UniversalFacade
}

// NewRailHandler creates the Rail handler. universal proxies the cross-product unified
// booking capability.
func NewRailHandler(facade *usecase.RailFacade, universal *usecase.UniversalFacade) *RailHandler {
	return &RailHandler{facade: facade, universal: universal}
}

// RegisterRoutes registers all endpoints of the Rail service.
func (h *RailHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/rail/rail-availability-search", f.RailAvailabilitySearch)
	registerPortHandler(mux, apiBasePath+"/rail/rail-exchange", f.RailExchange)
	registerPortHandler(mux, apiBasePath+"/rail/rail-exchange-quote", f.RailExchangeQuote)
	registerPortHandler(mux, apiBasePath+"/rail/rail-refund", f.RailRefund)
	registerPortHandler(mux, apiBasePath+"/rail/rail-refund-quote", f.RailRefundQuote)
	registerPortHandler(mux, apiBasePath+"/rail/rail-seat-map", f.RailSeatMap)

	// Productized aliases: booking is delegated to the UniversalRecord engine. Rail
	// cancellation uniformly goes through UniversalRecordCancel.
	registerPortHandler(mux, apiBasePath+"/rail/book", h.universal.RailCreateReservation)
}
