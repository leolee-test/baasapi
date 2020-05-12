package registries

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type registryCreatePayload struct {
	Name           string
	Type           int
	URL            string
	Authentication bool
	Username       string
	Password       string
}

func (payload *registryCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Name) {
		return baasapi.Error("Invalid registry name")
	}
	if govalidator.IsNull(payload.URL) {
		return baasapi.Error("Invalid registry URL")
	}
	if payload.Authentication && (govalidator.IsNull(payload.Username) || govalidator.IsNull(payload.Password)) {
		return baasapi.Error("Invalid credentials. Username and password must be specified when authentication is enabled")
	}
	if payload.Type != 1 && payload.Type != 2 && payload.Type != 3 {
		return baasapi.Error("Invalid registry type. Valid values are: 1 (Quay.io), 2 (Azure container registry) or 3 (custom registry)")
	}
	return nil
}

func (handler *Handler) registryCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload registryCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	registries, err := handler.RegistryService.Registries()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve registries from the database", err}
	}
	for _, r := range registries {
		if r.URL == payload.URL {
			return &httperror.HandlerError{http.StatusConflict, "A registry with the same URL already exists", baasapi.ErrRegistryAlreadyExists}
		}
	}

	registry := &baasapi.Registry{
		Type:            baasapi.RegistryType(payload.Type),
		Name:            payload.Name,
		URL:             payload.URL,
		Authentication:  payload.Authentication,
		Username:        payload.Username,
		Password:        payload.Password,
		AuthorizedUsers: []baasapi.UserID{},
		AuthorizedTeams: []baasapi.TeamID{},
	}

	err = handler.RegistryService.CreateRegistry(registry)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the registry inside the database", err}
	}

	hideFields(registry)
	return response.JSON(w, registry)
}
