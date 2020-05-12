package baask8sgroups

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// DELETE request on /api/baask8s_groups/:id
func (handler *Handler) baask8sGroupDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sGroupID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s group identifier route variable", err}
	}

	if baask8sGroupID == 1 {
		return &httperror.HandlerError{http.StatusForbidden, "Unable to remove the default 'Unassigned' group", baasapi.ErrCannotRemoveDefaultGroup}
	}

	_, err = handler.Baask8sGroupService.Baask8sGroup(baasapi.Baask8sGroupID(baask8sGroupID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s group with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s group with the specified identifier inside the database", err}
	}

	err = handler.Baask8sGroupService.DeleteBaask8sGroup(baasapi.Baask8sGroupID(baask8sGroupID))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove the baask8s group from the database", err}
	}

	baask8ss, err := handler.Baask8sService.Baask8ss()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	}

	for _, baask8s := range baask8ss {
		if baask8s.GroupID == baasapi.Baask8sGroupID(baask8sGroupID) {
			baask8s.GroupID = baasapi.Baask8sGroupID(1)
			err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update baask8s", err}
			}
		}
	}

	return response.Empty(w)
}
