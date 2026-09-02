package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/OscarAreiza/lms-library/membership-service/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/membership-service/internal/domain/membership"
	"github.com/OscarAreiza/lms-library/membership-service/internal/domain/shared"
	"github.com/OscarAreiza/lms-library/membership-service/internal/infrastructure/http/middleware"
	"github.com/OscarAreiza/lms-library/membership-service/internal/infrastructure/http/response"
)

// StudentHandler implements the /students endpoints (HU-02).
type StudentHandler struct {
	createStudent *usecase.CreateStudent
	getStudent    *usecase.GetStudent
}

func NewStudentHandler(createStudent *usecase.CreateStudent, getStudent *usecase.GetStudent) *StudentHandler {
	return &StudentHandler{createStudent: createStudent, getStudent: getStudent}
}

type createStudentRequest struct {
	FullName   string `json:"fullName"`
	DocumentID string `json:"documentId"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type studentResponse struct {
	ID             string  `json:"id"`
	FullName       string  `json:"fullName"`
	DocumentID     string  `json:"documentId"`
	Email          string  `json:"email"`
	Phone          *string `json:"phone,omitempty"`
	SuspendedUntil *string `json:"suspendedUntil,omitempty"`
	DeactivatedAt  *string `json:"deactivatedAt,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

func toStudentResponse(s *membership.Student) studentResponse {
	resp := studentResponse{
		ID:         s.ID,
		FullName:   s.FullName,
		DocumentID: s.DocumentID,
		Email:      s.Email.String(),
		CreatedAt:  s.CreatedAt.Format(timeFormat),
		UpdatedAt:  s.UpdatedAt.Format(timeFormat),
	}
	if s.Phone != "" {
		resp.Phone = &s.Phone
	}
	if s.SuspendedUntil != nil {
		v := s.SuspendedUntil.Format(timeFormat)
		resp.SuspendedUntil = &v
	}
	if s.DeactivatedAt != nil {
		v := s.DeactivatedAt.Format(timeFormat)
		resp.DeactivatedAt = &v
	}
	return resp
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

// Create — POST /students (HU-02, FR-003, FR-004).
func (h *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())

	var req createStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", correlationID)
		return
	}

	student, err := h.createStudent.Execute(r.Context(), req.FullName, req.DocumentID, req.Email, req.Phone)
	switch {
	case errors.Is(err, usecase.ErrDocumentIDAlreadyExists):
		response.Error(w, http.StatusConflict, "DOCUMENT_ID_ALREADY_EXISTS", "This document ID is already registered", correlationID)
		return
	case errors.Is(err, shared.ErrInvalidEmail):
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid email format", correlationID)
		return
	case err != nil:
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), correlationID)
		return
	}

	response.JSON(w, http.StatusCreated, toStudentResponse(student))
}

// Get — GET /students/{id}. Used by circulation-service to check a student's
// suspension status before registering a loan (HU-06) — Membership can't be
// queried by SQL join anymore now that Loan lives in a different database.
func (h *StudentHandler) Get(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())
	id := chi.URLParam(r, "id")

	student, err := h.getStudent.Execute(r.Context(), id)
	switch {
	case errors.Is(err, membership.ErrStudentNotFound):
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Student not found", correlationID)
		return
	case err != nil:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", correlationID)
		return
	}

	response.JSON(w, http.StatusOK, toStudentResponse(student))
}
