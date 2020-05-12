package teams

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// Handler is the HTTP handler used to handle team operations.
type Handler struct {
	*mux.Router
	TeamService            baasapi.TeamService
	TeamMembershipService  baasapi.TeamMembershipService
	ResourceControlService baasapi.ResourceControlService
}

// NewHandler creates a handler to manage team operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/groups",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.teamCreate))).Methods(http.MethodPost)
	h.Handle("/groups",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.teamList))).Methods(http.MethodGet)
	h.Handle("/groups/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.teamInspect))).Methods(http.MethodGet)
	h.Handle("/groups/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.teamUpdate))).Methods(http.MethodPut)
	h.Handle("/groups/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.teamDelete))).Methods(http.MethodDelete)
	h.Handle("/groups/{id}/memberships",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.teamMemberships))).Methods(http.MethodGet)

	return h
}
