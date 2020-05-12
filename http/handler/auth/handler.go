package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

const (
	// ErrInvalidCredentials is an error raised when credentials for a user are invalid
	ErrInvalidCredentials = baasapi.Error("Invalid credentials")
	// ErrAuthDisabled is an error raised when trying to access the authentication baask8ss
	// when the server has been started with the --no-auth flag
	ErrAuthDisabled = baasapi.Error("Authentication is disabled")
)

// Handler is the HTTP handler used to handle authentication operations.
type Handler struct {
	*mux.Router
	authDisabled          bool
	UserService           baasapi.UserService
	CryptoService         baasapi.CryptoService
	JWTService            baasapi.JWTService
	LDAPService           baasapi.LDAPService
	SettingsService       baasapi.SettingsService
	TeamService           baasapi.TeamService
	TeamMembershipService baasapi.TeamMembershipService
	ExtensionService      baasapi.ExtensionService
}

// NewHandler creates a handler to manage authentication operations.
func NewHandler(bouncer *security.RequestBouncer, rateLimiter *security.RateLimiter, authDisabled bool) *Handler {
	h := &Handler{
		Router:       mux.NewRouter(),
		authDisabled: authDisabled,
	}

	h.Handle("/auth/oauth/validate",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.validateOAuth)))).Methods(http.MethodPost)
	h.Handle("/auth",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.authenticate)))).Methods(http.MethodPost, "OPTIONS")

	return h
}
