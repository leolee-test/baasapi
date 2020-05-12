package registries

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// GET request on /api/registries/:id
func (handler *Handler) registryInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	registryID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid registry identifier route variable", err}
	}

	registry, err := handler.RegistryService.Registry(baasapi.RegistryID(registryID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a registry with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a registry with the specified identifier inside the database", err}
	}

	err = handler.requestBouncer.RegistryAccess(r, registry)
	if err != nil {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to access registry", baasapi.ErrBaask8sAccessDenied}
	}

	hideFields(registry)
	return response.JSON(w, registry)
}
