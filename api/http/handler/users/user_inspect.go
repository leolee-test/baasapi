package users

import (
	"net/http"
	"fmt"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// GET request on /api/users/:id
func (handler *Handler) userInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	userID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid user identifier route variable", err}
	}

	user, err := handler.UserService.User(baasapi.UserID(userID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a user with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a user with the specified identifier inside the database", err}
	}

	hideFields(user)
	return response.JSON(w, user)
}

// GET request on /api/users/byname/:id
func (handler *Handler) userByUserName(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	username, err := request.RetrieveRouteVariableValue(r, "username")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid username name variable", err}
	}
	fmt.Println(username)


	tokenData, err := security.RetrieveTokenData(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve user authentication token", err}
	}

	if tokenData.Role != baasapi.AdministratorRole && tokenData.Username != username {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to retrieve user information", baasapi.ErrUnauthorized}
	}

	user, err := handler.UserService.UserByUsername(username)
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a user with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a user with the specified identifier inside the database", err}
	}

	hideFields(user)
	return response.JSON(w, user)
}
