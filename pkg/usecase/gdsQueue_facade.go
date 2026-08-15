// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the GdsQueue service (gdsqueue): methods map 1:1 to the GdsQueueServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	gdsqueuexsd "github.com/shuiyihan12/uapi-go/pkg/generated/gdsqueue"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	gdsQueuesvc "github.com/shuiyihan12/uapi-go/pkg/services/gdsQueue"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// GdsQueueFacade orchestrates the Travelport GdsQueue service use cases; methods map 1:1 to the service PortType.
type GdsQueueFacade struct {
	manager *manager.ServiceManager
}

// NewGdsQueueFacade creates the GdsQueue use-case layer.
func NewGdsQueueFacade(serviceManager *manager.ServiceManager) *GdsQueueFacade {
	return &GdsQueueFacade{manager: serviceManager}
}

// getService lazily retrieves the GdsQueue service client, handling nil manager and lookup failures uniformly.
func (f *GdsQueueFacade) getService() (*gdsQueuesvc.GdsQueueService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*gdsQueuesvc.GdsQueueService](f.manager, "gdsQueue")
	if err != nil {
		return nil, fmt.Errorf("failed to get gdsQueue service: %w", err)
	}
	return svc, nil
}

// GdsEnterQueue corresponds to GdsQueueServicePort.GdsEnterQueue。
func (f *GdsQueueFacade) GdsEnterQueue(ctx context.Context, req *gdsqueuexsd.GdsEnterQueueReq) (*gdsqueuexsd.GdsEnterQueueRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.GdsEnterQueue(ctx, req)
}

// GdsExitQueue corresponds to GdsQueueServicePort.GdsExitQueue。
func (f *GdsQueueFacade) GdsExitQueue(ctx context.Context, req *gdsqueuexsd.GdsExitQueueReq) (*gdsqueuexsd.GdsExitQueueRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.GdsExitQueue(ctx, req)
}

// GdsNextOnQueue corresponds to GdsQueueServicePort.GdsNextOnQueue。
func (f *GdsQueueFacade) GdsNextOnQueue(ctx context.Context, req *gdsqueuexsd.GdsNextOnQueueReq) (*gdsqueuexsd.GdsNextOnQueueRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.GdsNextOnQueue(ctx, req)
}

// GdsQueueAgentList corresponds to GdsQueueServicePort.GdsQueueAgentList。
func (f *GdsQueueFacade) GdsQueueAgentList(ctx context.Context, req *gdsqueuexsd.GdsQueueAgentListReq) (*gdsqueuexsd.GdsQueueAgentListRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.GdsQueueAgentList(ctx, req)
}

// GdsQueueCount corresponds to GdsQueueServicePort.GdsQueueCount。
func (f *GdsQueueFacade) GdsQueueCount(ctx context.Context, req *gdsqueuexsd.GdsQueueCountReq) (*gdsqueuexsd.GdsQueueCountRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.GdsQueueCount(ctx, req)
}

// GdsQueueList corresponds to GdsQueueServicePort.GdsQueueList。
func (f *GdsQueueFacade) GdsQueueList(ctx context.Context, req *gdsqueuexsd.GdsQueueListReq) (*gdsqueuexsd.GdsQueueListRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.GdsQueueList(ctx, req)
}

// GdsQueuePlace corresponds to GdsQueueServicePort.GdsQueuePlace。
func (f *GdsQueueFacade) GdsQueuePlace(ctx context.Context, req *gdsqueuexsd.GdsQueuePlaceReq) (*gdsqueuexsd.GdsQueuePlaceRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.GdsQueuePlace(ctx, req)
}

// GdsQueueRemove corresponds to GdsQueueServicePort.GdsQueueRemove。
func (f *GdsQueueFacade) GdsQueueRemove(ctx context.Context, req *gdsqueuexsd.GdsQueueRemoveReq) (*gdsqueuexsd.GdsQueueRemoveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.GdsQueueRemove(ctx, req)
}
