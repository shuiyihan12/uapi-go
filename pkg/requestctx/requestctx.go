// Package requestctx carries per-HTTP-request caller configuration:
//   - Authorization: the Travelport auth header supplied by the caller in the
//     request headers, passed through to UAPI verbatim by the gateway.
//   - Region: the region (americas / apac / emea) selected by the caller in
//     the request headers, used to build the Production endpoint.
//
// Neither belongs in process-level global configuration: a single gateway
// serving multiple Travelport accounts / regions requires each request to
// carry its own values. They flow from the HTTP entry point through context
// down to the SOAP calls in pkg/client.
package requestctx

import (
	"context"
	"strings"
)

type ctxKey string

const (
	authKey   ctxKey = "uapi_authorization"
	regionKey ctxKey = "uapi_region"
)

// HeaderAuthorization is the HTTP header carrying the caller's Travelport
// auth header.
const HeaderAuthorization = "Authorization"

// HeaderRegion is the HTTP header in which the caller selects the region.
const HeaderRegion = "X-UAPI-Region"

// validProductionRegions lists the region subdomains available in the
// Production environment (case-insensitive). Source: Travelport UAPI
// endpoint documentation (Production only; the sharedUprofile special path
// is ignored).
var validProductionRegions = map[string]string{
	"americas": "americas",
	"apac":     "apac",
	"emea":     "emea",
}

// WithAuthorization stores the caller-supplied Authorization value in the
// context.
func WithAuthorization(ctx context.Context, auth string) context.Context {
	return context.WithValue(ctx, authKey, auth)
}

// Authorization returns the Authorization value from the context, or an
// empty string if absent.
func Authorization(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(authKey).(string); ok {
		return v
	}
	return ""
}

// WithRegion stores the caller-supplied region in the context.
func WithRegion(ctx context.Context, region string) context.Context {
	return context.WithValue(ctx, regionKey, region)
}

// Region returns the raw (un-normalized) region from the context, or an
// empty string if absent.
func Region(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(regionKey).(string); ok {
		return v
	}
	return ""
}

// NormalizeRegion validates and normalizes a region name, accepting only the
// Production regions americas / apac / emea (case-insensitive,
// whitespace-trimmed). An empty return value means "no region specified; use
// the default endpoint (UAPI_ENDPOINT or the apac production environment)".
func NormalizeRegion(region string) string {
	region = strings.TrimSpace(strings.ToLower(region))
	if region == "" {
		return ""
	}
	if v, ok := validProductionRegions[region]; ok {
		return v
	}
	return ""
}
