package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OscarAreiza/lms-library/backend/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/catalog"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/membership"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/service"
	"github.com/OscarAreiza/lms-library/backend/internal/infrastructure/http/middleware"
	"github.com/OscarAreiza/lms-library/backend/internal/infrastructure/http/response"
)

const loanTimeFormat = "2006-01-02T15:04:05Z07:00"

// LoanHandler implements the /loans endpoints (HU-06).
type LoanHandler struct {
	registerLoan *usecase.RegisterLoan
}

func NewLoanHandler(registerLoan *usecase.RegisterLoan) *LoanHandler {
	return &LoanHandler{registerLoan: registerLoan}
}

type createLoanRequest struct {
	StudentID string `json:"studentId"`
	BookID    string `json:"bookId"`
}

type loanResponse struct {
	ID         string  `json:"id"`
	StudentID  string  `json:"studentId"`
	BookID     string  `json:"bookId"`
	LoanDate   string  `json:"loanDate"`
	DueDate    string  `json:"dueDate"`
	ReturnDate *string `json:"returnDate,omitempty"`
	Status     string  `json:"status"`
	WasLate    *bool   `json:"wasLate,omitempty"`
}

func toLoanResponse(l *circulation.Loan) loanResponse {
	resp := loanResponse{
		ID:        l.ID,
		StudentID: l.StudentID,
		BookID:    l.BookID,
		LoanDate:  l.LoanDate.Format(loanTimeFormat),
		DueDate:   l.DueDate.Format(loanTimeFormat),
		Status:    l.Status,
		WasLate:   l.WasLate,
	}
	if l.ReturnDate != nil {
		v := l.ReturnDate.Format(loanTimeFormat)
		resp.ReturnDate = &v
	}
	return resp
}

// Create — POST /loans (HU-06, FR-013, FR-014, FR-015).
func (h *LoanHandler) Create(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())

	var req createLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", correlationID)
		return
	}

	loan, err := h.registerLoan.Execute(r.Context(), req.StudentID, req.BookID)
	switch {
	case errors.Is(err, service.ErrStudentSuspended):
		response.Error(w, http.StatusConflict, "STUDENT_SUSPENDED",
			"This student cannot receive a new loan until their suspension ends", correlationID)
		return
	case errors.Is(err, service.ErrLoanLimitReached):
		response.Error(w, http.StatusConflict, "LOAN_LIMIT_REACHED",
			"This student already has the maximum of 2 active loans", correlationID)
		return
	case errors.Is(err, catalog.ErrNoCopiesAvailable):
		response.Error(w, http.StatusConflict, "NO_COPIES_AVAILABLE",
			"There are no available copies of this book", correlationID)
		return
	case errors.Is(err, membership.ErrStudentNotFound):
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Student not found", correlationID)
		return
	case errors.Is(err, catalog.ErrBookNotFound):
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Book not found", correlationID)
		return
	case err != nil:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", correlationID)
		return
	}

	response.JSON(w, http.StatusCreated, toLoanResponse(loan))
}
