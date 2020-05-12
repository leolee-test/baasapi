package bolt

import (
	"log"
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/baask8s"
	"github.com/baasapi/baasapi/api/bolt/baasmsp"
	"github.com/baasapi/baasapi/api/bolt/extension"
	"github.com/baasapi/baasapi/api/bolt/migrator"
	"github.com/baasapi/baasapi/api/bolt/registry"
	"github.com/baasapi/baasapi/api/bolt/resourcecontrol"
	"github.com/baasapi/baasapi/api/bolt/schedule"
	"github.com/baasapi/baasapi/api/bolt/settings"
	"github.com/baasapi/baasapi/api/bolt/tag"
	"github.com/baasapi/baasapi/api/bolt/team"
	"github.com/baasapi/baasapi/api/bolt/teammembership"
	"github.com/baasapi/baasapi/api/bolt/template"
	"github.com/baasapi/baasapi/api/bolt/user"
	"github.com/baasapi/baasapi/api/bolt/version"
	"github.com/baasapi/baasapi/api/bolt/webhook"
)

const (
	databaseFileName = "baasapi.db"
)

// Store defines the implementation of baasapi.DataStore using
// BoltDB as the storage system.
type Store struct {
	path                   string
	db                     *bolt.DB
	checkForDataMigration  bool
	fileService            baasapi.FileService
	Baask8sService         *baask8s.Service
	BaasmspService         *baasmsp.Service
	ExtensionService       *extension.Service
	RegistryService        *registry.Service
	ResourceControlService *resourcecontrol.Service
	SettingsService        *settings.Service
	TagService             *tag.Service
	TeamMembershipService  *teammembership.Service
	TeamService            *team.Service
	TemplateService        *template.Service
	UserService            *user.Service
	VersionService         *version.Service
	WebhookService         *webhook.Service
	ScheduleService        *schedule.Service
}

// NewStore initializes a new Store and the associated services
func NewStore(storePath string, fileService baasapi.FileService) (*Store, error) {
	store := &Store{
		path:        storePath,
		fileService: fileService,
	}

	databasePath := path.Join(storePath, databaseFileName)
	databaseFileExists, err := fileService.FileExists(databasePath)
	if err != nil {
		return nil, err
	}

	if !databaseFileExists {
		store.checkForDataMigration = false
	} else {
		store.checkForDataMigration = true
	}

	return store, nil
}

// Open opens and initializes the BoltDB database.
func (store *Store) Open() error {
	databasePath := path.Join(store.path, databaseFileName)
	db, err := bolt.Open(databasePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	store.db = db

	return store.initServices()
}

// Init creates the default data set.
func (store *Store) Init() error {
	//groups, err := store.Baask8sGroupService.Baask8sGroups()
	//if err != nil {
	//	return err
	//}

	//if len(groups) == 0 {
	//	unassignedGroup := &baasapi.Baask8sGroup{
	//		Name:            "Unassigned",
	//		Description:     "Unassigned baask8ss",
	//		Labels:          []baasapi.Pair{},
	//		AuthorizedUsers: []baasapi.UserID{},
	//		AuthorizedTeams: []baasapi.TeamID{},
	//		Tags:            []string{},
	//	}

	//	return store.Baask8sGroupService.CreateBaask8sGroup(unassignedGroup)
	//}

	return nil
}

// Close closes the BoltDB database.
func (store *Store) Close() error {
	if store.db != nil {
		return store.db.Close()
	}
	return nil
}

// MigrateData automatically migrate the data based on the DBVersion.
func (store *Store) MigrateData() error {
	if !store.checkForDataMigration {
		return store.VersionService.StoreDBVersion(baasapi.DBVersion)
	}

	version, err := store.VersionService.DBVersion()
	if err == baasapi.ErrObjectNotFound {
		version = 0
	} else if err != nil {
		return err
	}

	if version < baasapi.DBVersion {
		migratorParams := &migrator.Parameters{
			DB:                     store.db,
			DatabaseVersion:        version,
			Baask8sService:         store.Baask8sService,
			BaasmspService:         store.BaasmspService,
			ExtensionService:       store.ExtensionService,
			ResourceControlService: store.ResourceControlService,
			SettingsService:        store.SettingsService,
			TemplateService:        store.TemplateService,
			UserService:            store.UserService,
			VersionService:         store.VersionService,
			FileService:            store.fileService,
		}
		migrator := migrator.NewMigrator(migratorParams)

		log.Printf("Migrating database from version %v to %v.\n", version, baasapi.DBVersion)
		err = migrator.Migrate()
		if err != nil {
			log.Printf("An error occurred during database migration: %s\n", err)
			return err
		}
	}

	return nil
}

func (store *Store) initServices() error {

	baask8sService, err := baask8s.NewService(store.db)
	if err != nil {
		return err
	}
	store.Baask8sService = baask8sService

	baasmspService, err := baasmsp.NewService(store.db)
	if err != nil {
		return err
	}
	store.BaasmspService = baasmspService

	extensionService, err := extension.NewService(store.db)
	if err != nil {
		return err
	}
	store.ExtensionService = extensionService

	registryService, err := registry.NewService(store.db)
	if err != nil {
		return err
	}
	store.RegistryService = registryService

	resourcecontrolService, err := resourcecontrol.NewService(store.db)
	if err != nil {
		return err
	}
	store.ResourceControlService = resourcecontrolService

	settingsService, err := settings.NewService(store.db)
	if err != nil {
		return err
	}
	store.SettingsService = settingsService

	tagService, err := tag.NewService(store.db)
	if err != nil {
		return err
	}
	store.TagService = tagService

	teammembershipService, err := teammembership.NewService(store.db)
	if err != nil {
		return err
	}
	store.TeamMembershipService = teammembershipService

	teamService, err := team.NewService(store.db)
	if err != nil {
		return err
	}
	store.TeamService = teamService

	templateService, err := template.NewService(store.db)
	if err != nil {
		return err
	}
	store.TemplateService = templateService

	userService, err := user.NewService(store.db)
	if err != nil {
		return err
	}
	store.UserService = userService

	versionService, err := version.NewService(store.db)
	if err != nil {
		return err
	}
	store.VersionService = versionService

	webhookService, err := webhook.NewService(store.db)
	if err != nil {
		return err
	}
	store.WebhookService = webhookService

	scheduleService, err := schedule.NewService(store.db)
	if err != nil {
		return err
	}
	store.ScheduleService = scheduleService

	return nil
}
