// Package client wraps Travelport UAPI SOAP calls with auth pass-through
// and trace injection.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/shuiyihan12/uapi-go/internal/logging"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
)

// SOAPClient wraps SOAP communication with Travelport UAPI services:
// request orchestration, auth pass-through and response parsing. The
// endpoint and Authorization are no longer fixed at construction: every Call
// reads the request-level region / auth from the context, dynamically builds
// the Production endpoint
// (https://<region>.universal-api.travelport.com/B2BGateway/connect/uAPI/<Service>)
// and forwards the Authorization header.
type SOAPClient struct {
	httpClient   *http.Client
	baseEndpoint string // default endpoint prefix (UAPI_ENDPOINT or the apac production environment); used when no region is specified
	serviceName  string // UAPI SOAP service name (e.g. AirService), appended to the endpoint
	timeout      time.Duration
	transport    http.RoundTripper
	logger       logging.Logger
	soapAction   string
}

// soapResponseEnvelope parses only the raw XML of the SOAP envelope Body and
// leaves further deserialization to the caller.
type soapResponseEnvelope struct {
	Body struct {
		Contents []byte `xml:",innerxml"`
	} `xml:"Body"`
}

// soapFault models the SOAP <Fault> node, used to detect GDS failures in
// responses.
type soapFault struct {
	XMLName     xml.Name    `xml:"Fault"`
	FaultCode   string      `xml:"faultcode"`
	FaultString string      `xml:"faultstring"`
	Detail      faultDetail `xml:"detail"`
}

// faultDetail carries the SOAP Fault <detail> content, where Travelport
// usually places ErrorInfo.
type faultDetail struct {
	ErrorInfo errorInfo `xml:"ErrorInfo"`
}

// errorInfo models Travelport's <common_v55_0:ErrorInfo>.
type errorInfo struct {
	Code          string `xml:"Code"`
	Service       string `xml:"Service"`
	Type          string `xml:"Type"`
	Description   string `xml:"Description"`
	TransactionID string `xml:"TransactionId"`
	TraceID       string `xml:"TraceId"`
}

// SOAPFaultError represents a business/system error returned by Travelport
// via a SOAP Fault. Unlike network errors and timeouts, these originate
// from GDS business processing (e.g. unavailable inventory, invalid
// parameters).
type SOAPFaultError struct {
	FaultCode     string // SOAP faultcode, e.g. Server.Business / Client.xxx
	FaultString   string // SOAP faultstring, e.g. WARNING ROOM RATE NOT AVAILABLE *
	Code          string // detail/ErrorInfo/Code, e.g. 5436
	Type          string // detail/ErrorInfo/Type, e.g. Business
	Service       string // detail/ErrorInfo/Service, e.g. HTLSVC
	Description   string // detail/ErrorInfo/Description
	TransactionID string
	TraceID       string
}

// Error implements the error interface.
func (e *SOAPFaultError) Error() string {
	if e.Code != "" && e.Description != "" {
		return fmt.Sprintf("SOAP fault %s [%s]: %s", e.FaultCode, e.Code, e.Description)
	}
	if e.Code != "" {
		return fmt.Sprintf("SOAP fault %s [%s]", e.FaultCode, e.Code)
	}
	return fmt.Sprintf("SOAP fault %s: %s", e.FaultCode, e.FaultString)
}

// Category returns the error category used to decide the HTTP status, log
// level and retry policy. It is derived from Travelport's <ErrorInfo>/Type
// field (not the numeric code):
//   - system: system failure, retryable (502 / Error log);
//   - business: business condition (e.g. unavailable inventory), not
//     retryable (422 / Warn log);
//   - client: caller error (invalid parameters etc.), not retryable
//     (422 / Warn log).
func (e *SOAPFaultError) Category() string {
	switch strings.ToLower(strings.TrimSpace(e.Type)) {
	case "system":
		return "system"
	case "client":
		return "client"
	case "business":
		return "business"
	}
	// Without a Type, fall back to the faultcode: Client.* means a client
	// error.
	if strings.HasPrefix(e.FaultCode, "Client") {
		return "client"
	}
	return "system"
}

// IsSystem reports whether the error is a system error (retryable).
func (e *SOAPFaultError) IsSystem() bool {
	return e.Category() == "system"
}

// Retryable reports whether the error is retryable; system errors are,
// business/client errors are not.
func (e *SOAPFaultError) Retryable() bool {
	return e.IsSystem()
}

