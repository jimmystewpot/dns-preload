#!/usr/bin/make
SHELL  := /bin/bash


TOOL := dns-preload
export PATH = $(shell echo $$PATH):/usr/bin:/usr/local/bin:/usr/local/sbin:/usr/sbin:/bin:/sbin:/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/build/bin:/home/runner/go:/home/runner/go/bin:
BINPATH := bin
GO_DIR := src/github.com/jimmystewpot/dns-preload/
DOCKER_IMAGE := golang:1.21-bookworm
SYNK_IMAGE := snyk/snyk:golang
INTERACTIVE := $(shell [ -t 0 ] && echo 1)
TEST_DIRS := ./...

get-golang:
	docker pull ${DOCKER_IMAGE}

get-synk:
	docker pull ${SYNK_IMAGE}

get-sonarcloud:
	docker pull sonarsource/sonar-scanner-cli

.PHONY: clean
clean:
	@echo $(shell docker images -qa -f 'dangling=true'|egrep '[a-z0-9]+' && docker rmi $(shell docker images -qa -f 'dangling=true'))

lint:
	@echo ""
	@echo "***** linting ${TOOL} with golangci-lint *****"
ifdef INTERACTIVE
	golangci-lint run -v $(TEST_DIRS)
else
	golangci-lint run --out-format checkstyle -v $(TEST_DIRS) 1> reports/checkstyle-lint.xml
endif
.PHONY: lint

#
# build the software
#
build: get-golang
	@echo ""
	@echo "***** Building ${TOOL} *****"
	@docker run \
		--rm \
		-v $(CURDIR):/build/$(GO_DIR) \
		--workdir /build/$(GO_DIR) \
		-e GOPATH=/build \
		-e PATH=$(PATH) \
		-t ${DOCKER_IMAGE} \
		make build-all

build-all: deps lint test dns-preload

test-all: deps lint test

deps:
	@echo ""
	@echo "***** Installing dependencies for ${TOOL} *****"
	go clean --cache
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

dns-preload:
	@echo ""
	@echo "***** Building ${TOOL} *****"
	git config --global --add safe.directory /build/src/github.com/jimmystewpot/dns-preload
	git status
	go build -race -trimpath -ldflags="-s -w" -o $(BINPATH)/$(TOOL) ./cmd/$(TOOL)
	@echo ""

linux-arm64:
	@echo ""
	@echo "***** Building ${TOOL} for Linux ARM64 *****"
	GOOS=linux GOARCH=arm64 go build -race -ldflags="-s -w" -o $(BINPATH)/$(TOOL) ./cmd/$(TOOL)
	@echo ""

linux-x64:
	@echo ""
	@echo "***** Building ${TOOL} for Linux x86-64 *****"
	GOOS=linux GOARCH=amd64 go build -race -ldflags="-s -w" -o $(BINPATH)/$(TOOL) ./cmd/$(TOOL)
	@echo ""

linux-arm32:
	@echo ""
	@echo "***** Building ${TOOL} for Linux ARM32 *****"
	GOOS=linux GOARCH=arm go build -race -ldflags="-s -w" -o $(BINPATH)/$(TOOL) ./cmd/$(TOOL)
	@echo ""

test:
	@echo ""
	@echo "***** Testing ${TOOL} *****"
	go test -a -v -race -coverprofile=reports/coverage.txt -covermode=atomic -json $(TEST_DIRS) 1> reports/testreport.json
	@echo ""


test-synk: get-synk
	@echo ""
	@echo "***** Testing vulnerabilities using Synk *****"
	@docker run \
		--rm \
		-v $(CURDIR):/build/$(GO_DIR) \
		--workdir /build/$(GO_DIR) \
		-e SNYK_TOKEN=${SYNK_TOKEN} \
		-e MONITOR=true \
		-t ${SYNK_IMAGE}

test-sonarcloud: get-sonarcloud
	@echo ""
	@echo "***** Doing code analysis using SonarCloud *****"
	@docker run \
		--rm \
		-v $(CURDIR):/build/$(GO_DIR) \
		--workdir /build/$(GO_DIR) \
		-e SNYK_TOKEN=${SYNK_TOKEN} \
		-e MONITOR=true \
		-t ${SYNK_IMAGE} \
		sonarsource/sonar-scanner-cli