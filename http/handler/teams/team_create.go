package teams

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type teamCreatePayload struct {
	Name string
}

func (payload *teamCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Name) {
		return baasapi.Error("Invalid team name")
	}
	return nil
}

func (handler *Handler) teamCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload teamCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	team, err := handler.TeamService.TeamByName(payload.Name)
	if err != nil && err != baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve teams from the database", err}
	}
	if team != nil {
		return &httperror.HandlerError{http.StatusConflict, "A team with the same name already exists", baasapi.ErrTeamAlreadyExists}
	}

	team = &baasapi.Team{
		Name: payload.Name,
	}

	err = handler.TeamService.CreateTeam(team)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the team inside the database", err}
	}

	return response.JSON(w, team)
}
