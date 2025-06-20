SHELL=/usr/bin/env bash -o pipefail

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

BUILD_NAME?=
IMAGE_NAME?=
GO_PKG=github.com/yunling101/ControllerManager
GO_CMD=$(wildcard cmd/*)

TAG?=$(shell git rev-parse --short HEAD)
VERSION?=$(shell cat VERSION | tr -d " \t\n\r")
BUILD_DATE?=$(shell date +"%Y%m%d-%T")
BUILD_BRANCH?=$(shell git branch --show-current)
BUILD_HASH?=$(shell git rev-parse HEAD)

GO_BUILD_LDFLAGS=\
	-w \
	-s \
	-X $(GO_PKG)/common.Version=$(VERSION) \
	-X $(GO_PKG)/common.BuildDate=$(BUILD_DATE) \
	-X $(GO_PKG)/common.Branch=$(BUILD_BRANCH) \
	-X $(GO_PKG)/common.CCGitHash=$(BUILD_HASH)

GO_BUILD_RECIPE=\
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	CGO_ENABLED=0 \
	go build -ldflags="$(GO_BUILD_LDFLAGS)"

${BUILD_NAME}:
	@echo -e "\033[34mbuilding start ${BUILD_NAME}... \033[0m"
	@$(GO_BUILD_RECIPE) -o ./output/${BUILD_NAME} cmd/${BUILD_NAME}/${BUILD_NAME}.go

${IMAGE_NAME}:
	@bash scripts/config-generate.sh docker
	@docker build -f cmd/${IMAGE_NAME}/Dockerfile -t yunling101/$(shell echo ${IMAGE_NAME} | tr 'A-Z' 'a-z'):${VERSION} .
	@bash scripts/config-generate.sh docker-rm

cmdChannel:
	@docker build -f cmd/cmdChannel/Dockerfile -t yunling101/cmdchannel:latest .

.PHONY: build
build:
	@$(foreach dir, $(GO_CMD), \
		echo -e "\033[34mbuilding start $(dir)... \033[0m"; \
		$(GO_BUILD_RECIPE) -o ./output/$(notdir $(dir)) cmd/$(notdir $(dir))/$(notdir $(dir)).go; \
	)

.PHONY: docker
docker:
	@$(foreach dir, $(GO_CMD), \
		echo -e "\033[33mbuilding image $(dir)... \033[0m"; \
		docker build -f $(dir)/Dockerfile -t yunling101/$(shell echo $(notdir $(dir)) | tr 'A-Z' 'a-z'):${VERSION} .; \
	)

CONSUL?=$(shell cat third_party/consul/Dockerfile | grep "ARG VERSION" | tr -d 'ARG VERSION=\"')
ALERTMANAGER?=$(shell cat third_party/alertmanager/Dockerfile | grep "ARG VERSION" | tr -d 'ARG VERSION=\"')
PROMETHEUS?=$(shell cat third_party/prometheus/Dockerfile | grep "ARG VERSION" | tr -d 'ARG VERSION=\"')

.PHONY: consul
consul:
	@docker build -f third_party/consul/Dockerfile -t yunling101/consul:${CONSUL} .

.PHONY: alertmanager
alertmanager:
	@docker build -f third_party/alertmanager/Dockerfile -t yunling101/alertmanager:${ALERTMANAGER} .

.PHONY: prometheus
prometheus:
	@docker build -f third_party/prometheus/Dockerfile -t yunling101/prometheus:${PROMETHEUS} .
