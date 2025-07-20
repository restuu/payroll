package config

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// AppEnvKey is the environment variable key to determine the deployment environment.
const AppEnvKey = "APP_ENV"

// Environment-specific config files are embedded into the binary.
//
//go:embed *.config.yaml
var configFS embed.FS

func LoadConfig() (*Config, error) {
	var config Config
	v := viper.New()

	// 1. Determine environment and the corresponding config file.
	env := os.Getenv(AppEnvKey)
	if env == "" {
		env = "development" // Default to development to align with docker-compose
	}
	configFileName := fmt.Sprintf("%s.config.yaml", env)
	log.Printf("Loading config file %s", configFileName)

	// 2. Read the selected embedded configuration file.
	v.SetConfigType("yaml")

	fileBytes, err := configFS.ReadFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded config file '%s': %w. Make sure APP_ENV is set to a valid environment (e.g., development, staging, production)", configFileName, err)
	}

	if err := v.ReadConfig(bytes.NewReader(fileBytes)); err != nil {
		return nil, fmt.Errorf("failed to parse embedded config file '%s': %w", configFileName, err)
	}

	// Explicitly set the App Env from the determined environment.
	// This is more robust than relying on AutomaticEnv for this specific key,
	// especially when the 'app' key might not exist in the YAML file.
	v.Set("app.env", env)

	// Set default values
	v.SetDefault("app.name", "payroll")
	v.SetDefault("server.shutdown_timeout", 5*time.Second)

	// 3. Set up and apply environment variable overrides (highest priority).
	// e.g. DATABASE_USER will override config.database.user
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 4. Unmarshal the final configuration into the struct.
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	val := validator.New(validator.WithRequiredStructEnabled())

	if err := val.Struct(config); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &config, nil
}
