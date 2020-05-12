package teams

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type teamUpdatePayload struct {
	Name string
}

func (payload *teamUpdatePayload) Validate(r *http.Request) error {
	return nil
}

// PUT request on /api/teams/:id
func (handler *Handler) teamUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	teamID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid team identifier route variable", err}
	}

	var payload teamUpdatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	team, err := handler.TeamService.Team(baasapi.TeamID(teamID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a team with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a team with the specified identifier inside the database", err}
	}

	if payload.Name != "" {
		team.Name = payload.Name
	}

	err = handler.TeamService.UpdateTeam(team.ID, team)
	if err != nil {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to persist team changes inside the database", err}
	}

	return response.JSON(w, team)
}
