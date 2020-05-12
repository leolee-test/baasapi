package baask8ss

import (
	//"log"
	"net/http"
	//"runtime"
	//"strconv"
	"math/rand"
	"time"
	//"reflect"
	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	//"github.com/baasapi/baasapi/api/crypto"
	//"github.com/baasapi/baasapi/api/http/client"
)

const (
	// Baas deployment files
	BaaSDeploymentPath = "k8s/ansible/vars/namespaces"
)

type Network struct {
	Cas       []string
	Peers     []string
	Orderers  []string
	Zookeepers []string
	Kafkas     []string
}

type authenticatePayload struct {
	Version     string
	Networkname string
	Platform    string
	Owner       string
	Otype       string
	Orgnum      int
	Peernum     int
	Orderernum  int
	Kafkanum    int
	Zookeepernum int 
	Raftnum     int
	Peer_db     string
	Tls         string
	Logging_level string
	Networks     map[string]Network
	MSPs        []baasapi.MSP
	CHLs        []baasapi.CHL
	CCs         []baasapi.CC
	//Fabric struct{} `json:"fabric"`
	//Fabric    string `json:"fabric"`
}

func randomString(l int) string {
	//var pool ="0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz.-+_"
	//var pool ="0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz"
	var pool ="0123456789abcdefghijklmnopqrstuvwxyz"
    bytes := make([]byte, l)
    for i := 0; i < l; i++ {
        bytes[i] = pool[rand.Intn(len(pool))]
    }
    return string(bytes)
}

func (payload *authenticatePayload) Validate(r *http.Request) error {

	//username, err := request.RetrieveMultiPartFormValue(r, "Username", false)
	//log.Printf("http error: baask8s snapshot error (baask8s=%s)   (err=%s)\n", username, err)

	
	//var orgnuml = 1;
	//var orgnumr = 5;
	//var peernuml = 1;
	//var peernumr = 5;
	//var orderernuml = 1;
	//var orderernumr = 5;
	//var kafkanuml = 4;
	//var kafkanumr = 7;
	//var zookeepernuml = 3;
	//var zookeepernumr = 7;

	if govalidator.IsNull(payload.Networkname) {
		return baasapi.Error("Invalid network name")
	}
	if govalidator.IsNull(payload.Platform) {
		return baasapi.Error("Invalid platform")
	}
	if govalidator.IsNull(payload.Owner) {
		return baasapi.Error("Invalid Owner name")
	}
	//if govalidator.IsNull(payload.Networks) {
	//	return baasapi.Error("Invalid networks object")
	//}
	if govalidator.IsNull(payload.Otype) {
		return baasapi.Error("Invalid orderer type")
	}
	//if !govalidator.InRangeInt(payload.Orgnum, orgnuml, orgnumr) {
	//	return baasapi.Error("Invalid organization number")
	//}
	if !govalidator.InRangeInt(payload.Peernum, 1, 5 ) {
		return baasapi.Error("Invalid peer nodes number")
	}
	if !govalidator.InRangeInt(payload.Orderernum, 1, 5 ) {
		return baasapi.Error("Invalid orderer nodes number")
	}
	//if !govalidator.InRangeInt(payload.Kafkanum, 4, 7) {
	//	return baasapi.Error("Invalid kafka nodes number")
	//}
	//if !govalidator.InRangeInt(payload.Zookeepernum, 3, 7) {
	//	return baasapi.Error("Invalid zookeeper nodes number")
	//}
	//if !govalidator.InRangeInt(payload.Raftnum, 3, 7) {
	//	return baasapi.Error("Invalid raft nodes number")
	//}
	return nil
	
}

// POST request on /api/baask8ss
func (handler *Handler) baask8sCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	//if handler.authDisabled {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Cannot authenticate user. BaaSapi was started with the --no-auth flag", ErrAuthDisabled}
	//}

	//payload := &baask8sCreatePayload{}

	var payload authenticatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}


	rand.Seed(time.Now().UnixNano())
	networkID := randomString(64)

	var namespace = payload.Owner + "-" + networkID[0:13]
	var initfilen = namespace + ".input"


	
	//projectPath, err := handler.FileService.StoreYamlFileFromJSON(namespace, initfilen, payload, payload.Owner)
	_, err = handler.FileService.StoreYamlFileFromJSON(namespace, initfilen, payload, payload.Owner)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist parameters file on disk", err}
	}



	var ansible_env = " env=" +namespace+ " deploy_type=k8s"
	var ansible_extra = " mode=apply "
	var ansible_config = "/data/k8s/ansible/initparm.yml"
	err = handler.CAFilesManager.Deploy(payload.Owner, namespace, ansible_extra, ansible_env, ansible_config, true)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute ansible commands ", err}
	}



	var NetworkName = payload.Networkname
	var Platform = payload.Platform
	//var CreatedAt = time.Now().Unix()
	var CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	//var Tags = ["baas1.4"]



	baask8sID := handler.Baask8sService.GetNextIdentifier()
	baask8s := &baasapi.Baask8s{
		ID:               baasapi.Baask8sID(baask8sID),
		NetworkName:      NetworkName,
		NetworkID:        networkID,
		Owner:            payload.Owner,
		Platform:         Platform,
		CreatedAt:        CreatedAt,
		MSPs:            payload.MSPs,
		CHLs:            payload.CHLs,
		CCs:             payload.CCs,
		Applications:    []baasapi.BaasAPP{},
		Policys:         []baasapi.Policy{},
		AuthorizedUsers: []baasapi.UserID{},
		AuthorizedTeams: []baasapi.TeamID{},
		//Tags:            Tags,
		Status:          1,
		Namespace:       namespace,
		Snapshots:       []baasapi.Snapshot{},
	}
	//baask8s.Policys = append(baask8s.Policys, policy)

    err = handler.Baask8sService.CreateBaask8s(baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s inside the database", err}
	}


	//ansible_env = "mode=apply env=" +namespace+ " deploy_type=k8s"
	ansible_config = "/data/k8s/ansible/setupfabric.yml"
	err = handler.CAFilesManager.Deploy(payload.Owner, namespace, ansible_extra, ansible_env, ansible_config, false)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute ansible commands ", err}
	}
	
	var responseObject jsonResponse
	responseObject.Success = true
	responseObject.Message = "Network creation job has been submitted to server successfully"
	responseObject.Namespace = namespace
	//responseObject.Message = "Not authorized or jwt token was expired"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)

	//return response.JSON(w, baask8s)

	//return response.JSON(w, baask8s)
}

