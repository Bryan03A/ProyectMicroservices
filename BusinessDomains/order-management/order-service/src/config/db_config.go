// config/db_config.go
package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/order_db")
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}
}
