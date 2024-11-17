package db

import (
	"context"
	"moma-api/db/model"
)

type RateI interface {
	GetRate(ctx context.Context, fromCode, toCode string) (*model.Rate, error)
	AddRate(ctx context.Context, fromCode, toCode string, rate float32) error
}

type AccountDB interface {
	GetAccount(ctx context.Context, userID string) (*model.Account, error)
	CreateAccount(ctx context.Context, account *model.Account) error
}
