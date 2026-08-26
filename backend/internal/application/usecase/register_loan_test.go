package usecase_test

import (
	"context"
	"testing"

	"github.com/OscarAreiza/lms-library/backend/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/catalog"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/membership"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/service"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/shared"
)

// --- Fakes (11-quality/tdd-guide.md) — one per repository port RegisterLoan needs.

type fakeStudentRepo struct{ byID map[string]*membership.Student }

func (f *fakeStudentRepo) FindByID(_ context.Context, id string) (*membership.Student, error) {
	s, ok := f.byID[id]
	if !ok {
		return nil, membership.ErrStudentNotFound
	}
	return s, nil
}
func (f *fakeStudentRepo) FindByDocumentID(_ context.Context, _ string) (*membership.Student, error) {
	return nil, membership.ErrStudentNotFound
}
func (f *fakeStudentRepo) Search(_ context.Context, _ string, _, _ int) ([]*membership.Student, int, error) {
	return nil, 0, nil
}
func (f *fakeStudentRepo) Save(_ context.Context, s *membership.Student) error {
	f.byID[s.ID] = s
	return nil
}

type fakeBookRepo struct{ byID map[string]*catalog.Book }

func (f *fakeBookRepo) FindByID(_ context.Context, id string) (*catalog.Book, error) {
	b, ok := f.byID[id]
	if !ok {
		return nil, catalog.ErrBookNotFound
	}
	return b, nil
}
func (f *fakeBookRepo) FindByISBN(_ context.Context, _ string) (*catalog.Book, error) {
	return nil, catalog.ErrBookNotFound
}
func (f *fakeBookRepo) Search(_ context.Context, _, _ string, _, _ int) ([]*catalog.Book, int, error) {
	return nil, 0, nil
}
func (f *fakeBookRepo) Save(_ context.Context, b *catalog.Book) error {
	f.byID[b.ID] = b
	return nil
}

type fakeLoanRepo struct {
	active map[string]int
	saved  []*circulation.Loan
}

func (f *fakeLoanRepo) FindByID(_ context.Context, _ string) (*circulation.Loan, error) {
	return nil, circulation.ErrLoanNotFound
}
func (f *fakeLoanRepo) CountActiveByStudent(_ context.Context, studentID string) (int, error) {
	return f.active[studentID], nil
}
func (f *fakeLoanRepo) Search(_ context.Context, _ string, _ bool, _, _ string, _, _ int) ([]*circulation.Loan, int, error) {
	return nil, 0, nil
}
func (f *fakeLoanRepo) Save(_ context.Context, l *circulation.Loan) error {
	f.saved = append(f.saved, l)
	return nil
}

func newEligibleStudent(t *testing.T) *membership.Student {
	t.Helper()
	email, err := shared.NewEmail("jane@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, err := membership.NewStudent("Jane Doe", "1075300000", email, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return s
}

func newAvailableBook(t *testing.T) *catalog.Book {
	t.Helper()
	b, err := catalog.NewBook("Clean Code", "Robert C. Martin", "978-0132350884", "Software", 2008, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return b
}

func TestRegisterLoan_Succeeds(t *testing.T) {
	student := newEligibleStudent(t)
	book := newAvailableBook(t)
	students := &fakeStudentRepo{byID: map[string]*membership.Student{student.ID: student}}
	books := &fakeBookRepo{byID: map[string]*catalog.Book{book.ID: book}}
	loans := &fakeLoanRepo{active: map[string]int{}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	uc := usecase.NewRegisterLoan(svc)

	loan, err := uc.Execute(context.Background(), student.ID, book.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loan.Status != circulation.StatusActive {
		t.Fatalf("expected status ACTIVE, got %s", loan.Status)
	}
	if book.AvailableCopies != 0 {
		t.Fatalf("expected availableCopies to be decremented to 0, got %d", book.AvailableCopies)
	}
}

func TestRegisterLoan_RejectsSuspendedStudent(t *testing.T) {
	student := newEligibleStudent(t)
	student.Suspend(7)
	book := newAvailableBook(t)
	students := &fakeStudentRepo{byID: map[string]*membership.Student{student.ID: student}}
	books := &fakeBookRepo{byID: map[string]*catalog.Book{book.ID: book}}
	loans := &fakeLoanRepo{active: map[string]int{}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	uc := usecase.NewRegisterLoan(svc)

	_, err := uc.Execute(context.Background(), student.ID, book.ID)
	if err != service.ErrStudentSuspended {
		t.Fatalf("expected ErrStudentSuspended, got %v", err)
	}
}

func TestRegisterLoan_RejectsWhenLoanLimitReached(t *testing.T) {
	student := newEligibleStudent(t)
	book := newAvailableBook(t)
	students := &fakeStudentRepo{byID: map[string]*membership.Student{student.ID: student}}
	books := &fakeBookRepo{byID: map[string]*catalog.Book{book.ID: book}}
	loans := &fakeLoanRepo{active: map[string]int{student.ID: circulation.MaxActiveLoans}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	uc := usecase.NewRegisterLoan(svc)

	_, err := uc.Execute(context.Background(), student.ID, book.ID)
	if err != service.ErrLoanLimitReached {
		t.Fatalf("expected ErrLoanLimitReached, got %v", err)
	}
}

func TestRegisterLoan_RejectsWhenNoCopiesAvailable(t *testing.T) {
	student := newEligibleStudent(t)
	book, err := catalog.NewBook("Clean Code", "Robert C. Martin", "978-0132350884", "Software", 2008, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = book.LoanOneCopy() // exhaust the only copy
	students := &fakeStudentRepo{byID: map[string]*membership.Student{student.ID: student}}
	books := &fakeBookRepo{byID: map[string]*catalog.Book{book.ID: book}}
	loans := &fakeLoanRepo{active: map[string]int{}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	uc := usecase.NewRegisterLoan(svc)

	_, err = uc.Execute(context.Background(), student.ID, book.ID)
	if err != catalog.ErrNoCopiesAvailable {
		t.Fatalf("expected ErrNoCopiesAvailable, got %v", err)
	}
}
