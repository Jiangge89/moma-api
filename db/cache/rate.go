package cache

import (
	"context"
	"fmt"
	"moma-api/db/model"
	"sync"
	"time"
)

type RateCache struct {
	cache   map[string]model.Rate
	rwMutex sync.RWMutex
}

func (r RateCache) GetRate(ctx context.Context, fromCode, toCode string) (*model.Rate, error) {
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
	rateInfo := model.Rate{
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
		cache: make(map[string]model.Rate, 10000),
	}
}
