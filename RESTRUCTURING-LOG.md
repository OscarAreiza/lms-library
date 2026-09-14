
## HU-03 — Student Search, Editing & Deactivation (`membership-service`)

**Audited:** `internal/application/usecase/{search_students,update_student,deactivate_student}.go`,
`internal/domain/membership/port.go`.

**Finding — same ISP pattern as HU-02 (fixed):** `SearchStudents` only calls `Search`;
`UpdateStudent` and `DeactivateStudent` only call `FindByID`+`Save` — all three depended on
the fat 4-method `StudentRepository`.

**Fix applied:** added `StudentSearcher` (1 method) and `StudentEditor` (2 methods) narrow
ports; repointed the three use cases' constructors/fields to them. `SuspendStudent`
(same 2-method shape) is left depending on the fat interface for now — it belongs to HU-08's
scope, not HU-03's, and will get `StudentEditor` when HU-08 is revisited, to keep each fix
inside its own HU. Verified with `go build`/`go vet`/`go test` — all green.
