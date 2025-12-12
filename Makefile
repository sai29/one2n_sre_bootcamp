ifneq (,$(wildcard .env))
	include .env
	export
endif

APP_NAME=student_api
MIGRATIONS_DIR=./migrations

# ---- Docker image variables ----
IMAGE ?= student-api
VERSION ?= v0.1.0
GIT_SHA ?= $(shell git rev-parse --short HEAD)

.PHONY: build run test dev migrate migrate-down docker-up docker-down clean \
				build-prod build-debug run-prod run-debug tag push

build:
		go build -o $(APP_NAME)

run:
		go run ./cmd/api

test:
		go test ./... -v -cover

dev:
		air

migrate:
		migrate -path=$(MIGRATIONS_DIR) -database="$(STUDENT_API_DB_DSN)" up

migrate-down:
		migrate -path=$(MIGRATIONS_DIR) -database="$(STUDENT_API_DB_DSN)" down 1

docker-up:
		docker compose up -d

docker-down:
		docker compose down

clean:
		rm -f $(APP_NAME)
		docker compose down -v

build-prod:
		docker build --target prod-run -t $(IMAGE):$(VERSION) -t $(IMAGE):sha-$(GIT_SHA) .

build-debug:
		docker build --target debug-run -t $(IMAGE):debug .

run-prod:
		docker run -d --name $(IMAGE)-prod \
		--network one2n_sre_bootcamp_default \
		--env-file .env \
		-v "$(PWD)/.env:/app/.env:ro" \
		-p 4000:4000 \
		$(IMAGE):$(VERSION)

run-debug:
		docker run --rm -it --name $(IMAGE)-debug \
		--network one2n_sre_bootcamp_default \
		--env-file .env \
		-v "$(PWD)/.env:/app/.env:ro" \
		-p 4000:4000 \
		$(IMAGE):debug bash

tag:
		docker tag $(IMAGE):sha-$(GIT_SHA) $(IMAGE):$(VERSION)

push:
		docker push $(IMAGE):sha-$(GIT_SHA)
		docker push $(IMAGE):$(VERSION)
