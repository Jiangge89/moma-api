package sql

import (
	"context"
	"database/sql"
	"fmt"
	"moma-api/db"
	"time"
)

type RateDB struct {
	db *sql.DB
}

func (r RateDB) GetRate(ctx context.Context, fromCode, toCode string) (*db.Rate, error) {
	var rate db.Rate
	allFields := "id, from_code, to_code, rate, created_at, updated_at"
	queryCmd := "select " + allFields + " from rate where from_code = ? and to_code = ?"
	rows := r.db.QueryRowContext(ctx, queryCmd, fromCode, toCode)
	err := rows.Scan(&rate.ID, &rate.FromCode, &rate.ToCode, &rate.Rate, &rate.CreatedAt, &rate.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("not exist")

	return &rate, nil
}

func (r RateDB) AddRate(ctx context.Context, fromCode, toCode string, rate float32) error {
	nowTimestamp := time.Now().Nanosecond()
	queryCmd := "insert into rate(from_code, to_code, rate, created_at, updated_at) values (?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE rate = ?, updated_at = ?"
	_, err := r.db.ExecContext(ctx, queryCmd, fromCode, toCode, rate, nowTimestamp, nowTimestamp, rate, nowTimestamp)
	if err != nil {
		return err
	}

	return nil
}

func NewRateDB(db *sql.DB) RateDB {
	return RateDB{db: db}
}
