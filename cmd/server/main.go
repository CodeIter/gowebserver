package main

import (
	"log/slog"
	"os"

	"my-go-server/internal/config"
	"my-go-server/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	if err := server.Run(cfg); err != nil {
		slog.Error("server runtime error", "error", err)
		os.Exit(1)
	}
}
