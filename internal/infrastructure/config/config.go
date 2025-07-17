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

type Config struct {
	Debug    bool           `mapstructure:"debug" yaml:"debug"`
	App      AppConfig      `mapstructure:"app" yaml:"app"`
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Auth     AuthConfig     `mapstructure:"auth" yaml:"auth" validate:"required"`
}

type AppConfig struct {
	Name string `mapstructure:"name" yaml:"name"`
	Env  string `mapstructure:"env" yaml:"env"`
}

type ServerConfig struct {
	Port            int           `mapstructure:"port" yaml:"port"`
	Host            string        `mapstructure:"host" yaml:"host"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" yaml:"dbname"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"`
}

type AuthConfig struct {
	JWTSecret string        `mapstructure:"jwt_secret" yaml:"jwt_secret" validate:"required"`
	JWTExpiry time.Duration `mapstructure:"jwt_expiry" yaml:"jwt_expiry" validate:"required"`
}

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
