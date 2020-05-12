package baask8s

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "baask8ss"
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
func (service *Service) Baask8s(ID baasapi.Baask8sID) (*baasapi.Baask8s, error) {
	var baask8s baasapi.Baask8s
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &baask8s)
	if err != nil {
		return nil, err
	}

	return &baask8s, nil
}

// Baask8ss return an array containing all the baask8ss.
func (service *Service) Baask8ss() ([]baasapi.Baask8s, error) {
	var baask8ss = make([]baasapi.Baask8s, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var baask8s baasapi.Baask8s
			err := internal.UnmarshalObject(v, &baask8s)
			if err != nil {
				return err
			}
			baask8ss = append(baask8ss, baask8s)
		}

		return nil
	})

	return baask8ss, err
}

// CreateBaask8s assign an ID to a new baask8s and saves it.
func (service *Service) CreateBaask8s(baask8s *baasapi.Baask8s) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		// We manually manage sequences for baask8ss
		err := bucket.SetSequence(uint64(baask8s.ID))
		if err != nil {
			return err
		}

		data, err := internal.MarshalObject(baask8s)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(baask8s.ID)), data)
	})
}

// UpdateBaask8s updates an baask8s.
func (service *Service) UpdateBaask8s(ID baasapi.Baask8sID, baask8s *baasapi.Baask8s) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, baask8s)
}

// DeleteBaask8s deletes an baask8s.
func (service *Service) DeleteBaask8s(ID baasapi.Baask8sID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}

// GetNextIdentifier returns the next identifier for an baask8s.
func (service *Service) GetNextIdentifier() int {
	return internal.GetNextIdentifier(service.db, BucketName)
}


