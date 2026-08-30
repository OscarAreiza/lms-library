// Package service holds Domain Services — business logic that coordinates more than
// one bounded context's aggregate and therefore does not naturally belong inside a
// single Aggregate Root. See library-docs/02-domain/entities-and-rules.md, "Domain
// Services", and library-docs/09-microservices/services/01-library-api/decisions.md.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/OscarAreiza/lms-library/backend/internal/domain/catalog"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/event"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/membership"
)

// ErrStudentSuspended — INV-003 on Loan: a loan can only be created for a
// non-suspended student.
var ErrStudentSuspended = errors.New("student is suspended")

// ErrLoanLimitReached — INV-004 on Loan: max 2 simultaneous active loans per student.
var ErrLoanLimitReached = errors.New("loan limit reached")

// LoanRegistrationService coordinates Student (Membership), Book (Catalog), and Loan
// (Circulation) — this is the one place in the codebase where those three aggregates
// are checked/changed together. Neither Student, Book, nor Loan alone can enforce
// these cross-aggregate rules.
type LoanRegistrationService struct {
	students membership.StudentRepository
	books    catalog.BookRepository
	loans    circulation.LoanRepository
}

func NewLoanRegistrationService(
	students membership.StudentRepository,
	books catalog.BookRepository,
	loans circulation.LoanRepository,
) *LoanRegistrationService {
	return &LoanRegistrationService{students: students, books: books, loans: loans}
}

// RegisterLoan implements HU-06's acceptance criteria end to end.
func (s *LoanRegistrationService) RegisterLoan(ctx context.Context, studentID, bookID string) (*circulation.Loan, []any, error) {
	student, err := s.students.FindByID(ctx, studentID)
	if err != nil {
		return nil, nil, err
	}
	if !student.IsEligibleForLoan() { // INV-003
		return nil, nil, ErrStudentSuspended
	}

	activeCount, err := s.loans.CountActiveByStudent(ctx, studentID)
	if err != nil {
		return nil, nil, err
	}
	if activeCount >= circulation.MaxActiveLoans { // INV-004
		return nil, nil, ErrLoanLimitReached
	}

	book, err := s.books.FindByID(ctx, bookID)
	if err != nil {
		return nil, nil, err
	}
	if err := book.LoanOneCopy(); err != nil { // INV-002 / INV-001 on Book
		return nil, nil, err
	}

	loan := circulation.NewLoan(studentID, bookID) // INV-001 on Loan: dueDate = +7 days

	if err := s.books.Save(ctx, book); err != nil {
		return nil, nil, err
	}
	if err := s.loans.Save(ctx, loan); err != nil {
		return nil, nil, err
	}

	events := []any{
		event.LoanRegistered{
			LoanID: loan.ID, StudentID: loan.StudentID, BookID: loan.BookID,
			LoanDate: loan.LoanDate, DueDate: loan.DueDate, OccurredAt: time.Now().UTC(),
		},
	}
	return loan, events, nil
}

// RegisterReturn implements HU-07 and HU-08's acceptance criteria: closes the loan,
// restores availability, and — if the return was late — applies the flat suspension.
func (s *LoanRegistrationService) RegisterReturn(ctx context.Context, loanID string) (*circulation.Loan, []any, error) {
	loan, err := s.loans.FindByID(ctx, loanID)
	if err != nil {
		return nil, nil, err
	}

	returnTime := time.Now().UTC()
	wasLate, err := loan.RegisterReturn(returnTime) // INV-005
	if err != nil {
		return nil, nil, err
	}

	book, err := s.books.FindByID(ctx, loan.BookID)
	if err != nil {
		return nil, nil, err
	}
	if err := book.ReturnOneCopy(); err != nil {
		return nil, nil, err
	}
	if err := s.books.Save(ctx, book); err != nil {
		return nil, nil, err
	}

	events := []any{
		event.LoanReturned{
			LoanID: loan.ID, StudentID: loan.StudentID, BookID: loan.BookID,
			DueDate: loan.DueDate, ReturnDate: returnTime, IsLate: wasLate,
			OccurredAt: returnTime,
		},
	}

	if wasLate { // INV-006 on Loan: late return triggers a fixed 7-day suspension
		student, err := s.students.FindByID(ctx, loan.StudentID)
		if err != nil {
			return nil, nil, err
		}
		student.Suspend(circulation.SuspensionDays)
		if err := s.students.Save(ctx, student); err != nil {
			return nil, nil, err
		}
		events = append(events, event.StudentSuspended{
			StudentID: student.ID, Reason: "LATE_RETURN", SourceLoanID: loan.ID,
			SuspendedUntil: *student.SuspendedUntil, OccurredAt: returnTime,
		})
	}

	if err := s.loans.Save(ctx, loan); err != nil {
		return nil, nil, err
	}

	return loan, events, nil
}
