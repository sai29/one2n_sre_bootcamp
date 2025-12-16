# One2n SRE bootcamp

## Student API to practice SRE skills

A simple REST service for practicing infra/SRE concepts: env-based config, migrations, logging, versioned APIs, and container-based local setup.

### Requirements

- `Docker`
- `Docker Compose`
- `GNU Make`

You can verify if the required tools are installed by running:

`./scripts/install-deps.sh`

## Setup

1. Create a `.env` file and set 
  ```
  STUDENT_API_DB_DSN=postgres://student_api:pa55word@db:5432/student_api?sslmode=disable
  SERVER_PORT=4000
  GIN_MODE=release
  ```

### Start Application 

- `make dev`

This will:

- start Postgres

- run database migrations

- build the API image

- start the API container

API will be available at: http://localhost:4000

### Server logs

- `docker compose logs -f api`

### Stop the Server 

- `make down`

### Makefile

The primary supported workflow is:

- `make dev` — start the full local development environment

Other targets exist for experimentation, debugging, or to simulate prod runs.
See the `Makefile` for details.

### Production run local 

These commands simulate a production-style container locally.

1. First build the container with `make build-prod VERSION=version`, local semver tag (change version as needed).
2. Then run `make run-prod` to run the container.

### Debug image

The debug image includes a shell and is useful for inspecting runtime behavior locally.

1. First build the container with `make build-debug`.
2. Then run `make run-debug` to run the container and get shell access.

### Postman Collection
Import this file into Postman:  
[`student_api.postman_collection.json`](./postman/student_api.postman_collection.json)








