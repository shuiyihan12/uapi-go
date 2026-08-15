// Package uprofile provides the SOAP client implementation for the Travelport Uprofile
// service (the uprofile-mapped namespace).
// Its UprofileServicePort interface corresponds one-to-one with the WSDL *PortType operations.
// Request/response types come directly from the generator-produced uprofile package
// (script generated; do not edit by hand).
// This package keeps no hand-written request models; all infrastructure fields
// are injected uniformly by prepareReq before sending.
package uprofile

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/internal/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	uprofilexsd "github.com/shuiyihan12/uapi-go/pkg/generated/uprofile"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// UprofileServicePort mirrors the *PortType operations of the uprofile-mapped namespace one-to-one,
// making it easy to swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated uprofile package.
type UprofileServicePort interface {
	// ProfileChildSearch corresponds to the ProfileChildSearchReq operation of the Uprofile service.
	ProfileChildSearch(ctx context.Context, req *uprofilexsd.ProfileChildSearchReq) (*uprofilexsd.ProfileChildSearchRsp, error)
	// ProfileCreate corresponds to the ProfileCreateReq operation of the Uprofile service.
	ProfileCreate(ctx context.Context, req *uprofilexsd.ProfileCreateReq) (*uprofilexsd.ProfileCreateRsp, error)
	// ProfileCreateField corresponds to the ProfileCreateFieldReq operation of the Uprofile service.
	ProfileCreateField(ctx context.Context, req *uprofilexsd.ProfileCreateFieldReq) (*uprofilexsd.ProfileCreateFieldRsp, error)
	// ProfileCreateHierarchyLevel corresponds to the ProfileCreateHierarchyLevelReq operation of the Uprofile service.
	ProfileCreateHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileCreateHierarchyLevelReq) (*uprofilexsd.ProfileCreateHierarchyLevelRsp, error)
	// ProfileCreateTags corresponds to the ProfileCreateTagsReq operation of the Uprofile service.
	ProfileCreateTags(ctx context.Context, req *uprofilexsd.ProfileCreateTagsReq) (*uprofilexsd.ProfileCreateTagsRsp, error)
	// ProfileDelete corresponds to the ProfileDeleteReq operation of the Uprofile service.
	ProfileDelete(ctx context.Context, req *uprofilexsd.ProfileDeleteReq) (*uprofilexsd.ProfileDeleteRsp, error)
	// ProfileDeleteHierarchyLevel corresponds to the ProfileDeleteHierarchyLevelReq operation of the Uprofile service.
	ProfileDeleteHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileDeleteHierarchyLevelReq) (*uprofilexsd.ProfileDeleteHierarchyLevelRsp, error)
	// ProfileDeleteTag corresponds to the ProfileDeleteTagReq operation of the Uprofile service.
	ProfileDeleteTag(ctx context.Context, req *uprofilexsd.ProfileDeleteTagReq) (*struct{}, error)
	// ProfileModify corresponds to the ProfileModifyReq operation of the Uprofile service.
	ProfileModify(ctx context.Context, req *uprofilexsd.ProfileModifyReq) (*uprofilexsd.ProfileModifyRsp, error)
	// ProfileModifyBridgeBranches corresponds to the ProfileModifyBridgeBranchesReq operation of the Uprofile service.
	ProfileModifyBridgeBranches(ctx context.Context, req *uprofilexsd.ProfileModifyBridgeBranchesReq) (*uprofilexsd.ProfileModifyBridgeBranchesRsp, error)
	// ProfileModifyField corresponds to the ProfileModifyFieldReq operation of the Uprofile service.
	ProfileModifyField(ctx context.Context, req *uprofilexsd.ProfileModifyFieldReq) (*uprofilexsd.ProfileModifyFieldRsp, error)
	// ProfileModifyHierarchyLevel corresponds to the ProfileModifyHierarchyLevelReq operation of the Uprofile service.
	ProfileModifyHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileModifyHierarchyLevelReq) (*uprofilexsd.ProfileModifyHierarchyLevelRsp, error)
	// ProfileModifyTags corresponds to the ProfileModifyTagsReq operation of the Uprofile service.
	ProfileModifyTags(ctx context.Context, req *uprofilexsd.ProfileModifyTagsReq) (*uprofilexsd.ProfileModifyTagsRsp, error)
	// ProfileModifyTemplate corresponds to the ProfileModifyTemplateReq operation of the Uprofile service.
	ProfileModifyTemplate(ctx context.Context, req *uprofilexsd.ProfileModifyTemplateReq) (*uprofilexsd.ProfileModifyTemplateRsp, error)
	// ProfileRetrieve corresponds to the ProfileRetrieveReq operation of the Uprofile service.
	ProfileRetrieve(ctx context.Context, req *uprofilexsd.ProfileRetrieveReq) (*uprofilexsd.ProfileRetrieveRsp, error)
	// ProfileRetrieveAction corresponds to the ProfileRetrieveActionReq operation of the Uprofile service.
	ProfileRetrieveAction(ctx context.Context, req *uprofilexsd.ProfileRetrieveActionReq) (*uprofilexsd.ProfileRetrieveActionRsp, error)
	// ProfileRetrieveBridgeBranches corresponds to the ProfileRetrieveBridgeBranchesReq operation of the Uprofile service.
	ProfileRetrieveBridgeBranches(ctx context.Context, req *uprofilexsd.ProfileRetrieveBridgeBranchesReq) (*uprofilexsd.ProfileRetrieveBridgeBranchesRsp, error)
	// ProfileRetrieveHierarchy corresponds to the ProfileRetrieveHierarchyReq operation of the Uprofile service.
	ProfileRetrieveHierarchy(ctx context.Context, req *uprofilexsd.ProfileRetrieveHierarchyReq) (*uprofilexsd.ProfileRetrieveHierarchyRsp, error)
	// ProfileRetrieveHistory corresponds to the ProfileRetrieveHistoryReq operation of the Uprofile service.
	ProfileRetrieveHistory(ctx context.Context, req *uprofilexsd.ProfileRetrieveHistoryReq) (*uprofilexsd.ProfileRetrieveHistoryRsp, error)
	// ProfileRetrieveTemplate corresponds to the ProfileRetrieveTemplateReq operation of the Uprofile service.
	ProfileRetrieveTemplate(ctx context.Context, req *uprofilexsd.ProfileRetrieveTemplateReq) (*uprofilexsd.ProfileRetrieveTemplateRsp, error)
	// ProfileSearch corresponds to the ProfileSearchReq operation of the Uprofile service.
	ProfileSearch(ctx context.Context, req *uprofilexsd.ProfileSearchReq) (*uprofilexsd.ProfileSearchRsp, error)
	// ProfileSearchAction corresponds to the ProfileSearchActionReq operation of the Uprofile service.
	ProfileSearchAction(ctx context.Context, req *uprofilexsd.ProfileSearchActionReq) (*uprofilexsd.ProfileSearchActionRsp, error)
	// ProfileSearchField corresponds to the ProfileSearchFieldReq operation of the Uprofile service.
	ProfileSearchField(ctx context.Context, req *uprofilexsd.ProfileSearchFieldReq) (*uprofilexsd.ProfileSearchFieldRsp, error)
	// ProfileSearchTags corresponds to the ProfileSearchTagsReq operation of the Uprofile service.
	ProfileSearchTags(ctx context.Context, req *uprofilexsd.ProfileSearchTagsReq) (*uprofilexsd.ProfileSearchTagsRsp, error)
	// SingleProfileMigration corresponds to the SingleProfileMigrationReq operation of the Uprofile service.
	SingleProfileMigration(ctx context.Context, req *uprofilexsd.SingleProfileMigrationReq) (*uprofilexsd.SingleProfileMigrationRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// UprofileService is the SOAP implementation of UprofileServicePort.
type UprofileService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *UprofileService must satisfy the UprofileServicePort interface.
var _ UprofileServicePort = (*UprofileService)(nil)

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

// NewUprofileService builds a Uprofile service client from the given SOAP configuration and logger.
func NewUprofileService(config client.SOAPConfig, logger logging.Logger) (*UprofileService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "uprofile-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create uprofile service client: %w", err)
	}

	return &UprofileService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// ProfileChildSearch issues the ProfileChildSearchReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileChildSearch(ctx context.Context, req *uprofilexsd.ProfileChildSearchReq) (*uprofilexsd.ProfileChildSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileChildSearchRsp](s.client, ctx, "ProfileChildSearch", req)
}

