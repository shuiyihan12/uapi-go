// Package sharedUprofile provides the SOAP client implementation for the Travelport SharedUprofile
// service (the shareduprofile-mapped namespace).
// Its SharedUprofileServicePort interface corresponds one-to-one with the WSDL *PortType operations.
// Request/response types come directly from the generator-produced shareduprofile package
// (script generated; do not edit by hand).
// This package keeps no hand-written request models; all infrastructure fields
// are injected uniformly by prepareReq before sending.
package sharedUprofile

import (
	"context"
	"fmt"

	"github.com/shuiyihan12/uapi-go/internal/logging"
	"github.com/shuiyihan12/uapi-go/pkg/client"
	shareduprofilexsd "github.com/shuiyihan12/uapi-go/pkg/generated/shareduprofile"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// SharedUprofileServicePort mirrors the *PortType operations of the shareduprofile-mapped namespace one-to-one,
// making it easy to swap implementations and plug in test stubs.
// Request/response types come directly from the WSDL-generated shareduprofile package.
type SharedUprofileServicePort interface {
	// ProfileChildSearch corresponds to the ProfileChildSearchReq operation of the SharedUprofile service.
	ProfileChildSearch(ctx context.Context, req *shareduprofilexsd.ProfileChildSearchReq) (*shareduprofilexsd.ProfileChildSearchRsp, error)
	// ProfileCreate corresponds to the ProfileCreateReq operation of the SharedUprofile service.
	ProfileCreate(ctx context.Context, req *shareduprofilexsd.ProfileCreateReq) (*shareduprofilexsd.ProfileCreateRsp, error)
	// ProfileCreateField corresponds to the ProfileCreateFieldReq operation of the SharedUprofile service.
	ProfileCreateField(ctx context.Context, req *shareduprofilexsd.ProfileCreateFieldReq) (*shareduprofilexsd.ProfileCreateFieldRsp, error)
	// ProfileCreateTags corresponds to the ProfileCreateTagsReq operation of the SharedUprofile service.
	ProfileCreateTags(ctx context.Context, req *shareduprofilexsd.ProfileCreateTagsReq) (*shareduprofilexsd.ProfileCreateTagsRsp, error)
	// ProfileDelete corresponds to the ProfileDeleteReq operation of the SharedUprofile service.
	ProfileDelete(ctx context.Context, req *shareduprofilexsd.ProfileDeleteReq) (*shareduprofilexsd.ProfileDeleteRsp, error)
	// ProfileDeleteTag corresponds to the ProfileDeleteTagReq operation of the SharedUprofile service.
	ProfileDeleteTag(ctx context.Context, req *shareduprofilexsd.ProfileDeleteTagReq) (*struct{}, error)
	// ProfileModify corresponds to the ProfileModifyReq operation of the SharedUprofile service.
	ProfileModify(ctx context.Context, req *shareduprofilexsd.ProfileModifyReq) (*shareduprofilexsd.ProfileModifyRsp, error)
	// ProfileModifyField corresponds to the ProfileModifyFieldReq operation of the SharedUprofile service.
	ProfileModifyField(ctx context.Context, req *shareduprofilexsd.ProfileModifyFieldReq) (*shareduprofilexsd.ProfileModifyFieldRsp, error)
	// ProfileModifyTags corresponds to the ProfileModifyTagsReq operation of the SharedUprofile service.
	ProfileModifyTags(ctx context.Context, req *shareduprofilexsd.ProfileModifyTagsReq) (*shareduprofilexsd.ProfileModifyTagsRsp, error)
	// ProfileRetrieve corresponds to the ProfileRetrieveReq operation of the SharedUprofile service.
	ProfileRetrieve(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveReq) (*shareduprofilexsd.ProfileRetrieveRsp, error)
	// ProfileRetrieveHistory corresponds to the ProfileRetrieveHistoryReq operation of the SharedUprofile service.
	ProfileRetrieveHistory(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveHistoryReq) (*shareduprofilexsd.ProfileRetrieveHistoryRsp, error)
	// ProfileRetrieveParent corresponds to the ProfileRetrieveParentReq operation of the SharedUprofile service.
	ProfileRetrieveParent(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveParentReq) (*shareduprofilexsd.ProfileRetrieveParentRsp, error)
	// ProfileSearch corresponds to the ProfileSearchReq operation of the SharedUprofile service.
	ProfileSearch(ctx context.Context, req *shareduprofilexsd.ProfileSearchReq) (*shareduprofilexsd.ProfileSearchRsp, error)
	// ProfileSearchField corresponds to the ProfileSearchFieldReq operation of the SharedUprofile service.
	ProfileSearchField(ctx context.Context, req *shareduprofilexsd.ProfileSearchFieldReq) (*shareduprofilexsd.ProfileSearchFieldRsp, error)
	// ProfileSearchTags corresponds to the ProfileSearchTagsReq operation of the SharedUprofile service.
	ProfileSearchTags(ctx context.Context, req *shareduprofilexsd.ProfileSearchTagsReq) (*shareduprofilexsd.ProfileSearchTagsRsp, error)
	// SingleProfileMigration corresponds to the SingleProfileMigrationReq operation of the SharedUprofile service.
	SingleProfileMigration(ctx context.Context, req *shareduprofilexsd.SingleProfileMigrationReq) (*shareduprofilexsd.SingleProfileMigrationRsp, error)
	// UIMetaDataCreate corresponds to the UIMetaDataCreateReq operation of the SharedUprofile service.
	UIMetaDataCreate(ctx context.Context, req *shareduprofilexsd.UIMetaDataCreateReq) (*shareduprofilexsd.UIMetaDataCreateRsp, error)
	// UIMetaDataDelete corresponds to the UIMetaDataDeleteReq operation of the SharedUprofile service.
	UIMetaDataDelete(ctx context.Context, req *shareduprofilexsd.UIMetaDataDeleteReq) (*shareduprofilexsd.UIMetaDataDeleteRsp, error)
	// UIMetaDataModify corresponds to the UIMetaDataModifyReq operation of the SharedUprofile service.
	UIMetaDataModify(ctx context.Context, req *shareduprofilexsd.UIMetaDataModifyReq) (*shareduprofilexsd.UIMetaDataModifyRsp, error)
	// UIMetaDataRetrieve corresponds to the UIMetaDataRetrieveReq operation of the SharedUprofile service.
	UIMetaDataRetrieve(ctx context.Context, req *shareduprofilexsd.UIMetaDataRetrieveReq) (*shareduprofilexsd.UIMetaDataRetrieveRsp, error)
	// Close releases the underlying SOAP client connection resources.
	Close() error
}

// SharedUprofileService is the SOAP implementation of SharedUprofileServicePort.
type SharedUprofileService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// Compile-time assertion: *SharedUprofileService must satisfy the SharedUprofileServicePort interface.
var _ SharedUprofileServicePort = (*SharedUprofileService)(nil)

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

// NewSharedUprofileService builds a SharedUprofile service client from the given SOAP configuration and logger.
func NewSharedUprofileService(config client.SOAPConfig, logger logging.Logger) (*SharedUprofileService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "sharedUprofile-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create sharedUprofile service client: %w", err)
	}

	return &SharedUprofileService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// ProfileChildSearch issues the ProfileChildSearchReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileChildSearch(ctx context.Context, req *shareduprofilexsd.ProfileChildSearchReq) (*shareduprofilexsd.ProfileChildSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileChildSearchRsp](s.client, ctx, "ProfileChildSearch", req)
}

