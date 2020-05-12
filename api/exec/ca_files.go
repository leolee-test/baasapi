package exec

import (
	"log"
	"bytes"
	//"encoding/json"
	"os"
	"os/exec"
	"path"
	"strconv"
	//"runtime"

	"github.com/baasapi/baasapi/api"
)

const (
	// Baas deployment files
	BaaSDeploymentPath = "k8s/ansible/vars/namespaces"
	BinaryStorePath = "bin"
)

// CAFilesManager represents a service for managing cafiles.
type CAFilesManager struct {
	binaryPath       string
	dataPath         string
	fileService      baasapi.FileService
}

// NewCAFilesManager initializes a new CAFilesManager service.
func NewCAFilesManager(binaryPath, dataPath string, fileService baasapi.FileService) (*CAFilesManager, error) {
	manager := &CAFilesManager{
		binaryPath:       binaryPath,
		dataPath:         dataPath,
		fileService:      fileService,
	}

	err := manager.updateCAFilesCLIConfiguration(dataPath)
	if err != nil {
		return nil, err
	}

	return manager, nil
}

// Deploy executes the cafiles deploy command.
// Deploy(baask8s *Baask8s, creator string, networkID string) error
func (manager *CAFilesManager) Deploy(creator string, namespace string, ansible_extra string, ansible_env string, ansible_config string, flag bool) error {
	//stackFilePath := path.Join(BaaSDeploymentPath, creator, creator+"-"+networkID[0:13])
	stackFilePath := path.Join(BaaSDeploymentPath, namespace)

    
	log.Printf("!!!!!!!!!http success(baask8s=%s) (passowrd=%s) \n", stackFilePath, "config.jsondeploy")

	//bin_version := "v1.4"
	//command, args := prepareCACommandAndArgs(manager.binaryPath+"/"+manager.dataPath+"/"+BinaryStorePath+"/"+bin_version, stackFilePath)
    command, args := prepareCACommandAndArgs(BaaSDeploymentPath+"/../../", stackFilePath, ansible_extra, ansible_env, ansible_config)

	env := make([]string, 0)
	//env = append(env, "ANSIBLE_LOG_PATH=./vars/namespaces/logs/"+namespace+"/"+namespace+".log")
	env = append(env, "ANSIBLE_LOG_PATH=./vars/namespaces/logs/"+namespace+".log")
	//for _, envvar := range stack.Env {
	//	env = append(env, envvar.Name+"="+envvar.Value)
	//}
	//env := ""
	//return runCACommandAndCaptureStdErr(command, args, env, manager.binaryPath+"/"+manager.dataPath+"/k8s/ansible/", flag)
	return runCACommandAndCaptureStdErr(command, args, env, "/data/k8s/ansible/", flag)
	//stackFolder := path.Dir(stackFilePath)
	//return runCACommandAndCaptureStdErr(command, args, env, stackFolder)
}

func (manager *CAFilesManager) GetLogs(namespace string, nline int) (string, error) {
	//stackFilePath := path.Join(BaaSDeploymentPath, creator, creator+"-"+networkID[0:13])
	//stackFilePath := path.Join(BaaSDeploymentPath, namespace)

	//bin_version := "v1.4"
	//command, args := prepareCACommandAndArgs(manager.binaryPath+"/"+manager.dataPath+"/"+BinaryStorePath+"/"+bin_version, stackFilePath)
    command, args := prepareLogCommandAndArgs(BaaSDeploymentPath+"/../../", BaaSDeploymentPath, namespace, nline)

	//env := make([]string, 0)
	//env = append(env, "ANSIBLE_LOG_PATH=./vars/namespaces/logs/"+namespace+"/"+namespace+".log")
	//env = append(env, "ANSIBLE_LOG_PATH=./vars/namespaces/logs/"+namespace+".log")
	//for _, envvar := range stack.Env {
	//	env = append(env, envvar.Name+"="+envvar.Value)
	//}
	//env := ""
	//return runCACommandAndCaptureStdErr(command, args, env, manager.binaryPath+"/"+manager.dataPath+"/k8s/ansible/", flag)
	return runLogCommandAndCaptureStdErr(command, args, "/data/",namespace)


	//stackFolder := path.Dir(stackFilePath)
	//return runCACommandAndCaptureStdErr(command, args, env, stackFolder)
}

