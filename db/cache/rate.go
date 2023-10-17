package cache

import (
	"context"
	"fmt"
	"moma-api/db"
	"sync"
	"time"
)

type RateCache struct {
	cache   map[string]db.Rate
	rwMutex sync.RWMutex
}

func (r RateCache) GetRate(ctx context.Context, fromCode, toCode string) (*db.Rate, error) {
	r.rwMutex.RLock()

	rate, ok := r.cache[fromCode+toCode]

	r.rwMutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("not exist")
	}

	return &rate, nil
}

func (r RateCache) AddRate(ctx context.Context, fromCode, toCode string, rate float32) error {
	now := int64(time.Now().Nanosecond())
	rateInfo := db.Rate{
		FromCode:  fromCode,
		ToCode:    toCode,
		Rate:      rate,
		UpdatedAt: now,
		CreatedAt: now,
	}

	r.rwMutex.Lock()
	if item, ok := r.cache[fromCode+toCode]; ok {
		rateInfo.CreatedAt = item.CreatedAt
	}

	r.cache[fromCode+toCode] = rateInfo
	r.rwMutex.Unlock()

	return nil
}

func NewRateCache() *RateCache {
	return &RateCache{
		cache: make(map[string]db.Rate, 10000),
	}
}
