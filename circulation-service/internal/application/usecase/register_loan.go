package usecase

import (
	"context"

	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/service"
)

// StudentResolver translates the identifier an Administrator actually has on
// hand (the student's document ID) into the internal UUID the domain
// operates on — an Administrator has no way to know the UUID
// (library-docs/09-microservices/data-ownership-matrix.md).
type StudentResolver interface {
	ResolveByDocumentID(ctx context.Context, documentID string) (studentID string, err error)
}

// BookResolver translates an ISBN into the internal UUID the domain operates
// on, for the same reason.
type BookResolver interface {
	ResolveByISBN(ctx context.Context, isbn string) (bookID string, err error)
}

// RegisterLoan implements HU-06's acceptance criteria. The Administrator
// identifies the student and the book by their natural keys (document ID,
// ISBN), never the internal UUID. This use case resolves those to UUIDs
// first, then delegates the actual cross-aggregate coordination (Student
// eligibility, Book availability, Loan creation) to LoanRegistrationService,
// which stays keyed by UUID — per library-docs/02-domain/entities-and-rules.md
// ("Domain Services"), that coordination logic should not need to know or
// care which external identifier a caller used.
type RegisterLoan struct {
	students StudentResolver
	books    BookResolver
	service  *service.LoanRegistrationService
}

func NewRegisterLoan(students StudentResolver, books BookResolver, svc *service.LoanRegistrationService) *RegisterLoan {
	return &RegisterLoan{students: students, books: books, service: svc}
}

func (uc *RegisterLoan) Execute(ctx context.Context, studentDocumentID, bookISBN string) (*circulation.Loan, error) {
	studentID, err := uc.students.ResolveByDocumentID(ctx, studentDocumentID)
	if err != nil {
		return nil, err
	}

	bookID, err := uc.books.ResolveByISBN(ctx, bookISBN)
	if err != nil {
		return nil, err
	}

	return uc.service.RegisterLoan(ctx, studentID, bookID)
}
