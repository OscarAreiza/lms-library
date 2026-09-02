package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/OscarAreiza/lms-library/circulation-service/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/service"
)

// --- Fakes (11-quality/tdd-guide.md) standing in for the HTTP calls to
// membership-service and catalog-service.

type fakeStudentClient struct {
	eligible  map[string]bool
	suspended map[string]int
}

func (f *fakeStudentClient) IsEligible(_ context.Context, studentID string) (bool, error) {
	eligible, ok := f.eligible[studentID]
	if !ok {
		return true, nil
	}
	return eligible, nil
}
func (f *fakeStudentClient) Suspend(_ context.Context, studentID string, days int) error {
	if f.suspended == nil {
		f.suspended = map[string]int{}
	}
	f.suspended[studentID] = days
	f.eligible[studentID] = false
	return nil
}

type fakeBookClient struct {
	availableCopies map[string]int
}

func (f *fakeBookClient) LoanCopy(_ context.Context, bookID string) error {
	if f.availableCopies[bookID] <= 0 {
		return errNoCopiesAvailable{}
	}
	f.availableCopies[bookID]--
	return nil
}
func (f *fakeBookClient) ReturnCopy(_ context.Context, bookID string) error {
	f.availableCopies[bookID]++
	return nil
}

type errNoCopiesAvailable struct{}

func (errNoCopiesAvailable) Error() string { return "no copies available" }

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
func (f *fakeLoanRepo) Search(_ context.Context, status string, overdueOnly bool, _, _ string, _, _ int) ([]*circulation.Loan, int, error) {
	var matches []*circulation.Loan
	for _, l := range f.byID {
		if status != "" && l.Status != status {
			continue
		}
		if overdueOnly && !(l.Status == circulation.StatusActive && time.Now().UTC().After(l.DueDate)) {
			continue
		}
		matches = append(matches, l)
	}
	return matches, len(matches), nil
}
func (f *fakeLoanRepo) Save(_ context.Context, l *circulation.Loan) error {
	f.byID[l.ID] = l
	return nil
}

func setupLoanFixture(t *testing.T) (*service.LoanRegistrationService, *circulation.Loan, *fakeStudentClient, *fakeBookClient, *fakeLoanRepo) {
	t.Helper()
	students := &fakeStudentClient{eligible: map[string]bool{}}
	books := &fakeBookClient{availableCopies: map[string]int{"book-1": 1}}
	loans := &fakeLoanRepo{byID: map[string]*circulation.Loan{}, active: map[string]int{}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	loan, err := svc.RegisterLoan(context.Background(), "student-1", "book-1")
	if err != nil {
		t.Fatalf("unexpected error setting up loan fixture: %v", err)
	}
	loans.byID[loan.ID] = loan
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
	if books.availableCopies[loan.BookID] != 1 {
		t.Fatalf("expected availableCopies restored to 1, got %d", books.availableCopies[loan.BookID])
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
	if days := students.suspended[loan.StudentID]; days != circulation.SuspensionDays {
		t.Fatalf("expected the student suspended for %d days, got %d", circulation.SuspensionDays, days)
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
