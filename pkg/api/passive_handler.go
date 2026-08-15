// Package api provides the HTTP interface layer for the Passive service.
// In Travelport, passive has no standalone WSDL/portType of its own; its booking and
// cancellation can only be exposed through UniversalRecordService. Therefore this handler
// holds no service of its own and delegates requests directly to the Universal engine
// (/api/passive/book, /api/passive/cancel).
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// PassiveHandler handles the HTTP interface for passive (proxying the UniversalRecord
// engine).
type PassiveHandler struct {
	universal *usecase.UniversalFacade
}

// NewPassiveHandler creates the passive handler.
func NewPassiveHandler(universal *usecase.UniversalFacade) *PassiveHandler {
	return &PassiveHandler{universal: universal}
}

// RegisterRoutes registers the passive endpoints.
func (h *PassiveHandler) RegisterRoutes(mux *http.ServeMux) {
	registerPortHandler(mux, apiBasePath+"/passive/book", h.universal.PassiveCreateReservation)
	registerPortHandler(mux, apiBasePath+"/passive/cancel", h.universal.PassiveCancel)
}
