package baask8ss

import (
	//"log"
	"net/http"
	//"runtime"
	//"strconv"
	//"math/rand"
	//"time"
	//"reflect"
	//"github.com/asaskevich/govalidator"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	//"github.com/baasapi/baasapi/api"
	//"github.com/baasapi/baasapi/api/crypto"
	//"github.com/baasapi/baasapi/api/http/client"
)

const (
	// Baas deployment files
	//BaaSDeploymentPath = "k8s/ansible/vars/namespaces"
)

//type logPayload struct {
//	Namespace   string
//	Nline       int
//}

type logResponse struct {
	Success      bool
	Logstr       string
}

//func (payload *logPayload) Validate(r *http.Request) error {

//	return nil
	
//}

// POST request on /api/baask8ss
func (handler *Handler) baask8sCreateLog(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	//if handler.authDisabled {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Cannot authenticate user. BaaSapi was started with the --no-auth flag", ErrAuthDisabled}
	//}

	//payload := &baask8sCreatePayload{}

	//var payload logPayload
	namespace, err := request.RetrieveRouteVariableValue(r, "namespace")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid namespace variable", err}
	}
	nline, err := request.RetrieveNumericRouteVariableValue(r, "nline")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid line number variable", err}
	}

	strlog, err := handler.CAFilesManager.GetLogs(namespace, nline)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to restrieve log  ", err}
	}

	var logresponse logResponse
	logresponse.Success = true
	logresponse.Logstr = strlog

	return response.JSON(w, logresponse)

	//return response.JSON(w, baask8s)
}


