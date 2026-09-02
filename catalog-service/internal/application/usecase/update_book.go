package usecase

import (
	"context"

	"github.com/OscarAreiza/lms-library/catalog-service/internal/domain/catalog"
)

// UpdateBook implements HU-09's acceptance criteria. ISBN is intentionally not
// editable here — it is Book's stable business key
// (library-docs/02-domain/entities-and-rules.md, modeling note on Book) — so
// there is no "duplicate ISBN on edit" case to reject; the API contract
// (library-docs/07-api/contracts/openapi/library-api.yaml, UpdateBookRequest)
// only exposes title/author/category/year.
type UpdateBook struct {
	Books catalog.BookRepository
}

func NewUpdateBook(books catalog.BookRepository) *UpdateBook {
	return &UpdateBook{Books: books}
}

func (uc *UpdateBook) Execute(ctx context.Context, id, title, author, category string, year int) (*catalog.Book, error) {
	book, err := uc.Books.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := book.Update(title, author, category, year); err != nil {
		return nil, err
	}

	if err := uc.Books.Save(ctx, book); err != nil {
		return nil, err
	}
	return book, nil
}
