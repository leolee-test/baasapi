package endpointproxy

// TODO: legacy extension management

import (
	"strconv"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/baasapi/api"

	"net/http"
)

func (handler *Handler) proxyRequestsToStoridgeAPI(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	err = handler.requestBouncer.EndpointAccess(r, endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to access endpoint", baasapi.ErrEndpointAccessDenied}
	}

	var storidgeExtension *baasapi.EndpointExtension
	for _, extension := range endpoint.Extensions {
		if extension.Type == baasapi.StoridgeEndpointExtension {
			storidgeExtension = &extension
		}
	}

	if storidgeExtension == nil {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Storidge extension not supported on this endpoint", baasapi.ErrEndpointExtensionNotSupported}
	}

	proxyExtensionKey := string(endpoint.ID) + "_" + string(baasapi.StoridgeEndpointExtension)

	var proxy http.Handler
	proxy = handler.ProxyManager.GetLegacyExtensionProxy(proxyExtensionKey)
	if proxy == nil {
		proxy, err = handler.ProxyManager.CreateLegacyExtensionProxy(proxyExtensionKey, storidgeExtension.URL)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to create extension proxy", err}
		}
	}

	id := strconv.Itoa(endpointID)
	http.StripPrefix("/"+id+"/extensions/storidge", proxy).ServeHTTP(w, r)
	return nil
}