// SOAPConfig holds the connection, timeout and endpoint configuration used
// to build a SOAPClient.
//
// Authorization and Region are REQUEST-LEVEL configuration: the caller sends
// them in every HTTP request header, they flow through context into Call to
// dynamically build the endpoint, so no Authorization field is kept here.
// This project does not retry (GDS write operations are not idempotent and
// blind retries would issue duplicate tickets/bookings), hence no retry
// options.
type SOAPConfig struct {
	// BaseEndpoint is the default endpoint prefix (UAPI_ENDPOINT or the apac
	// production environment), e.g.
	// https://apac.universal-api.travelport.com/B2BGateway/connect/uAPI.
	// It is used when the request carries no X-UAPI-Region; the full endpoint
	// is this prefix + "/" + ServiceName.
	BaseEndpoint string
	// ServiceName is the UAPI SOAP service name (e.g. AirService /
	// UniversalRecordService), appended to the endpoint.
	ServiceName       string
	Timeout           time.Duration
	ConnectionTimeout time.Duration
	ReadTimeout       time.Duration
	// MaxIdleConns caps the total idle keep-alive connections kept warm
	// across all upstream hosts, avoiding a TCP+TLS re-handshake per
	// request. Zero means the default (100).
	MaxIdleConns int
	// MaxIdleConnsPerHost caps idle keep-alive connections per host; it
	// cannot exceed MaxIdleConns. Zero means the default (100).
	MaxIdleConnsPerHost int
	// SkipTLSVerify defaults to false (verify certificates). Only enabled
	// explicitly via UAPI_SKIP_TLS_VERIFY=1 for private environments with
	// self-signed certificates; never implicitly tied to Environment.
	SkipTLSVerify bool
	Logger        logging.Logger
	// SOAPAction is optional; when empty, no SOAPAction HTTP header is set.
	SOAPAction string
	// Transport is optional; when nil, NewSOAPClient creates the default
	// transport from the timeout/TLS settings. Mainly used to inject a
	// capturing transport in tests.
	Transport http.RoundTripper
}

// NewSOAPClient creates and initializes a SOAPClient from SOAPConfig.
func NewSOAPClient(config SOAPConfig) (*SOAPClient, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = 30 * time.Second
	}
	if config.MaxIdleConns <= 0 {
		config.MaxIdleConns = 100
	}
	if config.MaxIdleConnsPerHost <= 0 {
		config.MaxIdleConnsPerHost = 100
	}
	if config.Logger == nil {
		config.Logger = logging.NewDefaultLogger()
	}

	transport := config.Transport
	if transport == nil {
		transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: config.ConnectionTimeout,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.SkipTLSVerify,
			},
			MaxIdleConns:          config.MaxIdleConns,
			MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: config.ReadTimeout,
		}
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return &SOAPClient{
		httpClient:   httpClient,
		baseEndpoint: config.BaseEndpoint,
		serviceName:  config.ServiceName,
		timeout:      config.Timeout,
		transport:    transport,
		logger:       config.Logger,
		soapAction:   config.SOAPAction,
	}, nil
}

// Call performs one SOAP call.
//
// Before the call it reads the global trace ID from the context (generating
// one if absent) and:
//   - injects it into the request body's TraceId attribute (a standard
//     Travelport BaseCoreReq attribute, UUID v4 form without the TP-
//     prefix), correlating with the GDS via the TraceId echoed in
//     SOAPFaults;
//   - forwards the HTTP X-Trace-Id header so gateways/proxies can correlate;
//   - tags the request/response XML logs with the trace_id field;
//   - reads the request-level Authorization (carried by the caller in the
//     HTTP header) and Region from the context, dynamically builds the
//     Production endpoint and forwards the Authorization header (Travelport
//     gateway authentication).
func (c *SOAPClient) Call(ctx context.Context, action string, request interface{}) ([]byte, error) {
	ctx, traceID := trace.Ensure(ctx)
	logger := c.logger.WithContext(ctx)

	rawRegion := requestctx.Region(ctx)
	region := requestctx.NormalizeRegion(rawRegion)
	if rawRegion != "" && region == "" {
		// An invalid X-UAPI-Region would otherwise silently fall back to the
		// default endpoint — a typo (e.g. "apa") then lands requests on the
		// wrong geography and surfaces only as confusing upstream auth
		// failures. Warn explicitly instead.
		logger.Warn("invalid X-UAPI-Region header ignored; falling back to the default endpoint",
			zap.String("region", rawRegion),
			zap.String("default_endpoint", c.resolveEndpoint("")))
	}
	endpoint := c.resolveEndpoint(region)
	authorization := strings.TrimSpace(requestctx.Authorization(ctx))

	payload, err := c.buildEnvelope(request)
	if err != nil {
		return nil, err
	}

	logger.Raw(fmt.Sprintf("[GDS REQUEST] operation=%s trace_id=%s endpoint=%s\n%s", action, traceID, endpoint, string(payload)))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create SOAP request: %v", err)
	}
	httpReq.Header.Set("Accept", "text/xml")
	httpReq.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	if c.soapAction != "" {
		httpReq.Header.Set("SOAPAction", c.soapAction)
	}
	// Forward the caller's Travelport auth header verbatim (same semantics
	// as the Java SDK's AuthenticationConfigurer).
	if authorization != "" {
		httpReq.Header.Set("Authorization", authorization)
	}
	// Forward the global trace ID in an HTTP header for gateway/proxy-side
	// correlation.
	httpReq.Header.Set("X-Trace-Id", traceID)

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("SOAP call failed: %v", err)
	}
	defer httpRsp.Body.Close()

	responseBody, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read SOAP response: %v", err)
	}
	if httpRsp.StatusCode < http.StatusOK || httpRsp.StatusCode >= http.StatusMultipleChoices {
		logger.Error("GDS SOAP call returned non-2xx",
			zap.String("operation", action),
			zap.String("status", httpRsp.Status))
		return nil, fmt.Errorf("SOAP call failed with HTTP %s: %s", httpRsp.Status, string(responseBody))
	}

	logger.Raw(fmt.Sprintf("[GDS RESPONSE] operation=%s trace_id=%s\n%s", action, traceID, string(responseBody)))

	bodyContents, err := extractSOAPBody(responseBody)
	if err != nil {
		return nil, err
	}

	return bodyContents, nil
}

