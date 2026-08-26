package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OscarAreiza/lms-library/backend/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/catalog"
	"github.com/OscarAreiza/lms-library/backend/internal/infrastructure/http/middleware"
	"github.com/OscarAreiza/lms-library/backend/internal/infrastructure/http/response"
)

// BookHandler implements the /books endpoints (HU-04).
type BookHandler struct {
	createBook *usecase.CreateBook
}

func NewBookHandler(createBook *usecase.CreateBook) *BookHandler {
	return &BookHandler{createBook: createBook}
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
