package access_test

import (
	"testing"

	"github.com/OscarAreiza/lms-library/access-service/internal/domain/access"
)

func TestNewAdministrator_HashesPassword(t *testing.T) {
	// INV-001: the plaintext password is never persisted.
	admin, err := access.NewAdministrator("libadmin", "SecurePass123!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if admin.PasswordHash == "SecurePass123!" {
		t.Fatal("password must be hashed, not stored in plaintext")
	}
}

func TestAuthenticate_SucceedsWithCorrectPassword(t *testing.T) {
	admin, _ := access.NewAdministrator("libadmin", "SecurePass123!")

	if err := admin.Authenticate("SecurePass123!"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthenticate_FailsWithGenericErrorOnWrongPassword(t *testing.T) {
	// INV-002: never reveal whether the username or the password was wrong.
	admin, _ := access.NewAdministrator("libadmin", "SecurePass123!")

	if err := admin.Authenticate("WrongPassword"); err != access.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
