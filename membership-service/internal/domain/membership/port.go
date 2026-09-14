package membership

import "context"

// StudentRepository is the driven port for Student persistence.
type StudentRepository interface {
	FindByID(ctx context.Context, id string) (*Student, error)
	FindByDocumentID(ctx context.Context, documentID string) (*Student, error)
	Search(ctx context.Context, query string, page, limit int) (students []*Student, total int, err error)
	Save(ctx context.Context, s *Student) error
}

// StudentSearcher is the narrow port HU-03's SearchStudents actually needs —
// ISP: it never reads or writes a single Student, only searches. Any
// StudentRepository implementation satisfies this automatically.
type StudentSearcher interface {
	Search(ctx context.Context, query string, page, limit int) (students []*Student, total int, err error)
}

// StudentEditor is the narrow port HU-03's UpdateStudent and DeactivateStudent
// actually need — ISP: neither searches nor registers by document ID. Any
// StudentRepository implementation satisfies this automatically.
type StudentEditor interface {
	FindByID(ctx context.Context, id string) (*Student, error)
	Save(ctx context.Context, s *Student) error
}
