package baask8sgroups

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type baask8sGroupUpdateAccessPayload struct {
	AuthorizedUsers []int
	AuthorizedTeams []int
}

func (payload *baask8sGroupUpdateAccessPayload) Validate(r *http.Request) error {
	return nil
}

// PUT request on /api/baask8s_groups/:id/access
func (handler *Handler) baask8sGroupUpdateAccess(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sGroupID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s group identifier route variable", err}
	}

	var payload baask8sGroupUpdateAccessPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	baask8sGroup, err := handler.Baask8sGroupService.Baask8sGroup(baasapi.Baask8sGroupID(baask8sGroupID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s group with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s group with the specified identifier inside the database", err}
	}

	if payload.AuthorizedUsers != nil {
		authorizedUserIDs := []baasapi.UserID{}
		for _, value := range payload.AuthorizedUsers {
			authorizedUserIDs = append(authorizedUserIDs, baasapi.UserID(value))
		}
		baask8sGroup.AuthorizedUsers = authorizedUserIDs
	}

	if payload.AuthorizedTeams != nil {
		authorizedTeamIDs := []baasapi.TeamID{}
		for _, value := range payload.AuthorizedTeams {
			authorizedTeamIDs = append(authorizedTeamIDs, baasapi.TeamID(value))
		}
		baask8sGroup.AuthorizedTeams = authorizedTeamIDs
	}

	err = handler.Baask8sGroupService.UpdateBaask8sGroup(baask8sGroup.ID, baask8sGroup)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s group changes inside the database", err}
	}

	return response.JSON(w, baask8sGroup)
}
