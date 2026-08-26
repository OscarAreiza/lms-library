// Package usecase orchestrates the domain for each HU — it contains no business
// logic itself (that lives in the Aggregate), only coordination
// (library-docs/05-architecture/hexagonal-architecture.md).
package usecase

import (
	"context"
	"errors"

	"github.com/OscarAreiza/lms-library/backend/internal/domain/access"
)

// LoginResult is what the HTTP handler needs to build the response.
type LoginResult struct {
	AccessToken string
	ExpiresIn   int
}

// Login implements HU-01's acceptance criteria. Any failure — administrator
// not found, or wrong password — maps to the exact same
// access.ErrInvalidCredentials, per INV-002 on Administrator.
type Login struct {
	Administrators access.AdministratorRepository
	Tokens         access.TokenIssuer
}

func NewLogin(administrators access.AdministratorRepository, tokens access.TokenIssuer) *Login {
	return &Login{Administrators: administrators, Tokens: tokens}
}

func (uc *Login) Execute(ctx context.Context, username, password string) (*LoginResult, error) {
	admin, err := uc.Administrators.FindByUsername(ctx, username)
	if errors.Is(err, access.ErrAdministratorNotFound) {
		return nil, access.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := admin.Authenticate(password); err != nil {
		return nil, access.ErrInvalidCredentials
	}

	token, expiresIn, err := uc.Tokens.Issue(admin.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{AccessToken: token, ExpiresIn: expiresIn}, nil
}
