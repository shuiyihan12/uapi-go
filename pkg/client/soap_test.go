// Package client wraps Travelport UAPI SOAP calls with auth pass-through
// and trace injection.
package client

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
)

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
