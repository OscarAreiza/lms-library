// Command api is the entry point for circulation-service — see
// library-docs/09-microservices/services/05-circulation-service/README.md.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OscarAreiza/lms-library/circulation-service/internal/application/usecase"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/config"
	catalogclient "github.com/OscarAreiza/lms-library/circulation-service/internal/infrastructure/catalog"
	httpserver "github.com/OscarAreiza/lms-library/circulation-service/internal/infrastructure/http"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/infrastructure/http/handler"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/infrastructure/logger"
	membershipclient "github.com/OscarAreiza/lms-library/circulation-service/internal/infrastructure/membership"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/infrastructure/postgres"

	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/service"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		println("config error:", err.Error())
		return err
	}

	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		return err
	}
	defer func() { _ = zapLog.Sync() }()
	log := zapLog.Sugar()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Errorw("failed to connect to database", "error", err)
		return err
	}
	defer pool.Close()

	loanRepo := postgres.NewLoanRepository(pool)
	studentClient := membershipclient.NewClient(cfg.MembershipServiceURL, cfg.JWTSecret)
	bookClient := catalogclient.NewClient(cfg.CatalogServiceURL, cfg.JWTSecret)

	loanRegistrationService := service.NewLoanRegistrationService(studentClient, bookClient, loanRepo)
	registerLoanUseCase := usecase.NewRegisterLoan(loanRegistrationService)
	loanHandler := handler.NewLoanHandler(registerLoanUseCase)

	router := httpserver.NewRouter(httpserver.RouterConfig{
		DB:         pool,
		JWTSecret:  cfg.JWTSecret,
		CORSOrigin: cfg.CORSOrigin,
		Loans:      loanHandler,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Infow("circulation-service listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorw("server error", "error", err)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
