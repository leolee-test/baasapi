package security

import (
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"

	"net/http"
	"strings"
)

type (
	// RequestBouncer represents an entity that manages API request accesses
	RequestBouncer struct {
		jwtService            baasapi.JWTService
		userService           baasapi.UserService
		teamMembershipService baasapi.TeamMembershipService
		baask8sGroupService  baasapi.Baask8sGroupService
		authDisabled          bool
	}

	// RequestBouncerParams represents the required parameters to create a new RequestBouncer instance.
	RequestBouncerParams struct {
		JWTService            baasapi.JWTService
		UserService           baasapi.UserService
		TeamMembershipService baasapi.TeamMembershipService
		Baask8sGroupService  baasapi.Baask8sGroupService
		AuthDisabled          bool
	}

	// RestrictedRequestContext is a data structure containing information
	// used in RestrictedAccess
	RestrictedRequestContext struct {
		IsAdmin         bool
		IsTeamLeader    bool
		UserID          baasapi.UserID
		UserMemberships []baasapi.TeamMembership
	}
)

// NewRequestBouncer initializes a new RequestBouncer
func NewRequestBouncer(parameters *RequestBouncerParams) *RequestBouncer {
	return &RequestBouncer{
		jwtService:            parameters.JWTService,
		userService:           parameters.UserService,
		teamMembershipService: parameters.TeamMembershipService,
		baask8sGroupService:  parameters.Baask8sGroupService,
		authDisabled:          parameters.AuthDisabled,
	}
}

// PublicAccess defines a security check for public baask8ss.
// No authentication is required to access these baask8ss.
func (bouncer *RequestBouncer) PublicAccess(h http.Handler) http.Handler {
	h = mwSecureHeaders(h)
	return h
}

// AuthenticatedAccess defines a security check for private baask8ss.
// Authentication is required to access these baask8ss.
func (bouncer *RequestBouncer) AuthenticatedAccess(h http.Handler) http.Handler {
	h = bouncer.mwCheckAuthentication(h)
	h = mwSecureHeaders(h)
	return h
}

// RestrictedAccess defines a security check for restricted baask8ss.
// Authentication is required to access these baask8ss.
// The request context will be enhanced with a RestrictedRequestContext object
// that might be used later to authorize/filter access to resources.
func (bouncer *RequestBouncer) RestrictedAccess(h http.Handler) http.Handler {
	h = bouncer.mwUpgradeToRestrictedRequest(h)
	h = bouncer.AuthenticatedAccess(h)
	return h
}

// AdministratorAccess defines a chain of middleware for restricted baask8ss.
// Authentication as well as administrator role are required to access these baask8ss.
func (bouncer *RequestBouncer) AdministratorAccess(h http.Handler) http.Handler {
	h = mwCheckAdministratorRole(h)
	h = bouncer.AuthenticatedAccess(h)
	return h
}

// Baask8sAccess retrieves the JWT token from the request context and verifies
// that the user can access the specified baask8s.
// An error is returned when access is denied.
func (bouncer *RequestBouncer) Baask8sAccess(r *http.Request, baask8s *baasapi.Baask8s) error {
	tokenData, err := RetrieveTokenData(r)
	if err != nil {
		return err
	}

	if tokenData.Role == baasapi.AdministratorRole {
		return nil
	}

	memberships, err := bouncer.teamMembershipService.TeamMembershipsByUserID(tokenData.ID)
	if err != nil {
		return err
	}

	//group, err := bouncer.baask8sGroupService.Baask8sGroup(baask8s.GroupID)
	//if err != nil {
	//	return err
	//}

	//if !authorizedBaask8sAccess(baask8s, group, tokenData.ID, memberships) {
	if !authorizedBaask8sAccess(baask8s, tokenData.ID, memberships) {
		return baasapi.ErrBaask8sAccessDenied
	}

	return nil
}

