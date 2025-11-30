package db

import (
	"database/sql"
)

// NewTestPool creates a test DBPool from an existing sql.DB (e.g., sqlmock)
func NewTestPool(sqlDB *sql.DB) *SQLPool {
	return &SQLPool{DB: sqlDB}
}
