package webhooks

import (
	//"context"
	"net/http"
	//"strings"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// Acts on a passed in token UUID to restart the docker service
func (handler *Handler) webhookExecute(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	webhookToken, err := request.RetrieveRouteVariableValue(r, "token")

	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Invalid service id parameter", err}
	}

	webhook, err := handler.WebhookService.WebhookByToken(webhookToken)

	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a webhook with this token", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve webhook from the database", err}
	}

	resourceID := webhook.ResourceID
	baask8sID := webhook.Baask8sID
	webhookType := webhook.WebhookType

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}
	switch webhookType {
	case baasapi.ServiceWebhook:
		return handler.executeServiceWebhook(w, baask8s, resourceID)
	default:
		return &httperror.HandlerError{http.StatusInternalServerError, "Unsupported webhook type", baasapi.ErrUnsupportedWebhookType}
	}

}

func (handler *Handler) executeServiceWebhook(w http.ResponseWriter, baask8s *baasapi.Baask8s, resourceID string) *httperror.HandlerError {
	//dockerClient, err := handler.DockerClientFactory.CreateClient(baask8s, "")
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Error creating docker client", err}
	//}
	//defer dockerClient.Close()

	//service, _, err := dockerClient.ServiceInspectWithRaw(context.Background(), resourceID, dockertypes.ServiceInspectOptions{InsertDefaults: true})
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Error looking up service", err}
	//}

	//service.Spec.TaskTemplate.ForceUpdate++

	//service.Spec.TaskTemplate.ContainerSpec.Image = strings.Split(service.Spec.TaskTemplate.ContainerSpec.Image, "@sha")[0]
	//_, err = dockerClient.ServiceUpdate(context.Background(), resourceID, service.Version, service.Spec, dockertypes.ServiceUpdateOptions{QueryRegistry: true})
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Error updating service", err}
	//}
	return response.Empty(w)
}
