// Package api provides the HTTP interface layer for the Uprofile service, plus unified
// error response handling.
package api

import (
	"net/http"

	"github.com/shuiyihan12/uapi-go/pkg/usecase"
)

// UprofileHandler handles the HTTP interface for the Travelport Uprofile service.
type UprofileHandler struct {
	facade *usecase.UprofileFacade
}

// NewUprofileHandler creates the Uprofile handler.
func NewUprofileHandler(facade *usecase.UprofileFacade) *UprofileHandler {
	return &UprofileHandler{facade: facade}
}

// RegisterRoutes registers all endpoints of the Uprofile service.
func (h *UprofileHandler) RegisterRoutes(mux *http.ServeMux) {
	f := h.facade
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-child-search", f.ProfileChildSearch)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-create", f.ProfileCreate)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-create-field", f.ProfileCreateField)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-create-hierarchy-level", f.ProfileCreateHierarchyLevel)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-create-tags", f.ProfileCreateTags)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-delete", f.ProfileDelete)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-delete-hierarchy-level", f.ProfileDeleteHierarchyLevel)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-delete-tag", f.ProfileDeleteTag)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-modify", f.ProfileModify)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-modify-bridge-branches", f.ProfileModifyBridgeBranches)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-modify-field", f.ProfileModifyField)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-modify-hierarchy-level", f.ProfileModifyHierarchyLevel)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-modify-tags", f.ProfileModifyTags)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-modify-template", f.ProfileModifyTemplate)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-retrieve", f.ProfileRetrieve)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-retrieve-action", f.ProfileRetrieveAction)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-retrieve-bridge-branches", f.ProfileRetrieveBridgeBranches)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-retrieve-hierarchy", f.ProfileRetrieveHierarchy)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-retrieve-history", f.ProfileRetrieveHistory)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-retrieve-template", f.ProfileRetrieveTemplate)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-search", f.ProfileSearch)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-search-action", f.ProfileSearchAction)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-search-field", f.ProfileSearchField)
	registerPortHandler(mux, apiBasePath+"/uprofile/profile-search-tags", f.ProfileSearchTags)
	registerPortHandler(mux, apiBasePath+"/uprofile/single-profile-migration", f.SingleProfileMigration)
}
