# Circulation Service

Go REST API implementing the Circulation bounded context (library-docs/02-domain/domain-map.md)
— the **Core Domain**: loans, returns, and the suspension policy. Hexagonal Architecture —
see `library-docs/05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md`.

Part of the microservices split: this service owns the `loans` table exclusively (its own
PostgreSQL database, `circulation_db`). It validates the JWT issued by access-service using
the shared `JWT_SECRET`.

`LoanRegistrationService` coordinates Student (Membership), Book (Catalog), and Loan
(Circulation) — before the split this was one in-process transaction; now it's:

- `internal/infrastructure/membership/client.go` → HTTP calls to membership-service
  (`GET /students/{id}` for eligibility, `POST /students/{id}/suspend` for late returns)
- `internal/infrastructure/catalog/client.go` → HTTP calls to catalog-service
  (`POST /books/{id}/loan-copy` / `/return-copy`)

Both clients authenticate with a short-lived internal token signed with the shared
`JWT_SECRET` (v1 has no separate service-to-service auth scope — see the client files'
doc comments for the trade-off). There is no distributed transaction across these calls;
a mid-sequence failure is an accepted risk at this project's size.

This branch implements Return (triggering the suspension policy) and the Overdue report
(HU-08); Create lands on HU-06's own branch, and history/search on HU-07's, combining
when merged. `SuspensionDays` is 3 (lowered from 7 per product decision).

## Structure

```
cmd/api/                 → entry point (main.go)
internal/
├── domain/
│   ├── circulation/       → Loan aggregate, LoanRepository port
│   └── service/            → LoanRegistrationService (cross-service coordination)
├── application/usecase/  → ReturnLoan, OverdueLoans (HU-08)
├── config/                → environment variable loading
└── infrastructure/
    ├── http/               → chi router, middleware, handlers (primary adapters)
    ├── postgres/            → LoanRepository (secondary adapter)
    ├── membership/           → HTTP client for the Membership driven port
    ├── catalog/               → HTTP client for the Catalog driven port
    └── logger/                 → structured (zap) logger
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

Needs `membership-service` and `catalog-service` reachable at
`MEMBERSHIP_SERVICE_URL` / `CATALOG_SERVICE_URL` to register a return.

Or, with the rest of the stack, from the repo root:

```bash
docker compose up --build circulation-service
```

## Tests

```bash
go test ./...
```

## Correlations

* Domain rules → `library-docs/02-domain/entities-and-rules.md`
* API contract → `library-docs/07-api/contracts/openapi/library-api.yaml`
