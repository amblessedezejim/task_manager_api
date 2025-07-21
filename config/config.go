package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		panic(err)
	}

	DB.SetMaxIdleConns(25)
	DB.SetMaxOpenConns(25)

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database ", err)
	}

	// Create table if table doesn't exist
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS tasks (id INT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255) NOT NULL, description TEXT, completed BOOLEAN DEFAULT FALSE, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)")
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func GetServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	return port
}
