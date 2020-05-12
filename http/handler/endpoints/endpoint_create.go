package endpoints

import (
	"log"
	"net/http"
	"runtime"
	"strconv"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/crypto"
	"github.com/baasapi/baasapi/api/http/client"
)

type endpointCreatePayload struct {
	Name                   string
	URL                    string
	EndpointType           int
	PublicURL              string
	GroupID                int
	TLS                    bool
	TLSSkipVerify          bool
	TLSSkipClientVerify    bool
	TLSCACertFile          []byte
	TLSCertFile            []byte
	TLSKeyFile             []byte
	AzureApplicationID     string
	AzureTenantID          string
	AzureAuthenticationKey string
	Tags                   []string
}

func (payload *endpointCreatePayload) Validate(r *http.Request) error {
	name, err := request.RetrieveMultiPartFormValue(r, "Name", false)
	if err != nil {
		return baasapi.Error("Invalid endpoint name")
	}
	payload.Name = name

	endpointType, err := request.RetrieveNumericMultiPartFormValue(r, "EndpointType", false)
	if err != nil || endpointType == 0 {
		return baasapi.Error("Invalid endpoint type value. Value must be one of: 1 (Docker environment), 2 (Agent environment) or 3 (Azure environment)")
	}
	payload.EndpointType = endpointType

	groupID, _ := request.RetrieveNumericMultiPartFormValue(r, "GroupID", true)
	if groupID == 0 {
		groupID = 1
	}
	payload.GroupID = groupID

	var tags []string
	err = request.RetrieveMultiPartFormJSONValue(r, "Tags", &tags, true)
	if err != nil {
		return baasapi.Error("Invalid Tags parameter")
	}
	payload.Tags = tags
	if payload.Tags == nil {
		payload.Tags = make([]string, 0)
	}

	useTLS, _ := request.RetrieveBooleanMultiPartFormValue(r, "TLS", true)
	payload.TLS = useTLS

	if payload.TLS {
		skipTLSServerVerification, _ := request.RetrieveBooleanMultiPartFormValue(r, "TLSSkipVerify", true)
		payload.TLSSkipVerify = skipTLSServerVerification
		skipTLSClientVerification, _ := request.RetrieveBooleanMultiPartFormValue(r, "TLSSkipClientVerify", true)
		payload.TLSSkipClientVerify = skipTLSClientVerification

		if !payload.TLSSkipVerify {
			caCert, _, err := request.RetrieveMultiPartFormFile(r, "TLSCACertFile")
			if err != nil {
				return baasapi.Error("Invalid CA certificate file. Ensure that the file is uploaded correctly")
			}
			payload.TLSCACertFile = caCert
		}

		if !payload.TLSSkipClientVerify {
			cert, _, err := request.RetrieveMultiPartFormFile(r, "TLSCertFile")
			if err != nil {
				return baasapi.Error("Invalid certificate file. Ensure that the file is uploaded correctly")
			}
			payload.TLSCertFile = cert

			key, _, err := request.RetrieveMultiPartFormFile(r, "TLSKeyFile")
			if err != nil {
				return baasapi.Error("Invalid key file. Ensure that the file is uploaded correctly")
			}
			payload.TLSKeyFile = key
		}
	}

	switch baasapi.EndpointType(payload.EndpointType) {
	case baasapi.AzureEnvironment:
		azureApplicationID, err := request.RetrieveMultiPartFormValue(r, "AzureApplicationID", false)
		if err != nil {
			return baasapi.Error("Invalid Azure application ID")
		}
		payload.AzureApplicationID = azureApplicationID

		azureTenantID, err := request.RetrieveMultiPartFormValue(r, "AzureTenantID", false)
		if err != nil {
			return baasapi.Error("Invalid Azure tenant ID")
		}
		payload.AzureTenantID = azureTenantID

		azureAuthenticationKey, err := request.RetrieveMultiPartFormValue(r, "AzureAuthenticationKey", false)
		if err != nil {
			return baasapi.Error("Invalid Azure authentication key")
		}
		payload.AzureAuthenticationKey = azureAuthenticationKey
	default:
		url, err := request.RetrieveMultiPartFormValue(r, "URL", true)
		if err != nil {
			return baasapi.Error("Invalid endpoint URL")
		}
		payload.URL = url

		publicURL, _ := request.RetrieveMultiPartFormValue(r, "PublicURL", true)
		payload.PublicURL = publicURL
	}

	return nil
}

