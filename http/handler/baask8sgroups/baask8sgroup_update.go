package baask8sgroups

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type baask8sGroupUpdatePayload struct {
	Name                string
	Description         string
	AssociatedBaask8ss []baasapi.Baask8sID
	Tags                []string
}

func (payload *baask8sGroupUpdatePayload) Validate(r *http.Request) error {
	return nil
}

// PUT request on /api/baask8s_groups/:id
func (handler *Handler) baask8sGroupUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sGroupID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s group identifier route variable", err}
	}

	var payload baask8sGroupUpdatePayload
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

	if payload.Name != "" {
		baask8sGroup.Name = payload.Name
	}

	if payload.Description != "" {
		baask8sGroup.Description = payload.Description
	}

	if payload.Tags != nil {
		baask8sGroup.Tags = payload.Tags
	}

	err = handler.Baask8sGroupService.UpdateBaask8sGroup(baask8sGroup.ID, baask8sGroup)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s group changes inside the database", err}
	}

	baask8ss, err := handler.Baask8sService.Baask8ss()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	}

	for _, baask8s := range baask8ss {
		err = handler.updateBaask8sGroup(baask8s, baasapi.Baask8sGroupID(baask8sGroupID), payload.AssociatedBaask8ss)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update baask8s", err}
		}
	}

	return response.JSON(w, baask8sGroup)
}
