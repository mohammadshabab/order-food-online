package db

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mohammadshabab/order-food-online/config"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

// helper to create mock DB and set global Pool
func newMockPool(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	Pool = &SQLPool{DB: mockDB}

	cleanup := func() {
		mockDB.Close()
		Pool = nil
	}

	return mock, cleanup
}

func TestConnect(t *testing.T) {
	t.Run("connect returns error in test environment (expected)", func(t *testing.T) {
		cfg := &config.Config{
			DBUser:     "root",
			DBPassword: "pwd",
			DBHost:     "localhost",
			DBPort:     "3306",
			DBName:     "test",
		}

		err := Connect(cfg)
		assert.Error(t, err) // cannot connect to real DB in tests
	})
}

func TestSQLPool_Exec(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)
	t.Run("Exec success", func(t *testing.T) {
		mock, cleanup := newMockPool(t)
		defer cleanup()

		mock.ExpectExec("UPDATE users SET name").
			WithArgs("john", 10).
			WillReturnResult(sqlmock.NewResult(1, 1))

		res, err := Pool.Exec(context.Background(),
			"UPDATE users SET name = ? WHERE id = ?", "john", 10,
		)

		require.NoError(t, err)
		require.NotNil(t, res)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec error returns wrapped internal error", func(t *testing.T) {
		mock, cleanup := newMockPool(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM users").
			WithArgs(99).
			WillReturnError(errors.New("db error"))

		_, err := Pool.Exec(context.Background(),
			"DELETE FROM users WHERE id = ?", 99,
		)

		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "db error"))
	})
}

func TestSQLPool_Query(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)
	t.Run("Query success", func(t *testing.T) {
		mock, cleanup := newMockPool(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "alice").
			AddRow(2, "bob")

		mock.ExpectQuery("SELECT id, name FROM users").
			WillReturnRows(rows)

		r, err := Pool.Query(context.Background(), "SELECT id, name FROM users")
		require.NoError(t, err)
		require.NotNil(t, r)

		defer r.Close()

		count := 0
		for r.Next() {
			count++
		}

		assert.Equal(t, 2, count)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		mock, cleanup := newMockPool(t)
		defer cleanup()

		mock.ExpectQuery("SELECT .* FROM users").
			WillReturnError(errors.New("bad query"))

		_, err := Pool.Query(context.Background(), "SELECT * FROM users")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "bad query")
	})
}

func TestSQLPool_QueryRow(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)
	t.Run("QueryRow success", func(t *testing.T) {
		mock, cleanup := newMockPool(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(rows)

		row := Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users")

		var count int
		err := row.Scan(&count)

		require.NoError(t, err)
		assert.Equal(t, 5, count)
	})
}

func TestClose(t *testing.T) {
	t.Run("Close when pool is set", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)

		Pool = &SQLPool{DB: mockDB}

		assert.NotPanics(t, func() { Close() })
	})

	t.Run("Close when pool is nil", func(t *testing.T) {
		Pool = nil
		assert.NotPanics(t, func() { Close() })
	})
}
