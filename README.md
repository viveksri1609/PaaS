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
docker compose up -d
```

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

---

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

* Add automated tests for the API and worker flows.
* Replace hardcoded database and Docker settings with environment variables.
* Add app name and image validation.
* Logs and metrics endpoints implemented: `GET /apps/:id/logs` and `GET /apps/:id/metrics`.
* Add a UI for creating, listing, scaling, and deleting apps.
* Improve replica tracking so the primary app container and scale-out containers are handled separately.
* Add better status transitions and retry handling for failed deployments.
* Add cleanup and health checks around orphaned containers and database rows.

