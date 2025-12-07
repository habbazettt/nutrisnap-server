package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	MinIO    MinIOConfig
	JWT      JWTConfig
	Google   GoogleOAuthConfig
}

type ServerConfig struct {
	Port        string
	Environment string
	LogLevel    string
	BaseURL     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type MinIOConfig struct {
	Endpoint    string
	AccessKey   string
	SecretKey   string
	Bucket      string
	UseSSL      bool
	APIPort     string
	ConsolePort string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var cfg *Config

func Load() (*Config, error) {
	cfg = &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "3000"),
			Environment: getEnv("ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			BaseURL:     getEnv("BASE_URL", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		MinIO: MinIOConfig{
			Endpoint:    getEnv("MINIO_ENDPOINT", ""),
			AccessKey:   getEnv("MINIO_ACCESS_KEY", ""),
			SecretKey:   getEnv("MINIO_SECRET_KEY", ""),
			Bucket:      getEnv("MINIO_BUCKET", "nutrisnap"),
			UseSSL:      getEnvBool("MINIO_USE_SSL", false),
			APIPort:     getEnv("MINIO_API_PORT", "9010"),
			ConsolePort: getEnv("MINIO_CONSOLE_PORT", "9011"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "nutrisnap-secret-key-change-in-production"),
			AccessExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", 30*time.Minute),
			RefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			Issuer:        getEnv("JWT_ISSUER", "nutrisnap-api"),
		},
		Google: GoogleOAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/api/v1/auth/google/callback"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func Get() *Config {
	return cfg
}

func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return errors.New("DB_HOST is required")
	}
	if c.Database.User == "" {
		return errors.New("DB_USER is required")
	}
	if c.Database.Password == "" {
		return errors.New("DB_PASSWORD is required")
	}
	if c.Database.DBName == "" {
		return errors.New("DB_NAME is required")
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil {
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}
