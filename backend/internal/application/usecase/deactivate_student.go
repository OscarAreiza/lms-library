package usecase

import (
	"context"

	"github.com/OscarAreiza/lms-library/backend/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/membership"
)

// DeactivateStudent implements HU-03, Scenario 2 — this is the one place in the
// codebase where Membership and Circulation are coordinated together (like
// LoanRegistrationService coordinates Circulation/Catalog/Membership for loans),
// since a Student cannot know about its own active loans directly
// (library-docs/02-domain/entities-and-rules.md, Student.Deactivate).
type DeactivateStudent struct {
	Students membership.StudentRepository
	Loans    circulation.LoanRepository
}

func NewDeactivateStudent(students membership.StudentRepository, loans circulation.LoanRepository) *DeactivateStudent {
	return &DeactivateStudent{Students: students, Loans: loans}
}

func (uc *DeactivateStudent) Execute(ctx context.Context, id string) (*membership.Student, error) {
	student, err := uc.Students.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	activeLoans, err := uc.Loans.CountActiveByStudent(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := student.Deactivate(activeLoans > 0); err != nil {
		return nil, err
	}

	if err := uc.Students.Save(ctx, student); err != nil {
		return nil, err
	}
	return student, nil
}
