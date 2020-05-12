package endpoints

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// POST request on /api/endpoints/:id/snapshot
func (handler *Handler) endpointSnapshot(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	endpointID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid endpoint identifier route variable", err}
	}

	endpoint, err := handler.EndpointService.Endpoint(baasapi.EndpointID(endpointID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an endpoint with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an endpoint with the specified identifier inside the database", err}
	}

	if endpoint.Type == baasapi.AzureEnvironment {
		return &httperror.HandlerError{http.StatusBadRequest, "Snapshots not supported for Azure endpoints", err}
	}

	snapshot, snapshotError := handler.Snapshotter.CreateSnapshot(endpoint)

	latestEndpointReference, err := handler.EndpointService.Endpoint(endpoint.ID)
	if latestEndpointReference == nil {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an endpoint with the specified identifier inside the database", err}
	}

	latestEndpointReference.Status = baasapi.EndpointStatusUp
	if snapshotError != nil {
		latestEndpointReference.Status = baasapi.EndpointStatusDown
	}

	if snapshot != nil {
		latestEndpointReference.Snapshots = []baasapi.Snapshot{*snapshot}
	}

	err = handler.EndpointService.UpdateEndpoint(latestEndpointReference.ID, latestEndpointReference)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist endpoint changes inside the database", err}
	}

	return response.Empty(w)
}
