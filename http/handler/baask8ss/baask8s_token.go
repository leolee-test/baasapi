package baask8ss

import (
	"net/http"

	//"fmt"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	//"github.com/baasapi/baasapi/api"
	//"github.com/baasapi/baasapi/api/http/security"
)

//type tokenPayload struct {
//	Token     string     `json:"token"`
//}


//func (payload *tokenPayload) Validate(r *http.Request) error {
//	return nil;
//}

// GET request on /api/baask8ss/{id}/token
func (handler *Handler) baask8sTokenValidate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	token, err := request.RetrieveRouteVariableValue(r, "token")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid token name variable", err}
	}

	//var err error
	TokenData, err := handler.JWTService.ParseAndVerifyToken(token)
	if err != nil {
		//httperror.WriteError(w, http.StatusUnauthorized, "Invalid JWT token", err)
		return &httperror.HandlerError{http.StatusNotFound, "Invalid JWT token", err}
		//return
	}



	//securityContext, err := security.RetrieveRestrictedRequestContext(r)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	//}

	//filteredBaask8ss := security.FilterBaask8ss(baask8ss, securityContext)

	//for idx := range filteredBaask8ss {
	//	hideFields(&filteredBaask8ss[idx])
	//}

	return response.JSON(w, TokenData)
}
