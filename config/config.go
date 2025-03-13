package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Database struct {
		Host     string        `json:"host" validate:"required"`
		Port     string        `json:"port" validate:"required"`
		User     string        `json:"user" validate:"required"`
		Password string        `json:"password" validate:"required"`
		DbName   string        `json:"dbname" validate:"required"`
		SSLMode  string        `json:"sslmode" default:"disable"`
		MaxConns int           `json:"max_conns" default:"10"`
		MaxIdle  int           `json:"max_idle" default:"5"`
		Timeout  time.Duration `json:"timeout" default:"5s"`
	} `json:"database"`
	Jwt struct {
		Secret         string        `json:"secret" validate:"required"`
		ExpirationTime time.Duration `json:"expiration_time" default:"24h"`
		RefreshTime    time.Duration `json:"refresh_time" default:"168h"`
		Issuer         string        `json:"issuer" default:"product-tracker"`
		Audience       string        `json:"audience" default:"product-tracker-api"`
	} `json:"jwt"`
	Server struct {
		Port           string        `json:"port" validate:"required"`
		ReadTimeout    time.Duration `json:"read_timeout" default:"5s"`
		WriteTimeout   time.Duration `json:"write_timeout" default:"10s"`
		IdleTimeout    time.Duration `json:"idle_timeout" default:"120s"`
		MaxHeaderBytes int           `json:"max_header_bytes" default:"1048576"`
	} `json:"server"`
}

var (
	configInstance *Config
	configLoaded   bool
)

// LoadConfig loads and validates the configuration
func LoadConfig() error {
	if configLoaded {
		return nil
	}

	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_conns", 10)
	viper.SetDefault("database.max_idle", 5)
	viper.SetDefault("database.timeout", "5s")
	viper.SetDefault("jwt.expiration_time", "24h")
	viper.SetDefault("jwt.refresh_time", "168h")
	viper.SetDefault("jwt.issuer", "product-tracker")
	viper.SetDefault("jwt.audience", "product-tracker-api")
	viper.SetDefault("server.read_timeout", "5s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.max_header_bytes", 1048576)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Warning: Config file not found, using environment variables and defaults")
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	configInstance = &Config{}
	if err := viper.Unmarshal(configInstance); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	if err := validateConfig(configInstance); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	configLoaded = true
	return nil
}

// validateConfig validates the configuration values
func validateConfig(cfg *Config) error {
	if cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.Database.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if cfg.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if cfg.Database.DbName == "" {
		return fmt.Errorf("database name is required")
	}
	if cfg.Jwt.Secret == "" {
		return fmt.Errorf("jwt secret is required")
	}
	if cfg.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	return nil
}

// GetConfig returns the configuration instance
func GetConfig() *Config {
	if !configLoaded {
		if err := LoadConfig(); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}
	return configInstance
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DbName,
		c.Database.SSLMode,
	)
}
