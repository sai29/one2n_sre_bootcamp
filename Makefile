ifneq (,$(wildcard .env))
	include .env
	export
endif

APP_NAME=student_api
MIGRATIONS_DIR=./migrations

build:
		go build -o $(APP_NAME)

run:
		go run ./cmd/api

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