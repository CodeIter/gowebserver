package main

import (
	"log/slog"
	"os"

	"github.com/CodeIter/gowebserver/internal/config"
	"github.com/CodeIter/gowebserver/internal/server"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Set up structured logging
	slog.SetDefault(logger)

	// Load Configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	// Run the Server
	if err := server.Run(cfg); err != nil {
		slog.Error("server runtime error", "error", err)
		os.Exit(1)
	}
}
