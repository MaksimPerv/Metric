package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Init(connect string) error {
	var err error
	DB, err = sql.Open("pgx", connect)
	return err
}
