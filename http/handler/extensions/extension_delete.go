package extensions

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// DELETE request on /api/extensions/:id
func (handler *Handler) extensionDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	extensionIdentifier, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid extension identifier route variable", err}
	}
	extensionID := baasapi.ExtensionID(extensionIdentifier)

	extension, err := handler.ExtensionService.Extension(extensionID)
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a extension with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a extension with the specified identifier inside the database", err}
	}

	err = handler.ExtensionManager.DisableExtension(extension)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to delete extension", err}
	}

	err = handler.ExtensionService.DeleteExtension(extensionID)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to delete the extension from the database", err}
	}

	return response.Empty(w)
}
