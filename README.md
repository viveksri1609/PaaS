# Mini PaaS in Go

This project is a small learning-focused Platform as a Service prototype.
It is meant to help understand the moving parts behind systems like Azure App Service or Heroku.

You will find a Go API, a background worker, PostgreSQL persistence, Docker-based container execution, and Traefik routing.

---

# Objective

The goal of the project is to learn how a basic PaaS control plane works.

The platform currently focuses on:

* Creating apps through an API
* Storing app state in PostgreSQL
* Deploying containers for apps
* Routing traffic through Traefik
* Scaling apps up and down
* Checking container health
* Deleting apps cleanly

---

# Architecture

```text
           User / Client
                ↓
             Go API
                ↓
          PostgreSQL DB
                ↓
        Reconciler Worker
                ↓
          Docker Engine
                ↓
       Running Containers
                ↓
            Traefik
                ↓
        app.localhost
```

---

# Prerequisites

Install:

* Go 1.22 or newer
* Docker Desktop or Docker Engine
* curl or Postman

---

# Project Structure

```text
PaaS/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── db/
│   │   └── db.go
│   ├── docker/
│   │   ├── docker.go
│   │   └── health.go
│   ├── handlers/
│   │   ├── app_handler.go
│   │   └── scale_handler.go
│   ├── models/
│   │   ├── app.go
│   │   └── app_instance.go
│   └── reconciler/
│       ├── reconciler.go
│       ├── deploy.go
│       ├── health.go
│       ├── self_heal.go
│       ├── scale.go
│       ├── scale_up.go
│       └── scale_down.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── README.md
└── temp
```

---

# Setup

Install the Go dependencies used by the project:

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/moby/moby/client
go get github.com/moby/moby/api/types/container
go get github.com/moby/moby/api/types/network
```

Start the supporting services:

```bash
docker compose up --build -d
docker compose ps
```

The API image uses a Go build stage and a non-root distroless runtime; Compose reports it healthy when `/healthz` responds. The before/after image sizes were not measurable in this environment because the Docker daemon was unavailable (`docker image ls` could not connect); once Docker is running, record them with `docker image ls --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}'` before and after `docker compose build api`. To observe a failing health check while keeping the container running, suspend and resume the API process with `docker kill --signal=STOP mini-paas-api` and `docker kill --signal=CONT mini-paas-api`, then watch `docker compose ps`.

---

# Run The Project

Start the API server in one terminal:

```bash
go run ./cmd/api
```

Start the worker in a second terminal:

```bash
go run ./cmd/worker
```

The API listens on `localhost:8081`.

---

# How To Test

1. Create an app.

```bash
curl -X POST http://localhost:8081/apps \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "hello-app",
    "image": "nginx:latest"
  }'
```

2. List apps and confirm the record exists.

```bash
curl http://localhost:8081/apps
```

3. Wait for the worker to reconcile the app state and move it from pending to running.

4. Open the app through Traefik.

```bash
curl http://hello-app.localhost
```

5. Scale the app.

```bash
curl -X POST http://localhost:8081/apps/<app-id>/scale \
  -H 'Content-Type: application/json' \
  -d '{
    "replicas": 3
  }'
```

6. Confirm the app and its replicas are running.

```bash
curl http://localhost:8081/apps
```

7. Delete the app.

```bash
curl -X DELETE http://localhost:8081/apps/<app-id>
```

8. Run the test suite.

```bash
GOCACHE=/private/tmp/gocache go test ./...
```

9. Optional: watch worker logs while it reconciles state.

```bash
docker logs -f <worker-container-name>
```



# Logs & Metrics

The API provides endpoints to fetch container logs and runtime metrics for an app and its replicas.

- Get recent logs for an app (includes primary container + replicas):

```bash
curl "http://localhost:8081/apps/<app-id>/logs?tail=200"
```

Query params:

- `tail` (optional): number of log lines to return (default: `100`).

- Get runtime metrics (CPU%, memory usage, memory%, restarts, and status):

```bash
curl http://localhost:8081/apps/<app-id>/metrics
```

The metrics endpoint returns an array of per-container metrics for the app and its replicas.

---

# Current Notes

* The API handles app creation, listing, scaling, and deletion.
* The worker continuously deploys and reconciles app state.
* Docker is used directly to create and manage containers.
* Traefik handles local routing through `*.localhost`.


---

# Next Steps

1. Solidify the Docker-based control plane
   * finish reconciliation and error handling in `internal/reconciler/`
   * enhance container health checks and runtime status reporting
   * add focused tests for app creation, scaling, and deletion

2. Add a runtime abstraction layer
   * define a runtime interface for deploy, scale, delete, and status operations
   * keep the existing Docker backend while preparing future runtimes
   * reduce direct Docker engine coupling in the reconciler

3. Learn Kubernetes fundamentals
   * run a local cluster with Docker Desktop, `minikube`, or `kind`
   * practice `kubectl` commands, `Deployment`, `Service`, `Ingress`, and `Namespace`
   * learn `ConfigMap`, `Secret`, and pod readiness/health concepts

4. Implement Kubernetes support
   * add a Kubernetes backend, e.g. `internal/k8s/`
   * map apps to Kubernetes resources: `Deployment` + `Service` + `Ingress`
   * use `k8s.io/client-go` or `sigs.k8s.io/controller-runtime`
   * preserve the current API, DB model, and app lifecycle logic

5. Expand PaaS features on Kubernetes
   * support namespaces, autoscaling, env vars/config, and persistent storage
   * add ingress-based host routing and rollout/update behavior
   * surface app status, logs, and metrics through the API

Recommended learning resources:

* Kubernetes official tutorials: https://kubernetes.io/docs/tutorials/
* kind: https://kind.sigs.k8s.io/
* `kubectl` cheat sheet and local cluster practice

---