// ProfileCreate issues the ProfileCreateReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileCreate(ctx context.Context, req *shareduprofilexsd.ProfileCreateReq) (*shareduprofilexsd.ProfileCreateRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileCreateRsp](s.client, ctx, "ProfileCreate", req)
}

// ProfileCreateField issues the ProfileCreateFieldReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileCreateField(ctx context.Context, req *shareduprofilexsd.ProfileCreateFieldReq) (*shareduprofilexsd.ProfileCreateFieldRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileCreateFieldRsp](s.client, ctx, "ProfileCreateField", req)
}

// ProfileCreateTags issues the ProfileCreateTagsReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileCreateTags(ctx context.Context, req *shareduprofilexsd.ProfileCreateTagsReq) (*shareduprofilexsd.ProfileCreateTagsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileCreateTagsRsp](s.client, ctx, "ProfileCreateTags", req)
}

// ProfileDelete issues the ProfileDeleteReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileDelete(ctx context.Context, req *shareduprofilexsd.ProfileDeleteReq) (*shareduprofilexsd.ProfileDeleteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileDeleteRsp](s.client, ctx, "ProfileDelete", req)
}

// ProfileDeleteTag issues the ProfileDeleteTagReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileDeleteTag(ctx context.Context, req *shareduprofilexsd.ProfileDeleteTagReq) (*struct{}, error) {
	ctx = prepareReq(ctx, req)
	return callPort[struct{}](s.client, ctx, "ProfileDeleteTag", req)
}

