package baasmsps

import (
	"net/http"
	"log"
	//"strconv"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// DELETE request on /api/baasmsps/:id
func (handler *Handler) baasmspDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//if !handler.authorizeBaasmspManagement {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Baasmsp management is disabled", ErrBaasmspManagementDisabled}
	//}

	baasmspID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baasmsp identifier route variable", err}
	}

	baasmsp, err := handler.BaasmspService.Baasmsp(baasapi.BaasmspID(baasmspID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baasmsp with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baasmsp with the specified identifier inside the database", err}
	}

	//if baasmsp.TLSConfig.TLS {
	//	folder := strconv.Itoa(baasmspID)
	//	err = handler.FileService.DeleteTLSFiles(folder)
	//	if err != nil {
	//		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove TLS files from disk", err}
	//	}
	//}
	log.Printf("(baask8s=%s) (Namespace to be deleted too =) \n", baasmsp.Namespace)

	err = handler.BaasmspService.DeleteBaasmsp(baasapi.BaasmspID(baasmspID))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove baasmsp from the database", err}
	}

	//handler.ProxyManager.DeleteProxy(string(baasmspID))

	return response.Empty(w)
}
