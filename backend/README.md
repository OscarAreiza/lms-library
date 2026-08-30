# Backend — `library-api`

Go REST API implementing the LMS-LIBRARY-V1 domain: authentication, book catalog, student
registry, and the loan/return/suspension lifecycle. Hexagonal Architecture, one internal module
per bounded context — see `library-docs/05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md`.

## Structure

```
cmd/api/                 → entry point (main.go)
internal/
├── domain/
│   ├── shared/          → Value Objects shared across modules (Email)
│   ├── access/          → Administrator (Access — Generic)
│   ├── membership/       → Student (Membership — Supporting)
│   ├── catalog/          → Book (Catalog — Supporting)
│   ├── circulation/      → Loan (Circulation — Core Domain)
│   ├── event/             → In-process domain events
│   └── service/            → Cross-module Domain Services (LoanRegistrationService)
├── config/                → Environment variable loading
└── infrastructure/
    ├── http/               → chi router, middleware, handlers (primary adapters)
    ├── postgres/            → repository implementations (secondary adapters)
    └── logger/               → structured (zap) logger
migrations/                → SQL migrations (golang-migrate)
```

Each module's repository is the **only** code allowed to touch its own table — see
`library-docs/09-microservices/service-boundary-rules.md`.

## Tech Stack

* **Language:** Go 1.22
* **Router:** chi
* **Database driver:** pgx (PostgreSQL)
* **Auth:** JWT (golang-jwt)
* **Logging:** zap

## Development

```bash
go mod download
go run ./cmd/api/...
```

Or, with the rest of the stack, from the repo root:

```bash
docker compose up --build backend
```

## Tests

```bash
make test          # go test ./...
make test-cover     # with coverage report
```

## Correlations

* Domain rules → `library-docs/02-domain/entities-and-rules.md`
* Data model / migrations → `library-docs/06-data/models.md`
* API contract → `library-docs/07-api/contracts/openapi/library-api.yaml`
* Testing strategy → `library-docs/11-quality/testing-strategy.md`
