package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	"strings"
	//"log"
	//"os"
	//"time"
    //"bytes"
    "encoding/json"
	//"fmt"
	//"reflect"
    //"io/ioutil"
	//"path/filepath"

	//"k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/apimachinery/pkg/types"
	//"k8s.io/client-go/rest"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	//"github.com/baasapi/baasapi/api/http/security"
)


type jsondeployResponse struct {
    Success    bool    `json:"success"`
	Message    string  `json:"message"`
}


//type jsonCreateSVCResponse struct {
//    Success    bool    `json:"success"`
//	Message    string  `json:"message"`
//}

type SVCCreatePayload struct {
	Svcname                string     `json:"svcname"`
}

//type CHLPayload struct {
//	Allchannels            []baasapi.CHL
//	Currchannels           []baasapi.CHL
//}

	//{"success":true,
	//"secret":"",
	//"message":"Jim enrolled Successfully",
	//"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTU2OTg2NTcsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Ik9yZzEiLCJpYXQiOjE1NTU2NjI2NTd9.racjuDcqswHY2WS9gj4XLBBwW-ST_yb9dDTZAlbh33Q"
	//}  
	
	type jsonDataResponse struct {
		Success    bool    `json:"success"`
		Message    string  `json:"message"`
		Namespace  string  `json:"namespace"`
		Networkname string   `json:"networkname"`
		Data       []baasapi.BaasAPP       `json:"data"`
		//Token      string    `json:"token"`
		//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	}

func (payload *SVCCreatePayload) Validate(r *http.Request) error {
	return nil;
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sGetSvcs(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
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

	var responseObject jsonDataResponse
	responseObject.Success = true
	responseObject.Message = "Retrieved the Applications information successfully"
	responseObject.Namespace = baask8s.Namespace
	responseObject.Networkname = baask8s.NetworkName
	responseObject.Data      = baask8s.Applications
	//responseObject.Message = "Not authorized or jwt token was expired"
	//json.Unmarshal(bodyBytes, &responseObject)
	//return response.JSON(w, responseObject)

	return response.JSON(w, responseObject)
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sDeploySvcs(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	var payload SVCCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

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

	for i, _ := range baask8s.Applications {
		if ((payload.Svcname == baask8s.Applications[i].Name) && (baask8s.Applications[i].Status != "DELETED")) {

			var responseObject1 jsonResponse
			responseObject1.Success = false
			responseObject1.Message = "Service " + payload.Svcname + " is already deployed"
			return response.JSON(w, responseObject1)
			//return &httperror.HandlerError{http.StatusBadRequest, "Application is already deployed", err}
		}
	}

	mychns := CHLPayload{}
	
	mychns.Currchannels = baask8s.CHLs
	mychns.Allchannels = baask8s.CHLs

    if len(mychns.Currchannels) == 0 {
		return &httperror.HandlerError{http.StatusNotFound, "No channel found for this operations", baasapi.ErrObjectNotFound}
	//	channel.CHLName = "default"
	//	mychns.Currchannels = append(mychns.Currchannels, channel)
	}

    //response, err := http.Get("http://11.11.11.120:30500/users")
    //if err != nil {
    //    fmt.Printf("The HTTP request failed with error %s\n", err)
    //} else {
    //    data, _ := ioutil.ReadAll(response.Body)
    //    fmt.Println(string(data))
	//}
	//var data={};
    var jsonData []byte
	jsonData, err = json.Marshal(mychns)
    if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute ansible commands ", err}
	}


	var flag = false
	for i, _ := range baask8s.Applications {
		if ((payload.Svcname == baask8s.Applications[i].Name)) {
			flag = true;
			break;
		}
	}

	if (flag == false) {

		myapp := &baasapi.BaasAPP{
			Name:     payload.Svcname, 
			Status:   "DEPLOYING",
		}
		baask8s.Applications = append(baask8s.Applications, *myapp)
	
		err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
		}
	}
	
	//--extra-vars {"Allchannels":[{"Id":0,"CHLName":"tst4c","CreatedAt":1567312246,"ORGs":[]}]}

	var ansible_env = "mode="+ payload.Svcname+" env=" +baask8s.Namespace+ " deploy_type=k8s "
	var ansible_extra = string(jsonData)
	//ansible_env = ansible_env + " --extra-vars '{"allchannels":[{"Id":0,"CHLName":"tst4c","CreatedAt":1567312246,"ORGs":[]}]}'"
	//ansible_env = ansible_env + "\" --extra-vars '{\"allchannels\":[{\"Id\":0,\"CHLName\":\"tst4c\",\"CreatedAt\":1567312246,\"ORGs\":[]}]}'"
	var ansible_config = "/data/k8s/ansible/operatefabric.yml"
	err = handler.CAFilesManager.Deploy(baask8s.Owner, baask8s.Namespace, ansible_extra, ansible_env, ansible_config, true)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute ansible commands ", err}
	}


	for i := range baask8s.Applications {
		if (payload.Svcname == baask8s.Applications[i].Name) {
			baask8s.Applications[i].Status = "DEPLOYED"
		}
	}
	
	err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
	}
	
				var responseObject jsondeployResponse
				responseObject.Success = true
				responseObject.Message = "successfully deploy fabric "+payload.Svcname+ " application"
				//responseObject.Message = "Not authorized or jwt token was expired"
				//json.Unmarshal(bodyBytes, &responseObject)
				return response.JSON(w, responseObject)

}



