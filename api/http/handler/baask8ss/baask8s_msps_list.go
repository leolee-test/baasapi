package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	//"log"
	//"os"
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

type updateMSPsPayload struct {
	MSPName    string    `json:"mspname"`
    ORGName    string    `json:"orgname"`
	Peers      []string  `json:"peers"`
	Role       int       `json:"role"`
	//Token      string    `json:"token"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
}

type jsonMSPsResponse struct {
    Success    bool    `json:"success"`
	Message    string  `json:"message"`
	Namespace  string  `json:"namespace"`
	Networkname string   `json:"networkname"`
	MSPs       []baasapi.MSP       `json:"msps"`
	//Token      string    `json:"token"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
}

//"Id": 1,
//"MSPName": "ca1org1.orgorg1",
//"ORGs": [
//	{
//		"ORGName": "orgorg1",
//		"Anchor": "",
//		"Peers": [
//			"Anchor@peer1.orgorg1",
//			"Worker@peer2.orgorg1"
//		]
//	}
//],
//"Role": 1

func (payload *updateMSPsPayload) Validate(r *http.Request) error {
	return nil;
}

// GET request on /api/baask8ss/{id}/msps
func (handler *Handler) baask8sMspsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	//log.Println("retrieving baask8s ID: " + baask8sID)

	//var payload baask8sUpdateAccessPayload
	//err = request.DecodeAndValidateJSONPayload(r, &payload)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	//}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}


	var responseObject jsonMSPsResponse
	responseObject.Success = true
	responseObject.Message = "Retrieved the MSP information successfully"
	responseObject.Namespace = baask8s.Namespace
	responseObject.Networkname = baask8s.NetworkName
	responseObject.MSPs      = baask8s.MSPs
	//responseObject.Message = "Not authorized or jwt token was expired"
	//json.Unmarshal(bodyBytes, &responseObject)
	//return response.JSON(w, responseObject)

	return response.JSON(w, responseObject)
}



// GET request on /api/baask8ss/{id}/msps
func (handler *Handler) baask8sCCsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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


	return response.JSON(w, baask8s.CCs)
}

// POST request on /api/baask8ss/{id}/msps
func (handler *Handler) baask8sUpdateMSPs(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload updateMSPsPayload
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


	msporg := baasapi.MSPORG{
		ORGName:      payload.ORGName,
		Anchor:       "",
		Peers:        payload.Peers,
	}

    //msp := baasapi.MSP{
	//ID:               len(baask8s.MSPs)+1,
	//MSPName:          payload.MSPName,
	//Role:             payload.Role,
	//ORGs:             null,
	//}

	msp := baasapi.MSP{}
	msp.ID = len(baask8s.MSPs)+1
	msp.MSPName = payload.MSPName
	msp.Role = payload.Role
	msp.ORGs = []baasapi.MSPORG{}
	
	msp.ORGs = append(msp.ORGs, msporg)

	baask8s.MSPs = append(baask8s.MSPs, msp)

	err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
	}


	//log.Printf(" to get msps:", baask8s.CCs)

	var responseObject jsonResponse
	responseObject.Success = true
	//responseObject.Message = bodyString
	responseObject.Message = "Orgnizatons/peers have been sucessfully updated"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)

	//return response.JSON(w, baask8s.CCs)
}


