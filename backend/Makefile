# `make build VERSION=0.1.0`
VERSION ?= dev

LDFLAGS = -X 'main.Version=$(VERSION)'

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o ./bin/

.PHONY: run
run:
	air

.PHONY: gen-swagger
gen-swagger:
	swag init && swag fmt