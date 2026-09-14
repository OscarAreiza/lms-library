
## HU-05 — Inventory Control & Book Search (`catalog-service`)

**Merge note:** this branch was 32 commits behind `origin/main` and had drifted into a real
conflict — its own `docker-compose.yml`/`nginx.conf` were stale, catalog-only infra configs,
and git's rename detection paired `catalog-service/cmd/api/main.go` against
`access-service/cmd/api/main.go` as the same renamed file (both trace back to
`backend/cmd/api/main.go`), producing nested conflict markers in both files. Resolved by
reconstructing each service's `main.go` cleanly, taking `origin/main`'s full compose/nginx
config, and keeping HU-05's functional `BooksListPage.tsx` search UI over `origin/main`'s
placeholder. `access-service`/`membership-service` verified byte-identical to `origin/main`
after resolution.

**Audited:** `internal/application/usecase/search_books.go`, `internal/domain/catalog/port.go`.

**Finding — same ISP pattern as HU-03 (fixed):** `SearchBooks` only calls `Search`, but
depended on the fat 4-method `BookRepository`.

**Fix applied:** added a narrow `BookSearcher` port (1 method); repointed `SearchBooks`'s
field/constructor to it. Verified with `go build`/`go vet`/`go test` on all three touched
services (`access-service`, `membership-service`, `catalog-service`) — all green.
