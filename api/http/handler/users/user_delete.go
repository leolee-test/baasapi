package users

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// DELETE request on /api/users/:id
func (handler *Handler) userDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	userID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid user identifier route variable", err}
	}

	tokenData, err := security.RetrieveTokenData(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve user authentication token", err}
	}

	if tokenData.ID == baasapi.UserID(userID) {
		return &httperror.HandlerError{http.StatusForbidden, "Cannot remove your own user account. Contact another administrator", baasapi.ErrAdminCannotRemoveSelf}
	}

	user, err := handler.UserService.User(baasapi.UserID(userID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a user with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a user with the specified identifier inside the database", err}
	}

	if user.Role == baasapi.AdministratorRole {
		return handler.deleteAdminUser(w, user)
	}

	return handler.deleteUser(w, user)
}

func (handler *Handler) deleteAdminUser(w http.ResponseWriter, user *baasapi.User) *httperror.HandlerError {
	if user.Password == "" {
		return handler.deleteUser(w, user)
	}

	users, err := handler.UserService.Users()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve users from the database", err}
	}

	localAdminCount := 0
	for _, u := range users {
		if u.Role == baasapi.AdministratorRole && u.Password != "" {
			localAdminCount++
		}
	}

	if localAdminCount < 2 {
		return &httperror.HandlerError{http.StatusInternalServerError, "Cannot remove local administrator user", baasapi.ErrCannotRemoveLastLocalAdmin}
	}

	return handler.deleteUser(w, user)
}

func (handler *Handler) deleteUser(w http.ResponseWriter, user *baasapi.User) *httperror.HandlerError {
	err := handler.UserService.DeleteUser(baasapi.UserID(user.ID))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove user from the database", err}
	}

	err = handler.TeamMembershipService.DeleteTeamMembershipByUserID(baasapi.UserID(user.ID))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove user memberships from the database", err}
	}

	return response.Empty(w)
}
