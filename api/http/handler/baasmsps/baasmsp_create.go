package baasmsps

import (
	"log"
	"net/http"
	//"runtime"
	//"strconv"
	"math/rand"
	"time"
	"reflect"
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
	BaaSDeploymentPath = "baasdeployment"
)

type baasmspCreatePayload struct {
	Name                   string
	URL                    string
	BaasmspType            int
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

	
	var orgnuml = 1;
	var orgnumr = 5;
	//var peernuml = 1;
	//var peernumr = 5;
	//var orderernuml = 1;
	//var orderernumr = 5;
	//var kafkanuml = 4;
	//var kafkanumr = 7;
	//var zookeepernuml = 3;
	//var zookeepernumr = 7;
	
	log.Printf("http success: baask8s snapshot error (baask8s=%s) (passowrd=) \n", reflect.TypeOf(payload.Orgnum).Kind())
	log.Printf("http success: baask8s snapshot error (baask8s=%s) (passowrd=) \n", reflect.TypeOf(orgnumr).Kind())
	log.Printf("(baask8s=%s) (passowrd=) \n", payload.Orgnum)
	log.Printf("(baask8s=%s) (passowrd=) \n", orgnuml)
	log.Printf("(baask8s=%s) (passowrd=) \n", orgnumr)
	log.Printf("(baask8s=%s) (passowrd=) \n", payload)

	if govalidator.IsNull(payload.Networkname) {
		return baasapi.Error("Invalid network name")
	}
	if govalidator.IsNull(payload.Platform) {
		return baasapi.Error("Invalid platform")
	}
	if govalidator.IsNull(payload.Owner) {
		return baasapi.Error("Invalid network name")
	}
	if govalidator.IsNull(payload.Otype) {
		return baasapi.Error("Invalid orderer type")
	}
	if !govalidator.InRangeInt(payload.Orgnum, orgnuml, orgnumr) {
		return baasapi.Error("Invalid organization number")
	}
	if !govalidator.InRangeInt(payload.Peernum, 1, 5 ) {
		return baasapi.Error("Invalid peer nodes number")
	}
	if !govalidator.InRangeInt(payload.Orderernum, 1, 5 ) {
		return baasapi.Error("Invalid orderer nodes number")
	}
	if !govalidator.InRangeInt(payload.Kafkanum, 4, 7) {
		return baasapi.Error("Invalid kafka nodes number")
	}
	if !govalidator.InRangeInt(payload.Zookeepernum, 3, 7) {
		return baasapi.Error("Invalid zookeeper nodes number")
	}
	return nil
	
}

