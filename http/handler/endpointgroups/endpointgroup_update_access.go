package endpointgroups

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type endpointGroupUpdateAccessPayload struct {
	AuthorizedUsers []int
	AuthorizedTeams []int
}

func (payload *endpointGroupUpdateAccessPayload) Validate(r *http.Request) error {
	return nil
}

// PUT request on /api/endpoint_groups/:id/access
func (handler *Handler) endpointGroupUpdateAccess(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	endpointGroupID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid endpoint group identifier route variable", err}
	}

	var payload endpointGroupUpdateAccessPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	endpointGroup, err := handler.EndpointGroupService.EndpointGroup(baasapi.EndpointGroupID(endpointGroupID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an endpoint group with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an endpoint group with the specified identifier inside the database", err}
	}

	if payload.AuthorizedUsers != nil {
		authorizedUserIDs := []baasapi.UserID{}
		for _, value := range payload.AuthorizedUsers {
			authorizedUserIDs = append(authorizedUserIDs, baasapi.UserID(value))
		}
		endpointGroup.AuthorizedUsers = authorizedUserIDs
	}

	if payload.AuthorizedTeams != nil {
		authorizedTeamIDs := []baasapi.TeamID{}
		for _, value := range payload.AuthorizedTeams {
			authorizedTeamIDs = append(authorizedTeamIDs, baasapi.TeamID(value))
		}
		endpointGroup.AuthorizedTeams = authorizedTeamIDs
	}

	err = handler.EndpointGroupService.UpdateEndpointGroup(endpointGroup.ID, endpointGroup)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist endpoint group changes inside the database", err}
	}

	return response.JSON(w, endpointGroup)
}
