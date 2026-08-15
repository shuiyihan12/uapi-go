// Package api provides the HTTP interface layer for the Vehicle service, plus unified
// error response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// VehicleHandler handles the HTTP interface for the Travelport Vehicle service.
type VehicleHandler struct {
	facade    *usecase.VehicleFacade
	universal *usecase.UniversalFacade
}

// NewVehicleHandler creates the Vehicle handler. universal proxies the cross-product
// unified booking/cancellation capabilities.
func NewVehicleHandler(facade *usecase.VehicleFacade, universal *usecase.UniversalFacade) *VehicleHandler {
	return &VehicleHandler{facade: facade, universal: universal}
}

// RegisterRoutes registers all endpoints of the Vehicle service.
func (h *VehicleHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/vehicle/vehicle-keyword", f.VehicleKeyword)
	registerPortHandler(mux, apiBasePath+"/vehicle/vehicle-location", f.VehicleLocation)
	registerPortHandler(mux, apiBasePath+"/vehicle/vehicle-location-detail", f.VehicleLocationDetail)
	registerPortHandler(mux, apiBasePath+"/vehicle/vehicle-media-links", f.VehicleMediaLinks)
	registerPortHandler(mux, apiBasePath+"/vehicle/vehicle-retrieve", f.VehicleRetrieve)
	registerPortHandler(mux, apiBasePath+"/vehicle/vehicle-rules", f.VehicleRules)
	registerPortHandler(mux, apiBasePath+"/vehicle/vehicle-search-availability", f.VehicleSearchAvailability)
	registerPortHandler(mux, apiBasePath+"/vehicle/vehicle-upsell-search-availability", f.VehicleUpsellSearchAvailability)

	// Productized aliases: booking/cancellation is delegated to the UniversalRecord engine.
	registerPortHandler(mux, apiBasePath+"/vehicle/book", h.universal.VehicleCreateReservation)
	registerPortHandler(mux, apiBasePath+"/vehicle/cancel", h.universal.VehicleCancel)
}
