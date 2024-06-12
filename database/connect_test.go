package database

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	// Create a new mock SQL connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "user")
	os.Setenv("POSTGRES_PASSWORD", "pass")
	os.Setenv("POSTGRES_DB", "dbname")
	os.Setenv("DB_PORT", "5432")

	// Expect a database ping
	mock.ExpectPing()

	// Invoke the actual function that uses the mocked DB connection
	_, err = ConnectDatabaseUsing(db) // This is a hypothetical function that you would need to implement to use the passed DB instance
	assert.NoError(t, err)

	// Check if all expectations set on the mock database connection were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
