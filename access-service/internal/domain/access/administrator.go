// Package access implements the Access bounded context (library-docs/02-domain/domain-map.md):
// the single authenticated role that operates the system.
package access

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// BcryptCost matches NFR-004 (library-docs/04-requirements/non-functional.md): cost >= 12.
const BcryptCost = 12

// ErrInvalidCredentials is the single, generic error for any login failure — INV-002:
// login must never reveal whether the username or the password was wrong.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Administrator is the Access bounded context's Aggregate Root.
// See library-docs/02-domain/entities-and-rules.md, Entity: Administrator.
type Administrator struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewAdministrator hashes the plaintext password (INV-001: password is always hashed;
// the plaintext is never persisted, logged, or returned).
func NewAdministrator(username, plaintextPassword string) (*Administrator, error) {
	if username == "" {
		return nil, errors.New("username must not be empty")
	}
	if plaintextPassword == "" {
		return nil, errors.New("password must not be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), BcryptCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Administrator{
		Username:     username,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Authenticate verifies a plaintext password against the stored hash.
// It never distinguishes "wrong username" from "wrong password" to the caller —
// callers must map any error from this method (and a not-found lookup) to the
// same ErrInvalidCredentials response.
func (a *Administrator) Authenticate(plaintextPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(plaintextPassword)); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}
