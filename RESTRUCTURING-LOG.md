
## HU-04 — Book Registration (`catalog-service`)

**Audited:** `internal/domain/catalog/port.go`, `internal/application/usecase/create_book.go`.

**Finding — same ISP pattern as HU-02/HU-03 (fixed):** `CreateBook` only used `FindByISBN`
and `Save` from the fat 4-method `BookRepository` (`FindByID`/`Search` belong to HU-05/HU-09).

**Fix applied:** added a narrow `BookRegistrar` port (2 methods), repointed `CreateBook`'s
field/constructor to it. `adjust_book_availability.go` (HU-06's concern, lives here because
Catalog owns copy counts) is left on the fat interface for now — will get its own narrow
port when HU-06 is revisited. Verified with `go build`/`go vet`/`go test` — all green.
