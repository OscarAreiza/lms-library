# Restructuring log — SOLID & design patterns retrofit

> One entry per HU, added on that HU's own `feat/HU-0X-*` branch as it's revisited.
> Branches are never deleted — this file is the running record of what was
> actually checked/changed on each one, in execution order.

## HU-01 — Administrator Authentication (`access-service`)

**Audited:** SRP, OCP, LSP, ISP, DIP against `internal/domain/access/administrator.go`,
`internal/application/usecase/login.go`, `internal/infrastructure/http/handler/auth_handler.go`,
`internal/domain/access/port.go`, `cmd/api/main.go`.

**Result: no changes required.**
- DIP: `Login` depends on `AdministratorRepository`/`TokenIssuer` interfaces, never on
  concrete infrastructure types.
- SRP: each unit (usecase, handler, aggregate, composition root) does exactly one thing.
- OCP: `TokenIssuer` is swappable without touching `Login`.

**Related finding (not a SOLID issue, tracked separately):** `AdministratorRepository.Save`
and `access.NewAdministrator` already exist and are unit-tested, but nothing in production
code calls them — there is no create-administrator use case, handler, or route. The only
administrator today comes from a hardcoded seed migration. This is the actual gap behind
the "admin access control panel" requirement — tracked as its own story, **HU-10**, not as
part of this HU-01 retrofit.
