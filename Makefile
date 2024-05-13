GOCMD=go
BINARY_NAME=dockerright
MAJOR_VERSION?=0
MINOR_VERSION?=0
PATCH_VERSION?=8
VERSION=$(MAJOR_VERSION).$(MINOR_VERSION).$(PATCH_VERSION)
DOCKER_REGISTRY?=ghcr.io/bata94/
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

build:
	docker build --target prod --tag $(BINARY_NAME) .

release-git:
	git add .
	git commit -m "Release version $(VERSION)"
	git tag -a $(VERSION) -m "Release version $(VERSION)"
	git push
	git push --tags

release-docker:
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest

	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(MAJOR_VERSION)
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(MAJOR_VERSION)

	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(MAJOR_VERSION).$(MINOR_VERSION)
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(MAJOR_VERSION).$(MINOR_VERSION)

	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)

auto-release: fmt-go build release-git release-docker

lint: lint-go fmt-go lint-dockerfile

lint-dockerfile:
ifeq ($(shell test -e ./Dockerfile && echo -n yes),yes)
	$(eval CONFIG_OPTION = $(shell [ -e $(shell pwd)/.hadolint.yaml ] && echo "-v $(shell pwd)/.hadolint.yaml:/root/.config/hadolint.yaml" || echo "" ))
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--format checkstyle" || echo "" ))
	$(eval OUTPUT_FILE = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "| tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -i $(CONFIG_OPTION) hadolint/hadolint hadolint $(OUTPUT_OPTIONS) - < ./Dockerfile $(OUTPUT_FILE)
endif

lint-go:
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	golangci-lint run --enable-all --fix --color always

fmt-go:
	gofmt -w -s .
	gofumpt -w .
