package users

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

type userCreatePayload struct {
	Username string
	Password string
	Role     int
}

func (payload *userCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Username) || govalidator.Contains(payload.Username, " ") {
		return baasapi.Error("Invalid username. Must not contain any whitespace")
	}

	if payload.Role != 1 && payload.Role != 2 {
		return baasapi.Error("Invalid role value. Value must be one of: 1 (administrator) or 2 (regular user)")
	}
	return nil
}

// POST request on /api/users
func (handler *Handler) userCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload userCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	if !securityContext.IsAdmin && !securityContext.IsTeamLeader {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to create user", baasapi.ErrResourceAccessDenied}
	}

	if securityContext.IsTeamLeader && payload.Role == 1 {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to create administrator user", baasapi.ErrResourceAccessDenied}
	}

	user, err := handler.UserService.UserByUsername(payload.Username)
	if err != nil && err != baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve users from the database", err}
	}
	if user != nil {
		return &httperror.HandlerError{http.StatusConflict, "Another user with the same username already exists", baasapi.ErrUserAlreadyExists}
	}

	user = &baasapi.User{
		Username: payload.Username,
		Role:     baasapi.UserRole(payload.Role),
	}

	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve settings from the database", err}
	}

	if settings.AuthenticationMethod == baasapi.AuthenticationInternal {
		user.Password, err = handler.CryptoService.Hash(payload.Password)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to hash user password", baasapi.ErrCryptoHashFailure}
		}
	}

	err = handler.UserService.CreateUser(user)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist user inside the database", err}
	}

	hideFields(user)
	return response.JSON(w, user)
}
