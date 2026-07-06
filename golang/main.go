package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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
	fmt.Printf("Deployment check status: %s\n", deployments)
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

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	if len(nodes.Items) < 1 {
		return "node check failed", nil
	}
	return "Complete", nil
}

// getDeploymentStatus returns a list of deployments and their associated status
func getDeploymentStatus(clientset kubernetes.Interface) (string, error) {
	deploymentClient := clientset.AppsV1().Deployments(metav1.NamespaceAll)
	deployments, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	if len(deployments.Items) == 0 {
		return "No Deployments Found", nil
	}

	for _, d := range deployments.Items {
		fmt.Printf(" Name: %s | Replicas: %d Available / %d Desired\n",
			d.Name, d.Status.AvailableReplicas, *d.Spec.Replicas)
		if d.Status.AvailableReplicas != *d.Spec.Replicas {
			fmt.Println("\n^^^^^^ REPLICAS DO NOT MATCH ^^^^^^")
		}
	}

	return "Complete", nil
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
		fmt.Fprintf(w, "deployCheckStatus, %q", deployments)
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
