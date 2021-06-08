# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

SHELL = /bin/bash
GO-VER = go1.16

default: build


# #### GO Binary Management ####
.PHONY: deps-go-binary deps-goimports

deps-go-binary:
	echo "Expect: $(GO-VER)" && \
		echo "Actual: $$(go version)" && \
	 	go version | grep $(GO-VER) > /dev/null

HAS_GO_IMPORTS := $(shell command -v goimports;)

deps-goimports: deps-go-binary
ifndef HAS_GO_IMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif


# #### CLEAN ####
.PHONY: clean

clean: deps-go-binary 
	rm -rf build/*
	go clean --modcache


# #### DEPS ####
.PHONY: deps deps-counterfeiter deps-ginkgo deps-modules

deps-modules: deps-goimports deps-go-binary
	go mod download

deps-counterfeiter: deps-modules
	command -v counterfeiter >/dev/null 2>&1 || go get -u github.com/maxbrunsfeld/counterfeiter/v6

deps-ginkgo: deps-go-binary
	command -v ginkgo >/dev/null 2>&1 || go get -u github.com/onsi/ginkgo/ginkgo github.com/onsi/gomega

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

lint: deps-goimports
	git ls-files | grep '.go$$' | xargs goimports -l -w


# #### DEVOPS ####
.PHONY: set-pipeline
set-pipeline: ci/pipeline.yaml
	fly -t tie set-pipeline --config ci/pipeline.yaml --pipeline marketplace-cli
