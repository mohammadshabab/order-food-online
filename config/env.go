package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Env        string `env:"ENV, default=dev"`
	Service    string `env:"SERVICE, default=food-order-online"`
	DBHost     string `env:"DB_HOST, default=localhost"`
	DBPort     string `env:"DB_PORT, default=3306"`                // MariaDB default port
	DBUser     string `env:"DB_USER, default=root"`                // default MariaDB user
	DBPassword string `env:"DB_PASSWORD, default=mariadbpassword"` // default MariaDB password
	DBName     string `env:"DB_NAME, default=food_order"`
	DBMaxConns int    `env:"DB_MAX_CONNS, default=20"`
	DBMinConns int    `env:"DB_MIN_CONNS, default=2"`
	DBConnLife int    `env:"DB_CONN_LIFETIME_MIN, default=30"` // in minutes
	LogLevel   string `env:"LOG_LEVEL, default=info"`
	APIKey     string `env:"API_KEY, default=test"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(context.Background(), cfg); err != nil {
		return nil, err
	}
	// Ensure sensible defaults
	if cfg.DBConnLife <= 0 {
		cfg.DBConnLife = 30
	}

	return cfg, nil
}
