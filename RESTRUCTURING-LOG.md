
## HU-09 — Book Editing (`catalog-service`)

**Merge note:** real content conflicts in `catalog-service` (this branch's `UpdateBook` vs.
`origin/main`'s `CreateBook`/`LoanBookCopy`/`ReturnBookCopy`), same rename-detection artifact
as HU-05 pairing `catalog-service/cmd/api/main.go` against `access-service/cmd/api/main.go`,
and a genuine compile error surfaced only by `go vet` after resolving: `fakeBookRepository`
redeclared between `create_book_test.go` and `update_book_test.go` (both in package
`usecase_test`) — removed the duplicate, kept the richer fake from `create_book_test.go`.
Reused HU-05's already-resolved `docker-compose.yml`/`nginx.conf` (same domain scope).

**Audited:** `internal/application/usecase/update_book.go`, `internal/domain/catalog/port.go`.

**Finding — same ISP pattern as HU-03's UpdateStudent/DeactivateStudent (fixed):** `UpdateBook`
only used `FindByID`+`Save` from the fat 4-method `BookRepository`.

**Fix applied:** added a narrow `BookEditor` port (2 methods); repointed `UpdateBook`'s
field/constructor to it. Verified with `go build`/`go vet`/`go test` on all three affected
services (`access-service`, `membership-service`, `catalog-service`) — all green.

---

## Summary — all 9 HUs revisited

Every `feat/HU-0X-*` branch was brought up to date with `origin/main`, audited against SOLID,
fixed where a real violation existed, and pushed — no branch deleted, each with its own commit
history documenting exactly what was checked and changed:

| HU | Violation found | Fix |
|---|---|---|
| HU-01 | None | — |
| HU-02 | ISP | `StudentRegistrar` narrow port |
| HU-03 | ISP (×2) | `StudentSearcher`, `StudentEditor` narrow ports |
| HU-04 | ISP | `BookRegistrar` narrow port |
| HU-05 | ISP | `BookSearcher` narrow port |
| HU-06 | None (usecase layer) | — (see HU-07 for the related domain-service SRP fix) |
| HU-07 | SRP + ISP + DIP | Split `ReturnRegistrationService` out of `LoanRegistrationService`; `LoanSearcher` |
| HU-08 | ISP + SRP (own branch copy) | `LoanSearcher`; same `ReturnRegistrationService` split |
| HU-09 | ISP | `BookEditor` narrow port |

Recurring pattern across every fix: a use case depended on a repository's full CRUD interface
while only calling 1-2 of its methods — textbook ISP. The one deeper finding (HU-06/07/08)
was a domain service bundling two HUs' responsibilities into one class, forcing both to share
an oversized dependency surface — SRP, with ISP and DIP as symptoms.
