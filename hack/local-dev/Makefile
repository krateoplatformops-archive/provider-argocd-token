# Set the shell to bash always
SHELL := /bin/bash

# Look for a .env file, and if present, set make variables from it.
ifneq (,$(wildcard ./.env))
	include ./.env
	export $(shell sed 's/=.*//' ./.env)
endif

KIND_CLUSTER_NAME ?= local-dev
KUBECONFIG ?= $(HOME)/.kube/config

# Project name
PROJECT_NAME := provider-argocd-token
# Git Owner and Vendor Company
ORG_NAME := krateoplatformops
VENDOR := Kiratech


VERSION := $(shell git describe --tags --abbrev=0 | sed -n "s/^v\(.*\)$$/\1/p")
ifndef VERSION
VERSION := 0.0.0
endif

# Build Target
TARGET_OS := linux
TARGET_ARCH := amd64

# Github Container Registry
DOCKER_REGISTRY := ghcr.io/$(ORG_NAME)

# Tools
KIND=$(shell which kind)
KUBECTL=$(shell which kubectl)

.PHONY: generate tidy lint help build.vars 

.DEFAULT_GOAL := help

## help: Print this help
help: Makefile
	@echo
	@echo " Choose a command run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo


.PHONY: print.vars
## print.vars: Print all the build variables
print.vars:
	@echo VENDOR=$(VENDOR)
	@echo ORG_NAME=$(ORG_NAME)
	@echo PROJECT_NAME=$(PROJECT_NAME)
	@echo GITHUB_TOKEN=$(GITHUB_TOKEN)
	@echo TARGET_OS=$(TARGET_OS)
	@echo TARGET_ARCH=$(TARGET_ARCH)
	@echo VERSION=$(VERSION)
	@echo DOCKER_REGISTRY=$(DOCKER_REGISTRY)

## kind.up: Starts a KinD cluster for local development
kind.up:
	@$(KIND) get kubeconfig --name $(KIND_CLUSTER_NAME) >/dev/null 2>&1 || $(KIND) create cluster --name=$(KIND_CLUSTER_NAME)

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

## kind.down: Shuts down the KinD cluster
kind.down:
	@$(KIND) delete cluster --name=$(KIND_CLUSTER_NAME)


## install.provider: Install this provider using HELM
install.provider:
	@$(KUBECTL) crossplane install provider "$(DOCKER_REGISTRY)/crossplane-$(PROJECT_NAME):$(VERSION)"

## install.provider.helm: Install this provider using HELM
install.provider.argocd:
	@$(KUBECTL) create namespace argo-system || true
	@$(KUBECTL) create secret docker-registry cr-token \
	--namespace argo-system --docker-server=ghcr.io \
	--docker-password=$(GITHUB_TOKEN) --docker-username=$(ORG_NAME) || true
	@helm repo add krateo-runtime https://krateo.io/runtime-helm-charts || true
	@helm repo update
	helm install $(PROJECT_NAME) --namespace argo-system krateo-runtime/provider-argocd


.PHONY: cr.token
## cr.token: Create the secret with the backup container registry credentials
cr.token:
	@$(KUBECTL) create secret docker-registry cr-token \
	--namespace crossplane-system --docker-server=ghcr.io \
	--docker-password=$(GITHUB_TOKEN) --docker-username=$(ORG_NAME)

# proxy:
# 	kubectl port-forward svc/argocd-server -n argo-system 8080:443

