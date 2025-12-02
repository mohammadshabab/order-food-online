package promo

import (
	"compress/gzip"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mohammadshabab/order-food-online/internal/logger"
	"github.com/stretchr/testify/require"
)

func writeGzipFile(t *testing.T, dir, name string, lines []string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("cannot create gzip file: %v", err)
	}
	gz := gzip.NewWriter(f)

	for _, line := range lines {
		_, _ = gz.Write([]byte(line + "\n"))
	}
	gz.Close()
	f.Close()
	return path
}

func TestLoadCouponsWithContext(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := NewMockCache(ctrl)

	tmp := t.TempDir()

	// Create two gz files
	file1 := writeGzipFile(t, tmp, "a.gz", []string{"A1", "A2"})
	file2 := writeGzipFile(t, tmp, "b.gz", []string{"B1", "B2", ""})

	mockCache.EXPECT().Progress().AnyTimes()
	// We expect total = 2 files
	mockCache.EXPECT().SetTotalFiles(2)
	// For each file worker increments `IncrementLoaded`
	mockCache.EXPECT().IncrementLoaded().Times(2)

	// Expect coupon sets (order NOT guaranteed because workers are concurrent)
	mockCache.EXPECT().Set("A1", Coupon{Code: "A1"})
	mockCache.EXPECT().Set("A2", Coupon{Code: "A2"})
	mockCache.EXPECT().Set("B1", Coupon{Code: "B1"})
	mockCache.EXPECT().Set("B2", Coupon{Code: "B2"})

	// Final flags
	mockCache.EXPECT().MarkLoadedSuccessfully()
	mockCache.EXPECT().MarkReady()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cfg := LoaderConfig{
		Dir:         tmp,
		WorkerCount: 2,
	}

	err := LoadCouponsWithContext(ctx, cfg, mockCache)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Ensure files actually existed
	if _, err := os.Stat(file1); err != nil {
		t.Fatalf("file1 missing: %v", err)
	}
	if _, err := os.Stat(file2); err != nil {
		t.Fatalf("file2 missing: %v", err)
	}
}

func TestLoadOne(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.gz")

	// Create gz with 2 lines
	f, err := os.Create(file)
	require.NoError(t, err)
	gz := gzip.NewWriter(f)
	_, err = gz.Write([]byte("C1\nC2\n"))
	require.NoError(t, err)
	require.NoError(t, gz.Close())
	require.NoError(t, f.Close())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := NewMockCache(ctrl)

	// Expectations for loadOne
	mockCache.EXPECT().Set("C1", Coupon{Code: "C1"})
	mockCache.EXPECT().Set("C2", Coupon{Code: "C2"})

	err = loadOne(context.Background(), file, mockCache)
	require.NoError(t, err)
}
