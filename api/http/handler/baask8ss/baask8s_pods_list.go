package baask8ss

import (
	"net/http"

	//"flag"
	"encoding/json"
	//"fmt"
	//"reflect"
	//"log"
	//"os"
	"strings"
	"path/filepath"

	"k8s.io/api/core/v1"
	apps "k8s.io/api/apps/v1"
	//apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/types"
	//"k8s.io/client-go/rest"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	//"github.com/baasapi/baasapi/api/http/security"
)

type Baaspods struct {
	NetworkID        string              `json:"NetworkID"`
	Owner            string              `json:"Owner"`
	Platform         string              `json:"Platform"`
	NodeName         string              `json:"NodeName"`
	PodName          string              `json:"PodName"`
	//Tags             []string            `json:"Tags"`
	Status           v1.PodStatus         `json:"Status"`
	StartTime        *metav1.Time        `json:"startTime,omitempty" protobuf:"bytes,7,opt,name=startTime"`
	UID              types.UID           `json:"UID"`
	Namespace        string              `json:"Namespace"`
	ObjectMeta       metav1.Object       `json:"ObjectMeta"`
	//Containers       []v1.Container      `json:"Containers"`
}

type Baasdeployment struct {
	Name         string              `json:"name"`
	Createdtime  metav1.Time              `json:"createdtime"`
	Replicas     int32              `json:"replicas"`
	Namespace    string              `json:"namespace"`
	Type         string              `json:"type"`

}



type scalePayload struct {
	ReplicationControllerName     string         `json:"replicationcontrollername"`
	Replicas                      uint32         `json:"replicas"`
	Type                      string         `json:"type"`
}

//  patchStringValue specifies a patch operation for a string.
type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

//  patchStringValue specifies a patch operation for a uint32.
type patchUInt32Value struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}

func (payload *scalePayload) Validate(r *http.Request) error {
	return nil;
}

func scaleReplicationController(clientSet *kubernetes.Clientset, namespace string, type_svc string, replicasetName string, scale uint32) error {
	payload := []patchUInt32Value{{
		Op:    "replace",
		Path:  "/spec/replicas",
		Value: scale,
	}}
	payloadBytes, _ := json.Marshal(payload)

	if (type_svc == "deployments"){
	_, err := clientSet.
		AppsV1().
		Deployments(namespace).
		Patch(replicasetName, types.JSONPatchType, payloadBytes)
	return err
	} else {
		if (type_svc == "statefulsets"){
			_, err := clientSet.
			AppsV1().
			StatefulSets(namespace).
			Patch(replicasetName, types.JSONPatchType, payloadBytes)
		return err
		} else {
			if (type_svc == "daemonsets"){
			_, err := clientSet.
			AppsV1().
			DaemonSets(namespace).
			Patch(replicasetName, types.JSONPatchType, payloadBytes)
			return err
			} else {
				return baasapi.ErrScaleTypeError
			}
		}
	}
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sPatchScale(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	var payload scalePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	if (payload.Replicas > 5) {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to set scale for this resource to kubernetes API", baasapi.ErrScaleNumberError}
		//return baasapi.ErrScaleNumberError
	}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

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


	//todo: check if namespace exsit?

	

	//var (
		//  Leave blank for the default context in your kube config.
		//context = ""
	
		//  Name of the replication controller to scale, and the desired number of replicas.
		//replicationControllerName = "fabricsdks1l8yivsxahhk"
		//replicas                  = uint32(0)
	//)

	err = scaleReplicationController(clientset, baask8s.Namespace, payload.Type, payload.ReplicationControllerName, payload.Replicas)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to set scale for this resource to kubernetes API", err}
	}


	var responseObject jsonResponse
	responseObject.Success = true
	responseObject.Message = "Set scale successfully to the resource"
	//responseObject.Namespace = namespace
	//responseObject.Message = "Not authorized or jwt token was expired"
	//json.Unmarshal(bodyBytes, &responseObject)
	return response.JSON(w, responseObject)
	
}




