package webhook

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "webhooks"
)

// Service represents a service for managing webhook data.
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

//Webhooks returns an array of all webhooks
func (service *Service) Webhooks() ([]baasapi.Webhook, error) {
	var webhooks = make([]baasapi.Webhook, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var webhook baasapi.Webhook
			err := internal.UnmarshalObject(v, &webhook)
			if err != nil {
				return err
			}
			webhooks = append(webhooks, webhook)
		}

		return nil
	})

	return webhooks, err
}

// Webhook returns a webhook by ID.
func (service *Service) Webhook(ID baasapi.WebhookID) (*baasapi.Webhook, error) {
	var webhook baasapi.Webhook
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}

// WebhookByResourceID returns a webhook by the ResourceID it is associated with.
func (service *Service) WebhookByResourceID(ID string) (*baasapi.Webhook, error) {
	var webhook *baasapi.Webhook

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var w baasapi.Webhook
			err := internal.UnmarshalObject(v, &w)
			if err != nil {
				return err
			}

			if w.ResourceID == ID {
				webhook = &w
				break
			}
		}

		if webhook == nil {
			return baasapi.ErrObjectNotFound
		}

		return nil
	})

	return webhook, err
}

// WebhookByToken returns a webhook by the random token it is associated with.
func (service *Service) WebhookByToken(token string) (*baasapi.Webhook, error) {
	var webhook *baasapi.Webhook

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var w baasapi.Webhook
			err := internal.UnmarshalObject(v, &w)
			if err != nil {
				return err
			}

			if w.Token == token {
				webhook = &w
				break
			}
		}

		if webhook == nil {
			return baasapi.ErrObjectNotFound
		}

		return nil
	})

	return webhook, err
}

// DeleteWebhook deletes a webhook.
func (service *Service) DeleteWebhook(ID baasapi.WebhookID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}

// CreateWebhook assign an ID to a new webhook and saves it.
func (service *Service) CreateWebhook(webhook *baasapi.Webhook) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		id, _ := bucket.NextSequence()
		webhook.ID = baasapi.WebhookID(id)

		data, err := internal.MarshalObject(webhook)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(webhook.ID)), data)
	})
}
