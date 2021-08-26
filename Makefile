# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

SHELL = /bin/bash

default: build

# #### GO Binary Management ####
.PHONY: deps-go-binary deps-counterfeiter deps-ginkgo deps-golangci-lint

GO_VERSION := $(shell go version)
GO_VERSION_REQUIRED = go1.16
GO_VERSION_MATCHED := $(shell go version | grep $(GO_VERSION_REQUIRED))

deps-go-binary:
ifndef GO_VERSION
	$(error Go not installed)
endif
ifndef GO_VERSION_MATCHED
	$(error Required Go version is $(GO_VERSION_REQUIRED), but was $(GO_VERSION))
endif
	@:

HAS_COUNTERFEITER := $(shell command -v counterfeiter;)
HAS_GINKGO := $(shell command -v ginkgo;)
HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)

# If go get is run from inside the project directory it will add the dependencies
# to the go.mod file. To avoid that we import from another directory
deps-counterfeiter: deps-go-binary
ifndef HAS_COUNTERFEITER
	cd /; go get -u github.com/maxbrunsfeld/counterfeiter/v6
endif

deps-ginkgo: deps-go-binary
ifndef HAS_GINKGO
	cd /; go get github.com/onsi/ginkgo/ginkgo github.com/onsi/gomega
endif

deps-golangci-lint: deps-go-binary
ifndef HAS_GOLANGCI_LINT
	cd /; go get github.com/golangci/golangci-lint/cmd/golangci-lint
endif

# #### CLEAN ####
.PHONY: clean

clean: deps-go-binary 
	rm -rf build/*
	go clean --modcache


# #### DEPS ####
.PHONY: deps deps-counterfeiter deps-ginkgo deps-modules

deps-modules: deps-go-binary
	go mod download

deps: deps-modules deps-counterfeiter deps-ginkgo


# #### BUILD ####
.PHONY: build

SRC = $(shell find . -name "*.go" | grep -v "_test\." )
VERSION := $(or $(VERSION), dev)
LDFLAGS="-X github.com/vmware-labs/marketplace-cli/v2/cmd.Version=$(VERSION)"

build/mkpcli: $(SRC)
	go build -o build/mkpcli -ldflags ${LDFLAGS} ./main.go

build/mkpcli-darwin: $(SRC)
	GOARCH=amd64 GOOS=darwin go build -o build/mkpcli-darwin -ldflags ${LDFLAGS} ./main.go

build/mkpcli-linux: $(SRC)
	GOARCH=amd64 GOOS=linux go build -o build/mkpcli-linux -ldflags ${LDFLAGS} ./main.go

build: deps build/mkpcli

build-all: build/mkpcli-darwin build/mkpcli-linux

build-image: build/mkpcli-linux
	docker build . --tag harbor-repo.vmware.com/tanzu_isv_engineering/mkpcli:$(VERSION)

# #### TESTS ####
.PHONY: lint test test-features test-units

test-units: deps
	ginkgo -r -skipPackage external,features .

test-features: deps
	ginkgo -r -tags=feature features

test-external: deps
ifndef CSP_API_TOKEN
	$(error CSP_API_TOKEN must be defined to run external tests)
else
	ginkgo -r -tags=external external
endif

test: deps lint test-units test-features test-external

lint: deps-golangci-lint
	golangci-lint run

# #### DEVOPS ####
.PHONY: set-pipeline
set-pipeline: ci/pipeline.yaml
	fly -t tie set-pipeline --config ci/pipeline.yaml --pipeline marketplace-cli
