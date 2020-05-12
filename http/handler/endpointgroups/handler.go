package endpointgroups

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// Handler is the HTTP handler used to handle endpoint group operations.
type Handler struct {
	*mux.Router
	EndpointService      baasapi.EndpointService
	EndpointGroupService baasapi.EndpointGroupService
}

// NewHandler creates a handler to manage endpoint group operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/endpoint_groups",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointGroupCreate))).Methods(http.MethodPost)
	h.Handle("/endpoint_groups",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.endpointGroupList))).Methods(http.MethodGet)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointGroupInspect))).Methods(http.MethodGet)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointGroupUpdate))).Methods(http.MethodPut)
	h.Handle("/endpoint_groups/{id}/access",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointGroupUpdateAccess))).Methods(http.MethodPut)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointGroupDelete))).Methods(http.MethodDelete)

	return h
}

func (handler *Handler) checkForGroupUnassignment(endpoint baasapi.Endpoint, associatedEndpoints []baasapi.EndpointID) error {
	for _, id := range associatedEndpoints {
		if id == endpoint.ID {
			return nil
		}
	}

	endpoint.GroupID = baasapi.EndpointGroupID(1)
	return handler.EndpointService.UpdateEndpoint(endpoint.ID, &endpoint)
}

func (handler *Handler) checkForGroupAssignment(endpoint baasapi.Endpoint, groupID baasapi.EndpointGroupID, associatedEndpoints []baasapi.EndpointID) error {
	for _, id := range associatedEndpoints {

		if id == endpoint.ID {
			endpoint.GroupID = groupID
			return handler.EndpointService.UpdateEndpoint(endpoint.ID, &endpoint)
		}
	}
	return nil
}

func (handler *Handler) updateEndpointGroup(endpoint baasapi.Endpoint, groupID baasapi.EndpointGroupID, associatedEndpoints []baasapi.EndpointID) error {
	if endpoint.GroupID == groupID {
		return handler.checkForGroupUnassignment(endpoint, associatedEndpoints)
	} else if endpoint.GroupID == baasapi.EndpointGroupID(1) {
		return handler.checkForGroupAssignment(endpoint, groupID, associatedEndpoints)
	}
	return nil
}
