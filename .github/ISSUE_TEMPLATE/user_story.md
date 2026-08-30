---
name: User Story
about: Implement a User Story from the source-of-truth backlog
title: "[HU-XX] Short title"
labels: user-story
---

## User Story

<!-- Copy the exact HU from the docs repo — do not paraphrase -->

**Source:** `library-docs/04-requirements/user-stories.md` — HU-XX

**As** the administrator
**I want** ...
**so that** ...

## Acceptance Criteria

<!-- Copy the Gherkin scenarios from the source HU -->

```gherkin
Scenario 1: ...
  Given ...
  When  ...
  Then  ...
```

## Scope for this issue

- [ ] Module(s) touched: `access` / `membership` / `catalog` / `circulation`
- [ ] Backend endpoint(s): `...` (see `library-docs/07-api/contracts/openapi/library-api.yaml`)
- [ ] Frontend screen(s): `...` (see `library-docs/12-ux-ui/navigation-map.md`)

## Definition of Done

- [ ] Code reviewed and approved (1+ approval)
- [ ] Unit tests written and passing
- [ ] Acceptance criteria verified
- [ ] `04-requirements/traceability-matrix.md` status updated (in the docs repo) once merged
