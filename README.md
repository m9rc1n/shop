# Items + Reservations API

Simple Rest API representing items and reservations in the system.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

1. Install [Skaffold](https://skaffold.dev/)
2. Install [Kubernetes](https://kubernetes.io/)
3. Install [Helm](https://helm.sh/)
4. Install [Minikube](https://minikube.sigs.k8s.io/docs/start/)
5. Export env variables needed for Skaffold
   ```
   export SHOP_DB_USER=user &&
   export SHOP_DB_PASSWORD=super-secret &&
   export SHOP_DB_NAME=shop
   ```

### Running (locally)

1. Start minikube
   ```
   minikube start
   ```
2. Start skaffold to deploy helm charts to your cluster
   ```
   skaffold dev
   ```

### Debugging

You can use Jetbrains Run configurations

## Code Quality

### Linting

[Install](https://github.com/golangci/golangci-lint) to get linting advices.

#### To install liner
```
go get github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### To run linter
```
golangci-lint run -v --go=1.18 --timeout=240s 
```

### Testing
1. Install [GoMock](https://github.com/golang/mock)
2. Generate mocks
   ```
   make mock
   ```
3. Run tests
   ```
   make test
   ```

## Built With

* [Skaffold](https://skaffold.dev/)
* [OpenAPI](https://swagger.io/specification/)
* [Go](https://golang.org/)
* [Docker](https://www.docker.com/)
* [Kubernetes](https://kubernetes.io/)
* [Helm](https://helm.sh/)

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see
the [tags on this repository](https://github.com/your/project/tags).
