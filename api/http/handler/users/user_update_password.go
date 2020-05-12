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

type userUpdatePasswordPayload struct {
	Password    string
	NewPassword string
}

func (payload *userUpdatePasswordPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Password) {
		return baasapi.Error("Invalid current password")
	}
	if govalidator.IsNull(payload.NewPassword) {
		return baasapi.Error("Invalid new password")
	}
	return nil
}

// PUT request on /api/users/:id/passwd
func (handler *Handler) userUpdatePassword(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	var payload userUpdatePasswordPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	user, err := handler.UserService.User(baasapi.UserID(userID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a user with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a user with the specified identifier inside the database", err}
	}

	err = handler.CryptoService.CompareHashAndData(user.Password, payload.Password)
	if err != nil {
		return &httperror.HandlerError{http.StatusForbidden, "Specified password do not match actual password", baasapi.ErrUnauthorized}
	}

	user.Password, err = handler.CryptoService.Hash(payload.NewPassword)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to hash user password", baasapi.ErrCryptoHashFailure}
	}

	err = handler.UserService.UpdateUser(user.ID, user)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist user changes inside the database", err}
	}

	return response.Empty(w)
}
