package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/OscarAreiza/lms-library/backend/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/catalog"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/membership"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/service"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/shared"
)

// --- Fakes (11-quality/tdd-guide.md) — one per repository port ReturnLoan needs.

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
	byID   map[string]*circulation.Loan
	active map[string]int
}

func (f *fakeLoanRepo) FindByID(_ context.Context, id string) (*circulation.Loan, error) {
	l, ok := f.byID[id]
	if !ok {
		return nil, circulation.ErrLoanNotFound
	}
	return l, nil
}
func (f *fakeLoanRepo) CountActiveByStudent(_ context.Context, studentID string) (int, error) {
	return f.active[studentID], nil
}
func (f *fakeLoanRepo) Search(_ context.Context, _ string, _ bool, _, _ string, _, _ int) ([]*circulation.Loan, int, error) {
	return nil, 0, nil
}
func (f *fakeLoanRepo) Save(_ context.Context, l *circulation.Loan) error {
	f.byID[l.ID] = l
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

func setupLoanFixture(t *testing.T) (*service.LoanRegistrationService, *circulation.Loan, *fakeStudentRepo, *fakeBookRepo, *fakeLoanRepo) {
	t.Helper()
	student := newEligibleStudent(t)
	book, err := catalog.NewBook("Clean Code", "Robert C. Martin", "978-0132350884", "Software", 2008, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	students := &fakeStudentRepo{byID: map[string]*membership.Student{student.ID: student}}
	books := &fakeBookRepo{byID: map[string]*catalog.Book{book.ID: book}}
	loans := &fakeLoanRepo{byID: map[string]*circulation.Loan{}, active: map[string]int{}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	loan, _, err := svc.RegisterLoan(context.Background(), student.ID, book.ID)
	if err != nil {
		t.Fatalf("unexpected error setting up loan fixture: %v", err)
	}
	return svc, loan, students, books, loans
}

func TestReturnLoan_OnTime(t *testing.T) {
	svc, loan, _, books, _ := setupLoanFixture(t)
	uc := usecase.NewReturnLoan(svc)

	returned, err := uc.Execute(context.Background(), loan.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if returned.Status != circulation.StatusReturned {
		t.Fatalf("expected status RETURNED, got %s", returned.Status)
	}
	if *returned.WasLate {
		t.Fatal("expected an on-time return")
	}
	if books.byID[loan.BookID].AvailableCopies != 1 {
		t.Fatalf("expected availableCopies restored to 1, got %d", books.byID[loan.BookID].AvailableCopies)
	}
}

func TestReturnLoan_Late_SuspendsStudent(t *testing.T) {
	svc, loan, students, _, _ := setupLoanFixture(t)
	loan.DueDate = time.Now().UTC().Add(-1 * time.Hour) // force overdue for the test
	uc := usecase.NewReturnLoan(svc)

	returned, err := uc.Execute(context.Background(), loan.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !*returned.WasLate {
		t.Fatal("expected a late return")
	}
	if students.byID[loan.StudentID].IsEligibleForLoan() {
		t.Fatal("expected the student to be suspended after a late return")
	}
}

func TestReturnLoan_RejectsDoubleReturn(t *testing.T) {
	svc, loan, _, _, _ := setupLoanFixture(t)
	uc := usecase.NewReturnLoan(svc)

	if _, err := uc.Execute(context.Background(), loan.ID); err != nil {
		t.Fatalf("unexpected error on first return: %v", err)
	}
	if _, err := uc.Execute(context.Background(), loan.ID); err != circulation.ErrLoanAlreadyReturned {
		t.Fatalf("expected ErrLoanAlreadyReturned, got %v", err)
	}
}
