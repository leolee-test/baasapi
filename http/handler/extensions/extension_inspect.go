package extensions

import (
	"encoding/json"
	"net/http"

	"github.com/coreos/go-semver/semver"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/client"
)

// GET request on /api/extensions/:id
func (handler *Handler) extensionInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	extensionIdentifier, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid extension identifier route variable", err}
	}
	extensionID := baasapi.ExtensionID(extensionIdentifier)

	extensionData, err := client.Get(baasapi.ExtensionDefinitionsURL, 30)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve extension definitions", err}
	}

	var extensions []baasapi.Extension
	err = json.Unmarshal(extensionData, &extensions)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to parse external extension definitions", err}
	}

	var extension baasapi.Extension
	for _, p := range extensions {
		if p.ID == extensionID {
			extension = p
			if extension.DescriptionURL != "" {
				description, _ := client.Get(extension.DescriptionURL, 10)
				extension.Description = string(description)
			}
			break
		}
	}

	storedExtension, err := handler.ExtensionService.Extension(extensionID)
	if err == baasapi.ErrObjectNotFound {
		return response.JSON(w, extension)
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a extension with the specified identifier inside the database", err}
	}

	extension.Enabled = storedExtension.Enabled

	extensionVer := semver.New(extension.Version)
	pVer := semver.New(storedExtension.Version)

	if pVer.LessThan(*extensionVer) {
		extension.UpdateAvailable = true
	}

	return response.JSON(w, extension)
}
