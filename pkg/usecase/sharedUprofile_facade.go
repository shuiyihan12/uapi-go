// Package usecase provides the business orchestration layer for each Travelport service.
// This file wraps the SharedUprofile service (shareduprofile): methods map 1:1 to the SharedUprofileServicePort,
// each simply retrieves the service client and delegates the call without additional orchestration.
package usecase

import (
	"context"
	"errors"
	"fmt"

	shareduprofilexsd "github.com/shuiyihan12/uapi-go/pkg/generated/shareduprofile"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	sharedUprofilesvc "github.com/shuiyihan12/uapi-go/pkg/services/sharedUprofile"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// SharedUprofileFacade orchestrates the Travelport SharedUprofile service use cases; methods map 1:1 to the service PortType.
type SharedUprofileFacade struct {
	manager *manager.ServiceManager
}

// NewSharedUprofileFacade creates the SharedUprofile use-case layer.
func NewSharedUprofileFacade(serviceManager *manager.ServiceManager) *SharedUprofileFacade {
	return &SharedUprofileFacade{manager: serviceManager}
}

// getService lazily retrieves the SharedUprofile service client, handling nil manager and lookup failures uniformly.
func (f *SharedUprofileFacade) getService() (*sharedUprofilesvc.SharedUprofileService, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	svc, err := manager.Get[*sharedUprofilesvc.SharedUprofileService](f.manager, "sharedUprofile")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharedUprofile service: %w", err)
	}
	return svc, nil
}

// ProfileChildSearch corresponds to SharedUprofileServicePort.ProfileChildSearch。
func (f *SharedUprofileFacade) ProfileChildSearch(ctx context.Context, req *shareduprofilexsd.ProfileChildSearchReq) (*shareduprofilexsd.ProfileChildSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileChildSearch(ctx, req)
}

// ProfileCreate corresponds to SharedUprofileServicePort.ProfileCreate。
func (f *SharedUprofileFacade) ProfileCreate(ctx context.Context, req *shareduprofilexsd.ProfileCreateReq) (*shareduprofilexsd.ProfileCreateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileCreate(ctx, req)
}

// ProfileCreateField corresponds to SharedUprofileServicePort.ProfileCreateField。
func (f *SharedUprofileFacade) ProfileCreateField(ctx context.Context, req *shareduprofilexsd.ProfileCreateFieldReq) (*shareduprofilexsd.ProfileCreateFieldRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileCreateField(ctx, req)
}

// ProfileCreateTags corresponds to SharedUprofileServicePort.ProfileCreateTags。
func (f *SharedUprofileFacade) ProfileCreateTags(ctx context.Context, req *shareduprofilexsd.ProfileCreateTagsReq) (*shareduprofilexsd.ProfileCreateTagsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileCreateTags(ctx, req)
}

// ProfileDelete corresponds to SharedUprofileServicePort.ProfileDelete。
func (f *SharedUprofileFacade) ProfileDelete(ctx context.Context, req *shareduprofilexsd.ProfileDeleteReq) (*shareduprofilexsd.ProfileDeleteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileDelete(ctx, req)
}

// ProfileDeleteTag corresponds to SharedUprofileServicePort.ProfileDeleteTag。
func (f *SharedUprofileFacade) ProfileDeleteTag(ctx context.Context, req *shareduprofilexsd.ProfileDeleteTagReq) (*struct{}, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileDeleteTag(ctx, req)
}

// ProfileModify corresponds to SharedUprofileServicePort.ProfileModify。
func (f *SharedUprofileFacade) ProfileModify(ctx context.Context, req *shareduprofilexsd.ProfileModifyReq) (*shareduprofilexsd.ProfileModifyRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModify(ctx, req)
}

// ProfileModifyField corresponds to SharedUprofileServicePort.ProfileModifyField。
func (f *SharedUprofileFacade) ProfileModifyField(ctx context.Context, req *shareduprofilexsd.ProfileModifyFieldReq) (*shareduprofilexsd.ProfileModifyFieldRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModifyField(ctx, req)
}

