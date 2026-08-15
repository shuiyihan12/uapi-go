// Package api provides the HTTP interface layer for the SharedBooking service, plus
// unified error response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// SharedBookingHandler handles the HTTP interface for the Travelport SharedBooking service.
type SharedBookingHandler struct {
	facade *usecase.SharedBookingFacade
}

// NewSharedBookingHandler creates the SharedBooking handler.
func NewSharedBookingHandler(facade *usecase.SharedBookingFacade) *SharedBookingHandler {
	return &SharedBookingHandler{facade: facade}
}

// RegisterRoutes registers all endpoints of the SharedBooking service.
func (h *SharedBookingHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-air-exchange", f.BookingAirExchange)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-air-exchange-quote", f.BookingAirExchangeQuote)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-air-pnr-element", f.BookingAirPnrElement)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-air-segment", f.BookingAirSegment)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-display", f.BookingDisplay)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-end", f.BookingEnd)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-hotel-pnr-element", f.BookingHotelPnrElement)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-hotel-segment", f.BookingHotelSegment)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-pnr-element", f.BookingPnrElement)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-pricing", f.BookingPricing)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-retrieve-document", f.BookingRetrieveDocument)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-seat-assignment", f.BookingSeatAssignment)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-start", f.BookingStart)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-terminal", f.BookingTerminal)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-traveler", f.BookingTraveler)
	registerPortHandler(mux, apiBasePath+"/sharedBooking/booking-vehicle-pnr-element", f.BookingVehiclePnrElement)
}
