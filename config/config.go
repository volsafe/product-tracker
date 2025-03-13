package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var cfg *Config

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server" json:"server"`
	Database DatabaseConfig `yaml:"database" json:"database"`
	JWT      JWTConfig      `yaml:"jwt" json:"jwt"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port           string        `yaml:"port" json:"port"`
	ReadTimeout    time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout" json:"write_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
	MaxHeaderBytes int           `yaml:"max_header_bytes" json:"max_header_bytes"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	DbName   string `yaml:"dbname" json:"dbname"`
	SSLMode  string `yaml:"sslmode" json:"sslmode"`
}

// JWTConfig represents the JWT configuration
type JWTConfig struct {
	Secret         string        `yaml:"secret" json:"secret"`
	ExpirationTime time.Duration `yaml:"expiration_time" json:"expiration_time"`
}

// LoadConfig loads the configuration from config.yml and environment variables
func LoadConfig() (*Config, error) {
	// Default configuration
	cfg = &Config{
		Server: ServerConfig{
			Port:           "8080",
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1MB
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "pgsql",
			Password: "pgsql",
			DbName:   "consumers",
			SSLMode:  "disable",
		},
		JWT: JWTConfig{
			Secret:         "your-secret-key",
			ExpirationTime: 24 * time.Hour,
		},
	}

	// Try to read config.yaml from the config directory
	if data, err := os.ReadFile("config/config.yaml"); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			log.Printf("Warning: Failed to parse config.yaml: %v", err)
		}
	}

	// Override with environment variables
	cfg.Server.Port = getEnvOrDefault("SERVER_PORT", cfg.Server.Port)
	cfg.Database.Host = getEnvOrDefault("DB_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvOrDefault("DB_PORT", cfg.Database.Port)
	cfg.Database.User = getEnvOrDefault("DB_USER", cfg.Database.User)
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", cfg.Database.Password)
	cfg.Database.DbName = getEnvOrDefault("DB_NAME", cfg.Database.DbName)
	cfg.Database.SSLMode = getEnvOrDefault("DB_SSL_MODE", cfg.Database.SSLMode)
	cfg.JWT.Secret = getEnvOrDefault("JWT_SECRET", cfg.JWT.Secret)

	log.Printf("Loaded configuration - Database: %s@%s:%s/%s",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DbName)

	return cfg, nil
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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
