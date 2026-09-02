package usecase

import (
	"context"

	"github.com/OscarAreiza/lms-library/catalog-service/internal/domain/catalog"
)

// SearchBooks implements HU-05's acceptance criteria — an empty result set is a
// valid outcome (Scenario 2), never an error.
type SearchBooks struct {
	Books catalog.BookRepository
}

func NewSearchBooks(books catalog.BookRepository) *SearchBooks {
	return &SearchBooks{Books: books}
}

func (uc *SearchBooks) Execute(ctx context.Context, query, category string, page, limit int) ([]*catalog.Book, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.Books.Search(ctx, query, category, page, limit)
}
