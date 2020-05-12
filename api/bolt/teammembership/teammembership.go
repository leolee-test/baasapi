package teammembership

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "team_membership"
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

// TeamMembership returns a TeamMembership object by ID
func (service *Service) TeamMembership(ID baasapi.TeamMembershipID) (*baasapi.TeamMembership, error) {
	var membership baasapi.TeamMembership
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &membership)
	if err != nil {
		return nil, err
	}

	return &membership, nil
}

// TeamMemberships return an array containing all the TeamMembership objects.
func (service *Service) TeamMemberships() ([]baasapi.TeamMembership, error) {
	var memberships = make([]baasapi.TeamMembership, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var membership baasapi.TeamMembership
			err := internal.UnmarshalObject(v, &membership)
			if err != nil {
				return err
			}
			memberships = append(memberships, membership)
		}

		return nil
	})

	return memberships, err
}

// TeamMembershipsByUserID return an array containing all the TeamMembership objects where the specified userID is present.
func (service *Service) TeamMembershipsByUserID(userID baasapi.UserID) ([]baasapi.TeamMembership, error) {
	var memberships = make([]baasapi.TeamMembership, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var membership baasapi.TeamMembership
			err := internal.UnmarshalObject(v, &membership)
			if err != nil {
				return err
			}

			if membership.UserID == userID {
				memberships = append(memberships, membership)
			}
		}

		return nil
	})

	return memberships, err
}

// TeamMembershipsByTeamID return an array containing all the TeamMembership objects where the specified teamID is present.
func (service *Service) TeamMembershipsByTeamID(teamID baasapi.TeamID) ([]baasapi.TeamMembership, error) {
	var memberships = make([]baasapi.TeamMembership, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var membership baasapi.TeamMembership
			err := internal.UnmarshalObject(v, &membership)
			if err != nil {
				return err
			}

			if membership.TeamID == teamID {
				memberships = append(memberships, membership)
			}
		}

		return nil
	})

	return memberships, err
}

// UpdateTeamMembership saves a TeamMembership object.
func (service *Service) UpdateTeamMembership(ID baasapi.TeamMembershipID, membership *baasapi.TeamMembership) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, membership)
}

// CreateTeamMembership creates a new TeamMembership object.
func (service *Service) CreateTeamMembership(membership *baasapi.TeamMembership) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		id, _ := bucket.NextSequence()
		membership.ID = baasapi.TeamMembershipID(id)

		data, err := internal.MarshalObject(membership)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(membership.ID)), data)
	})
}

// DeleteTeamMembership deletes a TeamMembership object.
func (service *Service) DeleteTeamMembership(ID baasapi.TeamMembershipID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}

// DeleteTeamMembershipByUserID deletes all the TeamMembership object associated to a UserID.
func (service *Service) DeleteTeamMembershipByUserID(userID baasapi.UserID) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var membership baasapi.TeamMembership
			err := internal.UnmarshalObject(v, &membership)
			if err != nil {
				return err
			}

			if membership.UserID == userID {
				err := bucket.Delete(internal.Itob(int(membership.ID)))
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// DeleteTeamMembershipByTeamID deletes all the TeamMembership object associated to a TeamID.
func (service *Service) DeleteTeamMembershipByTeamID(teamID baasapi.TeamID) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var membership baasapi.TeamMembership
			err := internal.UnmarshalObject(v, &membership)
			if err != nil {
				return err
			}

			if membership.TeamID == teamID {
				err := bucket.Delete(internal.Itob(int(membership.ID)))
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
