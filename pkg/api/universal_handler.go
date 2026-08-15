// Package api provides the HTTP interface layer for the Universal service, plus unified
// error response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// UniversalHandler handles the HTTP interface for the Travelport Universal service.
type UniversalHandler struct {
	facade *usecase.UniversalFacade
}

// NewUniversalHandler creates the Universal handler.
func NewUniversalHandler(facade *usecase.UniversalFacade) *UniversalHandler {
	return &UniversalHandler{facade: facade}
}

// RegisterRoutes registers all endpoints of the Universal service.
func (h *UniversalHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/universal/ack-schedule-change", f.AckScheduleChange)
	registerPortHandler(mux, apiBasePath+"/universal/air-cancel", f.AirCancel)
	registerPortHandler(mux, apiBasePath+"/universal/air-create-reservation", f.AirCreateReservation)
	registerPortHandler(mux, apiBasePath+"/universal/air-merchandising-fulfillment", f.AirMerchandisingFulfillment)
	registerPortHandler(mux, apiBasePath+"/universal/hotel-cancel", f.HotelCancel)
	registerPortHandler(mux, apiBasePath+"/universal/hotel-create-reservation", f.HotelCreateReservation)
	registerPortHandler(mux, apiBasePath+"/universal/passive-cancel", f.PassiveCancel)
	registerPortHandler(mux, apiBasePath+"/universal/passive-create-reservation", f.PassiveCreateReservation)
	registerPortHandler(mux, apiBasePath+"/universal/provider-reservation-display-details", f.ProviderReservationDisplayDetails)
	registerPortHandler(mux, apiBasePath+"/universal/provider-reservation-divide", f.ProviderReservationDivide)
	registerPortHandler(mux, apiBasePath+"/universal/rail-create-reservation", f.RailCreateReservation)
	registerPortHandler(mux, apiBasePath+"/universal/saved-trip-create", f.SavedTripCreate)
	registerPortHandler(mux, apiBasePath+"/universal/saved-trip-delete", f.SavedTripDelete)
	registerPortHandler(mux, apiBasePath+"/universal/saved-trip-modify", f.SavedTripModify)
	registerPortHandler(mux, apiBasePath+"/universal/saved-trip-retrieve", f.SavedTripRetrieve)
	registerPortHandler(mux, apiBasePath+"/universal/saved-trip-search", f.SavedTripSearch)
	registerPortHandler(mux, apiBasePath+"/universal/universal-record-cancel", f.UniversalRecordCancel)
	registerPortHandler(mux, apiBasePath+"/universal/universal-record-import", f.UniversalRecordImport)
	registerPortHandler(mux, apiBasePath+"/universal/universal-record-modify", f.UniversalRecordModify)
	registerPortHandler(mux, apiBasePath+"/universal/universal-record-retrieve", f.UniversalRecordRetrieve)
	registerPortHandler(mux, apiBasePath+"/universal/universal-record-search", f.UniversalRecordSearch)
	registerPortHandler(mux, apiBasePath+"/universal/vehicle-cancel", f.VehicleCancel)
	registerPortHandler(mux, apiBasePath+"/universal/vehicle-create-reservation", f.VehicleCreateReservation)
}
