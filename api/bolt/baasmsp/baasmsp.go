package baasmsp

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "baasmsps"
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

// Baask8s returns an baask8s by ID.
func (service *Service) Baasmsp(ID baasapi.BaasmspID) (*baasapi.Baasmsp, error) {
	var baasmsp baasapi.Baasmsp
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &baasmsp)
	if err != nil {
		return nil, err
	}

	return &baasmsp, nil
}

// Baask8ss return an array containing all the baask8ss.
func (service *Service) Baasmsps() ([]baasapi.Baasmsp, error) {
	var baasmsps = make([]baasapi.Baasmsp, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var baasmsp baasapi.Baasmsp
			err := internal.UnmarshalObject(v, &baasmsp)
			if err != nil {
				return err
			}
			baasmsps = append(baasmsps, baasmsp)
		}

		return nil
	})

	return baasmsps, err
}

// CreateBaasmsp assign an ID to a new baask8s and saves it.
func (service *Service) CreateBaasmsp(baasmsp *baasapi.Baasmsp) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		// We manually manage sequences for baasmsps
		err := bucket.SetSequence(uint64(baasmsp.ID))
		if err != nil {
			return err
		}

		data, err := internal.MarshalObject(baasmsp)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(baasmsp.ID)), data)
	})
}

// DeleteBaasmsp deletes an baasmsp.
func (service *Service) DeleteBaasmsp(ID baasapi.BaasmspID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}

// GetNextIdentifier returns the next identifier for an baask8s.
func (service *Service) GetNextIdentifier() int {
	return internal.GetNextIdentifier(service.db, BucketName)
}