// POST request on /api/baask8ss
func (handler *Handler) baasmspCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	//if handler.authDisabled {
	//	return &httperror.HandlerError{http.StatusServiceUnavailable, "Cannot authenticate user. BaaSapi was started with the --no-auth flag", ErrAuthDisabled}
	//}

	//payload := &baasmspCreatePayload{}

	var payload authenticatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	//username, err := payload.Username;
	//if err != nil {
	//	log.Printf("http error: baask8s snapshot error (baask8s=%s)   (err=%s)\n", username, err)
	//}

	log.Printf("http success: baask8s snapshot error (baask8s=%s) (passowrd=%s) \n", payload.Owner, payload.Networkname)

	
	//stackStorePath := path.Join(BaaSDeploymentPath, "networkid")
	//err = handler.FileService.CreateDirectoryInStore(stackStorePath)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusBadRequest, "Cannot Create the Directory request from payload", err}
	//}

	//stackFolder := strconv.Itoa(int(stack.ID))
	rand.Seed(time.Now().UnixNano())
	networkID := randomString(64)

	var namespace = payload.Owner + "-" + networkID[0:13]
	
	projectPath, err := handler.FileService.StoreYamlFileFromJSON(namespace, "input.config", payload, payload.Owner)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist Compose file on disk", err}
	}

	log.Printf("http success: baask8s snapshot error (baask8s=%s) (passowrd=%s) \n", projectPath, payload.Owner)


	err = handler.CAFilesManager.Deploy(payload.Owner, networkID)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to execute cafiles commands ", err}
	}

	//tmpPath := path.Join(handler.FileService.fileStorePath, stackStorePath)
	//configFilePath := path.Join(tmpPath, "input.json")
	
	//err = handler.FileService.WriteJSONToFile(configFilePath, payload)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist Yaml file on disk", err}
	//}

	// TODO: add logic to install blockchain network using helm chart. 

	// TODO: if no helm tool installed, then do with manually

	





	var NetworkName = payload.Networkname
	var Platform = payload.Platform
	var CreatedAt = time.Now().Unix()
	//var Tags = ["baas1.4"]

	//baasmspType := baasapi.DockerEnvironment2
	baasmspID := handler.BaasmspService.GetNextIdentifier()
	baasmsp := &baasapi.Baasmsp{
		ID:               baasapi.BaasmspID(baasmspID),
		NetworkName:      NetworkName,
		NetworkID:        networkID,
		Owner:            payload.Owner,
		Platform:         Platform,
		CreatedAt:        CreatedAt,
		MSPs:            []baasapi.MSP{},
		//Tags:          Tags,
		Status:          1,
		Namespace:       "default",
		Snapshots:       []baasapi.Snapshot{},
	}

    err = handler.BaasmspService.CreateBaasmsp(baasmsp)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baasmsp inside the database", err}
	}


	
	//stack.ProjectPath = projectPath

	//doCleanUp := true
	//defer handler.cleanUp(stack, &doCleanUp)


	//baasmsp, baasmspCreationError := handler.createBaasmsp(payload)
	//if baasmspCreationError != nil {
	//	return baasmspCreationError
	//}

	return response.JSON(w, baasmsp)

	//return response.JSON(w, baasmsp)
}

func (handler *Handler) createBaasmsp(payload *baasmspCreatePayload) (*baasapi.Baasmsp, *httperror.HandlerError) {

	return handler.createUnsecuredBaasmsp(payload)
}

func (handler *Handler) createUnsecuredBaasmsp(payload *baasmspCreatePayload) (*baasapi.Baasmsp, *httperror.HandlerError) {


	//err := handler.snapshotAndPersistBaasmsp(baasmsp)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil

	//return baasmsp, nil
}

//func (handler *Handler) snapshotAndPersistBaasmsp(baasmsp *baasapi.Baasmsp) *httperror.HandlerError {
	//snapshot, err := handler.Snapshotter.CreateSnapshot(baask8s)
	//baask8s.Status = baasapi.Baask8sStatusUp
	//if err != nil {
	//	log.Printf("http error: baask8s snapshot error (baask8s=%s, URL=%s) (err=%s)\n", baask8s.Name, baask8s.URL, err)
	//	baask8s.Status = baasapi.Baask8sStatusDown
	//}

	//if snapshot != nil {
	//	baask8s.Snapshots = []baasapi.Snapshot{*snapshot}
	//}

	//err = handler.Baask8sService.CreateBaask8s(baask8s)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s inside the database", err}
	//}

//	return nil
//}

//func (handler *Handler) storeTLSFiles(baask8s *baasapi.Baask8s, payload *baask8sCreatePayload) *httperror.HandlerError {
//	folder := strconv.Itoa(int(baask8s.ID))

//	if !payload.TLSSkipVerify {
//		caCertPath, err := handler.FileService.StoreTLSFileFromBytes(folder, baasapi.TLSFileCA, payload.TLSCACertFile)
//		if err != nil {
//			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist TLS CA certificate file on disk", err}
//		}
//		baask8s.TLSConfig.TLSCACertPath = caCertPath
//	}
//
//	if !payload.TLSSkipClientVerify {
//		certPath, err := handler.FileService.StoreTLSFileFromBytes(folder, baasapi.TLSFileCert, payload.TLSCertFile)
//		if err != nil {
//			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist TLS certificate file on disk", err}
//		}
//		baask8s.TLSConfig.TLSCertPath = certPath
//
//		keyPath, err := handler.FileService.StoreTLSFileFromBytes(folder, baasapi.TLSFileKey, payload.TLSKeyFile)
//		if err != nil {
//			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist TLS key file on disk", err}
//		}
//		baask8s.TLSConfig.TLSKeyPath = keyPath
//	}

//	return nil
//}
