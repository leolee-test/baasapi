package dockerhub

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

func hideFields(dockerHub *baasapi.DockerHub) {
	dockerHub.Password = ""
}

// Handler is the HTTP handler used to handle DockerHub operations.
type Handler struct {
	*mux.Router
	DockerHubService baasapi.DockerHubService
}

// NewHandler creates a handler to manage Dockerhub operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/dockerhub",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.dockerhubInspect))).Methods(http.MethodGet)
	h.Handle("/dockerhub",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.dockerhubUpdate))).Methods(http.MethodPut)

	return h
}
