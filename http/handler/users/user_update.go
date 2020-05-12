package users

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

type userUpdatePayload struct {
	Password string
	Role     int
}

func (payload *userUpdatePayload) Validate(r *http.Request) error {
	if payload.Role != 0 && payload.Role != 1 && payload.Role != 2 {
		return baasapi.Error("Invalid role value. Value must be one of: 1 (administrator) or 2 (regular user)")
	}
	return nil
}

// PUT request on /api/users/:id
func (handler *Handler) userUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	userID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid user identifier route variable", err}
	}

	tokenData, err := security.RetrieveTokenData(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve user authentication token", err}
	}

	if tokenData.Role != baasapi.AdministratorRole && tokenData.ID != baasapi.UserID(userID) {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to update user", baasapi.ErrUnauthorized}
	}

	var payload userUpdatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	if tokenData.Role != baasapi.AdministratorRole && payload.Role != 0 {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to update user to administrator role", baasapi.ErrResourceAccessDenied}
	}

	user, err := handler.UserService.User(baasapi.UserID(userID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a user with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a user with the specified identifier inside the database", err}
	}

	if payload.Password != "" {
		user.Password, err = handler.CryptoService.Hash(payload.Password)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to hash user password", baasapi.ErrCryptoHashFailure}
		}
	}

	if payload.Role != 0 {
		user.Role = baasapi.UserRole(payload.Role)
	}

	err = handler.UserService.UpdateUser(user.ID, user)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist user changes inside the database", err}
	}

	return response.JSON(w, user)
}