// Remove executes the docker stack rm command.
//func (manager *SwarmStackManager) Remove(stack *baasapi.Stack, baask8s *baasapi.Baask8s) error {
//	command, args := prepareDockerCommandAndArgs(manager.binaryPath, manager.dataPath, baask8s)
//	args = append(args, "stack", "rm", stack.Name)
//	return runCommandAndCaptureStdErr(command, args, nil, "")
//}

func runCACommandAndCaptureStdErr(command string, args []string, env []string, workingDir string, flag bool) error {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	log.Printf("command%s) (args=%s) (env=%s) \n", command, args, env)

	log.Printf("http success: baask8s snapshot error (baask8s=%s) (passowrd=%s) \n", workingDir, "config.json22222run")

	cmd.Dir = workingDir

	if env != nil {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, env...)
	}

	if flag {
		err := cmd.Run()
		if err != nil {
			return baasapi.Error(stderr.String()+stdout.String())
		}
	} else {
		err := cmd.Start()
		if err != nil {
			return baasapi.Error(stderr.String()+stdout.String())
		}
	}
	//err := cmd.Run()
	


	return nil
}

func runLogCommandAndCaptureStdErr(command string, args []string, workingDir string, namespace string) (string, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	log.Printf("command%s) (args=%s) (env=%s) \n", command, args, namespace)

	log.Printf("http success: baask8s snapshot error (baask8s=%s) (passowrd=%s) \n", workingDir, "config.json22222run")

	cmd.Dir = workingDir

		err := cmd.Run()
		if err != nil {
			return "Error", baasapi.Error(stderr.String()+stdout.String())
		}
	//err := cmd.Run()
	


	return stdout.String(), nil
}

func prepareCACommandAndArgs(binaryPath, dataPath, ansible_extra, ansible_env, ansible_config string) (string, []string) {
	// Assume Linux as a default
	command := path.Join("", "ansible-playbook")

	args := make([]string, 0)
	//args = append(args, "generate")
	args = append(args, ansible_config, "-e", ansible_env)
	args = append(args, "-e", ansible_extra)
	//args = append(args, "--config", dataPath)
	//args = append(args, "-H", baask8s.URL)

	log.Printf("http success: endpreparecacommandpoint snapshot error (baask8s=%s) (passowrd=%s) \n", command, args)

	return command, args
}

func prepareLogCommandAndArgs(binaryPath, dataPath, namespace string, nline int) (string, []string) {
	// Assume Linux as a default
	command := path.Join("", "tail")

	args := make([]string, 0)
	//args = append(args, "generate")
	args = append(args, "-n", strconv.Itoa(nline))
	args = append(args, "/data/k8s/ansible/vars/namespaces/logs/"+namespace+".log")
	//args = append(args, dataPath+"/logs/"+namespace+".log")
	//args = append(args, "-H", baask8s.URL)

	log.Printf("http success: endpreparecacommandpoint snapshot error (baask8s=%s) (passowrd=%s) \n", command, args)

	return command, args
}

func (manager *CAFilesManager) updateCAFilesCLIConfiguration(dataPath string) error {
	configFilePath := path.Join(dataPath, "config.json")

	log.Printf("http success: baask8s snapshot error (baask8s=%s) (passowrd=%s) \n", configFilePath, "config.jsonupdate")
	//config, err := manager.retrieveConfigurationFromDisk(configFilePath)
	//if err != nil {
	//	return err
	//}

	//if config["HttpHeaders"] == nil {
	//	config["HttpHeaders"] = make(map[string]interface{})
	//}
	//headersObject := config["HttpHeaders"].(map[string]interface{})
	//headersObject["X-BaaSapiAgent-ManagerOperation"] = "1"
	//headersObject["X-BaaSapiAgent-Signature"] = signature
	//headersObject["X-BaaSapiAgent-PublicKey"] = manager.signatureService.EncodedPublicKey()

	//err = manager.fileService.WriteJSONToFile(configFilePath, config)
	//if err != nil {
	//	return err
	//}

	return nil
}

//func (manager *SwarmStackManager) retrieveConfigurationFromDisk(path string) (map[string]interface{}, error) {
//	var config map[string]interface{}
//
//	raw, err := manager.fileService.GetFileContent(path)
//	if err != nil {
//		return make(map[string]interface{}), nil
//	}
//
//	err = json.Unmarshal(raw, &config)
//	if err != nil {
//		return nil, err
//	}

//	return config, nil
//}
