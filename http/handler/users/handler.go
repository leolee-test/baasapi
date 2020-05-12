package users

import (
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

func hideFields(user *baasapi.User) {
	user.Password = ""
}

// Handler is the HTTP handler used to handle user operations.
type Handler struct {
	*mux.Router
	UserService            baasapi.UserService
	TeamService            baasapi.TeamService
	TeamMembershipService  baasapi.TeamMembershipService
	ResourceControlService baasapi.ResourceControlService
	CryptoService          baasapi.CryptoService
	SettingsService        baasapi.SettingsService
}

// NewHandler creates a handler to manage user operations.
func NewHandler(bouncer *security.RequestBouncer, rateLimiter *security.RateLimiter) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/users",
		//bouncer.PublicAccess(httperror.LoggerHandler(h.userCreate))).Methods(http.MethodPost)
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.userCreate))).Methods(http.MethodPost)
	h.Handle("/users",
		//bouncer.PublicAccess(httperror.LoggerHandler(h.userList))).Methods(http.MethodGet)
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.userList))).Methods(http.MethodGet)
	h.Handle("/users/{id}",
		//bouncer.PublicAccess(httperror.LoggerHandler(h.userInspect))).Methods(http.MethodGet)
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.userInspect))).Methods(http.MethodGet)
	h.Handle("/users/byname/{username}",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.userByUserName))).Methods(http.MethodGet)
	h.Handle("/users/{id}",
		//bouncer.PublicAccess(httperror.LoggerHandler(h.userUpdate))).Methods(http.MethodPut)
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.userUpdate))).Methods(http.MethodPut)
	h.Handle("/users/{id}",
		//bouncer.PublicAccess(httperror.LoggerHandler(h.userDelete))).Methods(http.MethodDelete)
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.userDelete))).Methods(http.MethodDelete)
	h.Handle("/users/{id}/memberships",
		//bouncer.PublicAccess(httperror.LoggerHandler(h.userMemberships))).Methods(http.MethodGet)
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.userMemberships))).Methods(http.MethodGet)
	h.Handle("/users/{id}/passwd",
		//bouncer.PublicAccess(bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.userUpdatePassword)))).Methods(http.MethodPut)
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.userUpdatePassword))).Methods(http.MethodPut)
	h.Handle("/users/admin/check",
		bouncer.PublicAccess(httperror.LoggerHandler(h.adminCheck))).Methods(http.MethodGet)
	h.Handle("/users/admin/init",
		bouncer.PublicAccess(httperror.LoggerHandler(h.adminInit))).Methods(http.MethodPost)

	return h
}
