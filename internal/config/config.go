package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/CodeIter/gowebserver/version"

	"github.com/joho/godotenv"
)

// Config holds the server configuration parameters.
type Config struct {
	Host            string
	Port            int
	ResourcesDir    string
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
	// Load .env file if it exists (no error if missing)
	loadEnv()

	// 1. Static Defaults
	cfg := &Config{
		Host:            "0.0.0.0",
		Port:            8000,
		ResourcesDir:    "./resources",
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
	flag.StringVar(&cfg.ResourcesDir, "resources", cfg.ResourcesDir, "External Resources directory, reserved for large files")
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

	if err := ensureDirectory(cfg.ResourcesDir); err != nil {
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

// loadEnv loads environment variables from .env file if it exists.
// Missing .env file is not an error.
func loadEnv() {
	_ = godotenv.Load()
}
