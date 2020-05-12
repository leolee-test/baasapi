package baask8ss

import (
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

const (
	// ErrBaask8sManagementDisabled is an error raised when trying to access the baask8ss management baask8ss
	// when the server has been started with the --external-baask8ss flag
	Errbaask8sManagementDisabled = baasapi.Error("baask8s management is disabled")
)

//func hideFields(baask8s *baasapi.Baask8s) {
//	baask8s.AzureCredentials = baasapi.AzureCredentials{}
//}

// Handler is the HTTP handler used to handle baask8s operations.
type Handler struct {
	*mux.Router
	requestBouncer              *security.RequestBouncer
	Baask8sService              baasapi.Baask8sService
	BaasmspService              baasapi.BaasmspService
	Baask8sGroupService        baasapi.Baask8sGroupService
	FileService                 baasapi.FileService
	JWTService                  baasapi.JWTService
	Snapshotter                 baasapi.Snapshotter
	CAFilesManager              baasapi.CAFilesManager
	JobService                  baasapi.JobService
}

// NewHandler creates a handler to manage baask8s operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
		requestBouncer:              bouncer,
	}

	h.Handle("/baask8ss",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sCreate))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{namespace}/log/{nline}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sCreateLog))).Methods(http.MethodGet)
	h.Handle("/baask8ss/settings",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sSettings))).Methods(http.MethodPost)
	h.Handle("/baask8ss/settings",
		bouncer.PublicAccess(httperror.LoggerHandler(h.getbaask8sSettings))).Methods(http.MethodGet)
	//h.Handle("/baask8ss/snapshot",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sSnapshots))).Methods(http.MethodPost)
	h.Handle("/baask8ss",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/backendpods/{namespace}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.backendPodsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/pods",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sPodsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/podsvcs",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sPodSVCsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/statefulsets",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sStatefulSetsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/daemonsets",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sDaemonSetsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/token/{token}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sTokenValidate))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/deployments",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sDeploymentsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/pods",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sPodOperations))).Methods(http.MethodPost)
	h.Handle("/baask8ss/backend/backendpods",
		bouncer.PublicAccess(httperror.LoggerHandler(h.backendPodOperations))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/scale",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sPatchScale))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/msps",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sMspsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/policys",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sPolicysList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/policys",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sPolicysAdd))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/policys",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sPolicysUpdate))).Methods(http.MethodPut)
	h.Handle("/baask8ss/{id}/policys/{name}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sPolicysDelete))).Methods(http.MethodDelete)
	h.Handle("/baask8ss/{id}/svcs",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sDeploySvcs))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/svcs/{srvname}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sDeleteSvcs))).Methods(http.MethodDelete)
	h.Handle("/baask8ss/{id}/svcs",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sGetSvcs))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/msps",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sUpdateMSPs))).Methods(http.MethodPost)
	//h.Handle("/baask8ss/{id}/channels",
	//	bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sCHLsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/ccs",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sCCsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sListByID))).Methods(http.MethodGet)
	//h.Handle("/baask8ss/{id}",
	//	bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sInspect))).Methods(http.MethodGet)
	//h.Handle("/baask8ss/{id}",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sUpdate))).Methods(http.MethodPut)
	h.Handle("/baask8ss/{id}/access",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.baask8sUpdateAccess))).Methods(http.MethodPut)
	h.Handle("/baask8ss/{id}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sDelete))).Methods(http.MethodDelete)
	h.Handle("/baask8ss/baasonly/{id}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sDeleteBaasOnly))).Methods(http.MethodDelete)
	h.Handle("/baask8ss/{id}/channels",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChannelsList))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/channels/{channelname}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChannelsListByName))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/channels/sync/{channelname}/{currentid}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChannelsListByNameSync))).Methods(http.MethodGet)
	//h.Handle("/baask8ss/{id}/enroll",
	//	bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sEnroll))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/channels",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChannelsCreate))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/chaincodesbypeer",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8schaincodesbyPeer))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/channelsbypeer",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8schannelsbyPeer))).Methods(http.MethodGet)
	h.Handle("/baask8ss/{id}/channelinfobypeer",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8schannelinfobyPeer))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/channels/{channelname}/peers",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChannelsJoin))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/channels/{channelname}/anchorpeers",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChannelsAnchor))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/channels/{channelname}/chaincodes",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChaincodesInstantiate))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/channels/{channelname}/chaincodes/{ccname}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChaincodesInvoke))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/channels/{channelname}/chaincodes/{ccname}/query",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChaincodesQuery))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/chaincodes",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sChaincodesInstall))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/channels/{channelname}/blocks/{blocknum}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8sBlocksQuery))).Methods(http.MethodGet)
	//h.Handle("/baask8ss/{id}/extensions",
	//	bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.baask8sExtensionAdd))).Methods(http.MethodPost)
	//h.Handle("/baask8ss/{id}/extensions/{extensionType}",
	//	bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.baask8sExtensionRemove))).Methods(http.MethodDelete)
	//h.Handle("/baask8ss/{id}/job",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sJob))).Methods(http.MethodPost)
	//h.Handle("/baask8ss/{id}/snapshot",
	//	bouncer.AdministratorAccess(httperror.LoggerHandler(h.baask8sSnapshot))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/ca/users",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersRegister))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/ca/enroll",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersEnroll))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/ca/reenroll",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersreEnroll))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/ca/listusers",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersList))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/ca/users/{orgname}/{username}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersUpdate))).Methods(http.MethodPut)
	h.Handle("/baask8ss/{id}/ca/users/{orgname}/{username}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersDelete))).Methods(http.MethodDelete)
	h.Handle("/baask8ss/{id}/ca/revoke",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersRevoke))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/ca/users/{orgname}/{username}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersListbyName))).Methods(http.MethodPost)
	h.Handle("/baask8ss/{id}/ca/users/CRL",
		bouncer.PublicAccess(httperror.LoggerHandler(h.baask8susersGenerateCRL))).Methods(http.MethodPost)
	return h
}