// ProfileCreate issues the ProfileCreateReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileCreate(ctx context.Context, req *uprofilexsd.ProfileCreateReq) (*uprofilexsd.ProfileCreateRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileCreateRsp](s.client, ctx, "ProfileCreate", req)
}

// ProfileCreateField issues the ProfileCreateFieldReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileCreateField(ctx context.Context, req *uprofilexsd.ProfileCreateFieldReq) (*uprofilexsd.ProfileCreateFieldRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileCreateFieldRsp](s.client, ctx, "ProfileCreateField", req)
}

// ProfileCreateHierarchyLevel issues the ProfileCreateHierarchyLevelReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileCreateHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileCreateHierarchyLevelReq) (*uprofilexsd.ProfileCreateHierarchyLevelRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileCreateHierarchyLevelRsp](s.client, ctx, "ProfileCreateHierarchyLevel", req)
}

// ProfileCreateTags issues the ProfileCreateTagsReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileCreateTags(ctx context.Context, req *uprofilexsd.ProfileCreateTagsReq) (*uprofilexsd.ProfileCreateTagsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileCreateTagsRsp](s.client, ctx, "ProfileCreateTags", req)
}

// ProfileDelete issues the ProfileDeleteReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileDelete(ctx context.Context, req *uprofilexsd.ProfileDeleteReq) (*uprofilexsd.ProfileDeleteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileDeleteRsp](s.client, ctx, "ProfileDelete", req)
}

