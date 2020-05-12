package stacks

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/proxy"
	"github.com/baasapi/baasapi/api/http/security"
)

// GET request on /api/stacks/:id
func (handler *Handler) stackInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	extendedStack := proxy.ExtendedStack{*stack, baasapi.ResourceControl{}}
	if !securityContext.IsAdmin && resourceControl == nil {
		return &httperror.HandlerError{http.StatusForbidden, "Access denied to resource", baasapi.ErrResourceAccessDenied}
	}

	if resourceControl != nil {
		if securityContext.IsAdmin || proxy.CanAccessStack(stack, resourceControl, securityContext.UserID, securityContext.UserMemberships) {
			extendedStack.ResourceControl = *resourceControl
		} else {
			return &httperror.HandlerError{http.StatusForbidden, "Access denied to resource", baasapi.ErrResourceAccessDenied}
		}
	}

	return response.JSON(w, extendedStack)
}
