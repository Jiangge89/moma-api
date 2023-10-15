package db

import "context"

type RateI interface {
	GetRate(ctx context.Context, fromCode, toCode string) (*Rate, error)
	AddRate(ctx context.Context, fromCode, toCode string, rate float32) error
}
