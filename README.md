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
## Deployment Methods (Learning Journey)

This section documents the evolution of deployment approaches used in this project.

## Kubernetes Deployment (Legacy - For Reference Only)

This section documents the initial deployment approach using raw Kubernetes manifests. 
**Current deployment method:** See [ArgoCD section](#deploying-with-argocd-gitops) below


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

### Vault Setup (Dev Mode) (Legacy - For reference only)

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



## Deploying the Stack Using Helm  (Legacy)

This section documents the approach using helm chart. 
**Current deployment method:** See [ArgoCD section](#deploying-with-argocd-gitops) below


---

### Prerequisites

* Kubernetes cluster (tested with Minikube)
* `kubectl`
* `helm`


---

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
## Deploying with ArgoCD (GitOps)

This project uses ArgoCD for GitOps-based deployments. All infrastructure and applications are managed via ArgoCD Applications defined in `helm/argocd/apps/`.





### Prerequisites

* Kubernetes cluster (tested with Minikube)
* `kubectl`
* ArgoCD CLI (optional, for CLI access)

### Install ArgoCD

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

Wait for ArgoCD to be ready.

### Access ArgoCD UI

Get the admin password:

```bash
 argocd admin initial-password -n argocd
```

Port-forward ArgoCD server:

```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

Access ArgoCD UI at https://localhost:8080 
- Username: `admin`
- Password: (from the command above)


### Deploy Root Application

The root application syncs all other applications from the Git repository:

```bash
kubectl apply -f helm/argocd/root-app.yaml
```

This will automatically deploy:
- Infrastructure (namespaces, Vault, PostgreSQL, External Secrets Operator)
- Applications (Student API, secrets, migrations)
- Observability stack (Prometheus, Loki, Grafana, Promtail, exporters)

### One-Time Manual Steps - Vault setup

**Note:** Vault is setup in memory for learning purposes.

Before deploying applications, complete these one-time setup steps:

#### 1. Initialize and Unseal Vault

After Vault is deployed (via ArgoCD), initialize it:

```bash
kubectl exec -n vault -it vault-0 -- vault operator init
```

Unseal Vault (run 3 times with different unseal keys):

```bash
kubectl exec -n vault -it vault-0 -- vault operator unseal <UNSEAL_KEY>
```


#### 2. Login to Vault
Use the root token from the initialization output:

```bash
kubectl exec -n vault -it vault-0 -- vault login <ROOT_TOKEN>
```
**Note**: Root token is used for simplicity in this learning setup. In production, you would create policies and use service-specific tokens.

#### 3. Create Vault KV Secret Engine

Create a KV v2 secret engine at path `kv`:

```bash
kubectl exec -n vault -it vault-0 -- vault secrets enable -path=kv kv-v2
```

#### 4. Store Database Credentials in Vault

```bash
kubectl exec -n vault -it vault-0 -- vault kv put kv/student-api/db \
  username=student_api \
  password=pa55word \
  database=student_api \
  POSTGRES_PASSWORD=pa55word
```


#### 5. Store Vault root Token for External Secrets Operator


Store the token as a Kubernetes Secret:

```bash
kubectl create secret generic vault-eso-token \
  -n student-api \
  --from-literal=token=<ROOT_TOKEN>
```
**Note**: Root token is used for simplicity in this learning setup. In production, you would create policies and use service-specific tokens.

---

### Application Structure

```
helm/argocd/apps/
├── infra/          # Infrastructure components (sync-wave: -2 to -1)
│   ├── namespaces.yaml
│   ├── postgres.yaml
│   ├── vault.yaml
│   └── external-secrets.yaml
├── apps/            # Application components (sync-wave: 0 to 2)
│   ├── student-api.yaml
│   ├── student-api-secrets.yaml
│   └── student-api-migrations.yaml
└── observability/  # Observability stack (sync-wave: 1 to 4)
    ├── prometheus-crds.yaml
    ├── kube-prometheus-stack.yaml
    ├── loki.yaml
    ├── promtail.yaml
    ├── postgres-exporter.yaml
    ├── blackbox-exporter.yaml
    └── probes.yaml
```

Applications are deployed in order based on `sync-wave` annotations to ensure dependencies are created first.

### Verification


Verify all pods are running:

```bash
kubectl get pods -A
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

## Observability Stack Setup

The observability stack provides comprehensive monitoring, logging, and visualization for the entire system. 
**Note**: The stack gets installed automatically when you apply the root Argo CD app.

### Components

- **Prometheus** - Metrics collection and storage
- **Loki** - Log aggregation and storage  
- **Grafana** - Visualization and dashboards
- **Promtail** - Log shipper (collects logs and sends to Loki)
- **Postgres Exporter** - Database metrics exporter
- **Blackbox Exporter** - Endpoint uptime and latency monitoring
- **Node Exporter** - Node-level system metrics
- **Kube-state-metrics** - Kubernetes object metrics

### Deployment

The observability stack is automatically deployed via ArgoCD when you deploy the root application. All components are deployed in the `observability` namespace.

**Deployment order (via sync-waves):**
1. Prometheus CRDs (sync-wave: 1)
2. Kube-prometheus-stack (sync-wave: 2) - Includes Prometheus, Grafana, Alertmanager, Node Exporter, Kube-state-metrics
3. Loki (sync-wave: 3)
4. Postgres Exporter (sync-wave: 3)
5. Blackbox Exporter (sync-wave: 3)
6. Promtail (sync-wave: 4)
7. Blackbox Probes (sync-wave: 4)

### Accessing Grafana

Port-forward Grafana service:

```bash
kubectl port-forward -n observability svc/kube-prometheus-stack-grafana 3000:80
```

Access Grafana at http://localhost:3000
- Username: `admin`
- Password: `admin` (default)

### Accessing Prometheus

Port-forward Prometheus service:

```bash
kubectl port-forward -n observability svc/kube-prometheus-stack-prometheus 9090:9090
```

Access Prometheus at http://localhost:9090

### Grafana Data Sources

Grafana is pre-configured with two data sources:
- **Prometheus** - For metrics queries (URL: `http://kube-prometheus-stack-prometheus.observability.svc.cluster.local:9090`)
- **Loki** - For log queries (URL: `http://loki:3100`)

### Monitoring Endpoints

Blackbox exporter monitors the following endpoints via Probe CRDs:

- **Student API** - `http://student-api.student-api.svc.cluster.local:4000/v1/healthcheck`
- **ArgoCD Server** - `https://argocd-server.argocd.svc.cluster.local`
- **Vault** - `http://vault.vault.svc.cluster.local:8200/v1/sys/health`

Probe configurations are defined in `manifests/probes/` and are automatically synced via ArgoCD as a separate Argo CD app.

### Log Collection

Promtail is configured to collect logs only from the `student-api` namespace and send them to Loki. Logs are filtered to include only application logs.

### Verification

Check that all observability components are running:

```bash
kubectl get pods -n observability
```

Verify Prometheus is scraping targets:

```bash
kubectl port-forward -n observability svc/kube-prometheus-stack-prometheus 9090:9090
# Open http://localhost:9090 → Status → Targets
```

Check probe metrics:

```bash
# In Prometheus UI, query:
probe_success{job="student-api-healthcheck"}
probe_success{job="vault"}
probe_success{job="argocd-server"}
# Should show probe_success=1 for all targets
```

---
