package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/baasapi/api"
)

type oauthPayload struct {
	Code string
}

func (payload *oauthPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Code) {
		return baasapi.Error("Invalid OAuth authorization code")
	}
	return nil
}

func (handler *Handler) authenticateThroughExtension(code, licenseKey string, settings *baasapi.OAuthSettings) (string, error) {
	//extensionURL := handler.ProxyManager.GetExtensionURL(baasapi.OAuthAuthenticationExtension)
	extensionURL := ""

	encodedConfiguration, err := json.Marshal(settings)
	if err != nil {
		return "", nil
	}

	req, err := http.NewRequest("GET", extensionURL+"/validate", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req.Header.Set("X-OAuth-Config", string(encodedConfiguration))
	req.Header.Set("X-OAuth-Code", code)
	req.Header.Set("X-BaaSapiExtension-License", licenseKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type extensionResponse struct {
		Username string `json:"Username,omitempty"`
		Err      string `json:"err,omitempty"`
		Details  string `json:"details,omitempty"`
	}

	var extResp extensionResponse
	err = json.Unmarshal(body, &extResp)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", baasapi.Error(extResp.Err + ":" + extResp.Details)
	}

	return extResp.Username, nil
}

func (handler *Handler) validateOAuth(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload oauthPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve settings from the database", err}
	}

	if settings.AuthenticationMethod != 3 {
		return &httperror.HandlerError{http.StatusForbidden, "OAuth authentication is not enabled", baasapi.Error("OAuth authentication is not enabled")}
	}

	extension, err := handler.ExtensionService.Extension(baasapi.OAuthAuthenticationExtension)
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Oauth authentication extension is not enabled", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a extension with the specified identifier inside the database", err}
	}

	username, err := handler.authenticateThroughExtension(payload.Code, extension.License.LicenseKey, &settings.OAuthSettings)
	if err != nil {
		log.Printf("[DEBUG] - OAuth authentication error: %s", err)
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to authenticate through OAuth", baasapi.ErrUnauthorized}
	}

	user, err := handler.UserService.UserByUsername(username)
	if err != nil && err != baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve a user with the specified username from the database", err}
	}

	if user == nil && !settings.OAuthSettings.OAuthAutoCreateUsers {
		return &httperror.HandlerError{http.StatusForbidden, "Account not created beforehand in BaaSapi and automatic user provisioning not enabled", baasapi.ErrUnauthorized}
	}

	if user == nil {
		user = &baasapi.User{
			Username: username,
			Role:     baasapi.StandardUserRole,
		}

		err = handler.UserService.CreateUser(user)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist user inside the database", err}
		}

		if settings.OAuthSettings.DefaultTeamID != 0 {
			membership := &baasapi.TeamMembership{
				UserID: user.ID,
				TeamID: settings.OAuthSettings.DefaultTeamID,
				Role:   baasapi.TeamMember,
			}

			err = handler.TeamMembershipService.CreateTeamMembership(membership)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist team membership inside the database", err}
			}
		}
	}

	return handler.writeToken(w, user)
}