// ProfileModify issues the ProfileModifyReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileModify(ctx context.Context, req *shareduprofilexsd.ProfileModifyReq) (*shareduprofilexsd.ProfileModifyRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileModifyRsp](s.client, ctx, "ProfileModify", req)
}

// ProfileModifyField issues the ProfileModifyFieldReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileModifyField(ctx context.Context, req *shareduprofilexsd.ProfileModifyFieldReq) (*shareduprofilexsd.ProfileModifyFieldRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileModifyFieldRsp](s.client, ctx, "ProfileModifyField", req)
}

// ProfileModifyTags issues the ProfileModifyTagsReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileModifyTags(ctx context.Context, req *shareduprofilexsd.ProfileModifyTagsReq) (*shareduprofilexsd.ProfileModifyTagsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileModifyTagsRsp](s.client, ctx, "ProfileModifyTags", req)
}

// ProfileRetrieve issues the ProfileRetrieveReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileRetrieve(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveReq) (*shareduprofilexsd.ProfileRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileRetrieveRsp](s.client, ctx, "ProfileRetrieve", req)
}

// ProfileRetrieveHistory issues the ProfileRetrieveHistoryReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileRetrieveHistory(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveHistoryReq) (*shareduprofilexsd.ProfileRetrieveHistoryRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileRetrieveHistoryRsp](s.client, ctx, "ProfileRetrieveHistory", req)
}

// ProfileRetrieveParent issues the ProfileRetrieveParentReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileRetrieveParent(ctx context.Context, req *shareduprofilexsd.ProfileRetrieveParentReq) (*shareduprofilexsd.ProfileRetrieveParentRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileRetrieveParentRsp](s.client, ctx, "ProfileRetrieveParent", req)
}

// ProfileSearch issues the ProfileSearchReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileSearch(ctx context.Context, req *shareduprofilexsd.ProfileSearchReq) (*shareduprofilexsd.ProfileSearchRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileSearchRsp](s.client, ctx, "ProfileSearch", req)
}

// ProfileSearchField issues the ProfileSearchFieldReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileSearchField(ctx context.Context, req *shareduprofilexsd.ProfileSearchFieldReq) (*shareduprofilexsd.ProfileSearchFieldRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileSearchFieldRsp](s.client, ctx, "ProfileSearchField", req)
}

// ProfileSearchTags issues the ProfileSearchTagsReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) ProfileSearchTags(ctx context.Context, req *shareduprofilexsd.ProfileSearchTagsReq) (*shareduprofilexsd.ProfileSearchTagsRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.ProfileSearchTagsRsp](s.client, ctx, "ProfileSearchTags", req)
}

// SingleProfileMigration issues the SingleProfileMigrationReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) SingleProfileMigration(ctx context.Context, req *shareduprofilexsd.SingleProfileMigrationReq) (*shareduprofilexsd.SingleProfileMigrationRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.SingleProfileMigrationRsp](s.client, ctx, "SingleProfileMigration", req)
}

// UIMetaDataCreate issues the UIMetaDataCreateReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) UIMetaDataCreate(ctx context.Context, req *shareduprofilexsd.UIMetaDataCreateReq) (*shareduprofilexsd.UIMetaDataCreateRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.UIMetaDataCreateRsp](s.client, ctx, "UIMetaDataCreate", req)
}

// UIMetaDataDelete issues the UIMetaDataDeleteReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) UIMetaDataDelete(ctx context.Context, req *shareduprofilexsd.UIMetaDataDeleteReq) (*shareduprofilexsd.UIMetaDataDeleteRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.UIMetaDataDeleteRsp](s.client, ctx, "UIMetaDataDelete", req)
}

// UIMetaDataModify issues the UIMetaDataModifyReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) UIMetaDataModify(ctx context.Context, req *shareduprofilexsd.UIMetaDataModifyReq) (*shareduprofilexsd.UIMetaDataModifyRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.UIMetaDataModifyRsp](s.client, ctx, "UIMetaDataModify", req)
}

// UIMetaDataRetrieve issues the UIMetaDataRetrieveReq SOAP call and returns the strongly typed response.
func (s *SharedUprofileService) UIMetaDataRetrieve(ctx context.Context, req *shareduprofilexsd.UIMetaDataRetrieveReq) (*shareduprofilexsd.UIMetaDataRetrieveRsp, error) {
	ctx = prepareReq(ctx, req)
	return callPort[shareduprofilexsd.UIMetaDataRetrieveRsp](s.client, ctx, "UIMetaDataRetrieve", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType
// that performs a single SOAP call and decodes it into the strongly typed
// response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *SharedUprofileService) Close() error {
	return s.client.Close()
}
