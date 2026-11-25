// Package client wraps Travelport UAPI SOAP calls with auth pass-through
// and trace injection.
package client

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shuiyihan12/uapi-go/pkg/generated/common55"
	"github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
)

type responseTransportFunc func(*http.Request) (*http.Response, error)

func (f responseTransportFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

// An upstream gateway can return an HTML 404 page with HTTP 200. This must
// not be misreported as an empty SOAP Body or silently retried.
func TestHTMLResponseWithHTTP200IsNotAnEmptySOAPBody(t *testing.T) {
	const endpoint = "https://apac.universal-api.travelport.com/B2BGateway/connect/uAPI/UniversalRecordService"
	const html = `<!doctype html><html lang="en"><head><title>HTTP Status 404 – Not Found</title></head><body><h1>HTTP Status 404 – Not Found</h1></body></html>`
	var logs bytes.Buffer
	calls := 0
	client, err := NewSOAPClient(SOAPConfig{
		ServiceName: "UniversalRecordService",
		Logger:      logging.NewStdLogAdapter(log.New(&logs, "", 0)),
		Transport: responseTransportFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if req.URL.String() != endpoint {
				t.Errorf("unexpected direct GDS endpoint: %s", req.URL)
			}
			return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Request: req,
				Header: http.Header{"Content-Type": []string{"text/html;charset=utf-8"}},
				Body:   io.NopCloser(strings.NewReader(html))}, nil
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	_, err = client.Call(requestctx.WithRegion(context.Background(), "apac"), "UniversalRecordRetrieve", testSOAPRequest{})
	if err == nil || !strings.Contains(err.Error(), "non-SOAP response") || !strings.Contains(err.Error(), `root="html"`) || strings.Contains(err.Error(), "body is empty") {
		t.Fatalf("misleading response error: %v", err)
	}
	for _, want := range []string{"http_status=200", "content_type=text/html", "endpoint=" + endpoint, "HTTP Status 404"} {
		if !strings.Contains(logs.String(), want) {
			t.Errorf("missing HTTP response evidence: %s", want)
		}
	}
	if calls != 1 {
		t.Fatalf("unexpected automatic retry: %d calls", calls)
	}
}

func TestExtractSOAPBodyDistinguishesNonSOAPAndEmptyEnvelope(t *testing.T) {
	for _, tc := range []struct{ body, errorText string }{
		{`<html><body>404</body></html>`, "non-SOAP response"},
		{`<html><Body><Response/></Body></html>`, "non-SOAP response"},
		{`<SOAP:Envelope xmlns:SOAP="http://schemas.xmlsoap.org/soap/envelope/"><SOAP:Body/></SOAP:Envelope>`, "SOAP response body is empty"},
	} {
		if _, err := extractSOAPBody([]byte(tc.body)); err == nil || !strings.Contains(err.Error(), tc.errorText) {
			t.Errorf("wanted %q, got %v", tc.errorText, err)
		}
	}
}

// testSOAPRequest is a minimal request body used to verify envelope
// serialization.
type testSOAPRequest struct {
	XMLName xml.Name `xml:"testReq"`
}

// captureTransport records outgoing request headers so tests can assert
// auth and Content-Type injection.
type captureTransport struct {
	header http.Header
}

// RoundTrip records the request headers and returns a successful response
// with a valid SOAP Body.
func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.header = req.Header.Clone()
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`<SOAP:Envelope xmlns:SOAP="http://schemas.xmlsoap.org/soap/envelope/"><SOAP:Body><testRsp>ok</testRsp></SOAP:Body></SOAP:Envelope>`)),
		Request:    req,
	}, nil
}

// TestAuthorizationPassedThroughFromContext verifies that the Authorization
// header supplied by the caller flows through the context to the downstream
// request.
func TestAuthorizationPassedThroughFromContext(t *testing.T) {
	base := &captureTransport{}

	ctx := requestctx.WithAuthorization(context.Background(), "Basic token")

	client, err := NewSOAPClient(SOAPConfig{
		BaseEndpoint: "https://example.com/uAPI",
		ServiceName:  "HotelService",
		Transport:    base,
	})
	if err != nil {
		t.Fatalf("new SOAP client: %v", err)
	}

	// Call reads the Authorization straight from ctx, with no
	// construction-time configuration involved.
	if _, err := client.Call(ctx, "TestOperation", testSOAPRequest{}); err != nil {
		t.Fatalf("call: %v", err)
	}

	if got := base.header.Get("Authorization"); got != "Basic token" {
		t.Fatalf("unexpected Authorization header: %q", got)
	}
	if got := base.header.Get("Content-Type"); got != "text/xml;charset=UTF-8" {
		t.Fatalf("unexpected Content-Type header: %q", got)
	}
}

