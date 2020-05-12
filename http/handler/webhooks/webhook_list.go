package webhooks

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type webhookListOperationFilters struct {
	ResourceID string `json:"ResourceID"`
	Baask8sID int    `json:"Baask8sID"`
}

// GET request on /api/webhooks?(filters=<filters>)
func (handler *Handler) webhookList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var filters webhookListOperationFilters
	err := request.RetrieveJSONQueryParameter(r, "filters", &filters, true)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: filters", err}
	}

	webhooks, err := handler.WebhookService.Webhooks()
	webhooks = filterWebhooks(webhooks, &filters)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve webhooks from the database", err}
	}

	return response.JSON(w, webhooks)
}

func filterWebhooks(webhooks []baasapi.Webhook, filters *webhookListOperationFilters) []baasapi.Webhook {
	if filters.Baask8sID == 0 && filters.ResourceID == "" {
		return webhooks
	}

	filteredWebhooks := make([]baasapi.Webhook, 0, len(webhooks))
	for _, webhook := range webhooks {
		if webhook.Baask8sID == baasapi.Baask8sID(filters.Baask8sID) && webhook.ResourceID == string(filters.ResourceID) {
			filteredWebhooks = append(filteredWebhooks, webhook)
		}
	}

	return filteredWebhooks
}
