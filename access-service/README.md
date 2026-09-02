# Access Service

Go REST API implementing the Access bounded context (library-docs/02-domain/domain-map.md):
authenticates the Administrator, the single role that operates the system. Hexagonal
Architecture — see `library-docs/05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md`.

Part of the microservices split: this service owns the `administrators` table exclusively
(its own PostgreSQL database, `access_db`) and issues the JWT that every other service
validates.

## Structure

```
cmd/api/                 → entry point (main.go)
internal/
├── domain/access/        → Administrator aggregate, ports
├── application/usecase/  → Login
├── config/                → environment variable loading
└── infrastructure/
    ├── http/               → chi router, middleware, handlers (primary adapters)
    ├── postgres/            → AdministratorRepository (secondary adapter)
    ├── auth/                 → JWT issuer
    └── logger/                → structured (zap) logger
migrations/                → SQL migrations (golang-migrate)
```

## Tech Stack

* **Language:** Go 1.25
* **Router:** chi
* **Database driver:** pgx (PostgreSQL)
* **Auth:** JWT (golang-jwt), bcrypt

## Development

```bash
go mod download
go run ./cmd/api/...
```

Or, with the rest of the stack, from the repo root:

```bash
docker compose up --build access-service
```

## Tests

```bash
go test ./...
```

## Correlations

* Domain rules → `library-docs/02-domain/entities-and-rules.md`
* API contract → `library-docs/07-api/contracts/openapi/library-api.yaml`
