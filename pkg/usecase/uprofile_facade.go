// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the Uprofile service (uprofile): methods map 1:1 to the UprofileServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	uprofilexsd "github.com/shuiyihan12/uapi-go/pkg/generated/uprofile"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	uprofilesvc "github.com/shuiyihan12/uapi-go/pkg/services/uprofile"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// UprofileFacade orchestrates the Travelport Uprofile service use cases; methods map 1:1 to the service PortType.
type UprofileFacade struct {
	manager *manager.ServiceManager
}

// NewUprofileFacade creates the Uprofile use-case layer.
func NewUprofileFacade(serviceManager *manager.ServiceManager) *UprofileFacade {
	return &UprofileFacade{manager: serviceManager}
}

// getService lazily retrieves the Uprofile service client, handling nil manager and lookup failures uniformly.
func (f *UprofileFacade) getService() (*uprofilesvc.UprofileService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*uprofilesvc.UprofileService](f.manager, "uprofile")
	if err != nil {
		return nil, fmt.Errorf("failed to get uprofile service: %w", err)
	}
	return svc, nil
}

// ProfileChildSearch corresponds to UprofileServicePort.ProfileChildSearch。
func (f *UprofileFacade) ProfileChildSearch(ctx context.Context, req *uprofilexsd.ProfileChildSearchReq) (*uprofilexsd.ProfileChildSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileChildSearch(ctx, req)
}

// ProfileCreate corresponds to UprofileServicePort.ProfileCreate。
func (f *UprofileFacade) ProfileCreate(ctx context.Context, req *uprofilexsd.ProfileCreateReq) (*uprofilexsd.ProfileCreateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileCreate(ctx, req)
}

// ProfileCreateField corresponds to UprofileServicePort.ProfileCreateField。
func (f *UprofileFacade) ProfileCreateField(ctx context.Context, req *uprofilexsd.ProfileCreateFieldReq) (*uprofilexsd.ProfileCreateFieldRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileCreateField(ctx, req)
}

// ProfileCreateHierarchyLevel corresponds to UprofileServicePort.ProfileCreateHierarchyLevel。
func (f *UprofileFacade) ProfileCreateHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileCreateHierarchyLevelReq) (*uprofilexsd.ProfileCreateHierarchyLevelRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileCreateHierarchyLevel(ctx, req)
}

// ProfileCreateTags corresponds to UprofileServicePort.ProfileCreateTags。
func (f *UprofileFacade) ProfileCreateTags(ctx context.Context, req *uprofilexsd.ProfileCreateTagsReq) (*uprofilexsd.ProfileCreateTagsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileCreateTags(ctx, req)
}

// ProfileDelete corresponds to UprofileServicePort.ProfileDelete。
func (f *UprofileFacade) ProfileDelete(ctx context.Context, req *uprofilexsd.ProfileDeleteReq) (*uprofilexsd.ProfileDeleteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileDelete(ctx, req)
}

// ProfileDeleteHierarchyLevel corresponds to UprofileServicePort.ProfileDeleteHierarchyLevel。
func (f *UprofileFacade) ProfileDeleteHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileDeleteHierarchyLevelReq) (*uprofilexsd.ProfileDeleteHierarchyLevelRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileDeleteHierarchyLevel(ctx, req)
}

// ProfileDeleteTag corresponds to UprofileServicePort.ProfileDeleteTag。
func (f *UprofileFacade) ProfileDeleteTag(ctx context.Context, req *uprofilexsd.ProfileDeleteTagReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileDeleteTag(ctx, req)
}

// ProfileModify corresponds to UprofileServicePort.ProfileModify。
func (f *UprofileFacade) ProfileModify(ctx context.Context, req *uprofilexsd.ProfileModifyReq) (*uprofilexsd.ProfileModifyRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModify(ctx, req)
}

// ProfileModifyBridgeBranches corresponds to UprofileServicePort.ProfileModifyBridgeBranches。
func (f *UprofileFacade) ProfileModifyBridgeBranches(ctx context.Context, req *uprofilexsd.ProfileModifyBridgeBranchesReq) (*uprofilexsd.ProfileModifyBridgeBranchesRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModifyBridgeBranches(ctx, req)
}

// ProfileModifyField corresponds to UprofileServicePort.ProfileModifyField。
func (f *UprofileFacade) ProfileModifyField(ctx context.Context, req *uprofilexsd.ProfileModifyFieldReq) (*uprofilexsd.ProfileModifyFieldRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModifyField(ctx, req)
}

