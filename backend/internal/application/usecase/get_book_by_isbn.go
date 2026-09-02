package usecase

import (
	"context"

	"github.com/OscarAreiza/lms-library/backend/internal/domain/catalog"
)

// GetBookByISBN exposes lookup-by-natural-key over HTTP — the ISBN is the
// identifier a human (the Administrator) actually has on hand, unlike the
// internal UUID. Used by circulation-service to resolve a loan request
// entered by ISBN, and by the frontend loan form to avoid asking the
// Administrator for a raw UUID.
type GetBookByISBN struct {
	Books catalog.BookRepository
}

func NewGetBookByISBN(books catalog.BookRepository) *GetBookByISBN {
	return &GetBookByISBN{Books: books}
}

func (uc *GetBookByISBN) Execute(ctx context.Context, isbn string) (*catalog.Book, error) {
	return uc.Books.FindByISBN(ctx, isbn)
}
