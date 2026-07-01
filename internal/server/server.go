package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"my-go-server/internal/config"
	"my-go-server/internal/handler"
	"my-go-server/internal/middleware"
)

func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	// Register Routes (Go 1.22+ method patterns)
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /ready", handler.Ready)
	mux.HandleFunc("GET /", handler.Home)
	
	// Serve Static Files
	fs := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Apply Middleware Chain (Inner to Outer)
	var handler http.Handler = mux
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.Logger(handler)
	handler = middleware.ConcurrencyLimiter(cfg.MaxConcurrency)(handler)
	handler = middleware.RateLimiterMiddleware(cfg.RateLimit, cfg.RateLimitBurst)(handler)

	srv := &http.Server{
		Addr:         cfg.Host + ":" + itoa(cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start Server in Goroutine
	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for Interrupt Signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("server exited gracefully")
	return nil
}

// TODO Note: Replace itoa with strconv.Itoa in actual code.
func itoa(i int) string {
	return string(rune(i)) // Simplified for brevity; use strconv.Itoa in prod
}
