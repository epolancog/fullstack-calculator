package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/epolancog/fullstack-calculator/backend/internal/calculator"
	"github.com/epolancog/fullstack-calculator/backend/internal/handler"
	"github.com/epolancog/fullstack-calculator/backend/internal/middleware"
	"github.com/epolancog/fullstack-calculator/backend/internal/validator"
)

func main() {
	// Create calculator with all operators registered
	calc := calculator.NewCalculator()

	// Create validator with supported operations
	v := validator.NewValidator(calc.SupportedOperations())

	// Create handler with dependencies injected
	h := handler.NewHandler(calc, v)

	// Create mux and register routes
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Wrap mux with middleware chain: Recovery → Logging → CORS → ContentType → mux
	var chain http.Handler = mux
	chain = middleware.ContentType(chain)
	chain = middleware.CORS(chain)
	chain = middleware.Logging(chain)
	chain = middleware.Recovery(chain)

	// Read port from env, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: chain,
	}

	// Graceful shutdown on SIGINT/SIGTERM
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server started", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
	slog.Info("server stopped")
}
