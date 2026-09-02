package usecase

import (
	"context"

	"github.com/OscarAreiza/lms-library/membership-service/internal/domain/membership"
)

// GetStudentByDocumentID exposes lookup-by-natural-key over HTTP — the
// document ID is the identifier a human (the Administrator) actually has on
// hand, unlike the internal UUID. Used by circulation-service to resolve a
// loan request entered by document ID, and by the frontend loan form to
// avoid asking the Administrator for a raw UUID.
type GetStudentByDocumentID struct {
	Students membership.StudentRepository
}

func NewGetStudentByDocumentID(students membership.StudentRepository) *GetStudentByDocumentID {
	return &GetStudentByDocumentID{Students: students}
}

func (uc *GetStudentByDocumentID) Execute(ctx context.Context, documentID string) (*membership.Student, error) {
	return uc.Students.FindByDocumentID(ctx, documentID)
}
