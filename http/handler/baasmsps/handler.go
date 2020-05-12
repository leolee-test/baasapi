package baasmsps

import (
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

const (
	// ErrBaask8sManagementDisabled is an error raised when trying to access the baask8ss management baask8ss
	// when the server has been started with the --external-baask8ss flag
	ErrbaasmspManagementDisabled = baasapi.Error("baasmsp management is disabled")
)

//func hideFields(baask8s *baasapi.Baask8s) {
//	baask8s.AzureCredentials = baasapi.AzureCredentials{}
//}

// Handler is the HTTP handler used to handle baask8s operations.
type Handler struct {
	*mux.Router
	requestBouncer              *security.RequestBouncer
	Baask8sService              baasapi.Baask8sService
	BaasmspService              baasapi.BaasmspService
	Baask8sGroupService        baasapi.Baask8sGroupService
	FileService                 baasapi.FileService
	Snapshotter                 baasapi.Snapshotter
	CAFilesManager              baasapi.CAFilesManager
	JobService                  baasapi.JobService
}

// NewHandler creates a handler to manage baask8s operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
		requestBouncer:              bouncer,
	}

	h.Handle("/baasmsps",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baasmspCreate))).Methods(http.MethodPost)
	//h.Handle("/baask8ss/snapshot",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sSnapshots))).Methods(http.MethodPost)
	h.Handle("/baasmsps",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baasmspList))).Methods(http.MethodGet)
	//h.Handle("/baasmsps/{id}/pods",
	//	bouncer.PublicAccess(httperror.LoggerHandler(h.baasmspPodsList))).Methods(http.MethodGet)
	//h.Handle("/baasmsps/{id}/pods",
	//	bouncer.PublicAccess(httperror.LoggerHandler(h.baasmspPodOperations))).Methods(http.MethodPost)
	//h.Handle("/baask8ss/{id}",
	//	bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sInspect))).Methods(http.MethodGet)
	//h.Handle("/baask8ss/{id}",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sUpdate))).Methods(http.MethodPut)
	//h.Handle("/baask8ss/{id}/access",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sUpdateAccess))).Methods(http.MethodPut)
	h.Handle("/baasmsps/{id}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baasmspDelete))).Methods(http.MethodDelete)
	//h.Handle("/baask8ss/{id}/extensions",
	//	bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.baask8sExtensionAdd))).Methods(http.MethodPost)
	//h.Handle("/baask8ss/{id}/extensions/{extensionType}",
	//	bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.baask8sExtensionRemove))).Methods(http.MethodDelete)
	//h.Handle("/baask8ss/{id}/job",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sJob))).Methods(http.MethodPost)
	//h.Handle("/baask8ss/{id}/snapshot",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sSnapshot))).Methods(http.MethodPost)
	return h
}