// GET request on /api/baask8ss
func (handler *Handler) baask8sPodsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	//baask8ss, err := handler.Baask8sService.Baask8ss()
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve baask8ss from the database", err}
	//}

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

	//securityContext, err := security.RetrieveRestrictedRequestContext(r)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	//}

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


	//todo: check if namespace exsit?


	pods, err := clientset.CoreV1().Pods(baask8s.Namespace).List(metav1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods("empty-23hxbvmbplj2z").List(metav1.ListOptions{})
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve pods from kubernetes API", err}
	}

	// print pods

	podsii := make([]*v1.Pod, len(pods.Items))
	for j := range pods.Items {
		podsii[j] = &pods.Items[j]
	}

	podsjj := make([]*Baaspods, len(podsii))
	for k := range podsii {
		podsjj[k]=new(Baaspods)
		podsjj[k].PodName = podsii[k].GetName()
		podsjj[k].NodeName = podsii[k].Spec.NodeName
		podsjj[k].NetworkID = baask8s.NetworkID
		podsjj[k].Owner = baask8s.Owner
		podsjj[k].Platform = baask8s.Platform
		podsjj[k].Status = podsii[k].Status
		//podsjj[k].StartTime = podsii[k].Status.StartTime
		podsjj[k].UID = podsii[k].GetUID() 
		podsjj[k].Namespace = podsii[k].GetNamespace()
		podsjj[k].ObjectMeta = podsii[k].GetObjectMeta()
		//podsjj[k].Containers = new(v1.Container)
		//podsjj[k].Containers = podsii[k].Spec.Containers

	}


	//securityContext, err := security.RetrieveRestrictedRequestContext(r)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	//}

	//filteredBaask8ss := security.FilterBaask8ss(baask8ss, securityContext)

	//for idx := range filteredBaask8ss {
	//	hideFields(&filteredBaask8ss[idx])
	//}

	//return response.JSON(w, podsii)
	return response.JSON(w, podsjj)
	
}

// GET request on /api/baask8ss/ns
//func (handler *Handler) baask8sListNS(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

//}


