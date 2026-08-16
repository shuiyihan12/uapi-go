// Package terminal provides the SOAP client implementation for the Travelport Terminal
// service (the terminal-mapped namespace).
// Its TerminalServicePort interface corresponds one-to-one with the WSDL *PortType operations.
// Request/response types come directly from the generator-produced terminal package
// (script generated; do not edit by hand).
// This package keeps no hand-written request models; all infrastructure fields
// are injected uniformly by prepareReq before sending.
package terminal

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	terminalxsd "github.com/shuiyihan12/uapi-go/pkg/generated/terminal"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// TerminalServicePort mirrors the *PortType operations of the terminal-mapped namespace one-to-one,
// making it easy to swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated terminal package.
type TerminalServicePort interface {
	// CreateTerminalSession corresponds to the CreateTerminalSessionReq operation of the Terminal service.
	CreateTerminalSession(ctx context.Context, req *terminalxsd.CreateTerminalSessionReq) (*terminalxsd.CreateTerminalSessionRsp, error)
	// EndTerminalSession corresponds to the EndTerminalSessionReq operation of the Terminal service.
	EndTerminalSession(ctx context.Context, req *terminalxsd.EndTerminalSessionReq) (*terminalxsd.EndTerminalSessionRsp, error)
	// Terminal corresponds to the TerminalReq operation of the Terminal service.
	Terminal(ctx context.Context, req *terminalxsd.TerminalReq) (*terminalxsd.TerminalRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// TerminalService is the SOAP implementation of TerminalServicePort.
type TerminalService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *TerminalService must satisfy the TerminalServicePort interface.
var _ TerminalServicePort = (*TerminalService)(nil)

// prepareReq ensures the call carries a trace_id and injects one as a fallback
// when the request body's TraceId is empty. Authorization/business fields such
// as the billing point of sale and TargetBranch are no longer injected by code;
// the caller (API user) must provide them explicitly in the request body.
// Authentication (Authorization) is supplied by the caller as an HTTP request
// header and passed through the context to the SOAP call (startup-time
// environment variables are no longer used). Injection follows a "fallback"
// strategy: InjectInfrastructure only fills the request body's TraceId with the
// trace_id when it is empty; a TraceId business value already set by the caller
// is not overwritten. Requests without an InjectInfrastructure implementation
// are skipped (the call is unaffected).
func prepareReq(ctx context.Context, req any) context.Context {
	ctx, traceID := trace.Ensure(ctx)
	if inj, ok := req.(interface{ InjectInfrastructure(traceID string) }); ok {
		inj.InjectInfrastructure(traceID)
	}
	return ctx
}

// NewTerminalService builds a Terminal service client from the given SOAP configuration and logger.
func NewTerminalService(config client.SOAPConfig, logger logging.Logger) (*TerminalService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "terminal-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create terminal service client: %w", err)
	}

	return &TerminalService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// CreateTerminalSession issues the CreateTerminalSessionReq SOAP call and returns the strongly typed response.
func (s *TerminalService) CreateTerminalSession(ctx context.Context, req *terminalxsd.CreateTerminalSessionReq) (*terminalxsd.CreateTerminalSessionRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[terminalxsd.CreateTerminalSessionRsp](s.client, ctx, "CreateTerminalSession", req)
}

// EndTerminalSession issues the EndTerminalSessionReq SOAP call and returns the strongly typed response.
func (s *TerminalService) EndTerminalSession(ctx context.Context, req *terminalxsd.EndTerminalSessionReq) (*terminalxsd.EndTerminalSessionRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[terminalxsd.EndTerminalSessionRsp](s.client, ctx, "EndTerminalSession", req)
}

// Terminal issues the TerminalReq SOAP call and returns the strongly typed response.
func (s *TerminalService) Terminal(ctx context.Context, req *terminalxsd.TerminalReq) (*terminalxsd.TerminalRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[terminalxsd.TerminalRsp](s.client, ctx, "Terminal", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType
// that performs a single SOAP call and decodes it into the strongly typed
// response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *TerminalService) Close() error {
	return s.client.Close()
}
