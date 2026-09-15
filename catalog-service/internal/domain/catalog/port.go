package catalog

import "context"

// BookRepository is the driven port for Book persistence.
type BookRepository interface {
	FindByID(ctx context.Context, id string) (*Book, error)
	FindByISBN(ctx context.Context, isbn string) (*Book, error)
	Search(ctx context.Context, query, category string, page, limit int) (books []*Book, total int, err error)
	Save(ctx context.Context, b *Book) error
}

// BookEditor is the narrow port HU-09's UpdateBook actually needs — ISP: it
// never looks a book up by ISBN and never searches. Any BookRepository
// implementation satisfies this automatically.
type BookEditor interface {
	FindByID(ctx context.Context, id string) (*Book, error)
	Save(ctx context.Context, b *Book) error
}
