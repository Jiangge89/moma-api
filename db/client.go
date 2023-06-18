package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Client interface {
}

func NewClient() *sql.DB {
	db, err := sql.Open("mysql", "root:duftee2023@tcp(127.0.0.1:3306)/moma_api?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println("mysql conn error", err)
		return nil
	}

	return db
}
