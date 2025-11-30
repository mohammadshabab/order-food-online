package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Set env variables
	os.Setenv("ENV", "prod")
	os.Setenv("SERVICE", "test-service")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "1234")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_NAME", "order_db")
	os.Setenv("DB_MAX_CONNS", "50")
	os.Setenv("DB_MIN_CONNS", "10")
	os.Setenv("DB_CONN_LIFETIME_MIN", "15")
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	require.NoError(t, err)

	require.Equal(t, "prod", cfg.Env)
	require.Equal(t, "test-service", cfg.Service)
	require.Equal(t, "db.example.com", cfg.DBHost)
	require.Equal(t, "1234", cfg.DBPort)
	require.Equal(t, "admin", cfg.DBUser)
	require.Equal(t, "secret", cfg.DBPassword)
	require.Equal(t, "order_db", cfg.DBName)
	require.Equal(t, 50, cfg.DBMaxConns)
	require.Equal(t, 10, cfg.DBMinConns)
	require.Equal(t, 15, cfg.DBConnLife)
	require.Equal(t, "debug", cfg.LogLevel)
}

func TestLoadConfig_Defaults(t *testing.T) {
	os.Clearenv()

	cfg, err := LoadConfig()
	require.NoError(t, err)

	require.Equal(t, "dev", cfg.Env)
	require.Equal(t, "food-order-online", cfg.Service)
	require.Equal(t, "localhost", cfg.DBHost)
	require.Equal(t, "3306", cfg.DBPort)
	require.Equal(t, "root", cfg.DBUser)
	require.Equal(t, "mariadbpassword", cfg.DBPassword)
	require.Equal(t, "food_order", cfg.DBName)
	require.Equal(t, 20, cfg.DBMaxConns)
	require.Equal(t, 2, cfg.DBMinConns)
	require.Equal(t, 30, cfg.DBConnLife) // default
	require.Equal(t, "info", cfg.LogLevel)
}

func TestLoadConfig_InvalidConnLife_ShouldFallback(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_CONN_LIFETIME_MIN", "0")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	require.NoError(t, err)

	require.Equal(t, 30, cfg.DBConnLife) // should fallback to default
}
