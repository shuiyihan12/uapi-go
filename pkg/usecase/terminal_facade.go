// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the Terminal service (terminal): methods map 1:1 to the TerminalServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	terminalxsd "github.com/shuiyihan12/uapi-go/pkg/generated/terminal"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	terminalsvc "github.com/shuiyihan12/uapi-go/pkg/services/terminal"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// TerminalFacade orchestrates the Travelport Terminal service use cases; methods map 1:1 to the service PortType.nalServicePort。
type TerminalFacade struct {
	manager *manager.ServiceManager
}

// NewTerminalFacade creates the Terminal use-case layer.
func NewTerminalFacade(serviceManager *manager.ServiceManager) *TerminalFacade {
	return &TerminalFacade{manager: serviceManager}
}

// getService lazily retrieves the Terminal service client, handling nil manager and lookup failures uniformly.
func (f *TerminalFacade) getService() (*terminalsvc.TerminalService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*terminalsvc.TerminalService](f.manager, "terminal")
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal service: %w", err)
	}
	return svc, nil
}

// CreateTerminalSession corresponds to TerminalServicePort.CreateTerminalSession。
func (f *TerminalFacade) CreateTerminalSession(ctx context.Context, req *terminalxsd.CreateTerminalSessionReq) (*terminalxsd.CreateTerminalSessionRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.CreateTerminalSession(ctx, req)
}

// EndTerminalSession corresponds to TerminalServicePort.EndTerminalSession。
func (f *TerminalFacade) EndTerminalSession(ctx context.Context, req *terminalxsd.EndTerminalSessionReq) (*terminalxsd.EndTerminalSessionRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.EndTerminalSession(ctx, req)
}

// Terminal corresponds to TerminalServicePort.Terminal。
func (f *TerminalFacade) Terminal(ctx context.Context, req *terminalxsd.TerminalReq) (*terminalxsd.TerminalRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.Terminal(ctx, req)
}
