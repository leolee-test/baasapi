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

type enrollauthenticatePayload struct {
	Username   string  `json:"username"`
	OrgName    string  `json:"orgName"`
	Password    string  `json:"password"`
	//{"username": "Jim", "orgName": "Org1"}
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
}
type AttrReg struct {
	Name       string    `json:"name"`
	Value      string    `json:"value"`
}

type reenrollPayload struct {
	Username   string  `json:"username"`
	OrgName    string  `json:"orgName"`
	Token      string    `json:"token"`
	//Password    string  `json:"password"`
	//{"username": "Jim", "orgName": "Org1"}
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
}
type updatePayload struct {
	//Reqjson    interface{} `json:"reqjson"`
	Reqjson    interface{} `json:"reqjson"`
	Token      string    `json:"token"`
}
type deletePayload struct {
	//Reqjson    interface{} `json:"reqjson"`
	Force      bool `json:"force"`
	Token      string    `json:"token"`
}
//type Restrictionreq struct {
//	RevokedBefore      string    `json:"revokedBefore"`
//	RevokedAfter       string    `json:"revokedAfter"`
//	ExpireBefore       string    `json:"expireBefore"`
//	ExpireAfter       string    `json:"expireAfter"`
//}
type crlPayload struct {
	Restrictionreq struct {
		RevokedBefore string `json:"revokedBefore"`
		RevokedAfter  string `json:"revokedAfter"`
		ExpireBefore  string `json:"expireBefore"`
		ExpireAfter   string   `json:"expireAfter"`
	} `json:"restrictionreq"`
	OrgName    string  `json:"orgName"`
	Token      string    `json:"token"`
}

type revokePayload struct {
	//Reqjson    interface{} `json:"reqjson"`
	Username   string `json:"username"`
	OrgName    string  `json:"orgName"`
	Reason     string  `json:"reason"`
	Token      string    `json:"token"`
}

type listPayload struct {
	OrgName    string  `json:"orgName"`
	Token      string    `json:"token"`
	//Password    string  `json:"password"`
	//{"username": "Jim", "orgName": "Org1"}
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
}
type regPayload struct {
	Newusername   string  `json:"newusername"`
	OrgName     string  `json:"orgName"`
	Password    string  `json:"password"`
	Role        string  `json:"role"`
	Token      string    `json:"token"`
	Attrs      []AttrReg  `json:"attrs"`
//    "orgName": "orgdbs",
//    "password": "12345678",
//    "role": "client",
//    "newusername": "test27777",
//    "attrs": [{
//                            "name": "hf.Registrar.Attributes",
//                            "value": "client"
//                        },
//                        {
//                            "name": "hf.AffiliationMgr",
//                            "value": "1"
//                        },
//                        {
//                            "name": "hf.Registrar.Roles",
//                            "value": "client"
//                       },
//                        {
//                            "name": "hf.Registrar.DelegateRoles",
//                            "value": "client"
//                        },
//                        {
//                            "name": "hf.Revoker",
//                            "value": "1"
//                        },
//                        {
//                            "name": "hf.IntermediateCA",
//                            "value": "1"
//                       },
//                        {
//                            "name": "hf.GenCRL",
//                            "value": "1"
//                        }
//                    ]
}
	//{"success":true,
	//"secret":"",
	//"message":"Jim enrolled Successfully",
	//"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTU2OTg2NTcsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Ik9yZzEiLCJpYXQiOjE1NTU2NjI2NTd9.racjuDcqswHY2WS9gj4XLBBwW-ST_yb9dDTZAlbh33Q"
	//}    
type sdkResponse struct {
		Success    bool    `json:"success"`
		Errors     string  `json:"errors"`
		Messages    string  `json:"messages"`
		Result  interface{}  `json:"result"`
	}

func (payload *enrollauthenticatePayload) Validate(r *http.Request) error {
		return nil;
}
func (payload *regPayload) Validate(r *http.Request) error {
	return nil;
}
func (payload *reenrollPayload) Validate(r *http.Request) error {
	return nil;
}
func (payload *listPayload) Validate(r *http.Request) error {
return nil;
}
func (payload *updatePayload) Validate(r *http.Request) error {
	return nil;
}
func (payload *deletePayload) Validate(r *http.Request) error {
	return nil;
}
func (payload *revokePayload) Validate(r *http.Request) error {
	return nil;
}
func (payload *crlPayload) Validate(r *http.Request) error {
	return nil;
}



// GET request on /api/baask8ss/{id}/ca/users
func (handler *Handler) baask8susersEnroll(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload enrollauthenticatePayload
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

	jsonValue, _ := json.Marshal(payload)
	sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	//fmt.Printf(reflect.TypeOf(jsonValue))

    jsonresponse, err := http.Post(sdk_url+"/enrollusers", "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "The HTTP request failed with error", err}
        //fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(jsonresponse.Body)
		var responseObject jsonResponse
		json.Unmarshal(data, &responseObject)


		return response.JSON(w, responseObject)
	}
	var responseObject jsonResponse
	responseObject.Success = false
	responseObject.Message = "The HTTP request failed with error"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)
}

