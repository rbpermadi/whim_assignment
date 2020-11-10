package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Init returns connector to Soteria database
func NewMySQL() *sql.DB {

	var db *sql.DB

	env := os.Getenv("ENV")

	dbUsername := os.Getenv("DATABASE_USERNAME")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbName := os.Getenv("DATABASE_NAME")
	dbPort := os.Getenv("DATABASE_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	if env == "development" || env == "staging" {
		fmt.Printf("Connecting to [USERNAME]:[PASSWORD]@tcp(%s:%v)/%s?parseTime=true\n", dbHost, dbPort, dbName)
	}

	dataSourceName := fmt.Sprintf("%s:%v@(%s:%v)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, _ = sql.Open("mysql", dataSourceName)

	if err := db.Ping(); err != nil {
		panic(err.Error())
	}

	if dp, err := strconv.Atoi(os.Getenv("DATABASE_POOL")); err == nil && dp > 0 {
		db.SetMaxIdleConns(dp)
	}

	db.SetConnMaxLifetime(time.Minute)

	return db
}
