# load .env file if present
-include .env

VERSION ?= dev
COMMIT ?= local
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS = -X 'main.Version=$(VERSION)' \
					-X 'main.CommitHash=$(COMMIT)' \
					-X 'main.BuildTime=$(BUILD_TIME)'

# `make build VERSION=0.1.0 COMMIT=$(git rev-parse --short HEAD)`
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o ./bin/server

.PHONY: run
run:
	air

.PHONY: gen-swagger
gen-swagger:
	swag init && swag fmt

# `make create-migration NAME=migration_name`
.PHONY: create-migration
create-migration:
	@if [ -z "$(NAME)" ]; then \
		echo "Missing NAME parameter. Usage: make create-migration NAME=migration_name"; \
		exit 1; \
	fi
	atlas migrate diff $(NAME) --env local

.PHONY: migrate-dev
migrate-dev:
	atlas migrate apply --env runtime --url "$(ATLAS_DATABASE_URL)"