package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
)

type Config struct {
	Dev            bool   `long:"dev" env:"DEV"`
	Port           int    `long:"port" env:"PORT" default:"8080"`
	DBURI          string `long:"db-uri" env:"DB_URI"`
	AuthToken      string `long:"auth-token" env:"AUTH_TOKEN"`
	YouTubeAPIKey  string `long:"youtube-api-key" env:"YOUTUBE_API_KEY"`
	DownloadsDir   string `long:"downloads-dir" env:"DOWNLOADS_DIR" default:"/var/lib/ytrssil/downloads"`
	FetchInterval  string `long:"fetch-interval" env:"FETCH_INTERVAL" default:"5m"`
	CleanupInterval string `long:"cleanup-interval" env:"CLEANUP_INTERVAL" default:"1h"`
	CleanupAge     string `long:"cleanup-age" env:"CLEANUP_AGE" default:"48h"`
}

func getenvOrDefault(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	}

	return defaultValue
}

// Parse parses all the supplied configurations and returns
func Parse() (Config, error) {
	var config Config
	parser := flags.NewParser(&config, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		return config, err
	}

	if config.AuthToken == "" {
		return config, fmt.Errorf("missing AUTH_TOKEN env var")
	}
	if config.YouTubeAPIKey == "" {
		return config, fmt.Errorf("missing YOUTUBE_API_KEY env var")
	}

	if config.DownloadsDir == "" {
		return config, fmt.Errorf("missing DOWNLOADS_DIR env var")
	}
	if err := os.MkdirAll(config.DownloadsDir, 0o755); err != nil {
		return config, fmt.Errorf("failed to create downloads directory: %w", err)
	}

	return config, nil
}

// GetFetchInterval returns the fetch interval as a time.Duration
func (c Config) GetFetchInterval() time.Duration {
	d, _ := time.ParseDuration(c.FetchInterval)
	return d
}

// GetCleanupInterval returns the cleanup interval as a time.Duration
func (c Config) GetCleanupInterval() time.Duration {
	d, _ := time.ParseDuration(c.CleanupInterval)
	return d
}

// GetCleanupAge returns the cleanup age as a time.Duration
func (c Config) GetCleanupAge() time.Duration {
	d, _ := time.ParseDuration(c.CleanupAge)
	return d
}

// TestConfig returns a mostly hardcoded configuration used for running tests
func TestConfig() Config {
	dbURI := getenvOrDefault("DB_URI", "postgres://ytrssil:ytrssil@localhost:5432/ytrssil?sslmode=disable")
	if !strings.Contains(dbURI, "sslmode") {
		dbURI = dbURI + "?sslmode=disable"
	}

	config := Config{
		Port:          8080,
		DBURI:         dbURI,
		AuthToken:     "foo",
		DownloadsDir:  "/tmp/ytrssil-test-downloads",
		FetchInterval: "5m",
		CleanupInterval: "1h",
		CleanupAge:    "48h",
	}

	return config
}
