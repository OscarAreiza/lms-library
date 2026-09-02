package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OscarAreiza/lms-library/membership-service/internal/infrastructure/http/handler"
	"github.com/OscarAreiza/lms-library/membership-service/internal/infrastructure/http/middleware"
)

// RouterConfig carries what the router needs to wire itself.
type RouterConfig struct {
	DB         *pgxpool.Pool
	JWTSecret  string
	CORSOrigin string
	Students   *handler.StudentHandler
}

// NewRouter builds the chi router with the base middleware stack, health
// endpoints, and the Membership bounded context's /students routes. JWT
// validation happens here using the secret shared with access-service —
// no network call to access-service is needed to validate a token.
func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CorrelationID)
	r.Use(middleware.CORS(cfg.CORSOrigin))

	health := handler.NewHealthHandler(cfg.DB)
	r.Get("/health", health.Liveness)
	r.Get("/health/ready", health.Readiness)

	r.Route("/api/v1", func(api chi.Router) {
		api.Group(func(protected chi.Router) {
			protected.Use(middleware.RequireAuth(cfg.JWTSecret))

			protected.Route("/students", func(students chi.Router) {
				students.Post("/", cfg.Students.Create)                                 // HU-02
				students.Get("/{id}", cfg.Students.Get)                                 // needed by circulation-service
				students.Get("/by-document/{documentId}", cfg.Students.GetByDocumentID) // loan registration by natural key
			})
		})
	})

	return r
}
