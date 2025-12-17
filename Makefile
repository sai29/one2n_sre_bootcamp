ifneq (,$(wildcard .env))
	include .env
	export
endif

APP_NAME=student_api
MIGRATIONS_DIR=./migrations
COMPOSE=docker compose
NETWORK := $(shell basename "$(PWD)")_default

# ---- Docker image variables ----
IMAGE ?= ghcr.io/sai29/one2n_sre_bootcamp
VERSION ?= v0.1.0
GIT_SHA ?= $(shell git rev-parse --short HEAD)

.PHONY: \
	db-up db-migrate api-build api-up dev down reset \
	local-build local-run local-dev local-migrate local-migrate-down \
	test clean build-prod build-debug run-prod run-debug tag push lint

db-up:
		${COMPOSE} up -d db

db-migrate:
		docker run --rm \
		--network ${NETWORK} \
		-v ${PWD}/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$$STUDENT_API_DB_DSN" up

api-build:
		${COMPOSE} build api

api-up:
		${COMPOSE} up -d api

dev: db-up db-migrate api-build api-up

lint:
		golangci-lint run

down:
		${COMPOSE} down

reset:
		${COMPOSE} down -v

local-build:
		go build -o $(APP_NAME) ./cmd/api

local-run:
		go run ./cmd/api

test:
		go test ./... -v -cover

local-dev:
		air

local-migrate:
		migrate -path=$(MIGRATIONS_DIR) -database="$(STUDENT_API_DB_DSN)" up

local-migrate-down:
		migrate -path=$(MIGRATIONS_DIR) -database="$(STUDENT_API_DB_DSN)" down 1

clean:
		rm -f $(APP_NAME)
		docker compose down -v

build-prod:
		docker build --target prod-run -t $(IMAGE):$(VERSION) -t $(IMAGE):sha-$(GIT_SHA) .

build-debug:
		docker build --target debug-run -t $(IMAGE):debug .

run-prod:
		docker run -d --name $(IMAGE)-prod \
		--network ${NETWORK} \
		--env-file .env \
		-v "$(PWD)/.env:/app/.env:ro" \
		-p 4000:4000 \
		$(IMAGE):$(VERSION)

run-debug:
		docker run --rm -it --name $(IMAGE)-debug \
		--network ${NETWORK} \
		--env-file .env \
		-v "$(PWD)/.env:/app/.env:ro" \
		-p 4000:4000 \
		$(IMAGE):debug bash

tag:
		docker tag $(IMAGE):sha-$(GIT_SHA) $(IMAGE):$(VERSION)

push:
		docker push $(IMAGE):sha-$(GIT_SHA)
		docker push $(IMAGE):$(VERSION)
