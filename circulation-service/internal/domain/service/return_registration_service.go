// Package service holds Domain Services — business logic that coordinates more
// than one bounded context and therefore does not naturally belong inside a
// single Aggregate Root. See library-docs/02-domain/entities-and-rules.md,
// "Domain Services".
//
// Before the microservices split, StudentClient/BookClient were in-process
// repository calls against the same database. Now they're HTTP calls to
// membership-service and catalog-service — the coordination logic here is
// unchanged, only the driven ports' implementations moved to the network
// (library-docs/09-microservices/service-boundary-rules.md).
package service

import (
	"context"
	"time"

	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/circulation"
)

// LoanStore is the narrow port onto Loan persistence that RegisterReturn
// needs — find the loan being returned, then persist its new state. Search
// (HU-07/HU-08's queries) and CountActiveByStudent (HU-06's loan-limit check)
// belong to other use cases, not to this service.
type LoanStore interface {
	FindByID(ctx context.Context, id string) (*circulation.Loan, error)
	Save(ctx context.Context, l *circulation.Loan) error
}

// StudentSuspender is the narrow port onto membership-service that
// RegisterReturn needs — only the HU-08 suspension side effect. Checking
// eligibility to take out a new loan (IsEligible) is HU-06's concern.
type StudentSuspender interface {
	Suspend(ctx context.Context, studentID string, days int) error
}

// BookReturner is the narrow port onto catalog-service that RegisterReturn
// needs — only restoring availability. Decrementing it (LoanCopy) is HU-06's
// concern.
type BookReturner interface {
	ReturnCopy(ctx context.Context, bookID string) error
}

// ReturnRegistrationService implements HU-07 and HU-08's acceptance criteria:
// closes the loan, restores Book availability, and — if the return was late —
// applies Student's flat suspension. This is the one place in the codebase
// where those three are checked/changed together for a return.
type ReturnRegistrationService struct {
	loans    LoanStore
	students StudentSuspender
	books    BookReturner
}

func NewReturnRegistrationService(loans LoanStore, students StudentSuspender, books BookReturner) *ReturnRegistrationService {
	return &ReturnRegistrationService{loans: loans, students: students, books: books}
}

// RegisterReturn implements HU-07 and HU-08's acceptance criteria: closes the
// loan, restores availability, and — if the return was late — applies the
// flat suspension.
func (s *ReturnRegistrationService) RegisterReturn(ctx context.Context, loanID string) (*circulation.Loan, error) {
	loan, err := s.loans.FindByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	returnTime := time.Now().UTC()
	wasLate, err := loan.RegisterReturn(returnTime) // INV-005
	if err != nil {
		return nil, err
	}

	if err := s.books.ReturnCopy(ctx, loan.BookID); err != nil {
		return nil, err
	}

	if wasLate { // INV-006 on Loan: late return triggers a fixed suspension
		if err := s.students.Suspend(ctx, loan.StudentID, circulation.SuspensionDays); err != nil {
			return nil, err
		}
	}

	if err := s.loans.Save(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}
