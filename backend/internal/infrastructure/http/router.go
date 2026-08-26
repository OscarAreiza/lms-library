package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OscarAreiza/lms-library/backend/internal/infrastructure/http/handler"
	"github.com/OscarAreiza/lms-library/backend/internal/infrastructure/http/middleware"
)

// RouterConfig carries what the router needs to wire itself. Per-module routes
// (students, books, loans) are added here incrementally as each HU is implemented —
// see library-docs/07-api/contracts/openapi/library-api.yaml for the full contract.
type RouterConfig struct {
	DB         *pgxpool.Pool
	JWTSecret  string
	CORSOrigin string
	Loans      *handler.LoanHandler
}

// NewRouter builds the chi router with the base middleware stack and health
// endpoints. Protected routes are mounted under /api/v1 behind middleware.RequireAuth
// as each module's handler is added.
func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CorrelationID)
	r.Use(middleware.CORS(cfg.CORSOrigin))

	health := handler.NewHealthHandler(cfg.DB)
	r.Get("/health", health.Liveness)
	r.Get("/health/ready", health.Readiness)

	r.Route("/api/v1", func(api chi.Router) {
		// Public: POST /auth/login (added when HU-01 is implemented).

		api.Group(func(protected chi.Router) {
			protected.Use(middleware.RequireAuth(cfg.JWTSecret))

			protected.Post("/loans", cfg.Loans.Create) // HU-06
		})
	})

	return r
}
