package database

import (
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() Config {
	return Config{
		Host:     getEnv("DATABASE_HOST", "localhost"),
		Port:     getEnv("DATABASE_PORT", "5432"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Name:     os.Getenv("DATABASE_NAME"),
		SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
	}
}

func (c Config) Enabled() bool {
	return c.User != "" && c.Name != ""
}

func (c Config) DSN() string {
	return c.dsn(c.Name)
}

func (c Config) MaintenanceDSN() string {
	return c.dsn("postgres")
}

func (c Config) dsn(databaseName string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		databaseName,
		c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
