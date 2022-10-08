package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"

	"github.com/m9rc1n/shop/pkg/log"
)

// Config holds data for application configuration
type Config struct {
	DB                   *Database `yaml:"database,omitempty"`
	ReservationsEndpoint string    `yaml:"reservationsEndpoint"`
}

// Database holds data for database configuration
type Database struct {
	Dsn string `yaml:"dsn,omitempty"`
}

// New returns config from environment variables.
func New(logger log.Logger) (*Config, error) {
	logger.Infof("loading configuration")
	err := godotenv.Load()
	if err != nil {
		panic("error loading .env file")
	}
	return &Config{
		ReservationsEndpoint: os.Getenv("RESERVATIONS_ENDPOINT"),
		DB: &Database{
			Dsn: fmt.Sprintf(
				"postgresql://%s:%s@%s:%s/%s?&sslmode=%s",
				os.Getenv("SHOP_DB_USER"),
				os.Getenv("SHOP_DB_PASSWORD"),
				os.Getenv("SHOP_DB_HOST"),
				os.Getenv("SHOP_DB_PORT"),
				os.Getenv("SHOP_DB_NAME"),
				os.Getenv("SHOP_DB_SSLMODE"),
			),
		},
	}, nil
}
