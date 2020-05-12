package upload

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type jsonResponse struct {
    Success    bool    `json:"success"`
	Secret     string  `json:"secret"`
	Message    string  `json:"message"`
	Namespace  string  `json:"namespace"`
	Token      string  `json:"token"`
}

// POST request on /api/upload/tls/{certificate:(?:ca|cert|key)}?folder=<folder>
func (handler *Handler) uploadTLS(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	certificate, err := request.RetrieveRouteVariableValue(r, "certificate")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid certificate route variable", err}
	}

	folder, err := request.RetrieveMultiPartFormValue(r, "folder", false)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: folder", err}
	}

	file, _, err := request.RetrieveMultiPartFormFile(r, "file")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid certificate file. Ensure that the certificate file is uploaded correctly", err}
	}

	var fileType baasapi.TLSFileType
	switch certificate {
	case "ca":
		fileType = baasapi.TLSFileCA
	case "cert":
		fileType = baasapi.TLSFileCert
	case "key":
		fileType = baasapi.TLSFileKey
	default:
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid certificate route value. Value must be one of: ca, cert or key", baasapi.ErrUndefinedTLSFileType}
	}

	_, err = handler.FileService.StoreTLSFileFromBytes(folder, fileType, file)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist certificate file on disk", err}
	}

	var responseObject jsonResponse
	responseObject.Success = true
	responseObject.Message = "Sucessfully upload the files to server"
	//responseObject.Message = "Not authorized or jwt token was expired"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)
	//return response.Empty(w)
}

// POST request on /api/upload/kubeconfig?folder=<folder>
func (handler *Handler) uploadKubeconfig(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	//folder, err := request.RetrieveMultiPartFormValue(r, "folder", false)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: folder", err}
	//}

	file, _, err := request.RetrieveMultiPartFormFile(r, "file")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid certificate file. Ensure that the certificate file is uploaded correctly", err}
	}

	var fileType baasapi.TLSFileType

	fileType = baasapi.TLSFileOther

	//_, err = handler.FileService.StoreTLSFileFromBytes(folder, fileType, file)
	_, err = handler.FileService.StoreKubeconfigFileFromBytes("vars", fileType, file)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist config file on disk", err}
	}

	var responseObject jsonResponse
	responseObject.Success = true
	responseObject.Message = "Sucessfully upload the files to server"
	//responseObject.Message = "Not authorized or jwt token was expired"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)
	//return response.Empty(w)
}

