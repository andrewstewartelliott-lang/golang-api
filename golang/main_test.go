package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	disco "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetKubernetesVersion(t *testing.T) {
	okClientset := fake.NewSimpleClientset()
	okClientset.Discovery().(*disco.FakeDiscovery).FakedServerVersion = &version.Info{GitVersion: "1.25.0-fake"}

	okVer, err := getKubernetesVersion(okClientset)
	assert.NoError(t, err)
	assert.Equal(t, "1.25.0-fake", okVer)

	badClientset := fake.NewSimpleClientset()
	badClientset.Discovery().(*disco.FakeDiscovery).FakedServerVersion = &version.Info{}

	badVer, err := getKubernetesVersion(badClientset)
	assert.NoError(t, err)
	assert.Equal(t, "", badVer)
}

func TestGetDeploymentStatus(t *testing.T) {
	okClientset := fake.NewSimpleClientset()
	okClientset.Discovery()
	deployments := okClientset.AppsV1().Deployments("default")

	// Create a Deployment
	replicas := int32(1)
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deploy",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas:     1,
			AvailableReplicas: 1,
		},
	}

	_, err := deployments.Create(context.Background(), deploy, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("create deployment: %v", err)
	}
	okDeployment, err := getDeploymentStatus(okClientset)
	assert.NoError(t, err)
	assert.Equal(t, "Complete", okDeployment)
}

func TestGetK8sApiStatus(t *testing.T) {
	okClientset := fake.NewSimpleClientset()
	okClientset.Discovery()
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kube-apiserver",
			Namespace: "kube-system",
		},
	}

	// Seed the fake client
	_, err := okClientset.CoreV1().Pods("kube-system").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	okapitest, err := getK8sApiStatus(okClientset)
	assert.NoError(t, err)
	assert.Equal(t, "Complete", okapitest)
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)
	res := rec.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	defer func(Body io.ReadCloser) {
		assert.NoError(t, Body.Close())
	}(res.Body)
	resp, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.Equal(t, "ok", string(resp))
}
