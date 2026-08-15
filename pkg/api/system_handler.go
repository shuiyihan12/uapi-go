// Package api provides the HTTP interface layer for the System service, plus unified
// error response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// SystemHandler handles the HTTP interface for the Travelport System service.
type SystemHandler struct {
	facade *usecase.SystemFacade
}

// NewSystemHandler creates the System handler.
func NewSystemHandler(facade *usecase.SystemFacade) *SystemHandler {
	return &SystemHandler{facade: facade}
}

// RegisterRoutes registers all endpoints of the System service.
func (h *SystemHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/system/ping", f.Ping)
	registerPortHandler(mux, apiBasePath+"/system/info", f.Info)
	registerPortHandler(mux, apiBasePath+"/system/time", f.Time)
	registerPortHandler(mux, apiBasePath+"/system/cache", f.ExternalCacheAccess)
}
