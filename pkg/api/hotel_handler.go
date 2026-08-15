// Package api provides the HTTP interface layer for hotel search and details,
// plus unified error response handling.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	hotelxsd "github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// Common constants for the interface layer: default request timeout, request body size
// limit, and API route prefix.
const (
	defaultRequestTimeout = 60 * time.Second
	maxRequestBodyBytes   = 1 << 20
	apiBasePath           = "/api"
)

// HotelFacade defines the minimal interface the handler depends on.
// Requests and responses of all operations (including the REST business surface
// SearchAvailability/Details) use WSDL-generated hotel package types, registered
// uniformly via the generic registerPortHandler — no hand-written models, no
// per-method boilerplate.
type HotelFacade interface {
	// SearchAvailability performs a hotel availability search.
	SearchAvailability(ctx context.Context, req *hotelxsd.HotelSearchAvailabilityReq) (*usecase.HotelSearchAvailabilityOutput, error)
	// Details fetches hotel details.
	Details(ctx context.Context, req *hotelxsd.HotelDetailsReq) (*usecase.HotelDetailsOutput, error)
	// MediaLinks fetches hotel media links.
	MediaLinks(ctx context.Context, req *hotelxsd.HotelMediaLinksReq) (*hotelxsd.HotelMediaLinksRsp, error)
	// Retrieve retrieves an existing hotel booking by record locator.
	Retrieve(ctx context.Context, req *hotelxsd.HotelRetrieveReq) (*hotelxsd.HotelRetrieveRsp, error)
	// Rules fetches hotel rate rules.
	Rules(ctx context.Context, req *hotelxsd.HotelRulesReq) (*hotelxsd.HotelRulesRsp, error)
	// UpsellSearch searches hotel upsell offers.
	UpsellSearch(ctx context.Context, req *hotelxsd.HotelUpsellDetailsReq) (*hotelxsd.HotelUpsellDetailsRsp, error)
	// Keywords searches hotels by keyword.
	Keywords(ctx context.Context, req *hotelxsd.HotelKeywordReq) (*hotelxsd.HotelKeywordRsp, error)
	// SuperShopper performs a hotel super shopper search.
	SuperShopper(ctx context.Context, req *hotelxsd.HotelSuperShopperReq) (*hotelxsd.HotelSuperShopperRsp, error)
}

// HotelHandler handles hotel-related HTTP endpoints.
type HotelHandler struct {
	facade    HotelFacade
	universal *usecase.UniversalFacade
}

// NewHotelHandler creates the hotel handler. universal proxies the cross-product unified
// booking/cancellation capabilities.
func NewHotelHandler(facade HotelFacade, universal *usecase.UniversalFacade) *HotelHandler {
	return &HotelHandler{facade: facade, universal: universal}
}

// RegisterRoutes registers the hotel endpoints.
// All routes (including /search and /details) are registered via the generic
// registerPortHandler: input validation is done in the facade layer, while timeout,
// trace injection, auth pass-through, and error convergence are carried uniformly by the
// generic handler.
func (h *HotelHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/hotel/search", f.SearchAvailability)
	registerPortHandler(mux, apiBasePath+"/hotel/details", f.Details)
	registerPortHandler(mux, apiBasePath+"/hotel/media-links", f.MediaLinks)
	registerPortHandler(mux, apiBasePath+"/hotel/retrieve", f.Retrieve)
	registerPortHandler(mux, apiBasePath+"/hotel/rules", f.Rules)
	registerPortHandler(mux, apiBasePath+"/hotel/upsell", f.UpsellSearch)
	registerPortHandler(mux, apiBasePath+"/hotel/keywords", f.Keywords)
	registerPortHandler(mux, apiBasePath+"/hotel/super-shopper", f.SuperShopper)

	// Productized aliases: booking/cancellation is delegated to the UniversalRecord engine.
	registerPortHandler(mux, apiBasePath+"/hotel/book", h.universal.HotelCreateReservation)
	registerPortHandler(mux, apiBasePath+"/hotel/cancel", h.universal.HotelCancel)
}

// decodeJSONBody decodes the request body into a single JSON object within the size limit
// and rejects unknown fields.
func decodeJSONBody(r *http.Request, v interface{}) error {
	r.Body = io.NopCloser(io.LimitReader(r.Body, maxRequestBodyBytes))
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

// writeJSON writes a JSON response with the given status code (without the trace header,
// for edge cases without ctx such as 405).
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// writeJSONWithTrace first writes the global trace ID response header X-Trace-Id (taken
// from the trace_id in context), then writes the JSON response with the given status code.
// All business paths (success/error) go through this function so that callers always get
// the trace_id of the current call from the response header for correlation with logs and
// GDS messages.
func writeJSONWithTrace(w http.ResponseWriter, status int, ctx context.Context, payload interface{}) {
	if tid := trace.ID(ctx); tid != "" {
		w.Header().Set(HeaderTraceID, tid)
	}
	writeJSON(w, status, payload)
}
