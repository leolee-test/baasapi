package templates

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

const (
	errTemplateManagementDisabled = baasapi.Error("Template management is disabled")
)

// Handler represents an HTTP API handler for managing templates.
type Handler struct {
	*mux.Router
	TemplateService baasapi.TemplateService
	SettingsService baasapi.SettingsService
}

// NewHandler returns a new instance of Handler.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	h.Handle("/templates",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.templateList))).Methods(http.MethodGet)
	h.Handle("/templates",
		bouncer.AdministratorAccess(h.templateManagementCheck(httperror.LoggerHandler(h.templateCreate)))).Methods(http.MethodPost)
	h.Handle("/templates/{id}",
		bouncer.AdministratorAccess(h.templateManagementCheck(httperror.LoggerHandler(h.templateInspect)))).Methods(http.MethodGet)
	h.Handle("/templates/{id}",
		bouncer.AdministratorAccess(h.templateManagementCheck(httperror.LoggerHandler(h.templateUpdate)))).Methods(http.MethodPut)
	h.Handle("/templates/{id}",
		bouncer.AdministratorAccess(h.templateManagementCheck(httperror.LoggerHandler(h.templateDelete)))).Methods(http.MethodDelete)
	return h
}

func (handler *Handler) templateManagementCheck(next http.Handler) http.Handler {
	return httperror.LoggerHandler(func(rw http.ResponseWriter, r *http.Request) *httperror.HandlerError {
		settings, err := handler.SettingsService.Settings()
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve settings from the database", err}
		}

		if settings.TemplatesURL != "" {
			return &httperror.HandlerError{http.StatusServiceUnavailable, "BaaSapi is configured to use external templates, template management is disabled", errTemplateManagementDisabled}
		}

		next.ServeHTTP(rw, r)
		return nil
	})
}
