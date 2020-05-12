package registries

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

func hideFields(registry *baasapi.Registry) {
	registry.Password = ""
	registry.ManagementConfiguration = nil
}

// Handler is the HTTP handler used to handle registry operations.
type Handler struct {
	*mux.Router
	requestBouncer   *security.RequestBouncer
	RegistryService  baasapi.RegistryService
	ExtensionService baasapi.ExtensionService
	FileService      baasapi.FileService
}

// NewHandler creates a handler to manage registry operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
	}

	h.Handle("/registries",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.registryCreate))).Methods(http.MethodPost)
	h.Handle("/registries",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.registryList))).Methods(http.MethodGet)
	h.Handle("/registries/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.registryInspect))).Methods(http.MethodGet)
	h.Handle("/registries/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.registryUpdate))).Methods(http.MethodPut)
	h.Handle("/registries/{id}/access",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.registryUpdateAccess))).Methods(http.MethodPut)
	h.Handle("/registries/{id}/configure",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.registryConfigure))).Methods(http.MethodPost)
	h.Handle("/registries/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.registryDelete))).Methods(http.MethodDelete)

	return h
}
