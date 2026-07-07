package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type DeploymentStatus struct {
	namespace string
	name      string
	available int32
	desired   int32
	errors    bool
}

func main() {

	kubeconfig := flag.String("kubeconfig", "", "path to kubeconfig, leave empty for in-cluster")
	listenAddr := flag.String("address", ":8080", "HTTP server listen address")

	flag.Parse()

	kConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		panic(err)
	}

	version, err := getKubernetesVersion(clientset)
	if err != nil {
		panic(err)
	}

	// gets deployments in all namespaces
	deployments, err := getDeploymentStatus(clientset)
	if err != nil {
		panic(err)
	}

	apistatus, err := getK8sApiStatus(clientset)
	if err != nil {
		panic(err)
	}

	fmt.Printf("API check status: %s\n", apistatus)
	fmt.Println("Deployments")
	fmt.Println(deployments)

	b, err := json.Marshal(deployments)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)

	fmt.Printf("Connected to Kubernetes %s\n", version)

	if err := startServer(*listenAddr, clientset); err != nil {
		panic(err)
	}
}

// getKubernetesVersion returns a string GitVersion of the Kubernetes server defined by the clientset.
//
// If it can't connect an error will be returned, which makes it useful to check connectivity.
func getKubernetesVersion(clientset kubernetes.Interface) (string, error) {
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	return version.String(), nil
}

// test the k8s api connection
func getK8sApiStatus(clientset kubernetes.Interface) (string, error) {
	found := false
	pods, err := clientset.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if len(pods.Items) < 1 {
		return "pod check failed", nil
	} else {
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, "kube-apiserver") {
				found = true
			}
		}
		if found == false {
			return "didn't find kube-apiserver pod", nil
		}
	}

	return "Complete", nil
}

// getDeploymentStatus returns a list of deployments and their associated status
func getDeploymentStatus(clientset kubernetes.Interface) ([]DeploymentStatus, error) {
	deploymentClient := clientset.AppsV1().Deployments(metav1.NamespaceAll)
	deployments, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	if len(deployments.Items) == 0 {
		panic("error getting deployments")
	}

	deploymentList := make([]DeploymentStatus, 0, len(deployments.Items))

	errors := false
	for _, d := range deployments.Items {
		if d.Status.AvailableReplicas != *d.Spec.Replicas {
			errors = true
		}
		// populate array of structs
		deploymentList = append(deploymentList, DeploymentStatus{
			namespace: d.Namespace,
			name:      d.Name,
			available: d.Status.AvailableReplicas,
			desired:   *d.Spec.Replicas,
			errors:    errors,
		})
	}

	return deploymentList, nil
}

// startServer launches an HTTP server with defined handlers and blocks until it's terminated or fails with an error.
//
// Expects a listenAddr to bind to.
func startServer(listenAddr string, clientset kubernetes.Interface) error {

	// gets deployments in all namespaces
	deployments, err := getDeploymentStatus(clientset)
	if err != nil {
		panic(err)
	}
	// gets pods and nodes as an api test
	apistatus, err := getK8sApiStatus(clientset)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/deploymentstatus", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(deployments)
	})
	http.HandleFunc("/apistatus", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "apiCheckStatus, %q", apistatus)
	})

	fmt.Printf("Server listening on %s\n", listenAddr)

	return http.ListenAndServe(listenAddr, nil)
}

// healthHandler responds with the health status of the application.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("ok"))
	if err != nil {
		fmt.Println("failed writing to response")
	}
}
