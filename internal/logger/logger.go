package logger

import (
	"log/slog"
	"os"
)

var (
	log *slog.Logger
	env string
)

func Init(serviceName, environment string, level slog.Level) {
	env = environment

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	log = slog.New(handler).With(
		slog.String("service", serviceName),
		slog.String("environment", environment),
	)
}

func Log() *slog.Logger {
	return log
}

func IsProd() bool {
	return env == "prod" || env == "production"
}

func IsDev() bool {
	return env == "dev" || env == "local"
}
