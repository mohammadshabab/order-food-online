package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mohammadshabab/order-food-online/config"
	"github.com/mohammadshabab/order-food-online/internal/apperrors"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

// SQLPool provides a small wrapper so existing code can continue using db.Pool.Exec/Query/QueryRow
type SQLPool struct {
	DB *sql.DB
}

var Pool *SQLPool

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open mysql connection: %w", err)
	}

	// Pool settings
	dbConn.SetMaxOpenConns(cfg.DBMaxConns)
	dbConn.SetMaxIdleConns(cfg.DBMinConns)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	if err := dbConn.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping mysql: %w", err)
	}

	Pool = &SQLPool{DB: dbConn}
	logger.Info(ctx, "MySQL connected successfully")
	return nil
}

func (p *SQLPool) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := p.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	if err != nil {
		logger.Error(ctx, "DB Exec failed",
			"query", query,
			"args", args,
			"duration", duration.String(),
			"error", err.Error(),
		)
		return nil, apperrors.Internal("DB Exec failed", err)
	}

	logger.Info(ctx, "DB Exec success",
		"query", query,
		"args", args,
		"duration", duration.String(),
	)
	return res, nil
}

func (p *SQLPool) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := p.DB.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	if err != nil {
		logger.Error(ctx, "DB Query failed",
			"query", query,
			"args", args,
			"duration", duration.String(),
			"error", err.Error(),
		)
		return nil, apperrors.Internal("DB Query failed", err)
	}

	logger.Info(ctx, "DB Query success",
		"query", query,
		"args", args,
		"duration", duration.String(),
	)
	return rows, nil
}

func (p *SQLPool) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	logger.Info(ctx, "DB QueryRow",
		"query", query,
		"args", args,
	)
	return p.DB.QueryRowContext(ctx, query, args...)
}

func Close() {
	if Pool != nil && Pool.DB != nil {
		Pool.DB.Close()
	}
}
