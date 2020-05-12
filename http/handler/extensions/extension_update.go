package extensions

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type extensionUpdatePayload struct {
	Version string
}

func (payload *extensionUpdatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Version) {
		return baasapi.Error("Invalid extension version")
	}

	return nil
}

func (handler *Handler) extensionUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	extensionIdentifier, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid extension identifier route variable", err}
	}
	extensionID := baasapi.ExtensionID(extensionIdentifier)

	var payload extensionUpdatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	extension, err := handler.ExtensionService.Extension(extensionID)
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a extension with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a extension with the specified identifier inside the database", err}
	}

	err = handler.ExtensionManager.UpdateExtension(extension, payload.Version)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update extension", err}
	}

	err = handler.ExtensionService.Persist(extension)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist extension status inside the database", err}
	}

	return response.Empty(w)
}
