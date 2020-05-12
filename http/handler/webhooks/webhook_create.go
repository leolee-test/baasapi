package webhooks

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/satori/go.uuid"
)

type webhookCreatePayload struct {
	ResourceID  string
	Baask8sID  int
	WebhookType int
}

func (payload *webhookCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.ResourceID) {
		return baasapi.Error("Invalid ResourceID")
	}
	if payload.Baask8sID == 0 {
		return baasapi.Error("Invalid Baask8sID")
	}
	if payload.WebhookType != 1 {
		return baasapi.Error("Invalid WebhookType")
	}
	return nil
}

func (handler *Handler) webhookCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload webhookCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	webhook, err := handler.WebhookService.WebhookByResourceID(payload.ResourceID)
	if err != nil && err != baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "An error occurred retrieving webhooks from the database", err}
	}
	if webhook != nil {
		return &httperror.HandlerError{http.StatusConflict, "A webhook for this resource already exists", baasapi.ErrWebhookAlreadyExists}
	}

	token, err := uuid.NewV4()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Error creating unique token", err}
	}

	webhook = &baasapi.Webhook{
		Token:       token.String(),
		ResourceID:  payload.ResourceID,
		Baask8sID:  baasapi.Baask8sID(payload.Baask8sID),
		WebhookType: baasapi.WebhookType(payload.WebhookType),
	}

	err = handler.WebhookService.CreateWebhook(webhook)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the webhook inside the database", err}
	}

	return response.JSON(w, webhook)
}
