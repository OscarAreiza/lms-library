
## HU-08 — Late-Return Penalty System (`circulation-service`)

**Merge note:** same rename/rename conflict pattern as HU-06/HU-07's merges. Resolved the
same way; reused the already-correct combined `docker-compose.yml`/`nginx.conf`.

**Audited:** `internal/application/usecase/overdue_loans.go`,
`internal/domain/service/loan_registration_service.go`, `internal/domain/circulation/port.go`.

**Finding 1 — ISP (fixed):** `OverdueLoans` only calls `Search`, same pattern as HU-07's
`SearchLoans`. Added `LoanSearcher` (1 method) to this branch's own copy of `port.go`;
repointed `OverdueLoans`'s field/constructor to it.

**Finding 2 — same SRP violation as HU-07 (fixed the same way):** this branch has its own
independent copy of `loan_registration_service.go` (each HU branch owns one, extracted
independently from `backend/`), with the identical `RegisterLoan`/`RegisterReturn` bundling
found and fixed on HU-07's branch. Applied the same split here: renamed to
`return_registration_service.go`, `ReturnRegistrationService` with narrow `LoanStore`,
`StudentSuspender`, `BookReturner` ports. `overdue_loans_test.go` needed no changes — it
already depended on the fake repo directly, which satisfies `LoanSearcher` structurally.
Verified with `go build`/`go vet`/`go test` on all four services — all green.
