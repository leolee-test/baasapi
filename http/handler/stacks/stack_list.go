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

type stackListOperationFilters struct {
	SwarmID    string `json:"SwarmID"`
	EndpointID int    `json:"EndpointID"`
}

// GET request on /api/stacks?(filters=<filters>)
func (handler *Handler) stackList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var filters stackListOperationFilters
	err := request.RetrieveJSONQueryParameter(r, "filters", &filters, true)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: filters", err}
	}

	stacks, err := handler.StackService.Stacks()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve stacks from the database", err}
	}
	stacks = filterStacks(stacks, &filters)

	resourceControls, err := handler.ResourceControlService.ResourceControls()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve resource controls from the database", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	filteredStacks := proxy.FilterStacks(stacks, resourceControls, securityContext.IsAdmin,
		securityContext.UserID, securityContext.UserMemberships)

	return response.JSON(w, filteredStacks)
}

func filterStacks(stacks []baasapi.Stack, filters *stackListOperationFilters) []baasapi.Stack {
	if filters.EndpointID == 0 && filters.SwarmID == "" {
		return stacks
	}

	filteredStacks := make([]baasapi.Stack, 0, len(stacks))
	for _, stack := range stacks {
		if stack.Type == baasapi.DockerComposeStack && stack.EndpointID == baasapi.EndpointID(filters.EndpointID) {
			filteredStacks = append(filteredStacks, stack)
		}
		if stack.Type == baasapi.DockerSwarmStack && stack.SwarmID == filters.SwarmID {
			filteredStacks = append(filteredStacks, stack)
		}
	}

	return filteredStacks
}
