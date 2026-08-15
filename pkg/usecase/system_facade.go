// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the System service (system_v32_0): methods map 1:1 to the SystemServicePort,
// It only retrieves the service client and delegates; infrastructure fields (trace ID,
// billing point of sale) are injected uniformly by the service layer before sending,
package usecase

import (
	"context"
	"errors"
	"fmt"

	systemxsd "github.com/shuiyihan12/uapi-go/pkg/generated/system"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	systemsvc "github.com/shuiyihan12/uapi-go/pkg/services/system"
)

// SystemFacade orchestrates the Travelport System service use cases; methods map 1:1 to the service PortType.
type SystemFacade struct {
	manager *manager.ServiceManager
}

// NewSystemFacade creates the System use-case layer.
func NewSystemFacade(serviceManager *manager.ServiceManager) *SystemFacade {
	return &SystemFacade{manager: serviceManager}
}

// getService lazily retrieves the System service client, handling nil manager and lookup failures uniformly.
func (f *SystemFacade) getService() (*systemsvc.SystemService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*systemsvc.SystemService](f.manager, "system")
	if err != nil {
		return nil, fmt.Errorf("failed to get system service: %w", err)
	}
	return svc, nil
}

// Ping probes service connectivity (SystemPingPortType).
func (f *SystemFacade) Ping(ctx context.Context, req *systemxsd.PingReq) (*systemxsd.PingRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	return svc.Ping(ctx, req)
}

// Info retrieves gateway/system information (SystemInfoPortType).
func (f *SystemFacade) Info(ctx context.Context, req *systemxsd.SystemInfoReq) (*systemxsd.SystemInfoRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	return svc.Info(ctx, req)
}

// Time retrieves the current GDS time (SystemTimePortType).
func (f *SystemFacade) Time(ctx context.Context, req *systemxsd.TimeReq) (*systemxsd.TimeRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	return svc.Time(ctx, req)
}

// ExternalCacheAccess accesses the external cache (ExternalCacheAccessPortType).
func (f *SystemFacade) ExternalCacheAccess(ctx context.Context, req *systemxsd.ExternalCacheAccessReq) (*systemxsd.ExternalCacheAccessRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	return svc.ExternalCacheAccess(ctx, req)
}
