package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func connectDB() error {
	_, err := sql.Open("mysql", "")
	if err != nil {
		panic(err)
	}
	return nil
}
