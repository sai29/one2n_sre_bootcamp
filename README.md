# One2n SRE bootcamp

## Student API to practice SRE skills

A simple REST service for practicing infra/SRE concepts: env-based config, migrations, logging, versioned APIs, and container-based local setup.

### Requirements

- `Go` (1.22+)
- `Docker` + `Docker Compose`
- `air` (for hot reload)
- `migrate` (DB schema migrations)

## Setup

### Install air  

- `curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh`

### Install migrate

- `curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvz`

### Steps

1. Run `go build` to install dependencies.
2. Run `docker compose up -d` to run the `docker` setup for setting up the Postgres DB.
3. Create a `.env` file and set 
  ```
  STUDENT_API_DB_DSN=postgres://user:pass@localhost:5432/student_api?sslmode=disable
  SERVER_PORT=8080
  ```
4. Run `migrate -path=./migrations -database=$STUDENT_API_DB_DSN up` to run the migrations.
5. Run `air` to start the server and run the API.

### Makefile

Check the `Makefile` all the commands, some are mentioned below.

```
make build
make run
make dev
make test
make migrate
make docker-up
```

### Production deployment

1. To run the docker container on prod first build the container with `make build-prod VERSION=version`, local semver tag (change version as needed)
2. Then run `make run-prod` to run the container
### Postman Collection
Import this file into Postman:  
[`student_api.postman_collection.json`](./postman/student_api.postman_collection.json)








