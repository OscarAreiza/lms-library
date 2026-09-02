package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/OscarAreiza/lms-library/catalog-service/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/catalog-service/internal/domain/catalog"
	"github.com/OscarAreiza/lms-library/catalog-service/internal/infrastructure/http/middleware"
	"github.com/OscarAreiza/lms-library/catalog-service/internal/infrastructure/http/response"
)

const timeFormat = "2006-01-02T15:04:05Z07:00"

// BookHandler implements the /books endpoints relevant to HU-09.
type BookHandler struct {
	updateBook *usecase.UpdateBook
}

func NewBookHandler(updateBook *usecase.UpdateBook) *BookHandler {
	return &BookHandler{updateBook: updateBook}
}

type updateBookRequest struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Category string `json:"category"`
	Year     int    `json:"year"`
}

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

// Update — PATCH /books/{id} (HU-09, FR-011, FR-012). ISBN is not in
// updateBookRequest on purpose — see usecase.UpdateBook's doc comment.
func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req updateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", correlationID)
		return
	}

	book, err := h.updateBook.Execute(r.Context(), id, req.Title, req.Author, req.Category, req.Year)
	switch {
	case errors.Is(err, catalog.ErrBookNotFound):
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Book not found", correlationID)
		return
	case err != nil:
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), correlationID)
		return
	}

	response.JSON(w, http.StatusOK, toBookResponse(book))
}
