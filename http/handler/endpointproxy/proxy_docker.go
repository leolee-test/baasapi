package endpointproxy

import (
	"errors"
	"strconv"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/baasapi/api"

	"net/http"
)

func (handler *Handler) proxyRequestsToDockerAPI(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	if endpoint.Status == baasapi.EndpointStatusDown {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Unable to query endpoint", errors.New("Endpoint is down")}
	}

	err = handler.requestBouncer.EndpointAccess(r, endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to access endpoint", baasapi.ErrEndpointAccessDenied}
	}

	var proxy http.Handler
	proxy = handler.ProxyManager.GetProxy(string(endpointID))
	if proxy == nil {
		proxy, err = handler.ProxyManager.CreateAndRegisterProxy(endpoint)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to create proxy", err}
		}
	}

	id := strconv.Itoa(endpointID)
	http.StripPrefix("/"+id+"/docker", proxy).ServeHTTP(w, r)
	return nil
}