// ProfileModifyHierarchyLevel corresponds to UprofileServicePort.ProfileModifyHierarchyLevel。
func (f *UprofileFacade) ProfileModifyHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileModifyHierarchyLevelReq) (*uprofilexsd.ProfileModifyHierarchyLevelRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModifyHierarchyLevel(ctx, req)
}

// ProfileModifyTags corresponds to UprofileServicePort.ProfileModifyTags。
func (f *UprofileFacade) ProfileModifyTags(ctx context.Context, req *uprofilexsd.ProfileModifyTagsReq) (*uprofilexsd.ProfileModifyTagsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModifyTags(ctx, req)
}

// ProfileModifyTemplate corresponds to UprofileServicePort.ProfileModifyTemplate。
func (f *UprofileFacade) ProfileModifyTemplate(ctx context.Context, req *uprofilexsd.ProfileModifyTemplateReq) (*uprofilexsd.ProfileModifyTemplateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModifyTemplate(ctx, req)
}

// ProfileRetrieve corresponds to UprofileServicePort.ProfileRetrieve。
func (f *UprofileFacade) ProfileRetrieve(ctx context.Context, req *uprofilexsd.ProfileRetrieveReq) (*uprofilexsd.ProfileRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieve(ctx, req)
}

// ProfileRetrieveAction corresponds to UprofileServicePort.ProfileRetrieveAction。
func (f *UprofileFacade) ProfileRetrieveAction(ctx context.Context, req *uprofilexsd.ProfileRetrieveActionReq) (*uprofilexsd.ProfileRetrieveActionRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieveAction(ctx, req)
}

// ProfileRetrieveBridgeBranches corresponds to UprofileServicePort.ProfileRetrieveBridgeBranches。
func (f *UprofileFacade) ProfileRetrieveBridgeBranches(ctx context.Context, req *uprofilexsd.ProfileRetrieveBridgeBranchesReq) (*uprofilexsd.ProfileRetrieveBridgeBranchesRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieveBridgeBranches(ctx, req)
}

// ProfileRetrieveHierarchy corresponds to UprofileServicePort.ProfileRetrieveHierarchy。
func (f *UprofileFacade) ProfileRetrieveHierarchy(ctx context.Context, req *uprofilexsd.ProfileRetrieveHierarchyReq) (*uprofilexsd.ProfileRetrieveHierarchyRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieveHierarchy(ctx, req)
}

// ProfileRetrieveHistory corresponds to UprofileServicePort.ProfileRetrieveHistory。
func (f *UprofileFacade) ProfileRetrieveHistory(ctx context.Context, req *uprofilexsd.ProfileRetrieveHistoryReq) (*uprofilexsd.ProfileRetrieveHistoryRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieveHistory(ctx, req)
}

// ProfileRetrieveTemplate corresponds to UprofileServicePort.ProfileRetrieveTemplate。
func (f *UprofileFacade) ProfileRetrieveTemplate(ctx context.Context, req *uprofilexsd.ProfileRetrieveTemplateReq) (*uprofilexsd.ProfileRetrieveTemplateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieveTemplate(ctx, req)
}

// ProfileSearch corresponds to UprofileServicePort.ProfileSearch。
func (f *UprofileFacade) ProfileSearch(ctx context.Context, req *uprofilexsd.ProfileSearchReq) (*uprofilexsd.ProfileSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileSearch(ctx, req)
}

// ProfileSearchAction corresponds to UprofileServicePort.ProfileSearchAction。
func (f *UprofileFacade) ProfileSearchAction(ctx context.Context, req *uprofilexsd.ProfileSearchActionReq) (*uprofilexsd.ProfileSearchActionRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileSearchAction(ctx, req)
}

// ProfileSearchField corresponds to UprofileServicePort.ProfileSearchField。
func (f *UprofileFacade) ProfileSearchField(ctx context.Context, req *uprofilexsd.ProfileSearchFieldReq) (*uprofilexsd.ProfileSearchFieldRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileSearchField(ctx, req)
}

// ProfileSearchTags corresponds to UprofileServicePort.ProfileSearchTags。
func (f *UprofileFacade) ProfileSearchTags(ctx context.Context, req *uprofilexsd.ProfileSearchTagsReq) (*uprofilexsd.ProfileSearchTagsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileSearchTags(ctx, req)
}

// SingleProfileMigration corresponds to UprofileServicePort.SingleProfileMigration。
func (f *UprofileFacade) SingleProfileMigration(ctx context.Context, req *uprofilexsd.SingleProfileMigrationReq) (*uprofilexsd.SingleProfileMigrationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.SingleProfileMigration(ctx, req)
}
