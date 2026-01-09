# One2n SRE bootcamp

## Student API to practice SRE skills

A simple REST service for practicing infra/SRE concepts: env-based config, migrations, logging, versioned APIs, and container-based local setup.

### Requirements

#### Local development
- Docker
- Docker Compose
- GNU Make

#### Kubernetes deployments
- Minikube
- kubectl
- Helm

You can verify if the required tools are installed by running:

`./scripts/install-deps.sh`

## Local Setup

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


### Kubernetes setup (local)

For Kubernetes-based deployments, a local Minikube cluster is used:

```bash
minikube start
```
---

## Kubernetes Deployment (Raw Manifests)


### Manifest Structure

```
k8s/
├── namespaces
│   └── student-api.yml
├── vault
│   └── vault.yml
├── external-secrets
│   ├── secret-store.yml
│   ├── db-external-secret.yml
│   └── vault-token.yml
├── database
│   └── database.yml
└── application
    └── student-api.yml
```

---

### Vault Setup (Dev Mode)

#### For the raw Kubernetes deployment, **Vault was deployed in dev mode** for learning purposes and to speed up the vault -> eso -> k8s secrets wiring process.
---

### Deploy Vault

```bash
kubectl apply -f k8s/vault/
```

Verify Vault is running:

```bash
kubectl get pods -n student-api
```

---

### Access Vault Locally

Port-forward the Vault service:

```bash
kubectl port-forward -n student-api svc/vault 8200:8200
```

In a new terminal:

```bash
export VAULT_ADDR=http://127.0.0.1:8200
vault login root
```

The static `root` token works because Vault is running in dev mode.

---

### Write Secrets to Vault

Database credentials are written manually as a one-time step.

```bash
vault kv put secret/student-api/db \
  username=student_api \
  password=pa55word
```

Verify:

```bash
vault kv get secret/student-api/db
```

---

### Deploy External Secrets Configuration

```bash
kubectl apply -f k8s/external-secrets/
```

This creates:

* `SecretStore` pointing to Vault
* `ExternalSecret` that syncs DB credentials into Kubernetes

---

### Deploy PostgreSQL

```bash
kubectl apply -f k8s/database/
```

PostgreSQL consumes the Kubernetes Secret created by External Secrets Operator.

---

### Deploy the REST API

```bash
kubectl apply -f k8s/application/
```

The API init container runs migrations using the same synced database credentials.

---


### Verification

Check that all pods are running:

``` bash
`kubectl get pods -n student-api`
```

Verify the API:

```bash
curl http://<node-ip>:<node-port>/v1/healthcheck
```


Expected response:

```bash
{"env":"development","status":"available"}
```


