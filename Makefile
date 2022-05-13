#!/usr/bin/make
SHELL  := /bin/bash


TOOL := dns-preload
export PATH = /usr/bin:/usr/local/bin:/usr/local/sbin:/usr/sbin:/bin:/sbin:/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/build/bin
BINPATH := bin
GO_DIR := src/github.com/jimmystewpot/dns-preload/
DOCKER_IMAGE := golang:1.18-bullseye
SYNK_IMAGE := snyk/snyk:golang
INTERACTIVE := $(shell [ -t 0 ] && echo 1)
TEST_DIRS := ./cmd/...

get-golang:
	docker pull ${DOCKER_IMAGE}

get-synk:
	docker pull ${SYNK_IMAGE}

.PHONY: clean
clean:
	@echo $(shell docker images -qa -f 'dangling=true'|egrep '[a-z0-9]+' && docker rmi $(shell docker images -qa -f 'dangling=true'))

lint:
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
	@docker run \
		--rm \
		-v $(CURDIR):/build/$(GO_DIR) \
		--workdir /build/$(GO_DIR) \
		-e GOPATH=/build \
		-e PATH=$(PATH) \
		-t ${DOCKER_IMAGE} \
		make build-all

build-all: deps test lint pdns-statsd-proxy

deps:
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.1

dns-preload:
	@echo ""
	@echo "***** Building $$TOOL *****"
	go build -race -ldflags="-s -w" -o $(BINPATH)/$(TOOL) ./cmd/$(TOOL)
	@echo ""

test:
	@echo ""
	@echo "***** Testing ${TOOL} *****"
	go test -a -v -race -coverprofile=reports/coverage.txt -covermode=atomic -json ./cmd/$(TOOL) 1> reports/testreport.json
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