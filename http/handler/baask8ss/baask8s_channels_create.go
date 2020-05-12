package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	"strings"
	//"log"
	//"os"
	"time"
    "bytes"
    "encoding/json"
	//"fmt"
	//"reflect"
    "io/ioutil"
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





type jsonCreateCHLResponse struct {
    Success    bool    `json:"success"`
	Message    string  `json:"message"`
}

type CHLCreatePayload struct {
	Token                  string     `json:"token"`
	ChannelName            string     `json:"channelName"`
	ChannelConfigPath      string     `json:"channelConfigPath"`
	Username               string     `json:"username"`
	Orgname                string     `json:"orgname"`
}

type CHLPayload struct {
	Allchannels            []baasapi.CHL
	Currchannels           []baasapi.CHL
}

	//{"success":true,
	//"secret":"",
	//"message":"Jim enrolled Successfully",
	//"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTU2OTg2NTcsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Ik9yZzEiLCJpYXQiOjE1NTU2NjI2NTd9.racjuDcqswHY2WS9gj4XLBBwW-ST_yb9dDTZAlbh33Q"
	//}    

func (payload *CHLCreatePayload) Validate(r *http.Request) error {
	return nil;
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sChannelsCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	var payload CHLCreatePayload
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

	//exsitstr := false
	for index, _ := range baask8s.CHLs {
		if (baask8s.CHLs[index].CHLName == payload.ChannelName ){
			var responseObject1 jsonCreateCHLResponse
			responseObject1.Success = false
			responseObject1.Message = "Channel Name = " + payload.ChannelName + " already exist"
			return response.JSON(w, responseObject1)
			
		}
	    
	}


	myorgs := []baasapi.MSPORG{}

	//baask8sID := handler.Baask8sService.GetNextIdentifier()
	//var CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	channel := baasapi.CHL{
		CHLName:          payload.ChannelName,
		CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
		ORGs:             myorgs,
	}

	mychns := CHLPayload{}
	
	mychns.Currchannels = baask8s.CHLs
	baask8s.CHLs = append(baask8s.CHLs, channel)
	mychns.Allchannels = baask8s.CHLs

    if len(mychns.Currchannels) == 0 {
		channel.CHLName = "default"
		mychns.Currchannels = append(mychns.Currchannels, channel)
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

	
	//--extra-vars {"Allchannels":[{"Id":0,"CHLName":"tst4c","CreatedAt":1567312246,"ORGs":[]}]}

	var ansible_env = "mode=apply env=" +baask8s.Namespace+ " deploy_type=k8s channelnametocreate="+payload.ChannelName
	var ansible_extra = string(jsonData)
	//ansible_env = ansible_env + " --extra-vars '{"allchannels":[{"Id":0,"CHLName":"tst4c","CreatedAt":1567312246,"ORGs":[]}]}'"
	//ansible_env = ansible_env + "\" --extra-vars '{\"allchannels\":[{\"Id\":0,\"CHLName\":\"tst4c\",\"CreatedAt\":1567312246,\"ORGs\":[]}]}'"
	var ansible_config = "/data/k8s/ansible/operatefabric.yml"
	err = handler.CAFilesManager.Deploy(baask8s.Owner, baask8s.Namespace, ansible_extra, ansible_env, ansible_config, true)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute ansible commands ", err}
	}
	
    //jsonData := map[string]string{"username": "Jim", "orgName": "Org1"}
	//jsonValue, _ := json.Marshal(jsonData)
	jsonValue, _ := json.Marshal(payload)
	sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"

	var responseObject jsonCreateCHLResponse



	var client http.Client


	req, err := http.NewRequest("POST", sdk_url+"/channels", bytes.NewBuffer(jsonValue))
	req.Header.Add("Authorization" , "Bearer "+payload.Token)
	req.Header.Set("Content-Type", "application/json")

	if err != nil { 

		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
	}
	resp, err := client.Do(req)

	if err != nil {

		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}		
	}

	defer resp.Body.Close()


	if resp.StatusCode == 200 { // OK
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		//bodyBytes, _ := ioutil.ReadAll(resp.Body)
		//bodyString := string(bodyBytes)
		//var data2 interface{}
		//var responseObject jsonResponse
		//json.Unmarshal([]byte(bodyString), &responseObject)
		json.Unmarshal(bodyBytes, &responseObject)


		if responseObject.Success {

		



			err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
			}
			//return response.JSON(w, responseObject)



		} else {
			if (strings.Contains(responseObject.Message, "it is currently at version")) {
				err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
				if err != nil {
					return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
				}
			}

		}


		//var responseObject2 jsonResponse
		//responseObject2.Success = false
		//responseObject2.Message = responseObject
		return response.JSON(w, responseObject)

	    
	
	
	} else {
			if resp.StatusCode == 401 {

				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
				}
				bodyString := string(bodyBytes)
				//type jsonResponse struct {
				//    Success    bool    `json:"success"`
				//	Secret     string  `json:"secret"`
				//	Message    string  `json:"message"`
				//	Token      string  `json:"token"`
				//}
				var responseObject jsonResponse
				responseObject.Success = false
				responseObject.Message = bodyString
				//responseObject.Message = "Not authorized or jwt token was expired"
				//json.Unmarshal(bodyBytes, &responseObject)
				return response.JSON(w, responseObject)
			}


	}


	return response.JSON(w, nil)
}
