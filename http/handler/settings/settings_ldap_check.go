package settings

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/filesystem"
)

type settingsLDAPCheckPayload struct {
	LDAPSettings baasapi.LDAPSettings
}

func (payload *settingsLDAPCheckPayload) Validate(r *http.Request) error {
	return nil
}

// PUT request on /settings/ldap/check
func (handler *Handler) settingsLDAPCheck(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload settingsLDAPCheckPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	if (payload.LDAPSettings.TLSConfig.TLS || payload.LDAPSettings.StartTLS) && !payload.LDAPSettings.TLSConfig.TLSSkipVerify {
		caCertPath, _ := handler.FileService.GetPathForTLSFile(filesystem.LDAPStorePath, baasapi.TLSFileCA)
		payload.LDAPSettings.TLSConfig.TLSCACertPath = caCertPath
	}

	err = handler.LDAPService.TestConnectivity(&payload.LDAPSettings)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to connect to LDAP server", err}
	}

	return response.Empty(w)
}
