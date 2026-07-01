package main

import (
	"log/slog"
	"os"
	"flag"

	"my-go-server/internal/config"
	"my-go-server/internal/server"
	"my-go-server/version"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// Handle global flags before config parsing
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version.Get().String())
		os.Exit(0)
	}

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