// ProfileDeleteHierarchyLevel issues the ProfileDeleteHierarchyLevelReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileDeleteHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileDeleteHierarchyLevelReq) (*uprofilexsd.ProfileDeleteHierarchyLevelRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileDeleteHierarchyLevelRsp](s.client, ctx, "ProfileDeleteHierarchyLevel", req)
}

// ProfileDeleteTag issues the ProfileDeleteTagReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileDeleteTag(ctx context.Context, req *uprofilexsd.ProfileDeleteTagReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "ProfileDeleteTag", req)
}

// ProfileModify issues the ProfileModifyReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileModify(ctx context.Context, req *uprofilexsd.ProfileModifyReq) (*uprofilexsd.ProfileModifyRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileModifyRsp](s.client, ctx, "ProfileModify", req)
}

// ProfileModifyBridgeBranches issues the ProfileModifyBridgeBranchesReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileModifyBridgeBranches(ctx context.Context, req *uprofilexsd.ProfileModifyBridgeBranchesReq) (*uprofilexsd.ProfileModifyBridgeBranchesRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileModifyBridgeBranchesRsp](s.client, ctx, "ProfileModifyBridgeBranches", req)
}

// ProfileModifyField issues the ProfileModifyFieldReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileModifyField(ctx context.Context, req *uprofilexsd.ProfileModifyFieldReq) (*uprofilexsd.ProfileModifyFieldRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileModifyFieldRsp](s.client, ctx, "ProfileModifyField", req)
}

// ProfileModifyHierarchyLevel issues the ProfileModifyHierarchyLevelReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileModifyHierarchyLevel(ctx context.Context, req *uprofilexsd.ProfileModifyHierarchyLevelReq) (*uprofilexsd.ProfileModifyHierarchyLevelRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileModifyHierarchyLevelRsp](s.client, ctx, "ProfileModifyHierarchyLevel", req)
}

// ProfileModifyTags issues the ProfileModifyTagsReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileModifyTags(ctx context.Context, req *uprofilexsd.ProfileModifyTagsReq) (*uprofilexsd.ProfileModifyTagsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileModifyTagsRsp](s.client, ctx, "ProfileModifyTags", req)
}

