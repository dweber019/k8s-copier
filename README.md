# K8s-copier
This operator aims to provide Kubernetes CRD's to copy resources from one namespace into another.

## Implemented resource types
- v1.Secret
- v1.ConfigMap

## Development
For development you should use `minikube` or any other possible kubernetes compatible implementation, as advised by the operator-sdk framework.

### Update the CRD
To update the CRD use the following command
```
make generate
make manifests
```

### Build and run
To install the CRD to the Kubernetes and run the operator outside of Kubernetes use
```
make install run
```

## Deployment
To update the docker image use this command
```
make docker-build docker-push IMG=docker.io/dweber019/k8s-copier:v0.0.1
```
After this you can run with
```
make deploy IMG=docker.io/dweber019/k8s-copier:v0.0.1
```

## Useful links
- [Kubebilder](https://book.kubebuilder.io)
- [Operator tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)

## Todo:
- convential commits