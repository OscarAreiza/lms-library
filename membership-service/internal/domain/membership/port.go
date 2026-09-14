package membership

import "context"

// StudentRepository is the driven port for Student persistence.
type StudentRepository interface {
	FindByID(ctx context.Context, id string) (*Student, error)
	FindByDocumentID(ctx context.Context, documentID string) (*Student, error)
	Search(ctx context.Context, query string, page, limit int) (students []*Student, total int, err error)
	Save(ctx context.Context, s *Student) error
}

// StudentRegistrar is the narrow port HU-02's CreateStudent actually needs —
// ISP: a consumer should not depend on Search/FindByID it never calls. Any
// StudentRepository implementation satisfies this automatically.
type StudentRegistrar interface {
	FindByDocumentID(ctx context.Context, documentID string) (*Student, error)
	Save(ctx context.Context, s *Student) error
}
