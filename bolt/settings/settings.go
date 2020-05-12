package settings

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName  = "settings"
	settingsKey = "SETTINGS"
)

// Service represents a service for managing baask8s data.
type Service struct {
	db *bolt.DB
}

// NewService creates a new instance of a service.
func NewService(db *bolt.DB) (*Service, error) {
	err := internal.CreateBucket(db, BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

// Settings retrieve the settings object.
func (service *Service) Settings() (*baasapi.Settings, error) {
	var settings baasapi.Settings

	err := internal.GetObject(service.db, BucketName, []byte(settingsKey), &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// UpdateSettings persists a Settings object.
func (service *Service) UpdateSettings(settings *baasapi.Settings) error {
	return internal.UpdateObject(service.db, BucketName, []byte(settingsKey), settings)
}
