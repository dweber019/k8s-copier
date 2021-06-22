# K8s-copier
This operator aims to provide Kubernetes CRD's to copy resources from one namespace into another.  
This new CRD is called `CopyResrouce` and the resulting resource is called target resource.

## Implemented resource types
- v1.Secret
- v1.ConfigMap

## Usage
You can find some usage examples in `config/samples/**`.

### Behavior
If you delete a CopyResource the target resource won't be deleted as it's possible that other implementation depend on it.

## Development setup
### Conventional commits
Execute the following terminal command in the root:
```
curl -o- https://raw.githubusercontent.com/craicoverflow/sailr/master/scripts/install.sh | bash
```
### Go
Follow the installation guide at https://golang.org/doc/install

## Development
For development, you should use `minikube` or any other possible kubernetes compatible implementation, as advised by the operator-sdk framework.

### Update the CRD
To update the CRD use the following command
```
make generate
make manifests
```
The files will be generated into `config/**`.

### Build and run
To install the CRD to the Kubernetes and run the operator outside of Kubernetes use
```
make install run
```

## Deployment
### Automated
The deployment and therefore publishing a Docker image is fully automated with GitHub workflows.

### Manually
Use the automated way over GitHub!!!
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
- [Conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)

## Todo:
- Write some tests
- Add a feature where you can define a template, which will take the source resource as input and produce any arbitrary target resource