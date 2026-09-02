# Setup Guide

How to get `lms-library` running locally: prerequisites, environment variables, the
Docker images that make up the stack, and how to run each service outside Docker for
development.

## Prerequisites

| Tool | Version | Required for |
|------|---------|--------------|
| [Docker](https://docs.docker.com/get-docker/) + Docker Compose | v2+ | Running the full stack |
| [Go](https://go.dev/dl/) | 1.25+ | Running/testing the backend outside Docker |
| [Node.js](https://nodejs.org/) | 20+ | Running/testing the frontend outside Docker |
| Git | any recent | Cloning the repo |

You do **not** need Go or Node installed to run the project — Docker Compose builds
everything. They're only needed if you want to run a service natively (faster reload
loop, debugger attached, etc.).

## 1. Clone and configure environment variables

```bash
git clone git@github.com:OscarAreiza/lms-library.git
cd lms-library
cp .env.example .env
```

Open `.env` and adjust values as needed. Variables read by `docker-compose.yml`:

| Variable | Default | Used by | Notes |
|----------|---------|---------|-------|
| `POSTGRES_USER` | `lms_user` | db, migrate, backend | |
| `POSTGRES_PASSWORD` | `lms_password` | db, migrate, backend | |
| `POSTGRES_DB` | `lms_db` | db, migrate, backend | |
| `JWT_SECRET` | *(none — required)* | backend | Compose fails fast if unset. Generate one with `openssl rand -base64 32` |
| `JWT_EXPIRY` | `1h` | backend | |
| `LOG_LEVEL` | `info` (`debug` in `.env.example`) | backend | |
| `CORS_ORIGIN` | — | backend | Not yet wired into compose (hardcoded `*` there); kept in `.env.example` for when it is |
| `VITE_API_BASE_URL` | `http://localhost:8080/api/v1` | frontend | Baked into the static bundle at **build time** — changing it after `docker compose up` requires a rebuild |

Never commit the real `.env` file — it's already covered by `.gitignore`.

## 2. Start the stack

```bash
docker compose up -d --build
```

This builds and starts everything in dependency order. First run takes a few minutes
(Go module download, npm install, Postgres image pull); afterwards Docker's layer cache
makes it much faster.

Verify it's up:

```bash
curl http://localhost:8080/health          # backend liveness
curl http://localhost:8080/health/ready    # backend readiness (checks DB connection)
```

Open the app: **http://localhost:3000**

## 3. What's running — services and images

| Service | Image / build | Port (host) | Role |
|---------|---------------|--------------|------|
| `db` | `postgres:16-alpine` | `5432` | PostgreSQL database. Data persists in the `lms-db-data` named volume. |
| `migrate` | `migrate/migrate:v4.17.1` | — | Runs SQL migrations from `backend/migrations/` against `db`, then exits. Backend won't start until this completes successfully. |
| `backend` | built from `backend/Dockerfile` (multi-stage: `golang:1.25-alpine` → `alpine:3.21`) | `8080` | Go REST API (`library-api`). Waits on `db` (healthy) and `migrate` (completed). |
| `frontend` | built from `frontend/Dockerfile` (multi-stage: `node:20-alpine` → `nginx:alpine`) | `3000` (mapped to container's `80`) | React SPA, served as a static build through NGINX. |

All services share the `lms-network` bridge network and reach each other by service
name (e.g. the backend connects to `db:5432`, not `localhost:5432`).

### Startup order

```
db (waits until healthy)
 └── migrate (waits until db is healthy, then runs, then exits)
      └── backend (waits until migrate completes and db is healthy)
           └── frontend (waits until backend container exists — not necessarily healthy)
```

### Health checks

- `db`: `pg_isready` every 5s.
- `backend`: `GET /health` every 10s (5s grace period on startup).

## 4. Common operations

```bash
# Follow logs for one service
docker compose logs -f backend

# Rebuild a single service after changing its code
docker compose up -d --build backend

# Rebuild the frontend after changing VITE_API_BASE_URL in .env
docker compose up -d --build frontend

# Stop everything (keeps the database volume)
docker compose down

# Stop and wipe the database too (fresh migrations next start)
docker compose down -v

# Run a psql shell against the database
docker compose exec db psql -U lms_user -d lms_db
```

## 5. Running services outside Docker (development)

### Backend

```bash
cd backend
go mod download
go run ./cmd/api/...
```

Needs a reachable Postgres — either point `DB_HOST`/`DB_PORT`/etc. env vars at the
Dockerized `db` (exposed on `localhost:5432`) or run `docker compose up -d db migrate`
first and let the native `go run` process connect to it.

Tests:

```bash
make test          # go test ./...
make test-cover     # with coverage report
```

### Frontend

```bash
cd frontend
npm install
cp .env.example .env   # set VITE_API_BASE_URL if the backend isn't on the default port
npm run dev
```

Vite's dev server proxies to whatever `VITE_API_BASE_URL` points to — typically the
Dockerized or natively-running backend at `http://localhost:8080/api/v1`.

## 6. Troubleshooting

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| `docker compose up` fails immediately with a `JWT_SECRET` error | `.env` wasn't created or `JWT_SECRET` is empty | `cp .env.example .env` and set a value |
| `backend` container restarts in a loop | `db` not healthy yet, or migrations failed | `docker compose logs migrate` and `docker compose logs db` |
| Frontend loads but API calls fail (CORS or 404) | `VITE_API_BASE_URL` baked in at build time doesn't match where the backend actually listens | Update `.env`, then `docker compose up -d --build frontend` (env change alone won't take effect — it's compiled into the bundle) |
| Port `3000`, `8080`, or `5432` already in use | Another process/project is bound to that port | Stop the conflicting process, or remap the host port in `docker-compose.yml` (e.g. `"3001:80"`) |
| Changes to backend/frontend code don't show up | Docker is using the cached build layer | `docker compose up -d --build <service>` |

## Related documentation

- Domain rules, requirements, API contracts, ADRs → [`library-docs`](https://github.com/code-corhuila/library-docs)
- Backend architecture details → [`backend/README.md`](backend/README.md)
- Frontend architecture details → [`frontend/README.md`](frontend/README.md)
