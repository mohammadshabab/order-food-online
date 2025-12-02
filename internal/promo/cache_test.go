package promo

import (
	"log/slog"
	"testing"
	"time"

	"github.com/mohammadshabab/order-food-online/internal/logger"
)

func TestCacheSetAndGet(t *testing.T) {
	c := NewCache()

	// First insert
	c.Set("ABC12345", Coupon{Code: "ABC12345"})
	cp, ok := c.Get("ABC12345")
	if !ok {
		t.Fatalf("expected coupon to exist")
	}
	if cp.FileCount != 1 {
		t.Fatalf("expected FileCount=1, got %d", cp.FileCount)
	}

	// Second insert increments FileCount
	c.Set("ABC12345", Coupon{Code: "ABC12345"})
	cp, _ = c.Get("ABC12345")
	if cp.FileCount != 2 {
		t.Fatalf("expected FileCount=2, got %d", cp.FileCount)
	}
}

func TestSetTotalFiles(t *testing.T) {
	c := NewCache()
	c.SetTotalFiles(5)

	// Progress without loaded
	if c.Progress() != 0 {
		t.Fatalf("expected progress=0")
	}
}

func TestIncrementLoadedAndReady(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)
	c := NewCache()
	c.SetTotalFiles(2)

	c.IncrementLoaded()
	if c.IsReady() {
		t.Fatalf("cache should not be ready yet")
	}

	c.IncrementLoaded() // now totalFiles == loaded
	if !c.IsReady() {
		t.Fatalf("cache should be ready")
	}

	if !c.LoadedSuccessfully() {
		t.Fatalf("expected loadedSuccessfully=true")
	}
}

func TestProgress(t *testing.T) {
	c := NewCache()
	c.SetTotalFiles(4)

	c.IncrementLoaded()
	c.IncrementLoaded()

	p := c.Progress()
	if p != 0.5 {
		t.Fatalf("expected progress=0.5, got=%f", p)
	}
}

func TestMarkReadyManually(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)
	c := NewCache()

	c.SetTotalFiles(5)
	c.IncrementLoaded() // loaded = 1
	c.MarkReady()       // forced early

	if !c.IsReady() {
		t.Fatalf("expected IsReady=true")
	}

	// Since loaded != totalFiles, LoadedSuccessfully should remain false
	if c.LoadedSuccessfully() {
		t.Fatalf("expected loadedSuccessfully=false")
	}
}

func TestWaitUntilReadySuccess(t *testing.T) {
	logger.Init("test-service", "test", slog.LevelInfo)
	c := NewCache()
	c.SetTotalFiles(1)

	go func() {
		time.Sleep(50 * time.Millisecond)
		c.IncrementLoaded() // Marks ready
	}()

	ok := c.WaitUntilReady(200 * time.Millisecond)
	if !ok {
		t.Fatalf("expected WaitUntilReady to return true")
	}
}

func TestWaitUntilReadyTimeout(t *testing.T) {
	c := NewCache()
	c.SetTotalFiles(1)

	ok := c.WaitUntilReady(50 * time.Millisecond)
	if ok {
		t.Fatalf("expected timeout and return false")
	}
}

func TestMarkLoadedSuccessfully(t *testing.T) {
	c := NewCache()

	if c.LoadedSuccessfully() {
		t.Fatalf("expected initial loadedSuccessfully=false")
	}

	c.MarkLoadedSuccessfully()

	if !c.LoadedSuccessfully() {
		t.Fatalf("expected loadedSuccessfully=true")
	}
}
