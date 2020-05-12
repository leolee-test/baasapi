package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	"strings"
	//"log"
	//"os"
    "bytes"
    "encoding/json"
    "fmt"
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

type joinauthenticatePayload struct {
	ORGName    string    `json:"orgname"`
	Peers      []string  `json:"peers"`
	Token      string    `json:"token"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
}
	//{"success":true,
	//"secret":"",
	//"message":"Jim enrolled Successfully",
	//"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTU2OTg2NTcsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Ik9yZzEiLCJpYXQiOjE1NTU2NjI2NTd9.racjuDcqswHY2WS9gj4XLBBwW-ST_yb9dDTZAlbh33Q"
	//}    

func (payload *joinauthenticatePayload) Validate(r *http.Request) error {
		return nil;
}


// GET request on /api/baask8ss
func (handler *Handler) baask8sChannelsJoin(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	fmt.Println(channelname)

	var payload joinauthenticatePayload
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

	for index, _ := range baask8s.CHLs {
		if (baask8s.CHLs[index].CHLName == channelname ){

			for index2, _ := range baask8s.CHLs[index].ORGs {
				if (baask8s.CHLs[index].ORGs[index2].ORGName == payload.ORGName) {

					for index3, _ := range payload.Peers {
						if (contains(baask8s.CHLs[index].ORGs[index2].Peers, payload.Peers[index3])) {
							var responseObject jsonResponse
							responseObject.Success = false
							responseObject.Message = "Peer Name = " + payload.Peers[index3] + " already exist"
							//responseObject.Message = "Not authorized or jwt token was expired"
							//json.Unmarshal(bodyBytes, &responseObject)
							return response.JSON(w, responseObject)

						}
					}
				} 
			}
			
		}
		
	}

		//fmt.Println(responseObject)


		var client http.Client

		//jsonData := map[string][]string{"peers": ["peer0org2.demo-test.baas.com"]}

		//jsonData["peers"] = append(jsonData["peers"], "peer1org2.demo-test.baas.com")


		jsonValue, _ := json.Marshal(payload)

		//fmt.Println(bytes.NewBuffer(jsonValue))
		//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("POST", sdk_url+"/channels/"+channelname+"/peers", bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()


		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)




			if (responseObject.Success || strings.Contains(responseObject.Message, "LedgerID already exists")) {


				// MSP represents a MSP
	    		//CHL struct {
				//ID       CHLID      `json:"Id"`
				//CHLName  string     `json:"CHLName"`
				//CreatedAt        int64               `json:"CreatedAt"`
				//ORGs     []MSPORG   `json:"ORGs"`
				//}
				flag1 := 0
				for index, _ := range baask8s.CHLs {
					if (baask8s.CHLs[index].CHLName == channelname ){
						myorgs := baasapi.MSPORG{}
						myorgs.ORGName = payload.ORGName
						myorgs.Anchor = ""
						myorgs.Peers = payload.Peers

						flag1 = 0
						for index2, _ := range baask8s.CHLs[index].ORGs {
							if (baask8s.CHLs[index].ORGs[index2].ORGName == payload.ORGName) {
								flag1 = 1
								baask8s.CHLs[index].ORGs[index2].Peers = append(baask8s.CHLs[index].ORGs[index2].Peers, payload.Peers...)
							} 
						}
						if (flag1 == 0){
							baask8s.CHLs[index].ORGs = append(baask8s.CHLs[index].ORGs, myorgs)
						}

						err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
						if err != nil {
							return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
						}

						
					}
					
				}







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

	    return response.JSON(w, resp.Body)
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
