
## HU-02 — Student Registration (`membership-service`)

**Audited:** `internal/domain/membership/port.go`, `internal/application/usecase/create_student.go`.

**Finding — ISP violation (fixed):** `CreateStudent` depended on the fat `StudentRepository`
interface (4 methods: `FindByID`, `FindByDocumentID`, `Search`, `Save`) while only using 2
(`FindByDocumentID`, `Save`) — forced to depend on methods owned by other HUs (`FindByID` for
HU-03, `Search` for HU-05).

**Fix applied:** added a narrow `StudentRegistrar` port (2 methods) in `port.go`, changed
`CreateStudent`'s field and constructor to depend on it instead. The concrete Postgres
adapter needed no changes — it already satisfies the narrower interface structurally.
Verified with `go build`/`go vet`/`go test` — all green.
