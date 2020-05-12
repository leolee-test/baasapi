package status

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// Handler is the HTTP handler used to handle status operations.
type Handler struct {
	*mux.Router
	Status *baasapi.Status
}

// NewHandler creates a handler to manage status operations.
func NewHandler(bouncer *security.RequestBouncer, status *baasapi.Status) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
		Status: status,
	}
	h.Handle("/status",
		bouncer.PublicAccess(httperror.LoggerHandler(h.statusInspect))).Methods(http.MethodGet)

	return h
}
