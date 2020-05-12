package baask8sgroups

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type baask8sGroupCreatePayload struct {
	Name                string
	Description         string
	AssociatedBaask8ss  []baasapi.Baask8sID
	Tags                []string
}

func (payload *baask8sGroupCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Name) {
		return baasapi.Error("Invalid baask8s group name")
	}
	if payload.Tags == nil {
		payload.Tags = []string{}
	}
	return nil
}

// POST request on /api/baask8s_groups
func (handler *Handler) baask8sGroupCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload baask8sGroupCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	baask8sGroup := &baasapi.Baask8sGroup{
		Name:            payload.Name,
		Description:     payload.Description,
		AuthorizedUsers: []baasapi.UserID{},
		AuthorizedTeams: []baasapi.TeamID{},
		Tags:            payload.Tags,
	}

	err = handler.Baask8sGroupService.CreateBaask8sGroup(baask8sGroup)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the baask8s group inside the database", err}
	}

	baask8ss, err := handler.Baask8sService.Baask8ss()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	}

	for _, baask8s := range baask8ss {
		if baask8s.GroupID == baasapi.Baask8sGroupID(1) {
			err = handler.checkForGroupAssignment(baask8s, baask8sGroup.ID, payload.AssociatedBaask8ss)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update baask8s", err}
			}
		}
	}

	return response.JSON(w, baask8sGroup)
}
