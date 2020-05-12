package proxy

import (
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

func createRegistryAuthenticationHeader(serverAddress string, accessContext *registryAccessContext) *registryAuthenticationHeader {
	var authenticationHeader *registryAuthenticationHeader

	if serverAddress == "" {
		authenticationHeader = &registryAuthenticationHeader{
			Username:      accessContext.dockerHub.Username,
			Password:      accessContext.dockerHub.Password,
			Serveraddress: "docker.io",
		}
	} else {
		var matchingRegistry *baasapi.Registry
		for _, registry := range accessContext.registries {
			if registry.URL == serverAddress &&
				(accessContext.isAdmin || (!accessContext.isAdmin && security.AuthorizedRegistryAccess(&registry, accessContext.userID, accessContext.teamMemberships))) {
				matchingRegistry = &registry
				break
			}
		}

		if matchingRegistry != nil {
			authenticationHeader = &registryAuthenticationHeader{
				Username:      matchingRegistry.Username,
				Password:      matchingRegistry.Password,
				Serveraddress: matchingRegistry.URL,
			}
		}
	}

	return authenticationHeader
}
