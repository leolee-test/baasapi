package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	//"log"
	//"os"
	"path"
	//"reflect"
	"encoding/json"
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

type policyPayload struct {
	Name                    string        `json:"name"` 
	Description             string        `json:"desc"` 
	EndorsementPolicyparm   interface{}   `json:"endorsementPolicyparm"`
}


func (payload *policyPayload) Validate(r *http.Request) error {
	return nil;
}

// GET request on /api/baask8ss/{id}/policys
func (handler *Handler) baask8sPolicysList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	var initfilen = "policy.json"

	if len(baask8s.Policys) == 0 {
		content, err := handler.FileService.GetFileContent(path.Join(handler.FileService.GetBinaryFolder(),"../k8s/ansible/vars/namespaces/"+baask8s.Namespace+"/fabric/keyfiles/", initfilen))
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to read file on disk", err}
		}

		var data interface{}
		//var payload policyPayload
		err = json.Unmarshal(content, &data)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to parse settings file. Please review your settings definition file", err}
		//	log.Println("Unable to parse templates file. Please review your template definition file.")
		//	return err
		}

		policy := baasapi.Policy{
			Name:      "Default",
			Description: "At least one endorsement from all the related groups",
			EndorsementPolicyparm: data,
		}
		baask8s.Policys = append(baask8s.Policys, policy)
		err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
		}
	}


	return response.JSON(w, baask8s.Policys)
}


// PUT request on /api/baask8ss/{id}/policys
func (handler *Handler) baask8sPolicysUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}
//
	var payload policyPayload
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

	flag := 0
	for index, _ := range baask8s.Policys {
		if baask8s.Policys[index].Name == payload.Name {
			baask8s.Policys[index].Description = payload.Description
			baask8s.Policys[index].EndorsementPolicyparm = payload.EndorsementPolicyparm
			flag = 1
		}
	}

	if flag == 0 {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an policy with the specified identifier inside the database", baasapi.ErrObjectNotFound}
	}



//
//	baask8s.MSPs = append(baask8s.MSPs, msp)
//
	err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
	}


	//log.Printf(" to get msps:", baask8s.CCs)

	var responseObject jsonResponse
	responseObject.Success = true
	//responseObject.Message = bodyString
	responseObject.Message = "Policy have been sucessfully updated"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)

	//return response.JSON(w, baask8s.CCs)
}

// POST request on /api/baask8ss/{id}/policys
func (handler *Handler) baask8sPolicysAdd(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}
//
	var payload policyPayload
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

	flag := 0
	for index, _ := range baask8s.Policys {
		if baask8s.Policys[index].Name == payload.Name {
			//baask8s.Policys[index].Description = payload.Description
			//baask8s.Policys[index].EndorsementPolicyparm = payload.EndorsementPolicyparm
			flag = 1
		}
	}

	if flag == 1 {
		var responseObject jsonResponse
		responseObject.Success = false
		//responseObject.Message = bodyString
		responseObject.Message = "Policy already exists"
		//json.Unmarshal(bodyBytes, &responseObject)
		return response.JSON(w, responseObject)
		//return &httperror.HandlerError{http.StatusNotFound, "Unable to find an policy with the specified identifier inside the database", ErrPolicyAlreadyExists}
	}


	var policy = baasapi.Policy{}
	policy.Name = payload.Name
	policy.Description = payload.Description
	policy.EndorsementPolicyparm = payload.EndorsementPolicyparm
//
	baask8s.Policys = append(baask8s.Policys, policy)
//
	err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist policy changes inside the database", err}
	}


	//log.Printf(" to get msps:", baask8s.CCs)

	var responseObject jsonResponse
	responseObject.Success = true
	//responseObject.Message = bodyString
	responseObject.Message = "Policy have been sucessfully added"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)

	//return response.JSON(w, baask8s.CCs)
}

// Delete request on /api/baask8ss/{id}/policys
func (handler *Handler) baask8sPolicysDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}
	name, err := request.RetrieveRouteVariableValue(r, "name")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid policy name variable", err}
	}
//

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an baask8s with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
	}

	//policyii := make([]*baasapi.Policy, len(baask8s.Policys))
	//for j := range baask8s.Policys {
	//	policyii[j] = &baask8s.Policys[j]
	//}

	flag := 0

	if len(baask8s.Policys) == 1 {
		return &httperror.HandlerError{http.StatusBadRequest, "Cannot delete the last one", baasapi.ErrPolicyLastOne}
	}

	//for index, _ := range baask8s.Policys {
	for index := 0; index < len(baask8s.Policys); index ++ {
		if baask8s.Policys[index].Name == name {
		//if policyii[index].Name == name {
			//baask8s.Policys[index] = baask8s.Policys[len(baask8s.Policys)-1]
			//baask8s.Policys[len(baask8s.Policys)-1] = baasapi.Policy{}
			//baask8s.Policys = delete(baask8s.Policys,index)
			//fooType := reflect.TypeOf(baask8s.Policys)
			//log.Printf(" to get pods:", fooType)
			//log.Printf(" to get pods:", index)
			//log.Printf(" to get pods2:", baask8s.Policys[:index])
			//log.Printf(" to get pods2:", baask8s.Policys[index+1:])
			//log.Printf(" to get pods3:", baask8s.Policys[0])
			//log.Printf(" to get pods4:", baask8s.Policys[1])


			baask8s.Policys = append(baask8s.Policys[:index], baask8s.Policys[index+1:]...)
			//baask8s.Policys = append(baask8s.Policys[0:index], baask8s.Policys[index+1:]...)
			//baask8s.Policys = RemoveIndex(baask8s.Policys,index)
			//baask8s.Policys[index].Description = payload.Description
			//baask8s.Policys[index].EndorsementPolicyparm = payload.EndorsementPolicyparm
			flag = 1
		}
	}

	//baask8s.Policys = policyii

	

	if flag == 0 {
		var responseObject jsonResponse
		responseObject.Success = false
		//responseObject.Message = bodyString
		responseObject.Message = "Policy does not exists"
		//json.Unmarshal(bodyBytes, &responseObject)
		return response.JSON(w, responseObject)
		//return &httperror.HandlerError{http.StatusNotFound, "Unable to find an policy with the specified identifier inside the database", ErrPolicyAlreadyExists}
	}



//
//	baask8s.Policys = append(baask8s.Policys, payload)
//
	err = handler.Baask8sService.UpdateBaask8s(baask8s.ID, baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist policy changes inside the database", err}
	}


	//log.Printf(" to get msps:", baask8s.CCs)

	var responseObject jsonResponse
	responseObject.Success = true
	//responseObject.Message = bodyString
	responseObject.Message = "Policy have been sucessfully deleted"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)

	//return response.JSON(w, baask8s.CCs)
}



