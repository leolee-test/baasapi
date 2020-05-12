package extension

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "extension"
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

// Extension returns a extension by ID
func (service *Service) Extension(ID baasapi.ExtensionID) (*baasapi.Extension, error) {
	var extension baasapi.Extension
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &extension)
	if err != nil {
		return nil, err
	}

	return &extension, nil
}

// Extensions return an array containing all the extensions.
func (service *Service) Extensions() ([]baasapi.Extension, error) {
	var extensions = make([]baasapi.Extension, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var extension baasapi.Extension
			err := internal.UnmarshalObject(v, &extension)
			if err != nil {
				return err
			}
			extensions = append(extensions, extension)
		}

		return nil
	})

	return extensions, err
}

// Persist persists a extension inside the database.
func (service *Service) Persist(extension *baasapi.Extension) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		data, err := internal.MarshalObject(extension)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(extension.ID)), data)
	})
}

// DeleteExtension deletes a Extension.
func (service *Service) DeleteExtension(ID baasapi.ExtensionID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}
