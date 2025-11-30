package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mohammadshabab/order-food-online/config"
	"github.com/mohammadshabab/order-food-online/internal/db"
	"github.com/mohammadshabab/order-food-online/internal/health"
	"github.com/mohammadshabab/order-food-online/internal/logger"
	"github.com/mohammadshabab/order-food-online/internal/middleware"

	"github.com/mohammadshabab/order-food-online/internal/order"
	"github.com/mohammadshabab/order-food-online/internal/product"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger
	var level slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	logger.Init(cfg.Service, cfg.Env, level)

	// Initialize DB pool
	if err := db.Connect(cfg); err != nil {
		logger.Log().Error("failed to connect to database", "error", err)
		log.Fatalf("db connect failed: %v", err)
	}
	defer db.Close()

	e := echo.New()

	// Apply API key middleware (required by OpenAPI spec for /order endpoint)
	e.Use(middleware.NewAPIKeyMiddleware(*cfg))

	// Setup health check routes
	health.Register(e)

	// Product module
	productRepo := product.NewMariaDBRepository()
	product.Setup(e, productRepo)

	// Order module
	orderRepo := order.NewMariaDBRepository()
	order.Setup(e, orderRepo)

	// Start server in a goroutine
	go func() {
		logger.Log().Info("starting server on :8080")
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			logger.Log().Error("server start error", "error", err)
		}
	}()

	// Wait for interrupt signal (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown with 10 second timeout
	logger.Log().Info("shutting down server gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Log().Error("server shutdown error", "error", err)
		log.Fatalf("server forced to shutdown: %v", err)
	}

	logger.Log().Info("server stopped successfully")
}
