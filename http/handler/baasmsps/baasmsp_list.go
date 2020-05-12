package baasmsps

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/response"
	//"github.com/baasapi/baasapi/api/http/security"
)

// GET request on /api/baasmsps
func (handler *Handler) baasmspList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baasmsps, err := handler.BaasmspService.Baasmsps()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baasmsps from the database", err}
	}

	//securityContext, err := security.RetrieveRestrictedRequestContext(r)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	//}

	//filteredBaasmsps := security.FilterBaasmsps(baasmsps, securityContext)

	//for idx := range filteredBaasmsps {
	//	hideFields(&filteredBaasmsps[idx])
	//}

	return response.JSON(w, baasmsps)
}
