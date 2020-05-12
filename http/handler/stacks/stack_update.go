package stacks

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/proxy"
	"github.com/baasapi/baasapi/api/http/security"
)

type updateComposeStackPayload struct {
	StackFileContent string
	Env              []baasapi.Pair
}

func (payload *updateComposeStackPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.StackFileContent) {
		return baasapi.Error("Invalid stack file content")
	}
	return nil
}

type updateSwarmStackPayload struct {
	StackFileContent string
	Env              []baasapi.Pair
	Prune            bool
}

func (payload *updateSwarmStackPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.StackFileContent) {
		return baasapi.Error("Invalid stack file content")
	}
	return nil
}

// PUT request on /api/stacks/:id?endpointId=<endpointId>
func (handler *Handler) stackUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	stackID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid stack identifier route variable", err}
	}

	stack, err := handler.StackService.Stack(baasapi.StackID(stackID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a stack with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a stack with the specified identifier inside the database", err}
	}

	resourceControl, err := handler.ResourceControlService.ResourceControlByResourceID(stack.Name)
	if err != nil && err != baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve a resource control associated to the stack", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	if !securityContext.IsAdmin {
		if !proxy.CanAccessStack(stack, resourceControl, securityContext.UserID, securityContext.UserMemberships) {
			return &httperror.HandlerError{http.StatusForbidden, "Access denied to resource", baasapi.ErrResourceAccessDenied}
		}
	}

	// TODO: this is a work-around for stacks created with BaaSapi version >= 1.17.1
	// The EndpointID property is not available for these stacks, this API endpoint
	// can use the optional EndpointID query parameter to associate a valid endpoint identifier to the stack.
	endpointID, err := request.RetrieveNumericQueryParameter(r, "endpointId", true)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: endpointId", err}
	}
	if endpointID != int(stack.EndpointID) {
		stack.EndpointID = baasapi.EndpointID(endpointID)
	}

	endpoint, err := handler.EndpointService.Endpoint(stack.EndpointID)
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find the endpoint associated to the stack inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find the endpoint associated to the stack inside the database", err}
	}

	updateError := handler.updateAndDeployStack(r, stack, endpoint)
	if updateError != nil {
		return updateError
	}

	err = handler.StackService.UpdateStack(stack.ID, stack)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the stack changes inside the database", err}
	}

	return response.JSON(w, stack)
}

func (handler *Handler) updateAndDeployStack(r *http.Request, stack *baasapi.Stack, endpoint *baasapi.Endpoint) *httperror.HandlerError {
	if stack.Type == baasapi.DockerSwarmStack {
		return handler.updateSwarmStack(r, stack, endpoint)
	}
	return handler.updateComposeStack(r, stack, endpoint)
}

func (handler *Handler) updateComposeStack(r *http.Request, stack *baasapi.Stack, endpoint *baasapi.Endpoint) *httperror.HandlerError {
	var payload updateComposeStackPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	stack.Env = payload.Env

	stackFolder := strconv.Itoa(int(stack.ID))
	_, err = handler.FileService.StoreStackFileFromBytes(stackFolder, stack.EntryPoint, []byte(payload.StackFileContent))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist updated Compose file on disk", err}
	}

	config, configErr := handler.createComposeDeployConfig(r, stack, endpoint)
	if configErr != nil {
		return configErr
	}

	err = handler.deployComposeStack(config)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, err.Error(), err}
	}

	return nil
}

func (handler *Handler) updateSwarmStack(r *http.Request, stack *baasapi.Stack, endpoint *baasapi.Endpoint) *httperror.HandlerError {
	var payload updateSwarmStackPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	stack.Env = payload.Env

	stackFolder := strconv.Itoa(int(stack.ID))
	_, err = handler.FileService.StoreStackFileFromBytes(stackFolder, stack.EntryPoint, []byte(payload.StackFileContent))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist updated Compose file on disk", err}
	}

	config, configErr := handler.createSwarmDeployConfig(r, stack, endpoint, payload.Prune)
	if configErr != nil {
		return configErr
	}

	err = handler.deploySwarmStack(config)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, err.Error(), err}
	}

	return nil
}
