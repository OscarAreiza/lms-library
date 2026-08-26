package usecase_test

import (
	"context"
	"testing"

	"github.com/OscarAreiza/lms-library/backend/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/catalog"
)

// fakeBookRepository — Fake test double (11-quality/tdd-guide.md).
type fakeBookRepository struct {
	byID map[string]*catalog.Book
}

func newFakeBookRepository() *fakeBookRepository {
	return &fakeBookRepository{byID: map[string]*catalog.Book{}}
}

func (f *fakeBookRepository) FindByID(_ context.Context, id string) (*catalog.Book, error) {
	b, ok := f.byID[id]
	if !ok {
		return nil, catalog.ErrBookNotFound
	}
	return b, nil
}

func (f *fakeBookRepository) FindByISBN(_ context.Context, _ string) (*catalog.Book, error) {
	return nil, catalog.ErrBookNotFound
}

func (f *fakeBookRepository) Search(_ context.Context, _, _ string, _, _ int) ([]*catalog.Book, int, error) {
	return nil, 0, nil
}

func (f *fakeBookRepository) Save(_ context.Context, b *catalog.Book) error {
	f.byID[b.ID] = b
	return nil
}

func seedFakeBook(t *testing.T, repo *fakeBookRepository) *catalog.Book {
	t.Helper()
	book, err := catalog.NewBook("Clean Code", "Robert C. Martin", "978-0132350884", "Software", 2008, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.Save(context.Background(), book); err != nil {
		t.Fatalf("unexpected error saving fixture: %v", err)
	}
	return book
}

func TestUpdateBook_Succeeds(t *testing.T) {
	repo := newFakeBookRepository()
	book := seedFakeBook(t, repo)
	uc := usecase.NewUpdateBook(repo)

	updated, err := uc.Execute(context.Background(), book.ID, "Clean Code (2nd ed.)", "Robert C. Martin", "Software Engineering", 2009)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "Clean Code (2nd ed.)" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}
}

func TestUpdateBook_KeepsISBNUnchanged(t *testing.T) {
	// FR-012: ISBN is immutable through this action — it isn't even a parameter.
	repo := newFakeBookRepository()
	book := seedFakeBook(t, repo)
	originalISBN := book.ISBN
	uc := usecase.NewUpdateBook(repo)

	updated, err := uc.Execute(context.Background(), book.ID, "New Title", "New Author", "New Category", 2020)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.ISBN != originalISBN {
		t.Fatalf("expected ISBN to remain %q, got %q", originalISBN, updated.ISBN)
	}
}

func TestUpdateBook_FailsWhenNotFound(t *testing.T) {
	repo := newFakeBookRepository()
	uc := usecase.NewUpdateBook(repo)

	_, err := uc.Execute(context.Background(), "does-not-exist", "Title", "Author", "Category", 2020)
	if err != catalog.ErrBookNotFound {
		t.Fatalf("expected ErrBookNotFound, got %v", err)
	}
}
