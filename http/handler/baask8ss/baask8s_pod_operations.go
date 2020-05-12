package baask8ss

import (
	"net/http"

	//"flag"
	//"fmt"
	//"log"
	//"os"
	"bytes"
	"io"
	"path/filepath"
	//"encoding/json"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/apimachinery/pkg/types"
	//"k8s.io/client-go/rest"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	//"github.com/baasapi/baasapi/api/http/security"
)

type podOperationsPayload struct {
	PodName    string  `json:"podname"`
	Namespace  string  `json:"namespace"`
	Action     string  `json:"action"`
	Nline      int64     `json:"nline"`
	
}

func (payload *podOperationsPayload) Validate(r *http.Request) error {
	return nil
}


// GET request on /api/baask8ss
func (handler *Handler) baask8sPodOperations(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	var payload podOperationsPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
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



	// Bootstrap k8s configuration from local 	Kubernetes config file
	//kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	kubeconfig := filepath.Join("/data/k8s/ansible/", "vars", "kubeconfig")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to using kubeconfig file", err}
	}

	// Create an rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to create an rest client to kubernetes API", err}
	}

	//var str = ""
    // todo: check if namespace exist???
	buf := new(bytes.Buffer)
	if payload.Action == "Log" {

		//podLogOpts := coreV1.PodLogOptions{}
		//podLogOpts := &v1.PodLogOptions{}

		//var line         = int64(8)
		//&v1.PodLogOptions{TailLines: &line}
		var (
			line         = payload.Nline
		//	line         = int64(50)
		//	bytes        = int64(64)
		//	timestamp    = metav1.Now()
		//	sinceseconds = int64(10)
		)



		//logreq := clientset.CoreV1().Pods("default").GetLogs(payload.PodName, &v1.PodLogOptions{TailLines: &line})
		logreq := clientset.CoreV1().Pods(baask8s.Namespace).GetLogs(payload.PodName, &v1.PodLogOptions{TailLines: &line})
		
		//logreq := clientset.CoreV1().Pods(baask8s.Namespace).GetLogs(podOperationsPayload.PodName, {})
		//&v1.PodLogOptions{}		

		podLogs, err := logreq.Stream()
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Error in opening stream from kubernetes API", err}
		}
		defer podLogs.Close()
	
		
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Error in copy information from podLogs to buf from kubernetes API", err}
		}
		//str := buf.String()

		return response.JSON(w, buf.String())
	}

	if payload.Action == "Delete" {

		//podLogOpts := coreV1.PodLogOptions{}
		//podLogOpts := &v1.PodLogOptions{}
		err = clientset.CoreV1().Pods(baask8s.Namespace).Delete(payload.PodName, &metav1.DeleteOptions{})
		//err = clientset.CoreV1().Pods("default").Delete(payload.PodName, &metav1.DeleteOptions{})
		//logreq := clientset.CoreV1().Pods(baask8s.Namespace).GetLogs(podOperationsPayload.PodName, {})
		//name string, options *metav1.DeleteOptions
		//&v1.PodLogOptions{}	
		
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Error in deleting pod from kubernetes API", err}
		 }
		
		var deleteresult = "pod deleted"
        return response.JSON(w, deleteresult)
	}

	return response.JSON(w, buf)
}

// GET request on /api/baask8ss
func (handler *Handler) backendPodOperations(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	var payload podOperationsPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}


	// Bootstrap k8s configuration from local 	Kubernetes config file
	//kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	kubeconfig := filepath.Join("/data/k8s/ansible/", "vars", "kubeconfig")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to using kubeconfig file", err}
	}

	// Create an rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to create an rest client to kubernetes API", err}
	}

	//var str = ""

    // todo: check if namespace exist???
	buf := new(bytes.Buffer)
	if payload.Action == "Log" {

		//podLogOpts := coreV1.PodLogOptions{}
		//podLogOpts := &v1.PodLogOptions{}

		//var line         = int64(8)
		//&v1.PodLogOptions{TailLines: &line}
		var (
			line         = payload.Nline
		//	line         = int64(50)
		//	bytes        = int64(64)
		//	timestamp    = metav1.Now()
		//	sinceseconds = int64(10)
		)



		//logreq := clientset.CoreV1().Pods("default").GetLogs(payload.PodName, &v1.PodLogOptions{TailLines: &line})
		logreq := clientset.CoreV1().Pods(payload.Namespace).GetLogs(payload.PodName, &v1.PodLogOptions{TailLines: &line})
		
		//logreq := clientset.CoreV1().Pods(baask8s.Namespace).GetLogs(podOperationsPayload.PodName, {})
		//&v1.PodLogOptions{}		

		podLogs, err := logreq.Stream()
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Error in opening stream from kubernetes API", err}
		}
		defer podLogs.Close()
	
		
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Error in copy information from podLogs to buf from kubernetes API", err}
		}
		//str := buf.String()
		//log.Println("Using kubeconfig file: ", buf.String())

		return response.JSON(w, buf.String())
	}

	if payload.Action == "Delete" {

		//podLogOpts := coreV1.PodLogOptions{}
		//podLogOpts := &v1.PodLogOptions{}
		err = clientset.CoreV1().Pods(payload.Namespace).Delete(payload.PodName, &metav1.DeleteOptions{})
		//err = clientset.CoreV1().Pods("default").Delete(payload.PodName, &metav1.DeleteOptions{})
		//logreq := clientset.CoreV1().Pods(baask8s.Namespace).GetLogs(podOperationsPayload.PodName, {})
		//name string, options *metav1.DeleteOptions
		//&v1.PodLogOptions{}	
		
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Error in deleting pod from kubernetes API", err}
		 }
		
		var deleteresult = "pod deleted"
        return response.JSON(w, deleteresult)
	}

	return response.JSON(w, buf)
}

