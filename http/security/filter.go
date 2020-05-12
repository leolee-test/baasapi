package security

import "github.com/baasapi/baasapi/api"

// FilterUserTeams filters teams based on user role.
// non-administrator users only have access to team they are member of.
func FilterUserTeams(teams []baasapi.Team, context *RestrictedRequestContext) []baasapi.Team {
	filteredTeams := teams

	if !context.IsAdmin {
		filteredTeams = make([]baasapi.Team, 0)
		for _, membership := range context.UserMemberships {
			for _, team := range teams {
				if team.ID == membership.TeamID {
					filteredTeams = append(filteredTeams, team)
					break
				}
			}
		}
	}

	return filteredTeams
}

// FilterLeaderTeams filters teams based on user role.
// Team leaders only have access to team they lead.
func FilterLeaderTeams(teams []baasapi.Team, context *RestrictedRequestContext) []baasapi.Team {
	filteredTeams := teams

	if context.IsTeamLeader {
		filteredTeams = make([]baasapi.Team, 0)
		for _, membership := range context.UserMemberships {
			for _, team := range teams {
				if team.ID == membership.TeamID && membership.Role == baasapi.TeamLeader {
					filteredTeams = append(filteredTeams, team)
					break
				}
			}
		}
	}

	return filteredTeams
}

// FilterUsers filters users based on user role.
// Non-administrator users only have access to non-administrator users.
func FilterUsers(users []baasapi.User, context *RestrictedRequestContext) []baasapi.User {
	filteredUsers := users

	if !context.IsAdmin {
		filteredUsers = make([]baasapi.User, 0)

		for _, user := range users {
			if user.Role != baasapi.AdministratorRole {
				filteredUsers = append(filteredUsers, user)
			}
		}
	}

	return filteredUsers
}

// FilterRegistries filters registries based on user role and team memberships.
// Non administrator users only have access to authorized registries.
func FilterRegistries(registries []baasapi.Registry, context *RestrictedRequestContext) []baasapi.Registry {
	filteredRegistries := registries
	if !context.IsAdmin {
		filteredRegistries = make([]baasapi.Registry, 0)

		for _, registry := range registries {
			if AuthorizedRegistryAccess(&registry, context.UserID, context.UserMemberships) {
				filteredRegistries = append(filteredRegistries, registry)
			}
		}
	}

	return filteredRegistries
}

// FilterTemplates filters templates based on the user role.
// Non-administrato template do not have access to templates where the AdministratorOnly flag is set to true.
func FilterTemplates(templates []baasapi.Template, context *RestrictedRequestContext) []baasapi.Template {
	filteredTemplates := templates

	if !context.IsAdmin {
		filteredTemplates = make([]baasapi.Template, 0)

		for _, template := range templates {
			if !template.AdministratorOnly {
				filteredTemplates = append(filteredTemplates, template)
			}
		}
	}

	return filteredTemplates
}

// FilterBaask8ss filters baask8ss based on user role and team memberships.
// Non administrator users only have access to authorized baask8ss (can be inherited via endoint groups).
//func FilterBaask8ss(baask8ss []baasapi.Baask8s, groups []baasapi.Baask8sGroup, context *RestrictedRequestContext) []baasapi.Baask8s {
//	filteredBaask8ss := baask8ss

//	if !context.IsAdmin {
//		filteredBaask8ss = make([]baasapi.Baask8s, 0)
//
//		for _, baask8s := range baask8ss {
//			baask8sGroup := getAssociatedGroup(&baask8s, groups)
//
//			if authorizedBaask8sAccess(&baask8s, baask8sGroup, context.UserID, context.UserMemberships) {
//				filteredBaask8ss = append(filteredBaask8ss, baask8s)
//			}
//		}
//	}
//
//	return filteredBaask8ss
//}

// FilterBaask8ss filters baask8s based on user role and team memberships.
// Non administrator users only have access to authorized baask8ss (can be inherited via endoint groups).
func FilterBaask8ss(baask8ss []baasapi.Baask8s, context *RestrictedRequestContext) []baasapi.Baask8s {
	filteredBaask8ss := baask8ss

	if !context.IsAdmin {
		filteredBaask8ss = make([]baasapi.Baask8s, 0)

		for _, baask8s := range baask8ss {
		//	baask8sGroup := getAssociatedGroup(&baask8s, groups)

			if authorizedBaask8sAccess(&baask8s, context.UserID, context.UserMemberships) {
				filteredBaask8ss = append(filteredBaask8ss, baask8s)
			}
		}
	}
	//if !context.IsAdmin {
	//	filteredBaask8ss = make([]baasapi.Baask8s, 0)

	//	for _, baask8s := range baask8ss {
	//		if baask8s.Owner != baasapi.AdministratorRole {
	//			filteredUsers = append(filteredUsers, user)
	//		}
	//	}
	//}

	return filteredBaask8ss
}



// FilterBaask8sGroups filters baask8s groups based on user role and team memberships.
// Non administrator users only have access to authorized baask8s groups.
func FilterBaask8sGroups(baask8sGroups []baasapi.Baask8sGroup, context *RestrictedRequestContext) []baasapi.Baask8sGroup {
	filteredBaask8sGroups := baask8sGroups

	if !context.IsAdmin {
		filteredBaask8sGroups = make([]baasapi.Baask8sGroup, 0)

		for _, group := range baask8sGroups {
			if authorizedBaask8sGroupAccess(&group, context.UserID, context.UserMemberships) {
				filteredBaask8sGroups = append(filteredBaask8sGroups, group)
			}
		}
	}

	return filteredBaask8sGroups
}

func getAssociatedGroup(baask8s *baasapi.Baask8s, groups []baasapi.Baask8sGroup) *baasapi.Baask8sGroup {
	//for _, group := range groups {
	//	if group.ID == baask8s.GroupID {
	//		return &group
	//	}
	//}
	return nil
}
