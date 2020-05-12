package http

import (
	"time"

	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/handler"
	"github.com/baasapi/baasapi/api/http/handler/auth"
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
	"github.com/baasapi/baasapi/api/http/handler/baask8ss"
	//"github.com/baasapi/baasapi/api/http/handler/baasmsps"
	"github.com/baasapi/baasapi/api/http/security"

	"net/http"
	"path/filepath"
)

// Server implements the baasapi.Server interface
type Server struct {
	BindAddress            string
	AssetsPath             string
	AuthDisabled           bool
	Baask8sManagement     bool
	Status                 *baasapi.Status
	ExtensionManager       baasapi.ExtensionManager
	CryptoService          baasapi.CryptoService
	SignatureService       baasapi.DigitalSignatureService
	JobScheduler           baasapi.JobScheduler
	Snapshotter            baasapi.Snapshotter
	Baask8sService         baasapi.Baask8sService
	BaasmspService         baasapi.BaasmspService
	Baask8sGroupService   baasapi.Baask8sGroupService
	FileService            baasapi.FileService
	GitService             baasapi.GitService
	JWTService             baasapi.JWTService
	LDAPService            baasapi.LDAPService
	ExtensionService       baasapi.ExtensionService
	RegistryService        baasapi.RegistryService
	ResourceControlService baasapi.ResourceControlService
	ScheduleService        baasapi.ScheduleService
	SettingsService        baasapi.SettingsService
	CAFilesManager         baasapi.CAFilesManager
	TagService             baasapi.TagService
	TeamService            baasapi.TeamService
	TeamMembershipService  baasapi.TeamMembershipService
	TemplateService        baasapi.TemplateService
	UserService            baasapi.UserService
	Handler                *handler.Handler
	SSL                    bool
	SSLCert                string
	SSLKey                 string
	JobService             baasapi.JobService
}

// Start starts the HTTP server
func (server *Server) Start() error {
	requestBouncerParameters := &security.RequestBouncerParams{
		JWTService:            server.JWTService,
		UserService:           server.UserService,
		TeamMembershipService: server.TeamMembershipService,
		Baask8sGroupService:  server.Baask8sGroupService,
		AuthDisabled:          server.AuthDisabled,
	}
	requestBouncer := security.NewRequestBouncer(requestBouncerParameters)


	rateLimiter := security.NewRateLimiter(10, 1*time.Second, 1*time.Hour)

	var authHandler = auth.NewHandler(requestBouncer, rateLimiter, server.AuthDisabled)
	authHandler.UserService = server.UserService
	authHandler.CryptoService = server.CryptoService
	authHandler.JWTService = server.JWTService
	authHandler.LDAPService = server.LDAPService
	authHandler.SettingsService = server.SettingsService
	authHandler.TeamService = server.TeamService
	authHandler.TeamMembershipService = server.TeamMembershipService
	authHandler.ExtensionService = server.ExtensionService

	var baask8sHandler = baask8ss.NewHandler(requestBouncer)
	baask8sHandler.Baask8sService = server.Baask8sService
	baask8sHandler.BaasmspService = server.BaasmspService
	baask8sHandler.JWTService = server.JWTService
	baask8sHandler.FileService = server.FileService
	baask8sHandler.Snapshotter = server.Snapshotter
	baask8sHandler.JobService = server.JobService
	baask8sHandler.CAFilesManager = server.CAFilesManager

	var fileHandler = file.NewHandler(filepath.Join(server.AssetsPath, "public"))

	var motdHandler = motd.NewHandler(requestBouncer)

	var extensionHandler = extensions.NewHandler(requestBouncer)
	extensionHandler.ExtensionService = server.ExtensionService
	extensionHandler.ExtensionManager = server.ExtensionManager

	var registryHandler = registries.NewHandler(requestBouncer)
	registryHandler.RegistryService = server.RegistryService
	registryHandler.ExtensionService = server.ExtensionService
	registryHandler.FileService = server.FileService

	var resourceControlHandler = resourcecontrols.NewHandler(requestBouncer)
	resourceControlHandler.ResourceControlService = server.ResourceControlService

	var schedulesHandler = schedules.NewHandler(requestBouncer)
	schedulesHandler.ScheduleService = server.ScheduleService
	schedulesHandler.Baask8sService = server.Baask8sService
	schedulesHandler.FileService = server.FileService
	schedulesHandler.JobService = server.JobService
	schedulesHandler.JobScheduler = server.JobScheduler
	schedulesHandler.SettingsService = server.SettingsService

	var settingsHandler = settings.NewHandler(requestBouncer)
	settingsHandler.SettingsService = server.SettingsService
	settingsHandler.LDAPService = server.LDAPService
	settingsHandler.FileService = server.FileService
	settingsHandler.JobScheduler = server.JobScheduler
	settingsHandler.ScheduleService = server.ScheduleService

	var tagHandler = tags.NewHandler(requestBouncer)
	tagHandler.TagService = server.TagService

	var teamHandler = teams.NewHandler(requestBouncer)
	teamHandler.TeamService = server.TeamService
	teamHandler.TeamMembershipService = server.TeamMembershipService

	var teamMembershipHandler = teammemberships.NewHandler(requestBouncer)
	teamMembershipHandler.TeamMembershipService = server.TeamMembershipService
	var statusHandler = status.NewHandler(requestBouncer, server.Status)

	var templatesHandler = templates.NewHandler(requestBouncer)
	templatesHandler.TemplateService = server.TemplateService
	templatesHandler.SettingsService = server.SettingsService

	var uploadHandler = upload.NewHandler(requestBouncer)
	uploadHandler.FileService = server.FileService

	var userHandler = users.NewHandler(requestBouncer, rateLimiter)
	userHandler.UserService = server.UserService
	userHandler.TeamService = server.TeamService
	userHandler.TeamMembershipService = server.TeamMembershipService
	userHandler.CryptoService = server.CryptoService
	userHandler.ResourceControlService = server.ResourceControlService
	userHandler.SettingsService = server.SettingsService

	server.Handler = &handler.Handler{
		AuthHandler:            authHandler,
		Baask8sHandler:         baask8sHandler,
		FileHandler:            fileHandler,
		MOTDHandler:            motdHandler,
		ExtensionHandler:       extensionHandler,
		RegistryHandler:        registryHandler,
		ResourceControlHandler: resourceControlHandler,
		SettingsHandler:        settingsHandler,
		StatusHandler:          statusHandler,
		TagHandler:             tagHandler,
		TeamHandler:            teamHandler,
		TeamMembershipHandler:  teamMembershipHandler,
		TemplatesHandler:       templatesHandler,
		UploadHandler:          uploadHandler,
		UserHandler:            userHandler,
		SchedulesHanlder:       schedulesHandler,
	}

	if server.SSL {
		return http.ListenAndServeTLS(server.BindAddress, server.SSLCert, server.SSLKey, server.Handler)
	}
	return http.ListenAndServe(server.BindAddress, server.Handler)
}
