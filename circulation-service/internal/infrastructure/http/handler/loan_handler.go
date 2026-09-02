package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/OscarAreiza/lms-library/circulation-service/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/infrastructure/http/middleware"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/infrastructure/http/response"
)

const loanTimeFormat = "2006-01-02T15:04:05Z07:00"

// LoanHandler implements the /loans endpoints (HU-07 — Return registers HU-06's
// counterpart action; List is HU-07's history/query half).
type LoanHandler struct {
	returnLoan  *usecase.ReturnLoan
	searchLoans *usecase.SearchLoans
}

func NewLoanHandler(returnLoan *usecase.ReturnLoan, searchLoans *usecase.SearchLoans) *LoanHandler {
	return &LoanHandler{returnLoan: returnLoan, searchLoans: searchLoans}
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

// Return — POST /loans/{id}/return (HU-07, FR-016, FR-017, FR-018).
func (h *LoanHandler) Return(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())
	id := chi.URLParam(r, "id")

	loan, err := h.returnLoan.Execute(r.Context(), id)
	switch {
	case errors.Is(err, circulation.ErrLoanNotFound):
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Loan not found", correlationID)
		return
	case errors.Is(err, circulation.ErrLoanAlreadyReturned):
		response.Error(w, http.StatusConflict, "LOAN_ALREADY_RETURNED", "This loan was already returned", correlationID)
		return
	case err != nil:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", correlationID)
		return
	}

	response.JSON(w, http.StatusOK, toLoanResponse(loan))
}

// List — GET /loans (HU-07 history; also used by HU-08 with ?overdue=true).
func (h *LoanHandler) List(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")
	overdue := r.URL.Query().Get("overdue") == "true"
	studentID := r.URL.Query().Get("studentId")
	bookID := r.URL.Query().Get("bookId")

	loans, total, err := h.searchLoans.Execute(r.Context(), status, overdue, studentID, bookID, page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", correlationID)
		return
	}

	items := make([]loanResponse, 0, len(loans))
	for _, l := range loans {
		items = append(items, toLoanResponse(l))
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": items,
		"meta": map[string]any{"total": total},
	})
}
