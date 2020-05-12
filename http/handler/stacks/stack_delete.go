package stacks

import (
	"net/http"
	"strconv"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/proxy"
	"github.com/baasapi/baasapi/api/http/security"
)

// DELETE request on /api/stacks/:id?external=<external>&endpointId=<endpointId>
// If the external query parameter is set to true, the id route variable is expected to be
// the name of an external stack as a string.
func (handler *Handler) stackDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	stackID, err := request.RetrieveRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid stack identifier route variable", err}
	}

	externalStack, _ := request.RetrieveBooleanQueryParameter(r, "external", true)
	if externalStack {
		return handler.deleteExternalStack(r, w, stackID)
	}

	id, err := strconv.Atoi(stackID)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid stack identifier route variable", err}
	}

	stack, err := handler.StackService.Stack(baasapi.StackID(id))
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
	// can use the optional EndpointID query parameter to set a valid endpoint identifier to be
	// used in the context of this request.
	endpointID, err := request.RetrieveNumericQueryParameter(r, "endpointId", true)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: endpointId", err}
	}
	endpointIdentifier := stack.EndpointID
	if endpointID != 0 {
		endpointIdentifier = baasapi.EndpointID(endpointID)
	}

	endpoint, err := handler.EndpointService.Endpoint(endpointIdentifier)
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find the endpoint associated to the stack inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find the endpoint associated to the stack inside the database", err}
	}

	err = handler.deleteStack(stack, endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, err.Error(), err}
	}

	err = handler.StackService.DeleteStack(baasapi.StackID(id))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove the stack from the database", err}
	}

	err = handler.FileService.RemoveDirectory(stack.ProjectPath)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove stack files from disk", err}
	}

	return response.Empty(w)
}

func (handler *Handler) deleteExternalStack(r *http.Request, w http.ResponseWriter, stackName string) *httperror.HandlerError {
	stack, err := handler.StackService.StackByName(stackName)
	if err != nil && err != baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to check for stack existence inside the database", err}
	}
	if stack != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "A stack with this name exists inside the database. Cannot use external delete method", baasapi.ErrStackNotExternal}
	}

	endpointID, err := request.RetrieveNumericQueryParameter(r, "endpointId", false)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: endpointId", err}
	}

	endpoint, err := handler.EndpointService.Endpoint(baasapi.EndpointID(endpointID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find the endpoint associated to the stack inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find the endpoint associated to the stack inside the database", err}
	}

	err = handler.requestBouncer.EndpointAccess(r, endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to access endpoint", baasapi.ErrEndpointAccessDenied}
	}

	stack = &baasapi.Stack{
		Name: stackName,
		Type: baasapi.DockerSwarmStack,
	}

	err = handler.deleteStack(stack, endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to delete stack", err}
	}

	return response.Empty(w)
}

func (handler *Handler) deleteStack(stack *baasapi.Stack, endpoint *baasapi.Endpoint) error {
	if stack.Type == baasapi.DockerSwarmStack {
		return handler.SwarmStackManager.Remove(stack, endpoint)
	}
	return handler.ComposeStackManager.Down(stack, endpoint)
}
