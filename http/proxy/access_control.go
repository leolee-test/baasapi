package proxy

import (
	"github.com/baasapi/baasapi/api"
)

type (
	// ExtendedStack represents a stack combined with its associated access control
	ExtendedStack struct {
		baasapi.Stack
		ResourceControl baasapi.ResourceControl `json:"ResourceControl"`
	}
)

// applyResourceAccessControlFromLabel returns an optionally decorated object as the first return value and the
// access level for the user (granted or denied) as the second return value.
// It will retrieve an identifier from the labels object. If an identifier exists, it will check for
// an existing resource control associated to it.
// Returns a decorated object and authorized access (true) when a resource control is found and the user can access the resource.
// Returns the original object and denied access (false) when no resource control is found.
// Returns the original object and denied access (false) when a resource control is found and the user cannot access the resource.
func applyResourceAccessControlFromLabel(labelsObject, resourceObject map[string]interface{}, labelIdentifier string,
	context *restrictedOperationContext) (map[string]interface{}, bool) {

	if labelsObject != nil && labelsObject[labelIdentifier] != nil {
		resourceIdentifier := labelsObject[labelIdentifier].(string)
		return applyResourceAccessControl(resourceObject, resourceIdentifier, context)
	}
	return resourceObject, false
}

// applyResourceAccessControl returns an optionally decorated object as the first return value and the
// access level for the user (granted or denied) as the second return value.
// Returns a decorated object and authorized access (true) when a resource control is found to the specified resource
// identifier and the user can access the resource.
// Returns the original object and authorized access (false) when no resource control is found for the specified
// resource identifier.
// Returns the original object and denied access (false) when a resource control is associated to the resource
// and the user cannot access the resource.
func applyResourceAccessControl(resourceObject map[string]interface{}, resourceIdentifier string,
	context *restrictedOperationContext) (map[string]interface{}, bool) {

	resourceControl := getResourceControlByResourceID(resourceIdentifier, context.resourceControls)
	if resourceControl == nil {
		return resourceObject, context.isAdmin
	}

	if context.isAdmin || resourceControl.Public || canUserAccessResource(context.userID, context.userTeamIDs, resourceControl) {
		resourceObject = decorateObject(resourceObject, resourceControl)
		return resourceObject, true
	}

	return resourceObject, false
}

// decorateResourceWithAccessControlFromLabel will retrieve an identifier from the labels object. If an identifier exists,
// it will check for an existing resource control associated to it. If a resource control is found, the resource object will be
// decorated. If no identifier can be found in the labels or no resource control is associated to the identifier, the resource
// object will not be changed.
func decorateResourceWithAccessControlFromLabel(labelsObject, resourceObject map[string]interface{}, labelIdentifier string,
	resourceControls []baasapi.ResourceControl) map[string]interface{} {

	if labelsObject != nil && labelsObject[labelIdentifier] != nil {
		resourceIdentifier := labelsObject[labelIdentifier].(string)
		resourceObject = decorateResourceWithAccessControl(resourceObject, resourceIdentifier, resourceControls)
	}

	return resourceObject
}

// decorateResourceWithAccessControl will check if a resource control is associated to the specified resource identifier.
// If a resource control is found, the resource object will be decorated, otherwise it will not be changed.
func decorateResourceWithAccessControl(resourceObject map[string]interface{}, resourceIdentifier string,
	resourceControls []baasapi.ResourceControl) map[string]interface{} {

	resourceControl := getResourceControlByResourceID(resourceIdentifier, resourceControls)
	if resourceControl != nil {
		return decorateObject(resourceObject, resourceControl)
	}
	return resourceObject
}

func canUserAccessResource(userID baasapi.UserID, userTeamIDs []baasapi.TeamID, resourceControl *baasapi.ResourceControl) bool {
	for _, authorizedUserAccess := range resourceControl.UserAccesses {
		if userID == authorizedUserAccess.UserID {
			return true
		}
	}

	for _, authorizedTeamAccess := range resourceControl.TeamAccesses {
		for _, userTeamID := range userTeamIDs {
			if userTeamID == authorizedTeamAccess.TeamID {
				return true
			}
		}
	}

	return resourceControl.Public
}

func decorateObject(object map[string]interface{}, resourceControl *baasapi.ResourceControl) map[string]interface{} {
	if object["BaaSapi"] == nil {
		object["BaaSapi"] = make(map[string]interface{})
	}

	baasapiMetadata := object["BaaSapi"].(map[string]interface{})
	baasapiMetadata["ResourceControl"] = resourceControl
	return object
}

func getResourceControlByResourceID(resourceID string, resourceControls []baasapi.ResourceControl) *baasapi.ResourceControl {
	for _, resourceControl := range resourceControls {
		if resourceID == resourceControl.ResourceID {
			return &resourceControl
		}
		for _, subResourceID := range resourceControl.SubResourceIDs {
			if resourceID == subResourceID {
				return &resourceControl
			}
		}
	}
	return nil
}

// CanAccessStack checks if a user can access a stack
func CanAccessStack(stack *baasapi.Stack, resourceControl *baasapi.ResourceControl, userID baasapi.UserID, memberships []baasapi.TeamMembership) bool {
	if resourceControl == nil {
		return false
	}

	userTeamIDs := make([]baasapi.TeamID, 0)
	for _, membership := range memberships {
		userTeamIDs = append(userTeamIDs, membership.TeamID)
	}

	if canUserAccessResource(userID, userTeamIDs, resourceControl) {
		return true
	}

	return resourceControl.Public
}

// FilterStacks filters stacks based on user role and resource controls.
func FilterStacks(stacks []baasapi.Stack, resourceControls []baasapi.ResourceControl, isAdmin bool,
	userID baasapi.UserID, memberships []baasapi.TeamMembership) []ExtendedStack {

	filteredStacks := make([]ExtendedStack, 0)

	userTeamIDs := make([]baasapi.TeamID, 0)
	for _, membership := range memberships {
		userTeamIDs = append(userTeamIDs, membership.TeamID)
	}

	for _, stack := range stacks {
		extendedStack := ExtendedStack{stack, baasapi.ResourceControl{}}
		resourceControl := getResourceControlByResourceID(stack.Name, resourceControls)
		if resourceControl == nil && isAdmin {
			filteredStacks = append(filteredStacks, extendedStack)
		} else if resourceControl != nil && (isAdmin || resourceControl.Public || canUserAccessResource(userID, userTeamIDs, resourceControl)) {
			extendedStack.ResourceControl = *resourceControl
			filteredStacks = append(filteredStacks, extendedStack)
		}
	}

	return filteredStacks
}
