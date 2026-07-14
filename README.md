# project

This repository contains the boilerplate projects for golang api

### Project

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
docker tag tyk:latest andrewstewartelliott/tyk:latest
docker push andrewstewartelliott/tyk:latest
```

Using default service account created by helm, need to grant cluster-role perms
```
kubectl apply -f yaml/clusterRole.yaml
```

To install via helm:
```
helm install tyk-app ./tyk-sre-assignment-chart
```
To port forward with kind:
```
kubectl port-forward service/tyk-app-tyk-sre-assignment-chart 8080:8080
```

To uninstall via helm:
```
helm uninstall tyk-app
```
