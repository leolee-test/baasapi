package baask8ss

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type baask8sUpdateAccessPayload struct {
	AuthorizedUsers []int
	AuthorizedTeams []int
}

func (payload *baask8sUpdateAccessPayload) Validate(r *http.Request) error {
	return nil
}

// PUT request on /api/baask8ss/:id/access
func (handler *Handler) baask8sUpdateAccess(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//if !handler.authorizeBaask8sManagement {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Baask8s management is disabled", ErrBaask8sManagementDisabled}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload baask8sUpdateAccessPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

	if payload.AuthorizedUsers != nil {
		authorizedUserIDs := []baasapi.UserID{}
		for _, value := range payload.AuthorizedUsers {
			authorizedUserIDs = append(authorizedUserIDs, baasapi.UserID(value))
		}
		baask8s.AuthorizedUsers = authorizedUserIDs
	}

	if payload.AuthorizedTeams != nil {
		authorizedTeamIDs := []baasapi.TeamID{}
		for _, value := range payload.AuthorizedTeams {
			authorizedTeamIDs = append(authorizedTeamIDs, baasapi.TeamID(value))
		}
		baask8s.AuthorizedTeams = authorizedTeamIDs
	}

	err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
	}

	return response.JSON(w, baask8s)
}
