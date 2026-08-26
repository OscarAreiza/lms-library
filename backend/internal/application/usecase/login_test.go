package usecase_test

import (
	"context"
	"testing"

	"github.com/OscarAreiza/lms-library/backend/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/access"
)

// fakeAdministratorRepository — Fake test double (11-quality/tdd-guide.md).
type fakeAdministratorRepository struct {
	byUsername map[string]*access.Administrator
}

func (f *fakeAdministratorRepository) FindByUsername(_ context.Context, username string) (*access.Administrator, error) {
	admin, ok := f.byUsername[username]
	if !ok {
		return nil, access.ErrAdministratorNotFound
	}
	return admin, nil
}

func (f *fakeAdministratorRepository) Save(_ context.Context, a *access.Administrator) error {
	f.byUsername[a.Username] = a
	return nil
}

// fakeTokenIssuer — Stub, controls the scenario.
type fakeTokenIssuer struct{}

func (fakeTokenIssuer) Issue(administratorID string) (string, int, error) {
	return "fake-token-for-" + administratorID, 3600, nil
}

func TestLogin_SucceedsWithCorrectCredentials(t *testing.T) {
	admin, _ := access.NewAdministrator("admin@lms.com", "admin123")
	repo := &fakeAdministratorRepository{byUsername: map[string]*access.Administrator{"admin@lms.com": admin}}
	uc := usecase.NewLogin(repo, fakeTokenIssuer{})

	result, err := uc.Execute(context.Background(), "admin@lms.com", "admin123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AccessToken == "" {
		t.Fatal("expected a non-empty access token")
	}
}

func TestLogin_FailsWithGenericErrorWhenUsernameUnknown(t *testing.T) {
	repo := &fakeAdministratorRepository{byUsername: map[string]*access.Administrator{}}
	uc := usecase.NewLogin(repo, fakeTokenIssuer{})

	_, err := uc.Execute(context.Background(), "unknown@lms.com", "whatever")
	if err != access.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_FailsWithGenericErrorWhenPasswordWrong(t *testing.T) {
	// INV-002: same error as "unknown username" — never distinguish the two.
	admin, _ := access.NewAdministrator("admin@lms.com", "admin123")
	repo := &fakeAdministratorRepository{byUsername: map[string]*access.Administrator{"admin@lms.com": admin}}
	uc := usecase.NewLogin(repo, fakeTokenIssuer{})

	_, err := uc.Execute(context.Background(), "admin@lms.com", "wrong-password")
	if err != access.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
