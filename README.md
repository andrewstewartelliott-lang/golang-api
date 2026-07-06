# tyk-sre-assignment

This repository contains the boilerplate projects for the SRE role interview assignments.

### Project

Location: https://github.com/TykTechnologies/tyk-sre-assignment/tree/main/golang

In order to build the project run:
```
go mod tidy & go build
```

To run it against a real Kubernetes API server:
```
./tyk-sre-assignment --kubeconfig '/path/to/your/kube/conf' --address ":8080"
```

To execute unit tests:
```
go test -v
```

To build local docker image:
```
docker build -t tyk:latest .
```

Using default service account created by helm, need to grant cluster-role perms
```
kubectl apply -f yaml/clusterRole.yaml
```

To install via helm:
```
helm install tyk-app ./tyk-sre-assignment-chart
```