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
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"

	assets "github.com/CodeIter/gowebserver"
	"github.com/CodeIter/gowebserver/internal/config"
	"github.com/CodeIter/gowebserver/internal/handler"
	"github.com/CodeIter/gowebserver/internal/middleware"
	"github.com/CodeIter/gowebserver/internal/template"
)

// Run starts the HTTP server with the provided configuration
// and handles graceful shutdown.
func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	// Load HTML templates for server-side views.
	if err := template.LoadTemplates(); err != nil {
		return err
	}

	// Register Routes (Go 1.22+ method patterns)
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /ready", handler.Ready)
	mux.HandleFunc("GET /about", handler.About)
	mux.HandleFunc("GET /go/{domain}/{path...}", handler.Redirector)
	mux.HandleFunc("GET /random/number/{min}/{max}", handler.RandomNumber)
	mux.HandleFunc("GET /random/password/{length}", handler.RandomPassword)

	// Serve Resources Files
	ResourcesFs := http.FileServer(http.Dir(cfg.ResourcesDir))
	mux.Handle("GET /resources/", http.StripPrefix("/resources/", ResourcesFs))

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

	// Default to 404 if ServeDir is not set
	serveDirHandler := http.NotFoundHandler()
	if cfg.ServeDir != "false" {
		serveDirHandler = http.StripPrefix("/", http.FileServer(http.Dir(cfg.ServeDir)))
	}

	// Create a combined handler: try public files first, fall back to home page
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		// Try serving from ServeDir if configured
		if cfg.ServeDir != "false" && r.URL.Path != "/" {
			// Check if file exists in ServeDir before serving
			filePath := filepath.Join(cfg.ServeDir, strings.TrimPrefix(r.URL.Path, "/"))
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				// XXX http.FileServer automatically return 404 if file not found, so we need to check if file exists before calling ServeHTTP
				serveDirHandler.ServeHTTP(w, r)
				return
			}
		}

		// Try serving from embedded public FS
		embeddedPath := path.Join("public", strings.TrimPrefix(r.URL.Path, "/"))
		if embeddedPath == "public/" {
			embeddedPath = "public/index.html"
		}
		if info, err := iofs.Stat(assets.EmbeddedFiles, embeddedPath); err == nil && !info.IsDir() && r.URL.Path != "/" && !strings.HasPrefix(path.Base(r.URL.Path), ".") {
			// XXX http.FileServer automatically redirects index.html to / .
			// So public directory could not have index.html file.
			publicFSHandler.ServeHTTP(w, r)
			return
		}

		// Redirect /index.html and similar to /
		indexFiles := []string{"/index.html", "/default.html"}
		if slices.Contains(indexFiles, r.URL.Path) {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		if r.URL.Path != "/" {
			// If the path is not found in public, return 404
			http.NotFound(w, r)
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
