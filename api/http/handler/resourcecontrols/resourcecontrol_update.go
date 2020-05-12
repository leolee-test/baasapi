package resourcecontrols

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

type resourceControlUpdatePayload struct {
	Public bool
	Users  []int
	Teams  []int
}

func (payload *resourceControlUpdatePayload) Validate(r *http.Request) error {
	if len(payload.Users) == 0 && len(payload.Teams) == 0 && !payload.Public {
		return baasapi.Error("Invalid resource control declaration. Must specify Users, Teams or Public")
	}
	return nil
}

// PUT request on /api/resource_controls/:id
func (handler *Handler) resourceControlUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	resourceControlID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid resource control identifier route variable", err}
	}

	var payload resourceControlUpdatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	resourceControl, err := handler.ResourceControlService.ResourceControl(baasapi.ResourceControlID(resourceControlID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a resource control with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a resource control with with the specified identifier inside the database", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	if !security.AuthorizedResourceControlAccess(resourceControl, securityContext) {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to update the resource control", baasapi.ErrResourceAccessDenied}
	}

	resourceControl.Public = payload.Public

	var userAccesses = make([]baasapi.UserResourceAccess, 0)
	for _, v := range payload.Users {
		userAccess := baasapi.UserResourceAccess{
			UserID:      baasapi.UserID(v),
			AccessLevel: baasapi.ReadWriteAccessLevel,
		}
		userAccesses = append(userAccesses, userAccess)
	}
	resourceControl.UserAccesses = userAccesses

	var teamAccesses = make([]baasapi.TeamResourceAccess, 0)
	for _, v := range payload.Teams {
		teamAccess := baasapi.TeamResourceAccess{
			TeamID:      baasapi.TeamID(v),
			AccessLevel: baasapi.ReadWriteAccessLevel,
		}
		teamAccesses = append(teamAccesses, teamAccess)
	}
	resourceControl.TeamAccesses = teamAccesses

	if !security.AuthorizedResourceControlUpdate(resourceControl, securityContext) {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to update the resource control", baasapi.ErrResourceAccessDenied}
	}

	err = handler.ResourceControlService.UpdateResourceControl(resourceControl.ID, resourceControl)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist resource control changes inside the database", err}
	}

	return response.JSON(w, resourceControl)
}
