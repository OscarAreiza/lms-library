package access

import (
	"context"
	"errors"
)

// ErrAdministratorNotFound is the port-level "not found" error — use cases
// depend on this, never on an infrastructure package's own error type
// (library-docs/05-architecture/hexagonal-architecture.md, Dependency Rule).
var ErrAdministratorNotFound = errors.New("administrator not found")

// AdministratorRepository is the driven (secondary) port — implemented by an
// infrastructure adapter (e.g. internal/infrastructure/postgres), never called
// directly from another bounded context (library-docs/09-microservices/service-boundary-rules.md).
type AdministratorRepository interface {
	FindByUsername(ctx context.Context, username string) (*Administrator, error)
	Save(ctx context.Context, a *Administrator) error
}

// TokenIssuer is a driven port for issuing a signed session token on successful
// authentication. Implemented by a JWT adapter in infrastructure.
type TokenIssuer interface {
	Issue(administratorID string) (token string, expiresInSeconds int, err error)
}
