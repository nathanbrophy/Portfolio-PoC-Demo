BIN ?= ./bin
EXE ?= server
OS ?= $(uname -s)
ARCH ?= $(uname -p)
VERSION ?= v1.0.0
IMAGE_REPO ?= nathanbrophy/example-server
IMAGE_REGISTRY ?= docker.io

IMAGE_TAG = $(IMAGE_REGISTRY)/$(IMAGE_REPO):$(VERSION)

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: test
test:
	go test ./pkg/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build: fmt vet test
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(ARCH) go build -o $(BIN)/$(EXE)

.PHONY: run
run: build
	go run main.go

.PHONY: docker-build
docker-build: fmt vet test
	docker build -t $(IMAGE_TAG) .

.PHONY: docker-push
docker-push: docker-build
	docker push $(IMAGE_TAG)
