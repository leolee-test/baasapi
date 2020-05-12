package registry

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "registries"
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

// Registry returns an registry by ID.
func (service *Service) Registry(ID baasapi.RegistryID) (*baasapi.Registry, error) {
	var registry baasapi.Registry
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &registry)
	if err != nil {
		return nil, err
	}

	return &registry, nil
}

// Registries returns an array containing all the registries.
func (service *Service) Registries() ([]baasapi.Registry, error) {
	var registries = make([]baasapi.Registry, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var registry baasapi.Registry
			err := internal.UnmarshalObject(v, &registry)
			if err != nil {
				return err
			}
			registries = append(registries, registry)
		}

		return nil
	})

	return registries, err
}

// CreateRegistry creates a new registry.
func (service *Service) CreateRegistry(registry *baasapi.Registry) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		id, _ := bucket.NextSequence()
		registry.ID = baasapi.RegistryID(id)

		data, err := internal.MarshalObject(registry)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(registry.ID)), data)
	})
}

// UpdateRegistry updates an registry.
func (service *Service) UpdateRegistry(ID baasapi.RegistryID, registry *baasapi.Registry) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, registry)
}

// DeleteRegistry deletes an registry.
func (service *Service) DeleteRegistry(ID baasapi.RegistryID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}
