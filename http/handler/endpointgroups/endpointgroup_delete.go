package endpointgroups

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// DELETE request on /api/endpoint_groups/:id
func (handler *Handler) endpointGroupDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	endpointGroupID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid endpoint group identifier route variable", err}
	}

	if endpointGroupID == 1 {
		return &httperror.HandlerError{http.StatusForbidden, "Unable to remove the default 'Unassigned' group", baasapi.ErrCannotRemoveDefaultGroup}
	}

	_, err = handler.EndpointGroupService.EndpointGroup(baasapi.EndpointGroupID(endpointGroupID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an endpoint group with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an endpoint group with the specified identifier inside the database", err}
	}

	err = handler.EndpointGroupService.DeleteEndpointGroup(baasapi.EndpointGroupID(endpointGroupID))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove the endpoint group from the database", err}
	}

	endpoints, err := handler.EndpointService.Endpoints()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve endpoints from the database", err}
	}

	for _, endpoint := range endpoints {
		if endpoint.GroupID == baasapi.EndpointGroupID(endpointGroupID) {
			endpoint.GroupID = baasapi.EndpointGroupID(1)
			err = handler.EndpointService.UpdateEndpoint(endpoint.ID, &endpoint)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update endpoint", err}
			}
		}
	}

	return response.Empty(w)
}
