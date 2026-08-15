// Package api provides the HTTP interface layer for hotel search and details,
// plus unified error response handling.
package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// ErrorResponse is the unified API error response structure.
// All failures (input validation, upstream business errors, system errors, timeouts)
// converge to this structure; callers can handle errors programmatically by code and
// correlate logs and GDS messages via trace_id.
type ErrorResponse struct {
	// Code is the machine-readable error code.
	Code string `json:"code"`
	// Message is the human-readable message.
	Message string `json:"message"`
	// TraceID is the global trace ID of this request.
	TraceID string `json:"trace_id,omitempty"`
	// Details carries additional error context.
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error codes: let callers distinguish error categories programmatically.
const (
	CodeInvalidRequest        = "INVALID_REQUEST"         // input validation failure (HTTP 400)
	CodeProviderBusinessError = "PROVIDER_BUSINESS_ERROR" // GDS business/client error (HTTP 422)
	CodeProviderSystemError   = "PROVIDER_SYSTEM_ERROR"   // GDS system error (HTTP 502)
	CodeUpstreamError         = "UPSTREAM_ERROR"          // upstream network/unknown error (HTTP 502)
	CodeTimeout               = "TIMEOUT"                 // request timeout (HTTP 504)
	CodeMethodNotAllowed      = "METHOD_NOT_ALLOWED"      // method not allowed (HTTP 405)
)

// writeError maps an internal error to the unified error response.
//   - input validation failure → 400 INVALID_REQUEST;
//   - Travelport business/client SOAP Fault (e.g. inventory unavailable, invalid
//     parameters) → 422 PROVIDER_BUSINESS_ERROR, exposing code/type/service;
//   - Travelport system SOAP Fault → 502 PROVIDER_SYSTEM_ERROR;
//   - upstream network or other unknown errors → 502 UPSTREAM_ERROR;
//   - request timeout (context.DeadlineExceeded) → 504 TIMEOUT.
func writeError(w http.ResponseWriter, ctx context.Context, err error) {
	resp := ErrorResponse{
		Code:    CodeUpstreamError,
		Message: err.Error(),
		TraceID: trace.ID(ctx),
	}
	status := http.StatusBadGateway

	var verr *usecase.ValidationError
	if errors.As(err, &verr) {
		status = http.StatusBadRequest
		resp.Code = CodeInvalidRequest
		if verr.Field != "" {
			resp.Details = map[string]interface{}{"field": verr.Field}
		}
		writeJSONWithTrace(w, status, ctx, resp)
		return
	}

	// Travelport SOAP Fault: distinguish business vs. system by <ErrorInfo>/Type and
	// flatten provider fields directly.
	var fault *client.SOAPFaultError
	if errors.As(err, &fault) {
		if fault.IsSystem() {
			status = http.StatusBadGateway
			resp.Code = CodeProviderSystemError
		} else {
			status = http.StatusUnprocessableEntity
			resp.Code = CodeProviderBusinessError
		}
		payload := map[string]interface{}{
			"code":        fault.Code,
			"type":        fault.Type,
			"service":     fault.Service,
			"fault_code":  fault.FaultCode,
			"description": fault.Description,
			"trace_id":    trace.ID(ctx),
		}
		if fault.TransactionID != "" {
			payload["transaction_id"] = fault.TransactionID
		}
		writeJSONWithTrace(w, status, ctx, payload)
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		status = http.StatusGatewayTimeout
		resp.Code = CodeTimeout
	}

	writeJSONWithTrace(w, status, ctx, resp)
}

// writeValidationError writes an input validation error (HTTP 400).
func writeValidationError(w http.ResponseWriter, ctx context.Context, message, field string) {
	details := map[string]interface{}{}
	if field != "" {
		details["field"] = field
	}
	writeJSONWithTrace(w, http.StatusBadRequest, ctx, ErrorResponse{
		Code:    CodeInvalidRequest,
		Message: message,
		TraceID: trace.ID(ctx),
		Details: details,
	})
}

// writeMethodNotAllowed writes a method-not-allowed error (HTTP 405).
func writeMethodNotAllowed(w http.ResponseWriter) {
	writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
		Code:    CodeMethodNotAllowed,
		Message: "method not allowed",
	})
}
