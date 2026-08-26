# Pull Request

## Linked User Story

**HU:** HU-XX — <!-- required: no PR merges without a linked HU or a chore/hotfix justification -->
**Issue:** Closes #

## What changed

<!-- One or two sentences — the "why", not a line-by-line diff description -->

## Type of change

- [ ] `feat` — new functionality
- [ ] `fix` — bug fix
- [ ] `chore` — tooling/deps/infra, no user-facing behavior
- [ ] `docs` — documentation only
- [ ] `refactor` — no behavior change
- [ ] `hotfix` — urgent fix to `main`

## Checklist (per `library-docs/00-governance/definition-of-done.md`)

- [ ] Branch is named `[type]/HU-XX-description` (or `chore/description` if not HU-scoped)
- [ ] Commits follow Conventional Commits (`library-docs/00-governance/git-conventions.md`)
- [ ] `go test ./...` / frontend tests pass locally
- [ ] Lint passes (`golangci-lint run ./...` / `npm run lint`)
- [ ] No cross-module table access introduced (`library-docs/09-microservices/service-boundary-rules.md`)
- [ ] API changes reflected in `library-docs/07-api/contracts/openapi/library-api.yaml`
- [ ] This PR is under ~400 LOC (excluding tests) — if not, explain why it wasn't split

## How to test

<!-- Exact commands/steps a reviewer runs to verify this -->