// buildEnvelope serializes the request body into a full SOAP envelope.
// Travelport UAPI request auth travels in the HTTP Authorization header and
// the SOAP <Header> is empty (matching the official samples); the trace ID
// rides on the request body's TraceId attribute rather than a custom SOAP
// header.
func (c *SOAPClient) buildEnvelope(request interface{}) ([]byte, error) {
	body, err := xml.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SOAP request body: %v", err)
	}

	envelope := bytes.NewBuffer(nil)
	envelope.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	envelope.WriteString(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">`)
	envelope.WriteString(`<soapenv:Header/>`)
	envelope.WriteString(`<soapenv:Body>`)
	envelope.Write(body)
	envelope.WriteString(`</soapenv:Body>`)
	envelope.WriteString(`</soapenv:Envelope>`)
	return envelope.Bytes(), nil
}

// extractSOAPBody parses the SOAP envelope and extracts the Body contents,
// returning a SOAPFaultError when a Fault is present.
func extractSOAPBody(responseBody []byte) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(responseBody))
	decoder.CharsetReader = charset.NewReaderLabel

	var envelope soapResponseEnvelope
	if err := decoder.Decode(&envelope); err != nil {
		return nil, fmt.Errorf("failed to parse SOAP envelope: %v", err)
	}
	if len(bytes.TrimSpace(envelope.Body.Contents)) == 0 {
		return nil, fmt.Errorf("SOAP response body is empty")
	}
	if faultErr := parseSOAPFault(envelope.Body.Contents); faultErr != nil {
		return nil, faultErr
	}
	return envelope.Body.Contents, nil
}

// parseSOAPFault tries to parse the Body contents as a SOAP Fault and
// returns nil when it cannot.
func parseSOAPFault(bodyContents []byte) error {
	var fault soapFault
	decoder := xml.NewDecoder(bytes.NewReader(bodyContents))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&fault); err != nil || fault.XMLName.Local != "Fault" {
		return nil
	}

	return &SOAPFaultError{
		FaultCode:     strings.TrimSpace(fault.FaultCode),
		FaultString:   strings.TrimSpace(fault.FaultString),
		Code:          strings.TrimSpace(fault.Detail.ErrorInfo.Code),
		Type:          strings.TrimSpace(fault.Detail.ErrorInfo.Type),
		Service:       strings.TrimSpace(fault.Detail.ErrorInfo.Service),
		Description:   strings.TrimSpace(fault.Detail.ErrorInfo.Description),
		TransactionID: strings.TrimSpace(fault.Detail.ErrorInfo.TransactionID),
		TraceID:       strings.TrimSpace(fault.Detail.ErrorInfo.TraceID),
	}
}

// resolveEndpoint builds the full UAPI endpoint from the request-level
// region and the default base.
//
// With a valid region (americas / apac / emea) it uses the Production
// template:
//
//	https://<region>.universal-api.travelport.com/B2BGateway/connect/uAPI/<ServiceName>
//
// Otherwise it falls back to the default base (UAPI_ENDPOINT or the apac
// production environment) + "/" + ServiceName.
func (c *SOAPClient) resolveEndpoint(region string) string {
	if region != "" {
		return "https://" + region + ".universal-api.travelport.com/B2BGateway/connect/uAPI/" + c.serviceName
	}
	return strings.TrimRight(c.baseEndpoint, "/") + "/" + c.serviceName
}

// Close releases the underlying transport's connections.
func (c *SOAPClient) Close() error {
	if c.transport != nil {
		if t, ok := c.transport.(interface{ CloseIdleConnections() }); ok {
			t.CloseIdleConnections()
		}
	}
	return nil
}
