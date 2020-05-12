package teams

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// GET request on /api/teams/:id/memberships
func (handler *Handler) teamMemberships(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	teamID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid team identifier route variable", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	if !security.AuthorizedTeamManagement(baasapi.TeamID(teamID), securityContext) {
		return &httperror.HandlerError{http.StatusForbidden, "Access denied to team", baasapi.ErrResourceAccessDenied}
	}

	memberships, err := handler.TeamMembershipService.TeamMembershipsByTeamID(baasapi.TeamID(teamID))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve associated team memberships from the database", err}
	}

	return response.JSON(w, memberships)
}
