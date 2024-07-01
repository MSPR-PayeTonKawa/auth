package database

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const maxRetries = 10
const retryDelay = time.Second

func ConnectDatabase() (*sql.DB, error) {
	return ConnectDatabaseUsing(nil, 0)
}

func ConnectDatabaseUsing(db *sql.DB, retryCount int) (*sql.DB, error) {
	if db == nil {
		var err error
		db, err = sql.Open("postgres", constructDSN())
		if err != nil {
			return nil, err
		}
	}

	err := db.Ping()
	if err != nil {
		if retryCount < maxRetries {
			time.Sleep(retryDelay)
			return ConnectDatabaseUsing(db, retryCount+1)
		}
		return nil, err
	}

	return db, nil
}

func constructDSN() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DB_PORT")
	return "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
}
