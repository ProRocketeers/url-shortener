# `make build VERSION=0.1.0 COMMIT=$(git rev-parse --short HEAD)`
VERSION ?= dev
COMMIT ?= local
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS = -X 'main.Version=$(VERSION)' \
					-X 'main.CommitHash=$(COMMIT)' \
					-X 'main.BuildTime=$(BUILD_TIME)'

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o ./bin/server

.PHONY: run
run:
	air

.PHONY: gen-swagger
gen-swagger:
	swag init && swag fmt