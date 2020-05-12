package stacks

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// Handler is the HTTP handler used to handle stack operations.
type Handler struct {
	stackCreationMutex *sync.Mutex
	stackDeletionMutex *sync.Mutex
	requestBouncer     *security.RequestBouncer
	*mux.Router
	FileService            baasapi.FileService
	GitService             baasapi.GitService
	StackService           baasapi.StackService
	EndpointService        baasapi.EndpointService
	ResourceControlService baasapi.ResourceControlService
	RegistryService        baasapi.RegistryService
	DockerHubService       baasapi.DockerHubService
	SwarmStackManager      baasapi.SwarmStackManager
	ComposeStackManager    baasapi.ComposeStackManager
}

// NewHandler creates a handler to manage stack operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:             mux.NewRouter(),
		stackCreationMutex: &sync.Mutex{},
		stackDeletionMutex: &sync.Mutex{},
		requestBouncer:     bouncer,
	}
	h.Handle("/stacks",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.stackCreate))).Methods(http.MethodPost)
	h.Handle("/stacks",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.stackList))).Methods(http.MethodGet)
	h.Handle("/stacks/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.stackInspect))).Methods(http.MethodGet)
	h.Handle("/stacks/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.stackDelete))).Methods(http.MethodDelete)
	h.Handle("/stacks/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.stackUpdate))).Methods(http.MethodPut)
	h.Handle("/stacks/{id}/file",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.stackFile))).Methods(http.MethodGet)
	h.Handle("/stacks/{id}/migrate",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.stackMigrate))).Methods(http.MethodPost)
	return h
}
