package user

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "users"
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

// User returns a user by ID
func (service *Service) User(ID baasapi.UserID) (*baasapi.User, error) {
	var user baasapi.User
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UserByUsername returns a user by username.
func (service *Service) UserByUsername(username string) (*baasapi.User, error) {
	var user *baasapi.User

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var u baasapi.User
			err := internal.UnmarshalObject(v, &u)
			if err != nil {
				return err
			}

			if u.Username == username {
				user = &u
				break
			}
		}

		if user == nil {
			return baasapi.ErrObjectNotFound
		}
		return nil
	})

	return user, err
}

// Users return an array containing all the users.
func (service *Service) Users() ([]baasapi.User, error) {
	var users = make([]baasapi.User, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var user baasapi.User
			err := internal.UnmarshalObject(v, &user)
			if err != nil {
				return err
			}
			users = append(users, user)
		}

		return nil
	})

	return users, err
}

// UsersByRole return an array containing all the users with the specified role.
func (service *Service) UsersByRole(role baasapi.UserRole) ([]baasapi.User, error) {
	var users = make([]baasapi.User, 0)
	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var user baasapi.User
			err := internal.UnmarshalObject(v, &user)
			if err != nil {
				return err
			}

			if user.Role == role {
				users = append(users, user)
			}
		}
		return nil
	})

	return users, err
}

// UpdateUser saves a user.
func (service *Service) UpdateUser(ID baasapi.UserID, user *baasapi.User) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, user)
}

// CreateUser creates a new user.
func (service *Service) CreateUser(user *baasapi.User) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		id, _ := bucket.NextSequence()
		user.ID = baasapi.UserID(id)

		data, err := internal.MarshalObject(user)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(user.ID)), data)
	})
}

// DeleteUser deletes a user.
func (service *Service) DeleteUser(ID baasapi.UserID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}