// GET request on /api/baask8ss
func (handler *Handler) baask8sDeleteSvcs(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	srvname, err := request.RetrieveRouteVariableValue(r, "srvname")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid service name variable", err}
	}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

	//for i, _ := range baask8s.Applications {
	//	if (srvname == baask8s.Applications[i].Name) {
//
//			var responseObject1 jsonResponse
//			responseObject1.Success = false
//			responseObject1.Message = "Service " + payload.Svcname + " is not deployed yet"
//			return response.JSON(w, responseObject1)
//			//return &httperror.HandlerError{http.StatusBadRequest, "Application is already deployed", err}
//		}
//	}

	mychns := CHLPayload{}
	
	mychns.Currchannels = baask8s.CHLs
	mychns.Allchannels = baask8s.CHLs

    if len(mychns.Currchannels) == 0 {
		return &httperror.HandlerError{http.StatusNotFound, "No channel found for this operations", baasapi.ErrObjectNotFound}
	//	channel.CHLName = "default"
	//	mychns.Currchannels = append(mychns.Currchannels, channel)
	}

    var jsonData []byte
	jsonData, err = json.Marshal(mychns)
    if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute ansible commands ", err}
	}


	var srvname_list = strings.Split(srvname, "_")
	ansible_mode, _ := srvname_list[0], srvname_list[1]

	var ansible_env = "mode="+ ansible_mode+"_destroy"+" env=" +baask8s.Namespace+ " deploy_type=k8s "
	var ansible_extra = string(jsonData)
	//ansible_env = ansible_env + " --extra-vars '{"allchannels":[{"Id":0,"CHLName":"tst4c","CreatedAt":1567312246,"ORGs":[]}]}'"
	//ansible_env = ansible_env + "\" --extra-vars '{\"allchannels\":[{\"Id\":0,\"CHLName\":\"tst4c\",\"CreatedAt\":1567312246,\"ORGs\":[]}]}'"
	var ansible_config = "/data/k8s/ansible/operatefabric.yml"
	err = handler.CAFilesManager.Deploy(baask8s.Owner, baask8s.Namespace, ansible_extra, ansible_env, ansible_config, true)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute ansible commands ", err}
	}


	for i := range baask8s.Applications {
		if (srvname == baask8s.Applications[i].Name) {
			baask8s.Applications[i].Status = "DELETED"
		}
	}
	
	err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
	}
	
				var responseObject jsondeployResponse
				responseObject.Success = true
				responseObject.Message = "successfully uninstall fabric "+srvname+ " application"
				//responseObject.Message = "Not authorized or jwt token was expired"
				//json.Unmarshal(bodyBytes, &responseObject)
				return response.JSON(w, responseObject)

}