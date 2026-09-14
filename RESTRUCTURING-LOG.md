
## HU-07 — Return Registration & History Tracking (`circulation-service`)

**Merge note:** same rename/rename conflict pattern as HU-06's merge (independent extraction
of `circulation-service` from `backend/` vs. `origin/main`'s independent extraction of the
other three services). Resolved the same way; reused HU-06's already-correct combined
`docker-compose.yml`/`nginx.conf`.

**Audited:** `internal/domain/service/loan_registration_service.go`,
`internal/application/usecase/{return_loan,search_loans}.go`, `internal/domain/circulation/port.go`.

**Finding — the SRP violation flagged (and deferred) during HU-06's audit — fixed here:**
`LoanRegistrationService` bundled `RegisterLoan` (HU-06) and `RegisterReturn` (HU-07/HU-08) in
one class, forcing `ReturnLoan` to depend on the full `StudentClient` (`IsEligible`+`Suspend`),
`BookClient` (`LoanCopy`+`ReturnCopy`), and 4-method `LoanRepository` — three-quarters of that
surface belonging to HU-06, not HU-07.

**Fix applied:** split the domain service in two. `return_registration_service.go` (renamed
from `loan_registration_service.go`, which this branch never used for `RegisterLoan` anyway —
HU-06's own usecase lives on a separate branch) now defines `ReturnRegistrationService` with
three narrow, single-purpose ports defined at point of use: `LoanStore` (`FindByID`+`Save`),
`StudentSuspender` (`Suspend` only), `BookReturner` (`ReturnCopy` only) — mirroring the same
pattern `RegisterLoan`'s own `StudentResolver`/`BookResolver` already used correctly on HU-06's
branch. Also added `LoanSearcher` (1 method) for `SearchLoans`, same ISP pattern as HU-03/HU-05.
`register_loan.go`/`LoanRegistrationService.RegisterLoan` will get the equivalent split
(`StudentEligibilityChecker`, `BookLoaner`, `ActiveLoanCounter`) when these branches are
reconciled at integration time — expected, since each HU branch owns an independent copy of
this shared domain-service file until then. Verified with `go build`/`go vet`/`go test` on all
four services — all green.
