package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OscarAreiza/lms-library/access-service/internal/infrastructure/http/handler"
	"github.com/OscarAreiza/lms-library/access-service/internal/infrastructure/http/middleware"
)

// RouterConfig carries what the router needs to wire itself.
type RouterConfig struct {
	DB         *pgxpool.Pool
	JWTSecret  string
	CORSOrigin string
	Auth       *handler.AuthHandler
}

// NewRouter builds the chi router with the base middleware stack, health
// endpoints, and the Access bounded context's public /auth/login route.
func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CorrelationID)
	r.Use(middleware.CORS(cfg.CORSOrigin))

	health := handler.NewHealthHandler(cfg.DB)
	r.Get("/health", health.Liveness)
	r.Get("/health/ready", health.Readiness)

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/login", cfg.Auth.Login) // HU-01 — public
	})

	return r
}
