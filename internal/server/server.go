package server

import (
	"context"
	"fmt"
	iofs "io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"slices"
	"strconv"
	"strings"
	"syscall"

	assets "my-go-server"
	"my-go-server/internal/config"
	"my-go-server/internal/handler"
	"my-go-server/internal/middleware"
)

// Run starts the HTTP server with the provided configuration
// and handles graceful shutdown.
func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	// Load HTML templates for server-side views.
	if err := handler.LoadTemplates(); err != nil {
		return err
	}

	// Register Routes (Go 1.22+ method patterns)
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /ready", handler.Ready)

	// Serve External Files
	externalFs := http.FileServer(http.Dir(cfg.ExternalDir))
	mux.Handle("GET /external/", http.StripPrefix("/external/", externalFs))

	// Serve Static Files from embedded filesystem
	subStatic, err := iofs.Sub(assets.EmbeddedFiles, "static")
	if err != nil {
		return err
	}
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.FS(subStatic)))
	mux.Handle("GET /static/", staticHandler)

	// Serve Public Files from embedded filesystem
	subPublic, err := iofs.Sub(assets.EmbeddedFiles, "public")
	if err != nil {
		return err
	}
	publicFSHandler := http.FileServer(http.FS(subPublic))

	// Create a combined handler: try public files first, fall back to home page
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Try serving from embedded public FS
		embeddedPath := path.Join("public", strings.TrimPrefix(r.URL.Path, "/"))
		if embeddedPath == "public/" {
			embeddedPath = "public/index.html"
		}
		if info, err := iofs.Stat(assets.EmbeddedFiles, embeddedPath); err == nil && !info.IsDir() && r.URL.Path != "/" && !strings.HasPrefix(path.Base(r.URL.Path), ".") {
			publicFSHandler.ServeHTTP(w, r)
			return
		}

		// Redirect /index.html and similar to /
		indexFiles := []string{"/index.html", "/default.html"}
		if slices.Contains(indexFiles, r.URL.Path) {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		// Path is "/" or not found in public, serve home page
		handler.Home(w, r)
	})

	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, `User-agent: *
Disallow: /admin/
Disallow: /private/
Allow: /admin/public/

# Sitemap location
Sitemap: https://%s/sitemap.xml
Host: %s
`, r.Host, r.Host)
	})

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

func itoa(i int) string {
	return strconv.Itoa(i)
}
