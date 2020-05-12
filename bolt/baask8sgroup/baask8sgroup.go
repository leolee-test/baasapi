package baask8sgroup

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "baask8s_groups"
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

// Baask8sGroup returns an baask8s group by ID.
func (service *Service) Baask8sGroup(ID baasapi.Baask8sGroupID) (*baasapi.Baask8sGroup, error) {
	var baask8sGroup baasapi.Baask8sGroup
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &baask8sGroup)
	if err != nil {
		return nil, err
	}

	return &baask8sGroup, nil
}

// UpdateBaask8sGroup updates an baask8s group.
func (service *Service) UpdateBaask8sGroup(ID baasapi.Baask8sGroupID, baask8sGroup *baasapi.Baask8sGroup) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, baask8sGroup)
}

// DeleteBaask8sGroup deletes an baask8s group.
func (service *Service) DeleteBaask8sGroup(ID baasapi.Baask8sGroupID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}

// Baask8sGroups return an array containing all the baask8s groups.
func (service *Service) Baask8sGroups() ([]baasapi.Baask8sGroup, error) {
	var baask8sGroups = make([]baasapi.Baask8sGroup, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var baask8sGroup baasapi.Baask8sGroup
			err := internal.UnmarshalObject(v, &baask8sGroup)
			if err != nil {
				return err
			}
			baask8sGroups = append(baask8sGroups, baask8sGroup)
		}

		return nil
	})

	return baask8sGroups, err
}

// CreateBaask8sGroup assign an ID to a new baask8s group and saves it.
func (service *Service) CreateBaask8sGroup(baask8sGroup *baasapi.Baask8sGroup) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		id, _ := bucket.NextSequence()
		baask8sGroup.ID = baasapi.Baask8sGroupID(id)

		data, err := internal.MarshalObject(baask8sGroup)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(baask8sGroup.ID)), data)
	})
}
