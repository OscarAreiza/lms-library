package usecase_test

import (
	"context"
	"testing"

	"github.com/OscarAreiza/lms-library/circulation-service/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/service"
)

// --- Fakes (11-quality/tdd-guide.md) standing in for the HTTP calls to
// membership-service and catalog-service.

type fakeStudentClient struct {
	eligible map[string]bool
}

// ResolveByDocumentID is a passthrough in these tests — the fixture data
// uses the same string for both the "document ID" and the internal ID, so
// the resolution step doesn't need its own separate mapping.
func (f *fakeStudentClient) ResolveByDocumentID(_ context.Context, documentID string) (string, error) {
	return documentID, nil
}
func (f *fakeStudentClient) IsEligible(_ context.Context, studentID string) (bool, error) {
	eligible, ok := f.eligible[studentID]
	if !ok {
		return true, nil
	}
	return eligible, nil
}
func (f *fakeStudentClient) Suspend(_ context.Context, _ string, _ int) error { return nil }

type fakeBookClient struct {
	availableCopies map[string]int
}

// ResolveByISBN is a passthrough in these tests, for the same reason as
// fakeStudentClient.ResolveByDocumentID above.
func (f *fakeBookClient) ResolveByISBN(_ context.Context, isbn string) (string, error) {
	return isbn, nil
}
func (f *fakeBookClient) LoanCopy(_ context.Context, bookID string) error {
	if f.availableCopies[bookID] <= 0 {
		return catalogErrNoCopiesAvailable
	}
	f.availableCopies[bookID]--
	return nil
}
func (f *fakeBookClient) ReturnCopy(_ context.Context, bookID string) error {
	f.availableCopies[bookID]++
	return nil
}

// Mirrors circulation-service/internal/infrastructure/catalog.ErrNoCopiesAvailable
// without importing that package, keeping this fake self-contained.
var catalogErrNoCopiesAvailable = errNoCopiesAvailable{}

type errNoCopiesAvailable struct{}

func (errNoCopiesAvailable) Error() string { return "no copies available" }

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

func TestRegisterLoan_Succeeds(t *testing.T) {
	students := &fakeStudentClient{eligible: map[string]bool{}}
	books := &fakeBookClient{availableCopies: map[string]int{"book-1": 1}}
	loans := &fakeLoanRepo{active: map[string]int{}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	uc := usecase.NewRegisterLoan(students, books, svc)

	loan, err := uc.Execute(context.Background(), "student-1", "book-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loan.Status != circulation.StatusActive {
		t.Fatalf("expected status ACTIVE, got %s", loan.Status)
	}
	if books.availableCopies["book-1"] != 0 {
		t.Fatalf("expected availableCopies to be decremented to 0, got %d", books.availableCopies["book-1"])
	}
}

func TestRegisterLoan_RejectsSuspendedStudent(t *testing.T) {
	students := &fakeStudentClient{eligible: map[string]bool{"student-1": false}}
	books := &fakeBookClient{availableCopies: map[string]int{"book-1": 1}}
	loans := &fakeLoanRepo{active: map[string]int{}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	uc := usecase.NewRegisterLoan(students, books, svc)

	_, err := uc.Execute(context.Background(), "student-1", "book-1")
	if err != service.ErrStudentSuspended {
		t.Fatalf("expected ErrStudentSuspended, got %v", err)
	}
}

func TestRegisterLoan_RejectsWhenLoanLimitReached(t *testing.T) {
	students := &fakeStudentClient{eligible: map[string]bool{}}
	books := &fakeBookClient{availableCopies: map[string]int{"book-1": 1}}
	loans := &fakeLoanRepo{active: map[string]int{"student-1": circulation.MaxActiveLoans}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	uc := usecase.NewRegisterLoan(students, books, svc)

	_, err := uc.Execute(context.Background(), "student-1", "book-1")
	if err != service.ErrLoanLimitReached {
		t.Fatalf("expected ErrLoanLimitReached, got %v", err)
	}
}

func TestRegisterLoan_RejectsWhenNoCopiesAvailable(t *testing.T) {
	students := &fakeStudentClient{eligible: map[string]bool{}}
	books := &fakeBookClient{availableCopies: map[string]int{"book-1": 0}}
	loans := &fakeLoanRepo{active: map[string]int{}}

	svc := service.NewLoanRegistrationService(students, books, loans)
	uc := usecase.NewRegisterLoan(students, books, svc)

	_, err := uc.Execute(context.Background(), "student-1", "book-1")
	if err != catalogErrNoCopiesAvailable {
		t.Fatalf("expected no-copies error, got %v", err)
	}
}
