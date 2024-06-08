package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// ConnectDatabase establishes a connection to the database and returns a pointer to the sql.DB object.
// It retrieves the database connection details from environment variables and constructs the connection string.
// Returns an error if the connection cannot be established.
func ConnectDatabase() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DB_PORT")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	log.Println(dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
