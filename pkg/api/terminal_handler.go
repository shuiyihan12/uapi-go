// Package api provides the HTTP interface layer for the Terminal service, plus unified
// error response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// TerminalHandler handles the HTTP interface for the Travelport Terminal service.
type TerminalHandler struct {
	facade *usecase.TerminalFacade
}

// NewTerminalHandler creates the Terminal handler.
func NewTerminalHandler(facade *usecase.TerminalFacade) *TerminalHandler {
	return &TerminalHandler{facade: facade}
}

// RegisterRoutes registers all endpoints of the Terminal service.
func (h *TerminalHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/terminal/create-terminal-session", f.CreateTerminalSession)
	registerPortHandler(mux, apiBasePath+"/terminal/end-terminal-session", f.EndTerminalSession)
	registerPortHandler(mux, apiBasePath+"/terminal/terminal", f.Terminal)
}
