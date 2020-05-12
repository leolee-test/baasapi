package teammemberships

import (
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

// Handler is the HTTP handler used to handle team membership operations.
type Handler struct {
	*mux.Router
	TeamMembershipService  baasapi.TeamMembershipService
	ResourceControlService baasapi.ResourceControlService
}

// NewHandler creates a handler to manage team membership operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/group_memberships",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.teamMembershipCreate))).Methods(http.MethodPost)
	h.Handle("/group_memberships",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.teamMembershipList))).Methods(http.MethodGet)
	h.Handle("/group_memberships/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.teamMembershipUpdate))).Methods(http.MethodPut)
	h.Handle("/group_memberships/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.teamMembershipDelete))).Methods(http.MethodDelete)

	return h
}
