MODULE = $(shell go list -m)
SHELL := /bin/bash
LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
GO_PKGS=$(foreach pkg, $(shell go list ./...), $(if $(findstring /vendor/, $(pkg)), , $(pkg)))
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')

ENV := local
ifdef $$APP_ENV
ENV := $$APP_ENV
endif

export PROJECT = github.com/m9rc1n/shop

build:
	env GOOS=linux GOARCH=amd64 go build -o bin/server $(PROJECT)/cmd
	env GOOS=linux GOARCH=amd64 go build -o bin/admin $(PROJECT)/cmd/admin
	chmod +x bin/server
	chmod +x bin/admin

build-mac:
	env GOOS=darwin GOARCH=amd64 go build -o bin/server $(PROJECT)/cmd
	env GOOS=darwin GOARCH=amd64 go build -o bin/admin $(PROJECT)/cmd/admin
	chmod +x bin/server
	chmod +x bin/admin

test:
	go test ./... -cover

generate-mock:
	mockgen -source=app/items-api/repository/repository.go -destination=app/items-api/repository/mock/repository_mock.go \\
	mockgen -source=app/items-api/service/service.go -destination=app/items-api/service/mock/service_mock.go
	mockgen -source=app/reservations-api/api/reservations.api.gen.go -destination=app/reservations-api/api/mock/reservations_mock.api.gen.go

tidy:
	go mod tidy

generate-openapi:
	oapi-codegen --config app/items-api/api/items.config.yml app/items-api/api/items.yml > app/items-api/api/items.api.gen.go \\
	oapi-codegen --config app/reservations-api/api/reservations.config.yml app/reservations-api/api/reservations.yml > app/reservations-api/api/reservations.api.gen.go

fmt:
	@go fmt $(GO_PKGS)
	@goimports -w -l $(GO_FILES)

$(LINT_TOOL):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.26.0

qc: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...
	staticcheck ./...