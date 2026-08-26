package usecase_test

import (
	"context"
	"testing"

	"github.com/OscarAreiza/lms-library/backend/internal/application/usecase"
)

func TestSearchBooks_ReturnsEmptyResultWithoutError(t *testing.T) {
	// HU-05, Scenario 2: no match is a valid outcome, not an error.
	repo := newFakeBookRepository()
	uc := usecase.NewSearchBooks(repo)

	books, total, err := uc.Execute(context.Background(), "does-not-exist", "", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(books) != 0 || total != 0 {
		t.Fatalf("expected an empty result, got %d books (total=%d)", len(books), total)
	}
}

func TestSearchBooks_MatchesByTitle(t *testing.T) {
	repo := newFakeBookRepository()
	createUC := usecase.NewCreateBook(repo)
	_, err := createUC.Execute(context.Background(), "Clean Code", "Robert C. Martin", "978-0132350884", "Software", 2008, 3)
	if err != nil {
		t.Fatalf("unexpected error seeding fixture: %v", err)
	}

	searchUC := usecase.NewSearchBooks(repo)
	books, total, err := searchUC.Execute(context.Background(), "Clean", "", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 || len(books) != 1 {
		t.Fatalf("expected 1 match, got %d (total=%d)", len(books), total)
	}
}
