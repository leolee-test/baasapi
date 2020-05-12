package baask8ss

import (
	"net/http"

	//"fmt"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// GET request on /api/baask8ss
func (handler *Handler) baask8sList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8ss, err := handler.Baask8sService.Baask8ss()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	}

	//fmt.Println(baask8ss[0].Namespace)

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	//fmt.Println(securityContext)

	//securityContext, err := security.RetrieveRestrictedRequestContext(r)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	//}

	filteredBaask8ss := security.FilterBaask8ss(baask8ss, securityContext)

	//for idx := range filteredBaask8ss {
	//	hideFields(&filteredBaask8ss[idx])
	//}

	return response.JSON(w, filteredBaask8ss)
}

// GET request on /api/baask8ss/{id}
func (handler *Handler) baask8sListByID(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	//log.Println("retrieving baask8s ID: " + baask8sID)

	//var payload baask8sUpdateAccessPayload
	//err = request.DecodeAndValidateJSONPayload(r, &payload)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	//}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

	//securityContext, err := security.RetrieveRestrictedRequestContext(r)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	//}

	//filteredBaask8ss := security.FilterBaask8ss(baask8ss, securityContext)

	//for idx := range filteredBaask8ss {
	//	hideFields(&filteredBaask8ss[idx])
	//}

	return response.JSON(w, baask8s)
}

