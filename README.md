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
5. Export env variables needed for Skaffold (example here)
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
3. Endpoints
   - Items API will be available at `localhost:8080`
   - Items API Swagger UI will be available at `localhost:8080/docs`
   - Reservations API will be available at `localhost:8081`
   - Reservations API Swagger UI will be available at `localhost:8081/docs`

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
To run tests you can use:
```
make test
```

### Generators

#### To generate stubs for Open API

1. Install [OAPI Codegen](https://github.com/deepmap/oapi-codegen)
2. Generate API stubs
   ```
   make generate-openapi
   ```

#### To generate mocks for interfaces

1. Install [GoMock](https://github.com/golang/mock)
2. Generate mocks
   ```
   make generate-mock
   ```

## Notes

The solution was implemented following specification-first approach.

Here you can find list of operations:
- Used skeleton of the project using [Project Layout](https://github.com/golang-standards/project-layout) with additional app folder.
- Created OpenAPI3 files, providing specification for Items API and Reservations API.
- Generated API client and server files using [Open API Generator](https://github.com/deepmap/oapi-codegen).
- Started minikube.
- Added Docker/Skaffold/Helm files for deployment locally.
- Implemented dummy stores and servers for Items/Reservations. 
- Added Swagger UI for both microservices
- Implemented service
- Implemented repository and DB connections + migration
- Added [GoMock](https://github.com/golang/mock) for mocking dependencies
- Implemented tests for stores/services/repository

### Packages 

```
app # directory with packages, which are valid only for specific programs (services)
   - items-api
      - api # Open API files and generated server stubs
      - store # Implementation of Open API store
      - service # Business Logic layer for glueing network and persistance layers
      - repository # Wraps DB calls and provides interface to access the storage
      - server # Provides business logic for server
   - reservations-api
      - api # Open API files and generated server and client stubs
      - server # Provides business logic for server

cmd # collection of entry files with main functions, also used as containers for dependency injection 
   - items-api
   - reservations-api

internal # packages used internally
   - config # configuration for the services

pkg # packages that can be shared with other projects (agnostic)
   - log # provides wrapper for logging

migration # files used for migration of the database

deployments # Dockerfiles and Helm charts
```

### Self assessment

In general, I'm quite happy with outcome. I would start with TDD approach, 
if I would have more experience using Open API Generator. 
It was really fun to implement servers using specification-first approach and I will consider it for future projects.
It was a bit difficult to complete all the tasks within 3-4 hours (took me bit more than 5h).

### To improve
1. CI & CD is not finished. The solution will work only locally.
2. Tests are too optimistic and are not covering many pessimistic use cases. Coverage is too low.
3. Call to reservation API needs to be improved.
   - currently using goroutines for starting async call
   - current approach will fail when there are operations in progress and server restarts
   - this will lead to items having not completed reservations
4. Possible solutions to mitigate 3.
   - Implementation of queuing system with dead-letter redelivery. This will help to survive server restarts and can repeat operations if needed.
     - SNS/SQS
     - RabbitMQ/ActiveMQ
   - Another solution could be to implement synchronization job services between Items API and Reservations API.
5. Docker, helm charts and code needs to be reviewed and checked for security issues.
   - For instance root access in docker containers should be blocked.
## Built With

* [Skaffold](https://skaffold.dev/)
* [OpenAPI](https://swagger.io/specification/)
* [Go](https://golang.org/)
* [Docker](https://www.docker.com/)
* [Kubernetes](https://kubernetes.io/)
* [Helm](https://helm.sh/)
* [Minikube](https://minikube.sigs.k8s.io/docs/start/)

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see
the [tags on this repository](https://github.com/your/project/tags).
