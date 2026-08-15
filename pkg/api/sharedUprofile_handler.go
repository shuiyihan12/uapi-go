// Package api provides the HTTP interface layer for the SharedUprofile service, plus
// unified error response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// SharedUprofileHandler handles the HTTP interface for the Travelport SharedUprofile
// service.
type SharedUprofileHandler struct {
	facade *usecase.SharedUprofileFacade
}

// NewSharedUprofileHandler creates the SharedUprofile handler.
func NewSharedUprofileHandler(facade *usecase.SharedUprofileFacade) *SharedUprofileHandler {
	return &SharedUprofileHandler{facade: facade}
}

// RegisterRoutes registers all endpoints of the SharedUprofile service.
func (h *SharedUprofileHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-child-search", f.ProfileChildSearch)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-create", f.ProfileCreate)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-create-field", f.ProfileCreateField)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-create-tags", f.ProfileCreateTags)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-delete", f.ProfileDelete)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-delete-tag", f.ProfileDeleteTag)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-modify", f.ProfileModify)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-modify-field", f.ProfileModifyField)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-modify-tags", f.ProfileModifyTags)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-retrieve", f.ProfileRetrieve)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-retrieve-history", f.ProfileRetrieveHistory)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-retrieve-parent", f.ProfileRetrieveParent)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-search", f.ProfileSearch)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-search-field", f.ProfileSearchField)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/profile-search-tags", f.ProfileSearchTags)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/single-profile-migration", f.SingleProfileMigration)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/ui-meta-data-create", f.UIMetaDataCreate)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/ui-meta-data-delete", f.UIMetaDataDelete)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/ui-meta-data-modify", f.UIMetaDataModify)
	registerPortHandler(mux, apiBasePath+"/sharedUprofile/ui-meta-data-retrieve", f.UIMetaDataRetrieve)
}
