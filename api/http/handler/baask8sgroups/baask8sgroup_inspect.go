package baask8sgroups

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// GET request on /api/baask8s_groups/:id
func (handler *Handler) baask8sGroupInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sGroupID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s group identifier route variable", err}
	}

	baask8sGroup, err := handler.Baask8sGroupService.Baask8sGroup(baasapi.Baask8sGroupID(baask8sGroupID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s group with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s group with the specified identifier inside the database", err}
	}

	return response.JSON(w, baask8sGroup)
}
