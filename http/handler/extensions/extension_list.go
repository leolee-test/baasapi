package extensions

import (
	"net/http"

	"github.com/coreos/go-semver/semver"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// GET request on /api/extensions?store=<store>
func (handler *Handler) extensionList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	storeDetails, _ := request.RetrieveBooleanQueryParameter(r, "store", true)

	extensions, err := handler.ExtensionService.Extensions()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve extensions from the database", err}
	}

	if storeDetails {
		definitions, err := handler.ExtensionManager.FetchExtensionDefinitions()
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve extensions", err}
		}

		for idx := range definitions {
			associateExtensionData(&definitions[idx], extensions)
		}

		extensions = definitions
	}

	return response.JSON(w, extensions)
}

func associateExtensionData(definition *baasapi.Extension, extensions []baasapi.Extension) {
	for _, extension := range extensions {
		if extension.ID == definition.ID {

			definition.Enabled = extension.Enabled
			definition.License.Company = extension.License.Company
			definition.License.Expiration = extension.License.Expiration
			definition.License.Valid = extension.License.Valid

			definitionVersion := semver.New(definition.Version)
			extensionVersion := semver.New(extension.Version)
			if extensionVersion.LessThan(*definitionVersion) {
				definition.UpdateAvailable = true
			}

			break
		}
	}
}
