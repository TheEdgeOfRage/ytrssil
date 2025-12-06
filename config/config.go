package config

import (
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type Config struct {
	Port      int    `long:"port" env:"PORT" default:"8080"`
	DBURI     string `long:"db-uri" env:"DB_URI"`
	AuthToken string `long:"auth-token" env:"AUTH_TOKEN"`
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
	return config, err
}

// TestConfig returns a mostly hardcoded configuration used for running tests
func TestConfig() Config {
	dbURI := getenvOrDefault("DB_URI", "postgres://ytrssil:ytrssil@localhost:5432/ytrssil?sslmode=disable")
	if !strings.Contains(dbURI, "sslmode") {
		dbURI = dbURI + "?sslmode=disable"
	}

	config := Config{
		Port:      8080,
		DBURI:     dbURI,
		AuthToken: "foo",
	}

	return config
}
