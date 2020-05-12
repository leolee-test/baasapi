package baask8ss

import (
	//"log"
	"net/http"
	"path"
	"encoding/json"
	//"runtime"
	//"strconv"
	//"math/rand"
	//"time"
	//"reflect"
	"github.com/asaskevich/govalidator"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	//"github.com/baasapi/baasapi/api/crypto"
	//"github.com/baasapi/baasapi/api/http/client"
)

//const (
	// Baas deployment files
	//BaaSDeploymentPath = "k8s/ansible/vars/namespaces"
//)

type settingsPayload struct {
	Owner       string                `json:"owner"`
	Settings     interface{}    `json:"settings"`
	//Settings     map[string]string    `json:"settings"`

}



func (payload *settingsPayload) Validate(r *http.Request) error {

	//username, err := request.RetrieveMultiPartFormValue(r, "Username", false)
	//log.Printf("http error: baask8s snapshot error (baask8s=%s)   (err=%s)\n", username, err)


	if govalidator.IsNull(payload.Owner) {
		return baasapi.Error("Invalid Owner name")
	}

	return nil
	
}

// POST request on /api/baask8ss/settings
func (handler *Handler) baask8sSettings(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	//if handler.authDisabled {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Cannot authenticate user. BaaSapi was started with the --no-auth flag", ErrAuthDisabled}
	//}

	//payload := &baask8sCreatePayload{}

	var payload settingsPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	var initfilen = "settings.input"


	
	//projectPath, err := handler.FileService.StoreYamlFileFromJSON(namespace, initfilen, payload, payload.Owner)
	_, err = handler.FileService.StoreYamlFileFromJSON("", initfilen, payload, payload.Owner)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist parameters file on disk", err}
	}

	//}

	return response.JSON(w, payload)

	//return response.JSON(w, baask8s)
}

// GET request on /api/baask8ss/settings
func (handler *Handler) getbaask8sSettings(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	//if handler.authDisabled {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Cannot authenticate user. BaaSapi was started with the --no-auth flag", ErrAuthDisabled}
	//}

	//payload := &baask8sCreatePayload{}

	//var payload settingsPayload
	//err := request.DecodeAndValidateJSONPayload(r, &payload)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	//}

	var initfilen = "settings.input"
	//yamlFilePath := path.Join(path.Join(service.fileStorePath, stackStorePath), fileName)

	//templatesJSON, err := handler.FileService.GetFileContent(path.Join(BaaSDeploymentPath, initfilen))
	templatesJSON, err := handler.FileService.GetFileContent(path.Join(handler.FileService.GetBinaryFolder(),"../k8s/ansible/vars/namespaces/", initfilen))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve settings definitions via filesystem", err}
	}


	var payload settingsPayload
	//var templates []baasapi.Template
	err = json.Unmarshal(templatesJSON, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to parse settings file. Please review your settings definition file", err}
	//	log.Println("Unable to parse templates file. Please review your template definition file.")
	//	return err
	}



	return response.JSON(w, payload.Settings)

	//return response.JSON(w, baask8s)
}
