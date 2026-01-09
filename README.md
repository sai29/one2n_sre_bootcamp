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
- Vault

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
---


## Deploying the Stack Using Helm


---

### Prerequisites

* Kubernetes cluster (tested with Minikube)
* `kubectl`
* `helm`


---

### Helm Chart Structure

```
helm/
├── namespaces     # Namespace creation (Helm-owned)
├── vault          # HashiCorp Vault (community chart + values)
├── postgres       # PostgreSQL (Bitnami chart, vendored)
├── secrets        # SecretStore + ExternalSecret (Vault → K8s wiring)
└── student-api    # REST API (Deployment, Service, ConfigMap)
```

Each chart has a **single responsibility** and clear boundaries.

---

### One-Time Manual Steps (Required)

Vault bootstrap actions are **intentionally not Helm-managed** because they are unsafe to replay.

1. Initialize Vault:

   ```
   kubectl exec -n student-api -it vault-0 -- vault operator init
   ```

2. Unseal Vault (run 3 times with different keys out of the 5 generated in the previous step):

   ```
   kubectl exec -n student-api -it vault-0 -- vault operator unseal
   ```

3. Create a Vault policy for External Secrets Operator:

   ```
   kubectl exec -n student-api -it vault-0 -- vault policy write eso-policy - <<EOF
   path "secret/data/student-api/*" {
     capabilities = ["read"]
   }
   EOF
   ```

4. Store the Vault token as a Kubernetes Secret:

   ```
   kubectl create secret generic vault-eso-token \
     -n student-api \
     --from-literal=token=<VAULT_TOKEN>
   ```

---

### Deployment Order (Authoritative)

Helm charts must be installed in the following order:

1. Namespace

   ```
   helm install namespaces ./helm/namespaces
   ```

2. Vault

   ```
   helm install vault hashicorp/vault \
     -n student-api \
     -f helm/vault/values.yaml
   ```

3. External Secrets Operator

   ```
   helm install external-secrets external-secrets/external-secrets \
     -n external-secrets
   ```

4. Secrets (Vault → Kubernetes wiring)

   ```
   helm install secrets ./helm/secrets -n student-api
   ```

5. PostgreSQL

   ```
   helm install postgres ./helm/postgres/postgres \
     -n student-api \
     -f helm/postgres/values.yml
   ```

6. REST API

   ```
   helm install student-api ./helm/student-api -n student-api
   ```

---

### Verification

Check that all pods are running:

```
kubectl get pods -n student-api
```

Verify the API:

```
curl http://<node-ip>:<node-port>/v1/healthcheck
```

Expected response:

```
{"env":"development","status":"available"}
```

---
