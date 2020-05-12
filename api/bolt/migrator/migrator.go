package migrator

import (
	"github.com/boltdb/bolt"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/baask8s"
	"github.com/baasapi/baasapi/api/bolt/baasmsp"
	"github.com/baasapi/baasapi/api/bolt/extension"
	"github.com/baasapi/baasapi/api/bolt/resourcecontrol"
	"github.com/baasapi/baasapi/api/bolt/settings"
	"github.com/baasapi/baasapi/api/bolt/template"
	"github.com/baasapi/baasapi/api/bolt/user"
	"github.com/baasapi/baasapi/api/bolt/version"
)

type (
	// Migrator defines a service to migrate data after a BaaSapi version update.
	Migrator struct {
		currentDBVersion       int
		db                     *bolt.DB
		baask8sService         *baask8s.Service
		baasmspService         *baasmsp.Service
		extensionService       *extension.Service
		resourceControlService *resourcecontrol.Service
		settingsService        *settings.Service
		templateService        *template.Service
		userService            *user.Service
		versionService         *version.Service
		fileService            baasapi.FileService
	}

	// Parameters represents the required parameters to create a new Migrator instance.
	Parameters struct {
		DB                     *bolt.DB
		DatabaseVersion        int
		Baask8sService         *baask8s.Service
		BaasmspService         *baasmsp.Service
		ExtensionService       *extension.Service
		ResourceControlService *resourcecontrol.Service
		SettingsService        *settings.Service
		TemplateService        *template.Service
		UserService            *user.Service
		VersionService         *version.Service
		FileService            baasapi.FileService
	}
)

// NewMigrator creates a new Migrator.
func NewMigrator(parameters *Parameters) *Migrator {
	return &Migrator{
		db:                     parameters.DB,
		currentDBVersion:       parameters.DatabaseVersion,
		baask8sService:         parameters.Baask8sService,
		baasmspService:         parameters.BaasmspService,
		extensionService:       parameters.ExtensionService,
		resourceControlService: parameters.ResourceControlService,
		settingsService:        parameters.SettingsService,
		templateService:        parameters.TemplateService,
		userService:            parameters.UserService,
		versionService:         parameters.VersionService,
		fileService:            parameters.FileService,
	}
}

// Migrate checks the database version and migrate the existing data to the most recent data model.
func (m *Migrator) Migrate() error {

	// BaaSapi 1.20.3
	if m.currentDBVersion < 18 {
		err := m.updateExtensionsToDBVersion18()
		if err != nil {
			return err
		}
	}


	return m.versionService.StoreDBVersion(baasapi.DBVersion)
}