// TestStandardOutboundSOAPHeaders verifies the fixed B2BGateway header set:
// SOAPAction must be present even when empty, and the caching/long-connection
// headers must be explicit.
func TestStandardOutboundSOAPHeaders(t *testing.T) {
	base := &captureTransport{}
	client, err := NewSOAPClient(SOAPConfig{
		BaseEndpoint: "https://example.com/uAPI",
		ServiceName:  "HotelService",
		Transport:    base,
		UserAgent:    "TINGLV-GOSDK-v0.55.5",
	})
	if err != nil {
		t.Fatalf("new SOAP client: %v", err)
	}
	if _, err := client.Call(context.Background(), "TestOperation", testSOAPRequest{}); err != nil {
		t.Fatalf("call: %v", err)
	}
	for header, want := range map[string]string{
		"Accept":        "text/xml",
		"Content-Type":  "text/xml;charset=UTF-8",
		"User-Agent":    "TINGLV-GOSDK-v0.55.5",
		"Cache-Control": "no-cache",
		"Pragma":        "no-cache",
		"Connection":    "keep-alive",
	} {
		if got := base.header.Get(header); got != want {
			t.Errorf("unexpected %s header: got %q, want %q", header, got, want)
		}
	}
	if values, ok := base.header["Soapaction"]; !ok || len(values) != 1 || values[0] != "" {
		t.Errorf("SOAPAction must be present and empty, got %q (present=%v)", values, ok)
	}
}

// TestEmptyUserAgentOmitsHeader verifies that no User-Agent is injected when
// configuration leaves it empty (no hardcoded fallback in code).
func TestEmptyUserAgentOmitsHeader(t *testing.T) {
	base := &captureTransport{}
	client, err := NewSOAPClient(SOAPConfig{
		BaseEndpoint: "https://example.com/uAPI",
		ServiceName:  "HotelService",
		Transport:    base,
	})
	if err != nil {
		t.Fatalf("new SOAP client: %v", err)
	}
	if _, err := client.Call(context.Background(), "TestOperation", testSOAPRequest{}); err != nil {
		t.Fatalf("call: %v", err)
	}
	if values, ok := base.header["User-Agent"]; ok {
		t.Errorf("User-Agent must be omitted when unconfigured, got %q", values)
	}
}

// TestResolveEndpointUsesRegion verifies that the region takes precedence
// when building the Production endpoint.
func TestResolveEndpointUsesRegion(t *testing.T) {
	client, err := NewSOAPClient(SOAPConfig{
		BaseEndpoint: "https://apac.universal-api.travelport.com/B2BGateway/connect/uAPI",
		ServiceName:  "AirService",
	})
	if err != nil {
		t.Fatalf("new SOAP client: %v", err)
	}

	if got := client.resolveEndpoint("emea"); got != "https://emea.universal-api.travelport.com/B2BGateway/connect/uAPI/AirService" {
		t.Fatalf("unexpected region endpoint: %q", got)
	}
	if got := client.resolveEndpoint(""); got != "https://apac.universal-api.travelport.com/B2BGateway/connect/uAPI/AirService" {
		t.Fatalf("unexpected default endpoint: %q", got)
	}
}

// TestSOAPClientCallPostsEnvelopeDirectly verifies that Call POSTs the
// envelope directly and extracts the response body correctly.
func TestSOAPClientCallPostsEnvelopeDirectly(t *testing.T) {
	var (
		requestCount int
		gotAuth      string
		gotEndpoint  string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if r.Method != http.MethodPost {
			t.Fatalf("SOAP client must post directly to service endpoint, got %s", r.Method)
		}
		gotAuth = r.Header.Get("Authorization")
		gotEndpoint = r.URL.Path
		if got := r.Header.Get("SOAPAction"); got != "http://localhost:8080/kestrel/HotelService" {
			t.Fatalf("unexpected SOAPAction header: %q", got)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		bodyText := string(body)
		if !strings.Contains(bodyText, "<soapenv:Envelope") || !strings.Contains(bodyText, "<testReq>") {
			t.Fatalf("unexpected SOAP envelope: %s", bodyText)
		}

		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(`<SOAP:Envelope xmlns:SOAP="http://schemas.xmlsoap.org/soap/envelope/"><SOAP:Body><testRsp>ok</testRsp></SOAP:Body></SOAP:Envelope>`))
	}))
	defer server.Close()

	// Auth arrives from the caller via context (mirroring the HTTP
	// Authorization request header).
	ctx := requestctx.WithAuthorization(context.Background(), "Basic token")

	client, err := NewSOAPClient(SOAPConfig{
		BaseEndpoint: server.URL,
		ServiceName:  "HotelService",
		SOAPAction:   "http://localhost:8080/kestrel/HotelService",
	})
	if err != nil {
		t.Fatalf("new SOAP client: %v", err)
	}

	response, err := client.Call(ctx, "TestOperation", testSOAPRequest{})
	if err != nil {
		t.Fatalf("call SOAP client: %v", err)
	}
	if strings.TrimSpace(string(response)) != "<testRsp>ok</testRsp>" {
		t.Fatalf("unexpected response body: %s", response)
	}
	if requestCount != 1 {
		t.Fatalf("expected exactly one direct POST, got %d requests", requestCount)
	}
	if gotAuth != "Basic token" {
		t.Fatalf("unexpected Authorization header: %q", gotAuth)
	}
	// No region specified: the endpoint should be BaseEndpoint +
	// "/HotelService".
	if gotEndpoint != "/HotelService" {
		t.Fatalf("unexpected endpoint path: %q", gotEndpoint)
	}
}

