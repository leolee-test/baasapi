package baask8ss

import (
	"net/http"
	//"log"
	//"strconv"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// DELETE request on /api/baask8ss/:id
func (handler *Handler) baask8sDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//if !handler.authorizeBaask8sManagement {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Baask8s management is disabled", ErrBaask8sManagementDisabled}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

	//if baask8s.TLSConfig.TLS {
	//	folder := strconv.Itoa(baask8sID)
	//	err = handler.FileService.DeleteTLSFiles(folder)
	//	if err != nil {
	//		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove TLS files from disk", err}
	//	}
	//}

	var ansible_env = "env=" +baask8s.Namespace+ " deploy_type=k8s"
	var ansible_extra = " mode=destroy "
	var ansible_config = "/data/k8s/ansible/setupfabric.yml"
	err = handler.CAFilesManager.Deploy(baask8s.Owner, baask8s.Namespace, ansible_extra, ansible_env, ansible_config, true)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute ansible commands ", err}
	}


	err = handler.Baask8sService.DeleteBaask8s(baasapi.Baask8sID(baask8sID))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove baask8s from the database", err}
	}

	//handler.ProxyManager.DeleteProxy(string(baask8sID))

	return response.Empty(w)
}

// DELETE request on /api/baask8ss/baasonly/:id
func (handler *Handler) baask8sDeleteBaasOnly(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//if !handler.authorizeBaask8sManagement {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Baask8s management is disabled", ErrBaask8sManagementDisabled}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}


	err = handler.Baask8sService.DeleteBaask8s(baasapi.Baask8sID(baask8sID))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove baask8s from the database", err}
	}

	//handler.ProxyManager.DeleteProxy(string(baask8sID))

	return response.Empty(w)
}