// GET request on /api/baask8ss/{id}/ca/users
func (handler *Handler) baask8susersRegister(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload regPayload
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

		var client http.Client
		jsonValue, _ := json.Marshal(payload)

		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("POST", sdk_url+"/users", bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
			}
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                return response.JSON(w, responseObject)
		
			} else {
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
					return response.JSON(w, responseObject)
			}
		}
		var responseObject jsonResponse
		responseObject.Success = false
		responseObject.Message = "HTTP error" 
		return response.JSON(w, responseObject)
}

// GET request on /api/baask8ss/{id}/ca/reenroll
func (handler *Handler) baask8susersreEnroll(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload reenrollPayload
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

		var client http.Client
		jsonValue, _ := json.Marshal(payload)

		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("POST", sdk_url+"/reenrollusers", bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
			}
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                return response.JSON(w, responseObject)
		
			} else {
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
					return response.JSON(w, responseObject)
			}
		}
		var responseObject jsonResponse
		responseObject.Success = false
		responseObject.Message = "HTTP error" 
		return response.JSON(w, responseObject)
}

// GET request on /api/baask8ss/{id}/ca/users
func (handler *Handler) baask8susersList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload listPayload
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

		var client http.Client
		jsonValue, _ := json.Marshal(payload)

		//req, err := http.NewRequest("GET", "http://11.11.11.120:30500/channels/mychannel/blocks/1?peer=peer0org1.demo-test.baas.com", nil)
		//req.Header.Add("Authorization" , "Bearer "+payload.Token)
		//req.Header.Set("Content-Type", "application/json")

		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("POST", sdk_url+"/listusers", bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
			}
			//bodyString := string(bodyBytes)
			var responseObject sdkResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                return response.JSON(w, responseObject)
		
			} else {
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
					var responseObject sdkResponse
					responseObject.Success = false
					responseObject.Messages = bodyString
					return response.JSON(w, responseObject)
			}
		}
		var responseObject sdkResponse
		responseObject.Success = false
		responseObject.Messages = "HTTP error" 
		return response.JSON(w, responseObject)
}

// PUT request on /api/baask8ss/{id}/ca/users/{orgname}/{username}",
func (handler *Handler) baask8susersUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload updatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	orgname, err := request.RetrieveRouteVariableValue(r, "orgname")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid org name variable", err}
	}
	username, err := request.RetrieveRouteVariableValue(r, "username")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid user name variable", err}
	}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

		var client http.Client
		jsonValue, _ := json.Marshal(payload)

		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("PUT", sdk_url+"/users/"+orgname+"/"+username, bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
			}
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                return response.JSON(w, responseObject)
		
			} else {
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
					return response.JSON(w, responseObject)
			}
		}
		var responseObject jsonResponse
		responseObject.Success = false
		responseObject.Message = "HTTP error" 
		return response.JSON(w, responseObject)
}

// GET request on /api/baask8ss/{id}/ca/users
func (handler *Handler) baask8susersDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload deletePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	orgname, err := request.RetrieveRouteVariableValue(r, "orgname")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid org name variable", err}
	}
	username, err := request.RetrieveRouteVariableValue(r, "username")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid user name variable", err}
	}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

		var client http.Client
		jsonValue, _ := json.Marshal(payload)

		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("DELETE", sdk_url+"/users/"+orgname+"/"+username, bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
			}
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                return response.JSON(w, responseObject)
		
			} else {
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
					return response.JSON(w, responseObject)
			}
		}
		var responseObject jsonResponse
		responseObject.Success = false
		responseObject.Message = "HTTP error" 
		return response.JSON(w, responseObject)
}

// GET request on /api/baask8ss/{id}/ca/users
func (handler *Handler) baask8susersRevoke(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload revokePayload
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

		var client http.Client
		jsonValue, _ := json.Marshal(payload)

		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("POST", sdk_url+"/revokeusers", bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
			}
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                return response.JSON(w, responseObject)
		
			} else {
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
					return response.JSON(w, responseObject)
			}
		}
		var responseObject jsonResponse
		responseObject.Success = false
		responseObject.Message = "HTTP error" 
		return response.JSON(w, responseObject)
}

// GET request on /api/baask8ss/{id}/ca/users
func (handler *Handler) baask8susersListbyName(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload reenrollPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	//fmt.Println(payload)

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

		var client http.Client
		jsonValue, _ := json.Marshal(payload)

		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("POST", sdk_url+"/listusersbyid", bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
			}
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                return response.JSON(w, responseObject)
		
			} else {
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
					return response.JSON(w, responseObject)
			}
		}
		var responseObject jsonResponse
		responseObject.Success = false
		responseObject.Message = "HTTP error" 
		return response.JSON(w, responseObject)
}

// GET request on /api/baask8ss/{id}/ca/users
func (handler *Handler) baask8susersGenerateCRL(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload crlPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	//fmt.Println(payload)

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

		var client http.Client
		jsonValue, _ := json.Marshal(payload)

		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		req, err := http.NewRequest("POST", sdk_url+"/generateCRL", bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "unable to read io from response", err}
			}
			//bodyString := string(bodyBytes)
			var responseObject jsonResponse
			json.Unmarshal(bodyBytes, &responseObject)


			if responseObject.Success {

                return response.JSON(w, responseObject)
		
			} else {
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
					return response.JSON(w, responseObject)
			}
		}
		var responseObject jsonResponse
		responseObject.Success = false
		responseObject.Message = "HTTP error" 
		return response.JSON(w, responseObject)
}


