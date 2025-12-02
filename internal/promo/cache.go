package promo

import (
	"sync"
	"time"

	"github.com/mohammadshabab/order-food-online/internal/logger"
)

//go:generate mockgen -source=cache.go -destination=mock_cache.go -package=promo
type Cache interface {
	Set(code string, cp Coupon)
	Get(code string) (Coupon, bool)
	SetTotalFiles(n int)
	IncrementLoaded()
	Progress() float64
	MarkReady()
	IsReady() bool
	WaitUntilReady(timeout time.Duration) bool
	MarkLoadedSuccessfully()
	LoadedSuccessfully() bool
}

type couponCache struct {
	mu                 sync.RWMutex
	store              map[string]Coupon
	totalFiles         int
	loaded             int
	ready              bool
	readyCh            chan struct{}
	loadedSuccessfully bool
}

func NewCache() Cache {
	return &couponCache{
		store:   make(map[string]Coupon),
		readyCh: make(chan struct{}),
	}
}

func (c *couponCache) Set(code string, cp Coupon) {
	c.mu.Lock()
	existing, ok := c.store[code]
	if ok {
		cp.FileCount = existing.FileCount + 1
	} else {
		cp.FileCount = 1
	}
	c.store[code] = cp
	c.mu.Unlock()
}

func (c *couponCache) Get(code string) (Coupon, bool) {
	c.mu.RLock()
	v, ok := c.store[code]
	c.mu.RUnlock()
	return v, ok
}

func (c *couponCache) SetTotalFiles(n int) {
	c.mu.Lock()
	c.totalFiles = n
	c.mu.Unlock()
}

func (c *couponCache) IncrementLoaded() {
	c.mu.Lock()
	c.loaded++
	done := c.loaded == c.totalFiles
	c.mu.Unlock()

	if done {
		c.MarkReady()
	}
}

func (c *couponCache) Progress() float64 {
	c.mu.RLock()
	total := c.totalFiles
	loaded := c.loaded
	c.mu.RUnlock()

	if total == 0 {
		return 0
	}
	return float64(loaded) / float64(total)
}

func (c *couponCache) MarkReady() {
	c.mu.Lock()
	if !c.ready {
		c.ready = true
		close(c.readyCh)

		if c.loaded == c.totalFiles {
			c.loadedSuccessfully = true
			logger.Log().Info("[Promo Loader] All coupon files loaded successfully")
		} else {
			logger.Log().Warn("[Promo Loader] Marked ready but not all files were loaded")
		}
	}
	c.mu.Unlock()
}

func (c *couponCache) IsReady() bool {
	c.mu.RLock()
	r := c.ready
	c.mu.RUnlock()
	return r
}

func (c *couponCache) WaitUntilReady(timeout time.Duration) bool {
	select {
	case <-c.readyCh:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (c *couponCache) MarkLoadedSuccessfully() {
	c.mu.Lock()
	c.loadedSuccessfully = true
	c.mu.Unlock()
}

func (c *couponCache) LoadedSuccessfully() bool {
	c.mu.RLock()
	v := c.loadedSuccessfully
	c.mu.RUnlock()
	return v
}