// GET request on /api/baask8ss
func (handler *Handler) backendPodsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	namespace, err := request.RetrieveRouteVariableValue(r, "namespace")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid namespace name variable", err}
	}


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


	//todo: check if namespace exsit?
	//log.Println(baask8s.Namespace)

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods("empty-23hxbvmbplj2z").List(metav1.ListOptions{})
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve pods from kubernetes API", err}
	}

	// print pods

	podsii := make([]*v1.Pod, len(pods.Items))
	for j := range pods.Items {
		podsii[j] = &pods.Items[j]
	}

	//podsjj := make([]*Baaspods, len(podsii))
	podsjj := []Baaspods{}
	for k := range podsii {

		strings.Contains("something", "some") // true
		if (strings.Contains(podsii[k].GetName(),"baasapi") ) {
		spod := Baaspods{}
		spod.PodName = podsii[k].GetName()
		spod.NodeName = podsii[k].Spec.NodeName
		spod.NetworkID = "backend"
		spod.Owner = "backend"
		spod.Platform = "backend"
		spod.Status = podsii[k].Status
		//podsjj[k].StartTime = podsii[k].Status.StartTime
		spod.UID = podsii[k].GetUID() 
		spod.Namespace = podsii[k].GetNamespace()
		spod.ObjectMeta = podsii[k].GetObjectMeta()

		podsjj = append(podsjj,spod)

		}
		//podsjj[k].Containers = new(v1.Container)
		//podsjj[k].Containers = podsii[k].Spec.Containers

	}



	//securityContext, err := security.RetrieveRestrictedRequestContext(r)
	//if err != nil {
	//	return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	//}

	//filteredBaask8ss := security.FilterBaask8ss(baask8ss, securityContext)

	//for idx := range filteredBaask8ss {
	//	hideFields(&filteredBaask8ss[idx])
	//}

	//return response.JSON(w, podsii)
	return response.JSON(w, podsjj)
	
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sPodSVCsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

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


	//todo: check if namespace exsit?


	services, err := clientset.CoreV1().Services(baask8s.Namespace).List(metav1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods("empty-23hxbvmbplj2z").List(metav1.ListOptions{})
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve svcs from kubernetes API", err}
	}

	// print pods



	//return response.JSON(w, podsjj)
	return response.JSON(w, services.Items)
	
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sDeploymentsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

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


	//todo: check if namespace exsit?


	deployments, err := clientset.AppsV1().Deployments(baask8s.Namespace).List(metav1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods("empty-23hxbvmbplj2z").List(metav1.ListOptions{})
	if err != nil {

		//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve svcs from kubernetes API", err}
	}

	statefulsets, err := clientset.AppsV1().StatefulSets(baask8s.Namespace).List(metav1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods("empty-23hxbvmbplj2z").List(metav1.ListOptions{})
	if err != nil {

		//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve StatefulSets from kubernetes API", err}
	}

	daemonsets, err := clientset.AppsV1().DaemonSets(baask8s.Namespace).List(metav1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods("empty-23hxbvmbplj2z").List(metav1.ListOptions{})
	if err != nil {

		//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve DaemonSets from kubernetes API", err}
	}


	deployments_all := []Baasdeployment{}

	//my_deployment := Baasdeployment{}
	//podsii := make([]*v1.Pod, len(pods.Items))
	//deployments := make([]*Baasdeployment, len(deployments)+len(statefulsets.Items)+len(deamonsets.Items))
	//deployments_all := make([]*Baasdeployment, len(deployments.Items))

	deploymentii := make([]*apps.Deployment, len(deployments.Items))
	for j := range deployments.Items {
		deploymentii[j] = &deployments.Items[j]
	}

    //fooType := reflect.TypeOf(deploymentii[0])
	//for g := 0; g < fooType.NumMethod(); g++ {
	//	method := fooType.Method(g)
	//	log.Printf(method.Name)
	//}

	for i := range deployments.Items {
		var my_deployment Baasdeployment

		
		my_deployment.Name = deployments.Items[i].GetName()
		my_deployment.Namespace = deployments.Items[i].GetNamespace()
		my_deployment.Createdtime = deployments.Items[i].GetCreationTimestamp()
		my_deployment.Replicas = deployments.Items[i].Status.ReadyReplicas
		my_deployment.Type = "deployments"
		deployments_all = append(deployments_all,my_deployment)
		
	}
	for j := range statefulsets.Items {
		//my_deployment=new(Baasdeployment)
		var my_deployment Baasdeployment
		//fooType = reflect.TypeOf(statefulsets.Items)
		//for g := 0; g < fooType.NumMethod(); g++ {
		//	method := fooType.Method(g)
		//	log.Printf(method.Name)
		//}
		
		my_deployment.Name = statefulsets.Items[j].GetName()
		my_deployment.Namespace = statefulsets.Items[j].GetNamespace()
		my_deployment.Createdtime = statefulsets.Items[j].GetCreationTimestamp()
		my_deployment.Replicas = statefulsets.Items[j].Status.ReadyReplicas
		my_deployment.Type = "statefulsets"
		deployments_all = append(deployments_all,my_deployment)
		
	}
	for k := range daemonsets.Items {
		//my_deployment=new(Baasdeployment)
		var my_deployment Baasdeployment
		
		my_deployment.Name = daemonsets.Items[k].GetName()
		my_deployment.Namespace = daemonsets.Items[k].GetNamespace()
		my_deployment.Createdtime = daemonsets.Items[k].GetCreationTimestamp()
		my_deployment.Replicas = daemonsets.Items[k].Status.NumberReady
		my_deployment.Type = "daemonsets"
		deployments_all = append(deployments_all,my_deployment)
		
	}


	//return response.JSON(w, podsjj)
	return response.JSON(w, deployments_all)
	
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sStatefulSetsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

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


	//todo: check if namespace exsit?


	statefulsets, err := clientset.AppsV1().StatefulSets(baask8s.Namespace).List(metav1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods("empty-23hxbvmbplj2z").List(metav1.ListOptions{})
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve StatefulSets from kubernetes API", err}
	}


	//return response.JSON(w, podsjj)
	return response.JSON(w, statefulsets.Items)
	
}

// GET request on /api/baask8ss
func (handler *Handler) baask8sDaemonSetsList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	baask8sID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid baask8s identifier route variable", err}
	}

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


	//todo: check if namespace exsit?


	daemonsets, err := clientset.AppsV1().DaemonSets(baask8s.Namespace).List(metav1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods("empty-23hxbvmbplj2z").List(metav1.ListOptions{})
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve DaemonSets from kubernetes API", err}
	}


	//return response.JSON(w, podsjj)
	return response.JSON(w, daemonsets.Items)
	
}

