// Package api provides the HTTP interface layer for the GdsQueue service, plus unified
// error response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// GdsQueueHandler handles the HTTP interface for the Travelport GdsQueue service.
type GdsQueueHandler struct {
	facade *usecase.GdsQueueFacade
}

// NewGdsQueueHandler creates the GdsQueue handler.
func NewGdsQueueHandler(facade *usecase.GdsQueueFacade) *GdsQueueHandler {
	return &GdsQueueHandler{facade: facade}
}

// RegisterRoutes registers all endpoints of the GdsQueue service.
func (h *GdsQueueHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/gdsQueue/enter-queue", f.GdsEnterQueue)
	registerPortHandler(mux, apiBasePath+"/gdsQueue/exit-queue", f.GdsExitQueue)
	registerPortHandler(mux, apiBasePath+"/gdsQueue/next-on-queue", f.GdsNextOnQueue)
	registerPortHandler(mux, apiBasePath+"/gdsQueue/queue-agent-list", f.GdsQueueAgentList)
	registerPortHandler(mux, apiBasePath+"/gdsQueue/queue-count", f.GdsQueueCount)
	registerPortHandler(mux, apiBasePath+"/gdsQueue/queue-list", f.GdsQueueList)
	registerPortHandler(mux, apiBasePath+"/gdsQueue/queue-place", f.GdsQueuePlace)
	registerPortHandler(mux, apiBasePath+"/gdsQueue/queue-remove", f.GdsQueueRemove)
}
