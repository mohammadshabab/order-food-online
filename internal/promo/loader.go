package promo

import (
	"bufio"
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mohammadshabab/order-food-online/internal/logger"
)

type LoaderConfig struct {
	Dir         string
	WorkerCount int
}

func LoadCouponsWithContext(ctx context.Context, cfg LoaderConfig, cache Cache) error {
	files, err := filepath.Glob(filepath.Join(cfg.Dir, "*.gz"))
	if err != nil {
		return err
	}

	cache.SetTotalFiles(len(files))
	if len(files) == 0 {
		cache.MarkReady()
		logger.Info(ctx, "[Promo Loader] No .gz files found, marking ready")
		return nil
	}

	logger.Info(ctx, "[Promo Loader] Found %d files, starting load...", len(files))

	jobs := make(chan string)
	var wg sync.WaitGroup
	worker := cfg.WorkerCount
	if worker <= 0 {
		worker = 4
	}

	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				if err := loadOne(ctx, path, cache); err != nil {
					logger.Error(ctx, "Failed to load coupon file", "file", path, "error", err)
				} else {
					logger.Debug(ctx, "[Promo Loader] Loaded file: %s", path)
				}
				cache.IncrementLoaded()

				logger.Debug(ctx, "[Promo Loader] Progress: %.2f%%", cache.Progress()*100)
			}
		}()
	}

	go func() {
		for _, f := range files {
			select {
			case <-ctx.Done():
				break
			case jobs <- f:
			}
		}
		close(jobs)
	}()

	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		logger.Error(ctx, "[Promo Loader] Timeout while loading coupons")
		return ctx.Err()
	case <-doneCh:
		cache.MarkLoadedSuccessfully()
		cache.MarkReady()
		logger.Info(ctx, "[Promo Loader] All files processed successfully")
		return nil
	}
}

func loadOne(ctx context.Context, path string, cache Cache) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	scanner := bufio.NewScanner(gz)
	for scanner.Scan() {
		code := strings.TrimSpace(scanner.Text())
		if code == "" {
			continue
		}

		cache.Set(code, Coupon{Code: code})
	}
	return scanner.Err()
}
