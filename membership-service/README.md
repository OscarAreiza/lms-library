# Membership Service

Go REST API implementing the Membership bounded context (library-docs/02-domain/domain-map.md):
owns student records — who is authorized to borrow, and their suspension status. Hexagonal
Architecture — see `library-docs/05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md`.

Part of the microservices split: this service owns the `students` table exclusively (its own
PostgreSQL database, `membership_db`). It validates the JWT issued by access-service using the
shared `JWT_SECRET` — no call to access-service is needed to check a token.

## Structure

```
cmd/api/                 → entry point (main.go)
internal/
├── domain/
│   ├── membership/        → Student aggregate, ports
│   └── shared/             → Value Objects (Email) — duplicated per service, no shared Go module
├── application/usecase/  → CreateStudent (HU-02)
├── config/                → environment variable loading
└── infrastructure/
    ├── http/               → chi router, middleware, handlers (primary adapters)
    ├── postgres/            → StudentRepository (secondary adapter)
    └── logger/                → structured (zap) logger
migrations/                → SQL migrations (golang-migrate)
```

## Tech Stack

* **Language:** Go 1.25
* **Router:** chi
* **Database driver:** pgx (PostgreSQL)

## Development

```bash
go mod download
go run ./cmd/api/...
```

Or, with the rest of the stack, from the repo root:

```bash
docker compose up --build membership-service
```

## Tests

```bash
go test ./...
```

## Correlations

* Domain rules → `library-docs/02-domain/entities-and-rules.md`
* API contract → `library-docs/07-api/contracts/openapi/library-api.yaml`
