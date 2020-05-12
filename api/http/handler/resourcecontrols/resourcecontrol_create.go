package resourcecontrols

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

type resourceControlCreatePayload struct {
	ResourceID     string
	Type           string
	Public         bool
	Users          []int
	Teams          []int
	SubResourceIDs []string
}

func (payload *resourceControlCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.ResourceID) {
		return baasapi.Error("Invalid resource identifier")
	}

	if govalidator.IsNull(payload.Type) {
		return baasapi.Error("Invalid type")
	}

	if len(payload.Users) == 0 && len(payload.Teams) == 0 && !payload.Public {
		return baasapi.Error("Invalid resource control declaration. Must specify Users, Teams or Public")
	}
	return nil
}

// POST request on /api/resource_controls
func (handler *Handler) resourceControlCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload resourceControlCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	var resourceControlType baasapi.ResourceControlType
	switch payload.Type {
	case "container":
		resourceControlType = baasapi.ContainerResourceControl
	case "service":
		resourceControlType = baasapi.ServiceResourceControl
	case "volume":
		resourceControlType = baasapi.VolumeResourceControl
	case "network":
		resourceControlType = baasapi.NetworkResourceControl
	case "secret":
		resourceControlType = baasapi.SecretResourceControl
	case "stack":
		resourceControlType = baasapi.StackResourceControl
	case "config":
		resourceControlType = baasapi.ConfigResourceControl
	default:
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid type value. Value must be one of: container, service, volume, network, secret, stack or config", baasapi.ErrInvalidResourceControlType}
	}

	rc, err := handler.ResourceControlService.ResourceControlByResourceID(payload.ResourceID)
	if err != nil && err != baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve resource controls from the database", err}
	}
	if rc != nil {
		return &httperror.HandlerError{http.StatusConflict, "A resource control is already associated to this resource", baasapi.ErrResourceControlAlreadyExists}
	}

	var userAccesses = make([]baasapi.UserResourceAccess, 0)
	for _, v := range payload.Users {
		userAccess := baasapi.UserResourceAccess{
			UserID:      baasapi.UserID(v),
			AccessLevel: baasapi.ReadWriteAccessLevel,
		}
		userAccesses = append(userAccesses, userAccess)
	}

	var teamAccesses = make([]baasapi.TeamResourceAccess, 0)
	for _, v := range payload.Teams {
		teamAccess := baasapi.TeamResourceAccess{
			TeamID:      baasapi.TeamID(v),
			AccessLevel: baasapi.ReadWriteAccessLevel,
		}
		teamAccesses = append(teamAccesses, teamAccess)
	}

	resourceControl := baasapi.ResourceControl{
		ResourceID:     payload.ResourceID,
		SubResourceIDs: payload.SubResourceIDs,
		Type:           resourceControlType,
		Public:         payload.Public,
		UserAccesses:   userAccesses,
		TeamAccesses:   teamAccesses,
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	if !security.AuthorizedResourceControlCreation(&resourceControl, securityContext) {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to create a resource control for the specified resource", baasapi.ErrResourceAccessDenied}
	}

	err = handler.ResourceControlService.CreateResourceControl(&resourceControl)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the resource control inside the database", err}
	}

	return response.JSON(w, resourceControl)
}
