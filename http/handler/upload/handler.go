package upload

import (
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

// Handler is the HTTP handler used to handle upload operations.
type Handler struct {
	*mux.Router
	FileService baasapi.FileService
}

// NewHandler creates a handler to manage upload operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/upload/tls/{certificate:(?:ca|cert|key)}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.uploadTLS))).Methods(http.MethodPost)
	h.Handle("/upload/kubeconfig",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.uploadKubeconfig))).Methods(http.MethodPost)
	return h
}
