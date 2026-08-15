// Package api provides the HTTP interface layer for hotel search and details,
// plus unified error response handling.
package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// HeaderTraceID is the name of the HTTP header that carries the global trace ID (trace_id).
// Callers pass the trace_id of this call chain via the X-Trace-Id request header, and the
// gateway writes the trace_id actually used for this call back into the response header so
// that callers can correlate it with logs and GDS messages.
const HeaderTraceID = "X-Trace-Id"

// registerPortHandler registers an HTTP endpoint for a strongly typed GDS PortType operation.
//
// call corresponds directly to a PortType method on the service layer (of the form
// func(ctx, *Req) (*Rsp, error)):
//   - the request body is decoded from JSON into *Req (unknown fields are rejected);
//   - on success, *Rsp is written out as-is;
//   - on failure, everything converges to the flat error structure produced by writeError.
//
// Thanks to generics, all PortType operations share the same timeout, trace injection, and
// error handling, avoiding per-method boilerplate.
//
// Global trace ID (trace_id) propagation follows industry best practices (the W3C Trace
// Context / OpenTelemetry boundary pass-through principle): the X-Trace-Id HTTP header is
// the preferred channel for carrying the trace_id and is provided by the caller at the
// entry point; when the caller does not supply it, the gateway generates one. The gateway
// carries it through context into logs, outbound HTTP headers, and (as a fallback) the
// request body TraceId attribute, but does not overwrite a TraceId business value the
// caller filled in explicitly in the request body (see BaseCoreReq.InjectInfrastructure
// for the fallback semantics).
func registerPortHandler[TReq any, TRsp any](
	mux *http.ServeMux,
	path string,
	call func(ctx context.Context, req *TReq) (*TRsp, error),
) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeMethodNotAllowed(w)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), defaultRequestTimeout)
		defer cancel()
		// trace_id comes first from the caller's X-Trace-Id HTTP header; when absent the
		// gateway generates one (for log/upstream correlation).
		if h := strings.TrimSpace(r.Header.Get(HeaderTraceID)); h != "" {
			ctx = trace.WithTraceID(ctx, h)
		} else {
			ctx, _ = trace.Ensure(ctx)
		}
		// Request-scoped auth and region: carried in HTTP headers by the caller and passed
		// through to the SOAP call via context.
		ctx = requestctx.WithAuthorization(ctx, r.Header.Get(requestctx.HeaderAuthorization))
		ctx = requestctx.WithRegion(ctx, r.Header.Get(requestctx.HeaderRegion))

		var req TReq
		if err := decodeJSONBody(r, &req); err != nil {
			writeValidationError(w, ctx, err.Error(), "")
			return
		}

		resp, err := call(ctx, &req)
		if err != nil {
			writeError(w, ctx, err)
			return
		}

		writeJSONWithTrace(w, http.StatusOK, ctx, resp)
	})
}