// ProfileModifyTags corresponds to SharedUprofileServicePort.ProfileModifyTags。
func (f *SharedUprofileFacade) ProfileModifyTags(ctx context.Context, req *shareduprofilexsd.ProfileModifyTagsReq) (*shareduprofilexsd.ProfileModifyTagsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileModifyTags(ctx, req)
}

// ProfileRetrieve corresponds to SharedUprofileServicePort.ProfileRetrieve。
func (f *SharedUprofileFacade) ProfileRetrieve(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveReq) (*shareduprofilexsd.ProfileRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieve(ctx, req)
}

// ProfileRetrieveHistory corresponds to SharedUprofileServicePort.ProfileRetrieveHistory。
func (f *SharedUprofileFacade) ProfileRetrieveHistory(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveHistoryReq) (*shareduprofilexsd.ProfileRetrieveHistoryRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieveHistory(ctx, req)
}

// ProfileRetrieveParent corresponds to SharedUprofileServicePort.ProfileRetrieveParent。
func (f *SharedUprofileFacade) ProfileRetrieveParent(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveParentReq) (*shareduprofilexsd.ProfileRetrieveParentRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileRetrieveParent(ctx, req)
}

// ProfileSearch corresponds to SharedUprofileServicePort.ProfileSearch。
func (f *SharedUprofileFacade) ProfileSearch(ctx context.Context, req *shareduprofilexsd.ProfileSearchReq) (*shareduprofilexsd.ProfileSearchRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileSearch(ctx, req)
}

// ProfileSearchField corresponds to SharedUprofileServicePort.ProfileSearchField。
func (f *SharedUprofileFacade) ProfileSearchField(ctx context.Context, req *shareduprofilexsd.ProfileSearchFieldReq) (*shareduprofilexsd.ProfileSearchFieldRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileSearchField(ctx, req)
}

// ProfileSearchTags corresponds to SharedUprofileServicePort.ProfileSearchTags。
func (f *SharedUprofileFacade) ProfileSearchTags(ctx context.Context, req *shareduprofilexsd.ProfileSearchTagsReq) (*shareduprofilexsd.ProfileSearchTagsRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.ProfileSearchTags(ctx, req)
}

// SingleProfileMigration corresponds to SharedUprofileServicePort.SingleProfileMigration。
func (f *SharedUprofileFacade) SingleProfileMigration(ctx context.Context, req *shareduprofilexsd.SingleProfileMigrationReq) (*shareduprofilexsd.SingleProfileMigrationRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.SingleProfileMigration(ctx, req)
}

// UIMetaDataCreate corresponds to SharedUprofileServicePort.UIMetaDataCreate。
func (f *SharedUprofileFacade) UIMetaDataCreate(ctx context.Context, req *shareduprofilexsd.UIMetaDataCreateReq) (*shareduprofilexsd.UIMetaDataCreateRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UIMetaDataCreate(ctx, req)
}

// UIMetaDataDelete corresponds to SharedUprofileServicePort.UIMetaDataDelete。
func (f *SharedUprofileFacade) UIMetaDataDelete(ctx context.Context, req *shareduprofilexsd.UIMetaDataDeleteReq) (*shareduprofilexsd.UIMetaDataDeleteRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UIMetaDataDelete(ctx, req)
}

// UIMetaDataModify corresponds to SharedUprofileServicePort.UIMetaDataModify。
func (f *SharedUprofileFacade) UIMetaDataModify(ctx context.Context, req *shareduprofilexsd.UIMetaDataModifyReq) (*shareduprofilexsd.UIMetaDataModifyRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UIMetaDataModify(ctx, req)
}

// UIMetaDataRetrieve corresponds to SharedUprofileServicePort.UIMetaDataRetrieve。
func (f *SharedUprofileFacade) UIMetaDataRetrieve(ctx context.Context, req *shareduprofilexsd.UIMetaDataRetrieveReq) (*shareduprofilexsd.UIMetaDataRetrieveRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.UIMetaDataRetrieve(ctx, req)
}
