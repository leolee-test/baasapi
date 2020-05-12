package baask8sgroups

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api/http/security"
)

// GET request on /api/baask8s_groups
func (handler *Handler) baask8sGroupList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sGroups, err := handler.Baask8sGroupService.Baask8sGroups()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8s groups from the database", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	baask8sGroups = security.FilterBaask8sGroups(baask8sGroups, securityContext)
	return response.JSON(w, baask8sGroups)
}
