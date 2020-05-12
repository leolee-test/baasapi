package settings

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

func hideFields(settings *baasapi.Settings) {
	settings.LDAPSettings.Password = ""
	settings.OAuthSettings.ClientSecret = ""
}

// Handler is the HTTP handler used to handle settings operations.
type Handler struct {
	*mux.Router
	SettingsService baasapi.SettingsService
	LDAPService     baasapi.LDAPService
	FileService     baasapi.FileService
	JobScheduler    baasapi.JobScheduler
	ScheduleService baasapi.ScheduleService
}

// NewHandler creates a handler to manage settings operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/settings",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.settingsInspect))).Methods(http.MethodGet)
	h.Handle("/settings",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.settingsUpdate))).Methods(http.MethodPut)
	h.Handle("/settings/public",
		bouncer.PublicAccess(httperror.LoggerHandler(h.settingsPublic))).Methods(http.MethodGet)
	h.Handle("/settings/authentication/checkLDAP",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.settingsLDAPCheck))).Methods(http.MethodPut)

	return h
}
