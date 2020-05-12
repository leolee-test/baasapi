package teams

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api/http/security"
)

// GET request on /api/teams
func (handler *Handler) teamList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	teams, err := handler.TeamService.Teams()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve teams from the database", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	filteredTeams := security.FilterUserTeams(teams, securityContext)

	return response.JSON(w, filteredTeams)
}
