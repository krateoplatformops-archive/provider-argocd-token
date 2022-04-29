# Set the shell to bash always
SHELL := /bin/bash

# Look for a .env file, and if present, set make variables from it.
ifneq (,$(wildcard ./.env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

KIND_CLUSTER_NAME ?= local-dev
KUBECONFIG ?= $(HOME)/.kube/config

VERSION := $(shell git describe --dirty --always --tags | sed 's/-/./2' | sed 's/-/./2')
ifndef VERSION
VERSION := 0.0.0
endif

BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
REPO_URL := $(shell git config --get remote.origin.url | sed "s/git@/https\:\/\//; s/\.com\:/\.com\//; s/\.git//")
LAST_COMMIT := $(shell git log -1 --pretty=%h)

PROJECT_NAME := provider-argocd-token
ORG_NAME := krateoplatformops
VENDOR := Kiratech

# Github Container Registry
DOCKER_REGISTRY := ghcr.io/$(ORG_NAME)

TARGET_OS := linux
TARGET_ARCH := amd64

# Tools
KIND=$(shell which kind)
LINT=$(shell which golangci-lint)
KUBECTL=$(shell which kubectl)
DOCKER=$(shell which docker)
SED=$(shell which sed)

.DEFAULT_GOAL := help

.PHONY: help
## help: Print this help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo


.PHONY: print.vars
## print.vars: Print all the build variables
print.vars:
	@echo VENDOR=$(VENDOR)
	@echo ORG_NAME=$(ORG_NAME)
	@echo PROJECT_NAME=$(PROJECT_NAME)
	@echo REPO_URL=$(REPO_URL)
	@echo LAST_COMMIT=$(LAST_COMMIT)
	@echo VERSION=$(VERSION)
	@echo BUILD_DATE=$(BUILD_DATE)
	@echo TARGET_OS=$(TARGET_OS)
	@echo TARGET_ARCH=$(TARGET_ARCH)
	@echo DOCKER_REGISTRY=$(DOCKER_REGISTRY)


.PHONY: dev
## dev: Run the controller in debug mode
dev: generate
	$(KUBECTL) apply -f package/crds/ -R
	go run cmd/main.go -d

.PHONY: generate
## generate: Generate all CRDs
generate: tidy
	go generate ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	$(LINT) run

.PHONY: kind.up
## kind.up: Starts a KinD cluster for local development
kind.up:
	@$(KIND) get kubeconfig --name $(KIND_CLUSTER_NAME) >/dev/null 2>&1 || $(KIND) create cluster --name=$(KIND_CLUSTER_NAME)

.PHONY: kind.down
## kind.down: Shuts down the KinD cluster
kind.down:
	@$(KIND) delete cluster --name=$(KIND_CLUSTER_NAME)

.PHONY: image.build
## image.build: Build the Docker image
image.build:
	@$(DOCKER) build -t "$(DOCKER_REGISTRY)/$(PROJECT_NAME):$(VERSION)" \
	--build-arg METRICS_PORT=9090 \
	--build-arg VERSION="$(VERSION)" \
	--build-arg BUILD_DATE="$(BUILD_DATE)" \
	--build-arg REPO_URL="$(REPO_URL)" \
	--build-arg LAST_COMMIT="$(LAST_COMMIT)" \
	--build-arg PROJECT_NAME="$(PROJECT_NAME)" \
	--build-arg VENDOR="$(VENDOR)" .
	@$(DOCKER) rmi -f $$(docker images -f "dangling=true" -q)

.PHONY: image.push
## image.push: Push the Docker image to the Github Registry
image.push:
	@$(DOCKER) push "$(DOCKER_REGISTRY)/$(PROJECT_NAME):$(VERSION)"


.PHONY: cr.secret
cr.secret:
	$(KUBECTL) create secret docker-registry cr-token \
	--namespace crossplane-system --docker-server=ghcr.io \
	--docker-password=$(GITHUB_TOKEN) --docker-username=$(ORG_NAME)


## install.crossplane: Install Crossplane into the local KinD cluster
install.crossplane:
	$(KUBECTL) create namespace crossplane-system || true
	helm repo add crossplane-stable https://charts.crossplane.io/stable
	helm repo update
	helm install crossplane --namespace crossplane-system crossplane-stable/crossplane


## install.argocd: Install ArgoCD into the local KinD cluster
install.argocd:
	$(KUBECTL) create namespace argo-system || true
	$(KUBECTL) apply -n argo-system -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml


.PHONY: install.provider
 ## install.provider: Install this provider
install.provider: cr.secret
	@$(SED) 's/VERSION/$(VERSION)/g' ./examples/provider.yaml | $(KUBECTL) apply -f -
