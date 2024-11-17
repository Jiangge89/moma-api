package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"moma-api/db"
	"moma-api/db/model"
	"time"
)

type accountDB struct {
	db *sql.DB
}

func NewAccountDB() (db.AccountDB, error) {
	sqlCli, err := NewClient()
	if err != nil {
		return nil, err
	}

	return &accountDB{db: sqlCli}, nil
}

func (adb *accountDB) GetAccount(ctx context.Context, userID string) (*model.Account, error) {
	query := "SELECT id, user_id, user_name, user_email, created_at, updated_at FROM account WHERE user_id = ?"
	row := adb.db.QueryRowContext(ctx, query, userID)

	// 创建 Account 实例并填充数据
	account := &model.Account{}
	if err := row.Scan(&account.ID, &account.UserID, &account.UserName, &account.UserEmail, &account.CreatedAt, &account.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.NotFound
		}
		return nil, fmt.Errorf("failed to query account by user_id: %v due to: %w", userID, err)
	}

	return account, nil
}

func (adb *accountDB) CreateAccount(ctx context.Context, account *model.Account) error {
	createdAt := time.Now().UnixMilli()
	updatedAt := createdAt
	query := "INSERT INTO accounts (user_id, user_name, user_email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := adb.db.ExecContext(ctx, query, account.UserID, account.UserName, account.UserEmail, updatedAt, createdAt)
	if err != nil {
		return fmt.Errorf("failed to create account: %+v due to: %w", account, err)
	}
	return nil
}