// ProfileModifyTemplate issues the ProfileModifyTemplateReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileModifyTemplate(ctx context.Context, req *uprofilexsd.ProfileModifyTemplateReq) (*uprofilexsd.ProfileModifyTemplateRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileModifyTemplateRsp](s.client, ctx, "ProfileModifyTemplate", req)
}

// ProfileRetrieve issues the ProfileRetrieveReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileRetrieve(ctx context.Context, req *uprofilexsd.ProfileRetrieveReq) (*uprofilexsd.ProfileRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileRetrieveRsp](s.client, ctx, "ProfileRetrieve", req)
}

// ProfileRetrieveAction issues the ProfileRetrieveActionReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileRetrieveAction(ctx context.Context, req *uprofilexsd.ProfileRetrieveActionReq) (*uprofilexsd.ProfileRetrieveActionRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileRetrieveActionRsp](s.client, ctx, "ProfileRetrieveAction", req)
}

// ProfileRetrieveBridgeBranches issues the ProfileRetrieveBridgeBranchesReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileRetrieveBridgeBranches(ctx context.Context, req *uprofilexsd.ProfileRetrieveBridgeBranchesReq) (*uprofilexsd.ProfileRetrieveBridgeBranchesRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileRetrieveBridgeBranchesRsp](s.client, ctx, "ProfileRetrieveBridgeBranches", req)
}

// ProfileRetrieveHierarchy issues the ProfileRetrieveHierarchyReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileRetrieveHierarchy(ctx context.Context, req *uprofilexsd.ProfileRetrieveHierarchyReq) (*uprofilexsd.ProfileRetrieveHierarchyRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileRetrieveHierarchyRsp](s.client, ctx, "ProfileRetrieveHierarchy", req)
}

// ProfileRetrieveHistory issues the ProfileRetrieveHistoryReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileRetrieveHistory(ctx context.Context, req *uprofilexsd.ProfileRetrieveHistoryReq) (*uprofilexsd.ProfileRetrieveHistoryRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileRetrieveHistoryRsp](s.client, ctx, "ProfileRetrieveHistory", req)
}

// ProfileRetrieveTemplate issues the ProfileRetrieveTemplateReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileRetrieveTemplate(ctx context.Context, req *uprofilexsd.ProfileRetrieveTemplateReq) (*uprofilexsd.ProfileRetrieveTemplateRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileRetrieveTemplateRsp](s.client, ctx, "ProfileRetrieveTemplate", req)
}

// ProfileSearch issues the ProfileSearchReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileSearch(ctx context.Context, req *uprofilexsd.ProfileSearchReq) (*uprofilexsd.ProfileSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileSearchRsp](s.client, ctx, "ProfileSearch", req)
}

// ProfileSearchAction issues the ProfileSearchActionReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileSearchAction(ctx context.Context, req *uprofilexsd.ProfileSearchActionReq) (*uprofilexsd.ProfileSearchActionRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileSearchActionRsp](s.client, ctx, "ProfileSearchAction", req)
}

// ProfileSearchField issues the ProfileSearchFieldReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileSearchField(ctx context.Context, req *uprofilexsd.ProfileSearchFieldReq) (*uprofilexsd.ProfileSearchFieldRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileSearchFieldRsp](s.client, ctx, "ProfileSearchField", req)
}

// ProfileSearchTags issues the ProfileSearchTagsReq SOAP call and returns the strongly typed response.
func (s *UprofileService) ProfileSearchTags(ctx context.Context, req *uprofilexsd.ProfileSearchTagsReq) (*uprofilexsd.ProfileSearchTagsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.ProfileSearchTagsRsp](s.client, ctx, "ProfileSearchTags", req)
}

// SingleProfileMigration issues the SingleProfileMigrationReq SOAP call and returns the strongly typed response.
func (s *UprofileService) SingleProfileMigration(ctx context.Context, req *uprofilexsd.SingleProfileMigrationReq) (*uprofilexsd.SingleProfileMigrationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[uprofilexsd.SingleProfileMigrationRsp](s.client, ctx, "SingleProfileMigration", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType
// that performs a single SOAP call and decodes it into the strongly typed
// response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *UprofileService) Close() error {
	return s.client.Close()
}
