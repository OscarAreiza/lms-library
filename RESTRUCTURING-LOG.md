
## HU-06 — Loan Registration (`circulation-service`)

**Merge note:** this branch extracted `circulation-service` from `backend/` in parallel with
`origin/main`'s independent extraction of the other three services, causing widespread
rename/rename conflicts (including `go.mod` and `config.go` being paired against
`membership-service`'s/`access-service`'s own files by git's rename detection). Resolved
domain-by-domain: `access-service`/`membership-service`/`catalog-service` verbatim from
`origin/main`, `circulation-service` verbatim from this branch, `backend/` removed,
`docker-compose.yml`/`nginx.conf` manually combined to carry all four domains. Verified with
`go build`/`go vet`/`go test` across all four services — all green.

**Audited:** `internal/application/usecase/register_loan.go`,
`internal/domain/service/loan_registration_service.go`, `internal/domain/circulation/port.go`.

**Result: no changes required at the usecase layer.** `RegisterLoan` already defines
`StudentResolver`/`BookResolver` as narrow, single-method interfaces at the point of use —
textbook ISP, done correctly from the start.

**Related finding (not fixed here, deferred to HU-07):** the domain service
`LoanRegistrationService` bundles `RegisterLoan` (HU-06) and `RegisterReturn` (HU-07/HU-08)
in one struct, forcing both to depend on the same fat `StudentClient` (`IsEligible` +
`Suspend`), `BookClient` (`LoanCopy` + `ReturnCopy`), and `circulation.LoanRepository`
(4 methods) — an SRP violation one level below the usecases. Splitting it touches HU-07's
code as much as HU-06's, so it's deferred to HU-07's own retrofit pass rather than done here.
