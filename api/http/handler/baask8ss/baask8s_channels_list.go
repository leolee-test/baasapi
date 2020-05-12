package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	//"log"
	//"os"
    //"bytes"
    //"encoding/json"
    //"fmt"
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

type jsonResponse struct {
    Success    bool    `json:"success"`
	Secret     string  `json:"secret"`
	Message    string  `json:"message"`
	Namespace  string  `json:"namespace"`
	Token      string  `json:"token"`
	Result     map[string]interface{} `json:"result"`
}

type syncResponse struct {
    Success    bool    `json:"success"`
	Currentid     int  `json:"currentid"`
	CHL        baasapi.CHL `json:"chl"`
}
	//{"success":true,
	//"secret":"",
	//"message":"Jim enrolled Successfully",
	//"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTU2OTg2NTcsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Ik9yZzEiLCJpYXQiOjE1NTU2NjI2NTd9.racjuDcqswHY2WS9gj4XLBBwW-ST_yb9dDTZAlbh33Q"
	//}    



// GET request on /api/baask8ss
func (handler *Handler) baask8sChannelsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {


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



	return response.JSON(w, baask8s.CHLs)
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sChannelsListByName(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	channelname, err := request.RetrieveRouteVariableValue(r, "channelname")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid channel name variable", err}
	}


	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

	for index, _ := range baask8s.CHLs {
		if (baask8s.CHLs[index].CHLName == channelname) {
			return response.JSON(w, baask8s.CHLs[index])
		} 
	}
	return &httperror.HandlerError{http.StatusNotFound, "channel not found", baasapi.Error("Unable to find a channel with the name --- "+channelname)}

}

// GET request on /api/baask8ss
func (handler *Handler) baask8sChannelsListByNameSync(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	channelname, err := request.RetrieveRouteVariableValue(r, "channelname")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid channel name variable", err}
	}
	currentid, err := request.RetrieveNumericRouteVariableValue(r, "currentid")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid currentid variable", err}
	}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

	for index, _ := range baask8s.CHLs {
		if (baask8s.CHLs[index].CHLName == channelname) {

			var responseObject syncResponse
			responseObject.Success = true
			responseObject.Currentid = currentid
			responseObject.CHL = baask8s.CHLs[index]

			//type syncResponse struct {
			//	Success    bool    `json:"success"`
			//	Currentid     int  `json:"currentid"`
			//	CHL        baasapi.CHL `json:"chl"`
			//}
			return response.JSON(w, responseObject)
		} 
	}
	return &httperror.HandlerError{http.StatusNotFound, "channel not found", baasapi.Error("Unable to find a channel with the name --- "+channelname)}

}
