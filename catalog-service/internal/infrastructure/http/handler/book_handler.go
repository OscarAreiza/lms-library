package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/OscarAreiza/lms-library/catalog-service/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/catalog-service/internal/domain/catalog"
	"github.com/OscarAreiza/lms-library/catalog-service/internal/infrastructure/http/middleware"
	"github.com/OscarAreiza/lms-library/catalog-service/internal/infrastructure/http/response"
)

// BookHandler implements the /books endpoints (HU-04, HU-05).
type BookHandler struct {
	createBook     *usecase.CreateBook
	searchBooks    *usecase.SearchBooks
	loanBookCopy   *usecase.LoanBookCopy
	returnBookCopy *usecase.ReturnBookCopy
}

func NewBookHandler(createBook *usecase.CreateBook, searchBooks *usecase.SearchBooks, loanBookCopy *usecase.LoanBookCopy, returnBookCopy *usecase.ReturnBookCopy) *BookHandler {
	return &BookHandler{createBook: createBook, searchBooks: searchBooks, loanBookCopy: loanBookCopy, returnBookCopy: returnBookCopy}
}

type createBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Category    string `json:"category"`
	Year        int    `json:"year"`
	TotalCopies int    `json:"totalCopies"`
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

type bookResponse struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	ISBN            string `json:"isbn"`
	Category        string `json:"category"`
	Year            int    `json:"year"`
	TotalCopies     int    `json:"totalCopies"`
	AvailableCopies int    `json:"availableCopies"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

func toBookResponse(b *catalog.Book) bookResponse {
	return bookResponse{
		ID:              b.ID,
		Title:           b.Title,
		Author:          b.Author,
		ISBN:            b.ISBN,
		Category:        b.Category,
		Year:            b.Year,
		TotalCopies:     b.TotalCopies,
		AvailableCopies: b.AvailableCopies,
		CreatedAt:       b.CreatedAt.Format(timeFormat),
		UpdatedAt:       b.UpdatedAt.Format(timeFormat),
	}
}

// Create — POST /books (HU-04, FR-007, FR-008).
func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())

	var req createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", correlationID)
		return
	}

	book, err := h.createBook.Execute(r.Context(), req.Title, req.Author, req.ISBN, req.Category, req.Year, req.TotalCopies)
	switch {
	case errors.Is(err, usecase.ErrISBNAlreadyExists):
		response.Error(w, http.StatusConflict, "ISBN_ALREADY_EXISTS", "This ISBN is already registered", correlationID)
		return
	case err != nil:
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), correlationID)
		return
	}

	response.JSON(w, http.StatusCreated, toBookResponse(book))
}

// List — GET /books (HU-05, FR-009, FR-010).
func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	search := r.URL.Query().Get("search")
	category := r.URL.Query().Get("category")

	books, total, err := h.searchBooks.Execute(r.Context(), search, category, page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", correlationID)
		return
	}

	items := make([]bookResponse, 0, len(books))
	for _, b := range books {
		items = append(items, toBookResponse(b))
	}

	// Not an error even when items is empty — HU-05, Scenario 2.
	response.JSON(w, http.StatusOK, map[string]any{
		"data": items,
		"meta": map[string]any{"total": total},
	})
}

// LoanCopy — POST /books/{id}/loan-copy. Called by circulation-service when
// registering a loan (HU-06); not part of the public UI-facing contract.
func (h *BookHandler) LoanCopy(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())
	id := chi.URLParam(r, "id")

	book, err := h.loanBookCopy.Execute(r.Context(), id)
	switch {
	case errors.Is(err, catalog.ErrBookNotFound):
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Book not found", correlationID)
		return
	case errors.Is(err, catalog.ErrNoCopiesAvailable):
		response.Error(w, http.StatusConflict, "NO_COPIES_AVAILABLE", "There are no available copies of this book", correlationID)
		return
	case err != nil:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", correlationID)
		return
	}

	response.JSON(w, http.StatusOK, toBookResponse(book))
}

// ReturnCopy — POST /books/{id}/return-copy. Called by circulation-service
// when registering a return (HU-07).
func (h *BookHandler) ReturnCopy(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())
	id := chi.URLParam(r, "id")

	book, err := h.returnBookCopy.Execute(r.Context(), id)
	switch {
	case errors.Is(err, catalog.ErrBookNotFound):
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Book not found", correlationID)
		return
	case err != nil:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", correlationID)
		return
	}

	response.JSON(w, http.StatusOK, toBookResponse(book))
}
