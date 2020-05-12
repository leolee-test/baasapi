package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	//"log"
	//"os"
    "bytes"
    "encoding/json"
    //"fmt"
	"io/ioutil"
	//"reflect"
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

//type jsonResponse struct {
//    Success    bool    `json:"success"`
//	Secret     string  `json:"secret"`
//	Message    string  `json:"message"`
//	Token      string  `json:"token"`
//}

type invokeccauthenticatePayload struct {
	ORGName    string    `json:"orgname"`
	Peers      []string  `json:"peers"`
	ChaincodeName string `json:"chaincodeName"`
	Args       []string  `json:"args"`
	Fcn        string `json:"fcn"`
	ChaincodeVersion string `json:"chaincodeVersion"`
	Token      string    `json:"token"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}
	//{"success":true,
	//"secret":"",
	//"message":"Jim enrolled Successfully",
	//"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTU2OTg2NTcsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Ik9yZzEiLCJpYXQiOjE1NTU2NjI2NTd9.racjuDcqswHY2WS9gj4XLBBwW-ST_yb9dDTZAlbh33Q"
	//}    

func (payload *invokeccauthenticatePayload) Validate(r *http.Request) error {
		return nil;
}


// GET request on /api/baask8ss
func (handler *Handler) baask8sChaincodesInvoke(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	var payload invokeccauthenticatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}


	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

		//fmt.Println(responseObject)


		var client http.Client

		//jsonData := map[string][]string{"peers": ["peer0org2.demo-test.baas.com"]}

		//jsonData["peers"] = append(jsonData["peers"], "peer1org2.demo-test.baas.com")


		jsonValue, _ := json.Marshal(payload)

		//fmt.Println(bytes.NewBuffer(jsonValue))
		//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("POST", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName, bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			//bodyBytes, err := ioutil.ReadAll(resp.Body)
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                // ... 

			}



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

	return response.JSON(w, baask8s.CHLs)
}