// POST request on /api/endpoints
func (handler *Handler) endpointCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	if !handler.authorizeEndpointManagement {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Endpoint management is disabled", ErrEndpointManagementDisabled}
	}

	payload := &endpointCreatePayload{}
	err := payload.Validate(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	endpoint, endpointCreationError := handler.createEndpoint(payload)
	if endpointCreationError != nil {
		return endpointCreationError
	}

	return response.JSON(w, endpoint)
}

func (handler *Handler) createEndpoint(payload *endpointCreatePayload) (*baasapi.Endpoint, *httperror.HandlerError) {
	if baasapi.EndpointType(payload.EndpointType) == baasapi.AzureEnvironment {
		return handler.createAzureEndpoint(payload)
	}

	if payload.TLS {
		return handler.createTLSSecuredEndpoint(payload)
	}
	return handler.createUnsecuredEndpoint(payload)
}

func (handler *Handler) createAzureEndpoint(payload *endpointCreatePayload) (*baasapi.Endpoint, *httperror.HandlerError) {
	credentials := baasapi.AzureCredentials{
		ApplicationID:     payload.AzureApplicationID,
		TenantID:          payload.AzureTenantID,
		AuthenticationKey: payload.AzureAuthenticationKey,
	}

	httpClient := client.NewHTTPClient()
	_, err := httpClient.ExecuteAzureAuthenticationRequest(&credentials)
	if err != nil {
		return nil, &httperror.HandlerError{http.StatusInternalServerError, "Unable to authenticate against Azure", err}
	}

	endpointID := handler.EndpointService.GetNextIdentifier()
	endpoint := &baasapi.Endpoint{
		ID:               baasapi.EndpointID(endpointID),
		Name:             payload.Name,
		URL:              "https://management.azure.com",
		Type:             baasapi.AzureEnvironment,
		GroupID:          baasapi.EndpointGroupID(payload.GroupID),
		PublicURL:        payload.PublicURL,
		AuthorizedUsers:  []baasapi.UserID{},
		AuthorizedTeams:  []baasapi.TeamID{},
		Extensions:       []baasapi.EndpointExtension{},
		AzureCredentials: credentials,
		Tags:             payload.Tags,
		Status:           baasapi.EndpointStatusUp,
		Snapshots:        []baasapi.Snapshot{},
	}

	err = handler.EndpointService.CreateEndpoint(endpoint)
	if err != nil {
		return nil, &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist endpoint inside the database", err}
	}

	return endpoint, nil
}

func (handler *Handler) createUnsecuredEndpoint(payload *endpointCreatePayload) (*baasapi.Endpoint, *httperror.HandlerError) {
	endpointType := baasapi.DockerEnvironment

	if payload.URL == "" {
		payload.URL = "unix:///var/run/docker.sock"
		if runtime.GOOS == "windows" {
			payload.URL = "npipe:////./pipe/docker_engine"
		}
	} else {
		agentOnDockerEnvironment, err := client.ExecutePingOperation(payload.URL, nil)
		if err != nil {
			return nil, &httperror.HandlerError{http.StatusInternalServerError, "Unable to ping Docker environment", err}
		}
		if agentOnDockerEnvironment {
			endpointType = baasapi.AgentOnDockerEnvironment
		}
	}

	endpointID := handler.EndpointService.GetNextIdentifier()
	endpoint := &baasapi.Endpoint{
		ID:        baasapi.EndpointID(endpointID),
		Name:      payload.Name,
		URL:       payload.URL,
		Type:      endpointType,
		GroupID:   baasapi.EndpointGroupID(payload.GroupID),
		PublicURL: payload.PublicURL,
		TLSConfig: baasapi.TLSConfiguration{
			TLS: false,
		},
		AuthorizedUsers: []baasapi.UserID{},
		AuthorizedTeams: []baasapi.TeamID{},
		Extensions:      []baasapi.EndpointExtension{},
		Tags:            payload.Tags,
		Status:          baasapi.EndpointStatusUp,
		Snapshots:       []baasapi.Snapshot{},
	}

	err := handler.snapshotAndPersistEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}