// RegistryAccess retrieves the JWT token from the request context and verifies
// that the user can access the specified registry.
// An error is returned when access is denied.
func (bouncer *RequestBouncer) RegistryAccess(r *http.Request, registry *baasapi.Registry) error {
	tokenData, err := RetrieveTokenData(r)
	if err != nil {
		return err
	}

	if tokenData.Role == baasapi.AdministratorRole {
		return nil
	}

	memberships, err := bouncer.teamMembershipService.TeamMembershipsByUserID(tokenData.ID)
	if err != nil {
		return err
	}

	if !AuthorizedRegistryAccess(registry, tokenData.ID, memberships) {
		return baasapi.ErrBaask8sAccessDenied
	}

	return nil
}

// mwSecureHeaders provides secure headers middleware for handlers.
func mwSecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

// mwUpgradeToRestrictedRequest will enhance the current request with
// a new RestrictedRequestContext object.
func (bouncer *RequestBouncer) mwUpgradeToRestrictedRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenData, err := RetrieveTokenData(r)
		if err != nil {
			httperror.WriteError(w, http.StatusForbidden, "Access denied", baasapi.ErrResourceAccessDenied)
			return
		}

		requestContext, err := bouncer.newRestrictedContextRequest(tokenData.ID, tokenData.Role)
		if err != nil {
			httperror.WriteError(w, http.StatusInternalServerError, "Unable to create restricted request context ", err)
			return
		}

		ctx := storeRestrictedRequestContext(r, requestContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// mwCheckAdministratorRole check the role of the user associated to the request
func mwCheckAdministratorRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenData, err := RetrieveTokenData(r)
		if err != nil || tokenData.Role != baasapi.AdministratorRole {
			httperror.WriteError(w, http.StatusForbidden, "Access denied", baasapi.ErrResourceAccessDenied)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// mwCheckAuthentication provides Authentication middleware for handlers
func (bouncer *RequestBouncer) mwCheckAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenData *baasapi.TokenData
		if !bouncer.authDisabled {
			var token string

			// Optionally, token might be set via the "token" query parameter.
			// For example, in websocket requests
			token = r.URL.Query().Get("token")

			// Get token from the Authorization header
			tokens, ok := r.Header["Authorization"]
			if ok && len(tokens) >= 1 {
				token = tokens[0]
				token = strings.TrimPrefix(token, "Bearer ")
			}

			if token == "" {
				httperror.WriteError(w, http.StatusUnauthorized, "Unauthorized", baasapi.ErrUnauthorized)
				return
			}

			var err error
			tokenData, err = bouncer.jwtService.ParseAndVerifyToken(token)
			if err != nil {
				httperror.WriteError(w, http.StatusUnauthorized, "Invalid JWT token", err)
				return
			}

			_, err = bouncer.userService.User(tokenData.ID)
			if err != nil && err == baasapi.ErrObjectNotFound {
				httperror.WriteError(w, http.StatusUnauthorized, "Unauthorized", baasapi.ErrUnauthorized)
				return
			} else if err != nil {
				httperror.WriteError(w, http.StatusInternalServerError, "Unable to retrieve users from the database", err)
				return
			}
		} else {
			tokenData = &baasapi.TokenData{
				Role: baasapi.AdministratorRole,
			}
		}

		ctx := storeTokenData(r, tokenData)
		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

func (bouncer *RequestBouncer) newRestrictedContextRequest(userID baasapi.UserID, userRole baasapi.UserRole) (*RestrictedRequestContext, error) {
	requestContext := &RestrictedRequestContext{
		IsAdmin: true,
		UserID:  userID,
	}

	if userRole != baasapi.AdministratorRole {
		requestContext.IsAdmin = false
		memberships, err := bouncer.teamMembershipService.TeamMembershipsByUserID(userID)
		if err != nil {
			return nil, err
		}

		isTeamLeader := false
		for _, membership := range memberships {
			if membership.Role == baasapi.TeamLeader {
				isTeamLeader = true
			}
		}

		requestContext.IsTeamLeader = isTeamLeader
		requestContext.UserMemberships = memberships
	}

	return requestContext, nil
}
