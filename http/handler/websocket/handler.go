package websocket

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// Handler is the HTTP handler used to handle websocket operations.
type Handler struct {
	*mux.Router
	Baask8sService    baasapi.Baask8sService
	SignatureService   baasapi.DigitalSignatureService
	requestBouncer     *security.RequestBouncer
	connectionUpgrader websocket.Upgrader
}

// NewHandler creates a handler to manage websocket operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:             mux.NewRouter(),
		connectionUpgrader: websocket.Upgrader{},
		requestBouncer:     bouncer,
	}
	h.PathPrefix("/websocket/exec").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.websocketExec)))
	return h
}