func (handler *Handler) createTLSSecuredEndpoint(payload *endpointCreatePayload) (*baasapi.Endpoint, *httperror.HandlerError) {
	tlsConfig, err := crypto.CreateTLSConfigurationFromBytes(payload.TLSCACertFile, payload.TLSCertFile, payload.TLSKeyFile, payload.TLSSkipClientVerify, payload.TLSSkipVerify)
	if err != nil {
		return nil, &httperror.HandlerError{http.StatusInternalServerError, "Unable to create TLS configuration", err}
	}

	agentOnDockerEnvironment, err := client.ExecutePingOperation(payload.URL, tlsConfig)
	if err != nil {
		return nil, &httperror.HandlerError{http.StatusInternalServerError, "Unable to ping Docker environment", err}
	}

	endpointType := baasapi.DockerEnvironment
	if agentOnDockerEnvironment {
		endpointType = baasapi.AgentOnDockerEnvironment
	}

	endpointID := handler.EndpointService.GetNextIdentifier()
	endpoint := &baasapi.Endpoint{
		ID:        baasapi.EndpointID(endpointID),
		Name:      payload.Name,
		URL:       payload.URL,
		Type:      endpointType,
		GroupID:   baasapi.EndpointGroupID(payload.GroupID),
		PublicURL: payload.PublicURL,
		TLSConfig: baasapi.TLSConfiguration{
			TLS:           payload.TLS,
			TLSSkipVerify: payload.TLSSkipVerify,
		},
		AuthorizedUsers: []baasapi.UserID{},
		AuthorizedTeams: []baasapi.TeamID{},
		Extensions:      []baasapi.EndpointExtension{},
		Tags:            payload.Tags,
		Status:          baasapi.EndpointStatusUp,
		Snapshots:       []baasapi.Snapshot{},
	}

	filesystemError := handler.storeTLSFiles(endpoint, payload)
	if err != nil {
		return nil, filesystemError
	}

	endpointCreationError := handler.snapshotAndPersistEndpoint(endpoint)
	if endpointCreationError != nil {
		return nil, endpointCreationError
	}

	return endpoint, nil
}

func (handler *Handler) snapshotAndPersistEndpoint(endpoint *baasapi.Endpoint) *httperror.HandlerError {
	snapshot, err := handler.Snapshotter.CreateSnapshot(endpoint)
	endpoint.Status = baasapi.EndpointStatusUp
	if err != nil {
		log.Printf("http error: endpoint snapshot error (endpoint=%s, URL=%s) (err=%s)\n", endpoint.Name, endpoint.URL, err)
		endpoint.Status = baasapi.EndpointStatusDown
	}

	if snapshot != nil {
		endpoint.Snapshots = []baasapi.Snapshot{*snapshot}
	}

	err = handler.EndpointService.CreateEndpoint(endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist endpoint inside the database", err}
	}

	return nil
}

func (handler *Handler) storeTLSFiles(endpoint *baasapi.Endpoint, payload *endpointCreatePayload) *httperror.HandlerError {
	folder := strconv.Itoa(int(endpoint.ID))

	if !payload.TLSSkipVerify {
		caCertPath, err := handler.FileService.StoreTLSFileFromBytes(folder, baasapi.TLSFileCA, payload.TLSCACertFile)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist TLS CA certificate file on disk", err}
		}
		endpoint.TLSConfig.TLSCACertPath = caCertPath
	}

	if !payload.TLSSkipClientVerify {
		certPath, err := handler.FileService.StoreTLSFileFromBytes(folder, baasapi.TLSFileCert, payload.TLSCertFile)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist TLS certificate file on disk", err}
		}
		endpoint.TLSConfig.TLSCertPath = certPath

		keyPath, err := handler.FileService.StoreTLSFileFromBytes(folder, baasapi.TLSFileKey, payload.TLSKeyFile)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist TLS key file on disk", err}
		}
		endpoint.TLSConfig.TLSKeyPath = keyPath
	}

	return nil
}
