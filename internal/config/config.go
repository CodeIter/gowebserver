package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"my-go-server/version"
)

// Config holds the server configuration parameters.
type Config struct {
	Host            string
	Port            int
	StaticDir       string
	PublicDir       string
	ViewsDir        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	MaxConcurrency  int
	RateLimit       int // requests per second
	RateLimitBurst  int // max tokens in bucket
}

// Load reads the configuration from environment variables,
// command line flags, and static defaults.
func Load() (*Config, error) {
	// 1. Static Defaults
	cfg := &Config{
		Host:            "0.0.0.0",
		Port:            8000,
		StaticDir:       "./static",
		PublicDir:       "./public",
		ViewsDir:        "./views",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		MaxConcurrency:  100,
		RateLimit:       10, // average rate (requests per second)
		RateLimitBurst:  10, // burst capacity (max tokens in bucket)
	}

	// 2. Environment Variables Override
	if h := os.Getenv("HOST"); h != "" {
		cfg.Host = h
	}
	if p := os.Getenv("PORT"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			cfg.Port = val
		}
	}

	// 3. Command Line Flags Override (Highest Priority)
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "Server host")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "Server port")
	flag.StringVar(&cfg.StaticDir, "static", cfg.StaticDir, "Static files directory mounted at /static route")
	flag.StringVar(&cfg.PublicDir, "public", cfg.PublicDir, "Public files directory mounted at / route")
	flag.StringVar(&cfg.ViewsDir, "views", cfg.ViewsDir, "View templates directory")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Println(version.Get().String())
		os.Exit(0)
	}

	// Validate Port Range
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", cfg.Port)
	}

	if err := ensureDirectory(cfg.StaticDir); err != nil {
		return nil, err
	}
	if err := ensureDirectory(cfg.PublicDir); err != nil {
		return nil, err
	}
	if err := ensureDirectory(cfg.ViewsDir); err != nil {
		return nil, err
	}

	return cfg, nil
}

func ensureDirectory(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path exists and is not a directory: %s", path)
		}
		return nil
	}
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0o755)
	}
	return err
}
