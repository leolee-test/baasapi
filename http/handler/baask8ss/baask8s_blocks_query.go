package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	"log"
	//"os"
    "bytes"
    "encoding/json"
    "fmt"
	"io/ioutil"
	"reflect"
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

type queryblkauthenticatePayload struct {
	ORGName    string    `json:"orgname"`
	Peers      []string  `json:"peers"`
	Blocknum   string    `json:"blocknum"`
	Hashstr    string    `json:"hashstr"`
	Token      string    `json:"token"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}

type querychannelsPayload struct {
	ORGName    string    `json:"orgname"`
	Peer       string    `json:"peer"`
	Token      string    `json:"token"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}

type querychannelinfoPayload struct {
	Channelname    string    `json:"channelname"`
	Peer       string    `json:"peer"`
	Token      string    `json:"token"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}

type jsonChannelsResponse struct {
	Channels    interface{}    `json:"channels"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}
type jsonChannelinfoResponse struct {
	Height     interface{}    `json:"height"`
	Channelname string        `json:"channelname"`
}

type jsonCCsResponse struct {
	Success    bool    `json:"success"`
	Message    string    `json:"message"`
	Result     []queryCC    `json:"result"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}
type queryCC struct {
	Name    string    `json:"name"`
	Version    string    `json:"version"`
	Path     string    `json:"path"`
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

func (payload *queryblkauthenticatePayload) Validate(r *http.Request) error {
		return nil;
}

func (payload *querychannelsPayload) Validate(r *http.Request) error {
	return nil;
}

func (payload *querychannelinfoPayload) Validate(r *http.Request) error {
	return nil;
}


// GET request on /api/baask8ss
func (handler *Handler) baask8sBlocksQuery(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	blocknum, err := request.RetrieveRouteVariableValue(r, "blocknum")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid block numble variable", err}
	}
	fmt.Println(blocknum)

	var payload queryblkauthenticatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	fmt.Println(payload)

	fmt.Println("Starting the application...")

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

		fmt.Println(jsonValue)
		fmt.Println(payload)
		//fmt.Println(bytes.NewBuffer(jsonValue))
		//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
		
		req, err := http.NewRequest("GET", "http://11.11.11.120:30500/channels/mychannel/blocks/1?peer=peer0org1.demo-test.baas.com", nil)
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		fmt.Println(resp.StatusCode)

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			var responseObject jsonChannelsResponse
			json.Unmarshal(bodyBytes, &responseObject)

			log.Printf("(baask8s=%s) (passowrd=) \n", err)
			fmt.Println(bodyString)

			//if responseObject.Success {
				fmt.Println("bodyString & data2...")
				fmt.Println(bodyString)
				//fmt.Println(responseObject.Message)


				fmt.Println(baask8s)

				fmt.Println(baask8s.CHLs)

				//fmt.Printf(reflect.TypeOf(baask8s.CHLs))
				log.Printf("(baask8s=%s) (passowrd=) \n", reflect.TypeOf(baask8s.CHLs))
				//log.Printf("(baask8s=%s) (passowrd=) \n", reflect.TypeOf(responseObject.Message))
				log.Printf("(baask8s=%s) (passowrd=) \n", len(baask8s.CHLs))

				return response.JSON(w, baask8s.CHLs)


				// MSP represents a MSP
	    		//CHL struct {
				//ID       CHLID      `json:"Id"`
				//CHLName  string     `json:"CHLName"`
				//CreatedAt        int64               `json:"CreatedAt"`
				//ORGs     []MSPORG   `json:"ORGs"`
				//}
		
				//myorgs := baasapi.MSPORG{}

				//channel := baasapi.CHL{
				//	ID:               1,
				//	CHLName:          "data2.channels",
				//	CreatedAt:        time.Now().Unix(),
				//	ORGs:             myorgs,
				//}
				//MSPORG struct {
				//	ORGName   string     `json:"ORGName"`
				//	Peers     []string   `json:"Peers"`
				//}

				//myorgs.ORGName = "Org1"
				//myorgs.Anchor = ""
				//myorgs.Peers = ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"] 
				//myorgs.Peers = payload.Peers


				//baask8s.CHLs[0].ORGs = append(baask8s.CHLs[0].ORGs, myorgs)
				//fmt.Println(baask8s.CHLs)


				//err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
				//if err != nil {
				//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
				//}



			//}



			//return response.JSON(w, responseObject)
		
		}

	return response.JSON(w, baask8s.CHLs)
}


// GET request on /api/baask8ss
func (handler *Handler) baask8schannelsbyPeer(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	//peer, err := request.RetrieveRouteVariableValue(r, "peer")
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusBadRequest, "Invalid peer name variable", err}
	//}
	//fmt.Println(channelname)

	var payload querychannelsPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	fmt.Println(payload)

	fmt.Println("Starting the application...")

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

		fmt.Println(jsonValue)
		fmt.Println(payload)
		//fmt.Println(bytes.NewBuffer(jsonValue))
		//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
		
		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		//req, err := http.NewRequest("GET", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn=query&args=%5B%22a%22%5D", nil)
		//req, err := http.NewRequest("POST", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn="+payload.Fcn+"&args="+payload.Args[0], nil)
		req, err := http.NewRequest("GET", sdk_url+"/channels/?peer="+payload.Peer, bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			var responseObject jsonChannelsResponse
			json.Unmarshal(bodyBytes, &responseObject)

			log.Printf("(baask8s=%s) (passowrd=) \n", err)
			fmt.Println(bodyString)

			//if responseObject.Success {
				fmt.Println("bodyString & data2...")
				fmt.Println(bodyString)
				//fmt.Println(responseObject.Message)


				fmt.Println(baask8s)

				fmt.Println(baask8s.CHLs)

				//fmt.Printf(reflect.TypeOf(baask8s.CHLs))
				log.Printf("(baask8s=%s) (passowrd=) \n", reflect.TypeOf(baask8s.CHLs))
				//log.Printf("(baask8s=%s) (passowrd=) \n", reflect.TypeOf(responseObject.Message))
				log.Printf("(baask8s=%s) (passowrd=) \n", len(baask8s.CHLs))

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

	return response.JSON(w, nil)
}

// GET request on /api/baask8ss
func (handler *Handler) baask8schannelinfobyPeer(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	//peer, err := request.RetrieveRouteVariableValue(r, "peer")
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusBadRequest, "Invalid peer name variable", err}
	//}
	//fmt.Println(channelname)

	var payload querychannelinfoPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	fmt.Println(payload)

	fmt.Println("Starting the application...")

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

		fmt.Println(jsonValue)
		fmt.Println(payload)
		//fmt.Println(bytes.NewBuffer(jsonValue))
		//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
		
		sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
		//req, err := http.NewRequest("GET", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn=query&args=%5B%22a%22%5D", nil)
		//req, err := http.NewRequest("POST", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn="+payload.Fcn+"&args="+payload.Args[0], nil)
		req, err := http.NewRequest("GET", sdk_url+"/channels/"+payload.Channelname+"?peer="+payload.Peer, bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization" , "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {}
		resp, err := client.Do(req)

		if err != nil {}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			var responseObject jsonChannelinfoResponse
			json.Unmarshal(bodyBytes, &responseObject)
			responseObject.Channelname = payload.Channelname

			log.Printf("(baask8s=%s) (passowrd=) \n", err)
			fmt.Println(bodyString)

			//if responseObject.Success {
				fmt.Println("bodyString & data2...")
				fmt.Println(bodyString)
				//fmt.Println(responseObject.Message)


				fmt.Println(baask8s)

				fmt.Println(baask8s.CHLs)

				//fmt.Printf(reflect.TypeOf(baask8s.CHLs))
				//log.Printf("(baask8s=%s) (passowrd=) \n", reflect.TypeOf(baask8s.CHLs))
				//log.Printf("(baask8s=%s) (passowrd=) \n", reflect.TypeOf(responseObject.Message))
				//log.Printf("(baask8s=%s) (passowrd=) \n", len(baask8s.CHLs))

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

	return response.JSON(w, nil)
}


// GET request on /api/baask8ss
func baask8schannelsSYNC() {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}


			log.Printf("(baask8s=%s) (passowrd=) \n", "adf")

			//var baask8ss = make([]baasapi.Baask8s, 0)
			//baask8ss, err := handler.Baask8sService.Baask8ss()
			//for index, _ := range baask8ss {
			//	for index2, _ := range baask8ss[index].CHLs {
			//		log.Printf("(baask8s=%s) (passowrd=) \n", baask8ss[index].CHLs[index2].CHLName)
			//	}
			//}

	//return nill
}
func (handler *Handler) baask8schaincodesbyPeer(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	//peer, err := request.RetrieveRouteVariableValue(r, "peer")
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusBadRequest, "Invalid peer name variable", err}
	//}
	//fmt.Println(channelname)

	var payload querychannelinfoPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	fmt.Println(payload)

	fmt.Println("Starting the application...")

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}
	var client http.Client

	//var 
	jsonValue, _ := json.Marshal(payload)

	fmt.Println(jsonValue)
	fmt.Println(payload)
	//fmt.Println(bytes.NewBuffer(jsonValue))
	//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
	
	sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	//req, err := http.NewRequest("GET", sdk_url+"/chaincodes?peer="+peername, nil)
	//req, err := http.NewRequest("GET", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn=query&args=%5B%22a%22%5D", nil)
	//req, err := http.NewRequest("POST", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn="+payload.Fcn+"&args="+payload.Args[0], nil)
	req, err := http.NewRequest("GET", sdk_url+"/channels/"+payload.Channelname+"/chaincodes"+"?peer="+payload.Peer, bytes.NewBuffer(jsonValue))
	req.Header.Add("Authorization" , "Bearer "+payload.Token)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {}
	resp, err := client.Do(req)

	if err != nil {}

	defer resp.Body.Close()

	if resp.StatusCode == 200 { // OK
		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil { fmt.Println(err) }
		//bodyString := string(bodyBytes)
		var responseObject jsonCCsResponse
		json.Unmarshal(bodyBytes, &responseObject)

		fmt.Println(responseObject.Result)

		for index, _ := range responseObject.Result {

			fmt.Println("test")
			fmt.Println(responseObject.Result[index].Name)
			fmt.Println(responseObject.Result[index].Version)
			fmt.Println(responseObject.Result[index].Path)

			//myorgs := []baasapi.MSPORG{}

			//ccs := []baasapi.CC{}

			//CC struct {
			//	//ID              CCID      `json:"id"`
			//	ID              int        `json:"id"`
			//	CCName          string     `json:"chaincodeName"`
			//	CHLName         string     `json:"chlName"`
			//	Version         string     `json:"chaincodeVersion"`
			//	EndorsementPolicyparm        interface{}     `json:"endorsementPolicyparm"`
			//	InstallORGs     []MSPORG   `json:"installORGs"`
			//	InstantiateORGs []MSPORG   `json:"instantiateORGs"`
			//	ChaincodeType   string     `json:"chaincodeType"`
			//}
		}
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
	//return nil
}
