// Package event defines the domain events raised in-process within library-api
// (no message broker in v1 — see library-docs/02-domain/domain-events.md and
// library-docs/05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md).
package event

import "time"

// BookRegistered is raised when a new book is added to the catalog (HU-04).
type BookRegistered struct {
	BookID      string
	ISBN        string
	Title       string
	TotalCopies int
	OccurredAt  time.Time
}

// StudentRegistered is raised when a new student is registered (HU-02).
type StudentRegistered struct {
	StudentID  string
	DocumentID string
	FullName   string
	OccurredAt time.Time
}

// LoanRegistered is raised when a loan is created (HU-06). Consumed by the catalog
// module to decrement availability.
type LoanRegistered struct {
	LoanID     string
	StudentID  string
	BookID     string
	LoanDate   time.Time
	DueDate    time.Time
	OccurredAt time.Time
}

// LoanReturned is raised when a return is registered (HU-07). Consumed by the
// catalog module (increment availability) and the membership module (apply
// suspension if late).
type LoanReturned struct {
	LoanID     string
	StudentID  string
	BookID     string
	DueDate    time.Time
	ReturnDate time.Time
	IsLate     bool
	OccurredAt time.Time
}

// StudentSuspended is raised as a policy reaction to a late LoanReturned (HU-08).
type StudentSuspended struct {
	StudentID      string
	Reason         string
	SourceLoanID   string
	SuspendedUntil time.Time
	OccurredAt     time.Time
}
