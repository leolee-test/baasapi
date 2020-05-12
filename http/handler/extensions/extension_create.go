package extensions

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type extensionCreatePayload struct {
	License string
}

func (payload *extensionCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.License) {
		return baasapi.Error("Invalid license")
	}

	return nil
}

func (handler *Handler) extensionCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload extensionCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	extensionIdentifier, err := strconv.Atoi(string(payload.License[0]))
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid license format", err}
	}
	extensionID := baasapi.ExtensionID(extensionIdentifier)

	extensions, err := handler.ExtensionService.Extensions()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve extensions status from the database", err}
	}

	for _, existingExtension := range extensions {
		if existingExtension.ID == extensionID && existingExtension.Enabled {
			return &httperror.HandlerError{http.StatusConflict, "Unable to enable extension", baasapi.ErrExtensionAlreadyEnabled}
		}
	}

	extension := &baasapi.Extension{
		ID: extensionID,
	}

	extensionDefinitions, err := handler.ExtensionManager.FetchExtensionDefinitions()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve extension definitions", err}
	}

	for _, def := range extensionDefinitions {
		if def.ID == extension.ID {
			extension.Version = def.Version
			break
		}
	}

	err = handler.ExtensionManager.EnableExtension(extension, payload.License)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to enable extension", err}
	}

	extension.Enabled = true

	err = handler.ExtensionService.Persist(extension)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist extension status inside the database", err}
	}

	return response.Empty(w)
}
