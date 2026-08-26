package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OscarAreiza/lms-library/backend/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/backend/internal/domain/access"
	"github.com/OscarAreiza/lms-library/backend/internal/infrastructure/http/middleware"
	"github.com/OscarAreiza/lms-library/backend/internal/infrastructure/http/response"
)

// AuthHandler implements HU-01 (library-docs/07-api/contracts/openapi/library-api.yaml, /auth/login).
// A primary (driving) adapter — it depends on the application layer, never the
// other way around (library-docs/05-architecture/hexagonal-architecture.md).
type AuthHandler struct {
	login *usecase.Login
}

func NewAuthHandler(login *usecase.Login) *AuthHandler {
	return &AuthHandler{login: login}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.FromContext(r.Context())

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", correlationID)
		return
	}

	result, err := h.login.Execute(r.Context(), req.Username, req.Password)
	if errors.Is(err, access.ErrInvalidCredentials) {
		// INV-002: same generic error whether the username or the password was wrong.
		response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Incorrect username or password", correlationID)
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", correlationID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"accessToken": result.AccessToken,
		"expiresIn":   result.ExpiresIn,
	})
}
