# Load .env if it exists
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

GIT_COMMIT=$(shell git rev-parse --short HEAD)

.PHONY: dev lint

lint:
	golangci-lint run
	hadolint Dockerfile

dev:
	go mod tidy
	godotenv -f .env go run main.go

build:
	docker build -t ${DOCKER_REPO}/pam-manager:latest .
	docker tag ${DOCKER_REPO}/pam-manager:latest ${DOCKER_REPO}/pam-manager:latest:${GIT_COMMIT}