// TestExtractSOAPBodyReturnsFaultError verifies that a response containing
// a SOAP Fault is parsed into an error.
func TestExtractSOAPBodyReturnsFaultError(t *testing.T) {
	_, err := extractSOAPBody([]byte(`<SOAP:Envelope xmlns:SOAP="http://schemas.xmlsoap.org/soap/envelope/"><SOAP:Body><SOAP:Fault><faultcode>soap:Client</faultcode><faultstring>Invalid request</faultstring></SOAP:Fault></SOAP:Body></SOAP:Envelope>`))
	if err == nil {
		t.Fatalf("expected SOAP fault to be returned as error")
	}
	if !strings.Contains(err.Error(), "Invalid request") {
		t.Fatalf("unexpected fault error: %v", err)
	}
}

// TestBuildEnvelopeNamespaceHoist locks the request envelope shape:
// the SOAP prefix is soapenv, the empty Header is omitted, and the body's
// namespaces are hoisted to the request root as readable prefixes (hotel,
// common) rather than repeated default xmlns declarations per element.
func TestBuildEnvelopeNamespaceHoist(t *testing.T) {
	client := &SOAPClient{}

	tb := common55.TypeBranchCode("P3735951")
	req := hotel.HotelDetailsReq{
		BaseHotelDetailsReq: hotel.BaseHotelDetailsReq{
			BaseReq: common55.BaseReq{
				BaseCoreReq: common55.BaseCoreReq{
					TargetBranch: &tb,
					TraceId:      strPtr("TP-f14c5698-a1ff-4613-b4b3-e47c6ba6f510"),
					BillingPointOfSaleInfo: common55.BillingPointOfSaleInfo{
						OriginApplication: "UAPI",
					},
				},
			},
			HotelProperty: hotel.HotelProperty{
				HotelChain: "WY",
				HotelCode:  "G5685",
			},
			HotelDetailsModifiers: &hotel.HotelDetailsModifiers{
				HotelStay: hotel.HotelStay{
					CheckinDate:  "2025-11-29",
					CheckoutDate: "2025-11-30",
				},
			},
		},
	}

	got, err := client.buildEnvelope(req)
	if err != nil {
		t.Fatalf("buildEnvelope returned error: %v", err)
	}
	out := string(got)

	checks := []struct {
		name string
		want string
	}{
		{"soapenv prefix on Envelope", `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">`},
		{"empty Header omitted", `<soapenv:Body>`},
		{"hotel + common prefixes declared once on request root", `<hotel:HotelDetailsReq xmlns:common="http://www.travelport.com/schema/common_v55_0" xmlns:hotel="http://www.travelport.com/schema/hotel_v55_0"`},
		{"common element uses prefix", `<common:BillingPointOfSaleInfo`},
		{"hotel element uses prefix", `<hotel:HotelProperty`},
		{"hotel child element uses prefix", `<hotel:CheckinDate>`},
	}
	for _, c := range checks {
		if !strings.Contains(out, c.want) {
			t.Errorf("check %q: expected output to contain %q\nfull output:\n%s", c.name, c.want, out)
		}
	}

	// The namespace prefixes are declared once on the request root, so no
	// element may carry a repeated default-namespace declaration.
	negatives := []struct {
		name string
		bad  string
	}{
		{"no per-element repeated common xmlns", `xmlns="http://www.travelport.com/schema/common_v55_0"`},
		{"no per-element repeated hotel xmlns", `xmlns="http://www.travelport.com/schema/hotel_v55_0"`},
	}
	for _, n := range negatives {
		if strings.Contains(out, n.bad) {
			t.Errorf("negative check %q: output must NOT contain %q\nfull output:\n%s", n.name, n.bad, out)
		}
	}
}

func strPtr(s string) *string { return &s }
