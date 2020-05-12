package baask8sgroups

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// Handler is the HTTP handler used to handle baask8s group operations.
type Handler struct {
	*mux.Router
	Baask8sService      baasapi.Baask8sService
	Baask8sGroupService baasapi.Baask8sGroupService
}

// NewHandler creates a handler to manage baask8s group operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/baask8s_groups",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sGroupCreate))).Methods(http.MethodPost)
	h.Handle("/baask8s_groups",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sGroupList))).Methods(http.MethodGet)
	h.Handle("/baask8s_groups/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sGroupInspect))).Methods(http.MethodGet)
	h.Handle("/baask8s_groups/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sGroupUpdate))).Methods(http.MethodPut)
	h.Handle("/baask8s_groups/{id}/access",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sGroupUpdateAccess))).Methods(http.MethodPut)
	h.Handle("/baask8s_groups/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sGroupDelete))).Methods(http.MethodDelete)

	return h
}

func (handler *Handler) checkForGroupUnassignment(baask8s baasapi.Baask8s, associatedBaask8ss []baasapi.Baask8sID) error {
	for _, id := range associatedBaask8ss {
		if id == baask8s.ID {
			return nil
		}
	}

	baask8s.GroupID = baasapi.Baask8sGroupID(1)
	return handler.Baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
}

func (handler *Handler) checkForGroupAssignment(baask8s baasapi.Baask8s, groupID baasapi.Baask8sGroupID, associatedBaask8ss []baasapi.Baask8sID) error {
	for _, id := range associatedBaask8ss {

		if id == baask8s.ID {
			baask8s.GroupID = groupID
			return handler.Baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
		}
	}
	return nil
}

func (handler *Handler) updateBaask8sGroup(baask8s baasapi.Baask8s, groupID baasapi.Baask8sGroupID, associatedBaask8ss []baasapi.Baask8sID) error {
	if baask8s.GroupID == groupID {
		return handler.checkForGroupUnassignment(baask8s, associatedBaask8ss)
	} else if baask8s.GroupID == baasapi.Baask8sGroupID(1) {
		return handler.checkForGroupAssignment(baask8s, groupID, associatedBaask8ss)
	}
	return nil
}
