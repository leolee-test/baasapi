package handler

import (
	"net/http"
	"strings"

	"github.com/baasapi/baasapi/api/http/handler/auth"
	"github.com/baasapi/baasapi/api/http/handler/baask8ss"
	"github.com/baasapi/baasapi/api/http/handler/extensions"
	"github.com/baasapi/baasapi/api/http/handler/file"
	"github.com/baasapi/baasapi/api/http/handler/motd"
	"github.com/baasapi/baasapi/api/http/handler/registries"
	"github.com/baasapi/baasapi/api/http/handler/resourcecontrols"
	"github.com/baasapi/baasapi/api/http/handler/schedules"
	"github.com/baasapi/baasapi/api/http/handler/settings"
	"github.com/baasapi/baasapi/api/http/handler/status"
	"github.com/baasapi/baasapi/api/http/handler/tags"
	"github.com/baasapi/baasapi/api/http/handler/teammemberships"
	"github.com/baasapi/baasapi/api/http/handler/teams"
	"github.com/baasapi/baasapi/api/http/handler/templates"
	"github.com/baasapi/baasapi/api/http/handler/upload"
	"github.com/baasapi/baasapi/api/http/handler/users"
	//"github.com/baasapi/baasapi/api/http/handler/webhooks"
	//"github.com/baasapi/baasapi/api/http/handler/websocket"
)

// Handler is a collection of all the service handlers.
type Handler struct {
	AuthHandler *auth.Handler

	Baask8sHandler         *baask8ss.Handler
	FileHandler            *file.Handler
	MOTDHandler            *motd.Handler
	ExtensionHandler       *extensions.Handler
	RegistryHandler        *registries.Handler
	ResourceControlHandler *resourcecontrols.Handler
	SettingsHandler        *settings.Handler
	StatusHandler          *status.Handler
	TagHandler             *tags.Handler
	TeamMembershipHandler  *teammemberships.Handler
	TeamHandler            *teams.Handler
	TemplatesHandler       *templates.Handler
	UploadHandler          *upload.Handler
	UserHandler            *users.Handler
	SchedulesHanlder       *schedules.Handler
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	(*w).Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*r).Method == "OPTIONS" {
				return
					}

	switch {
	case strings.HasPrefix(r.URL.Path, "/api/auth"):
		http.StripPrefix("/api", h.AuthHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/baask8s"):
		http.StripPrefix("/api", h.Baask8sHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/motd"):
		http.StripPrefix("/api", h.MOTDHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/extensions"):
		http.StripPrefix("/api", h.ExtensionHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/registries"):
		http.StripPrefix("/api", h.RegistryHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/resource_controls"):
		http.StripPrefix("/api", h.ResourceControlHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/settings"):
		http.StripPrefix("/api", h.SettingsHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/status"):
		http.StripPrefix("/api", h.StatusHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/tags"):
		http.StripPrefix("/api", h.TagHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/templates"):
		http.StripPrefix("/api", h.TemplatesHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/upload"):
		http.StripPrefix("/api", h.UploadHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/users"):
		http.StripPrefix("/api", h.UserHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/groups"):
		http.StripPrefix("/api", h.TeamHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/group_memberships"):
		http.StripPrefix("/api", h.TeamMembershipHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/schedules"):
		http.StripPrefix("/api", h.SchedulesHanlder).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/"):
		h.FileHandler.ServeHTTP(w, r)
	}
}
