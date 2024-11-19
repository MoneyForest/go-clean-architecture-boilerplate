package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/dependency"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/controller/http/handler"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/controller/http/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
	defer cancel()

	dependency, err := dependency.Inject(ctx)
	if err != nil {
		return err
	}

	r := chi.NewRouter()

	// ミドルウェアの設定
	r.Use(middleware.Logger)
	r.Use(middleware.Recover)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// ハンドラーの初期化
	userHandler := &handler.UserHandler{
		UserInteractor: dependency.UserInteractor,
	}
	healthHandler := &handler.HealthHandler{}

	// ルーティングの設定
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.List)
			r.Post("/", userHandler.Create)
			r.Get("/{id}", userHandler.Get)
			r.Put("/{id}", userHandler.Update)
			r.Delete("/{id}", userHandler.Delete)
		})
		r.Route("/health", func(r chi.Router) {
			r.Get("/", healthHandler.Health)
			r.Get("/deep", healthHandler.Deep)
		})
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", dependency.Environment.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Failed to gracefully shutdown: %v\n", err)
		}
	}()

	log.Printf("Server starting on port %s\n", dependency.Environment.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
