package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// Setup the Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the database connection and handler
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := NewHandler(db)

	// Define the endpoint and attach the handler
	router.POST("/login", handler.Login)

	// Create a request to the endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"user_id": "user1", "password": "pass123"}`)))
	req.Header.Set("Content-Type", "application/json")

	// Mock the database response
	rows := sqlmock.NewRows([]string{"user_id", "password_hash"}).AddRow("user1", "$2a$10$7s1tgjVVfOw/lxP8CPSU9OBFz8z0otj8JEB1ijROnMCRS/d4RoEOm") // bcrypt hash for "pass123"
	mock.ExpectPrepare("SELECT id, password_hash FROM users WHERE user_id = ?").ExpectQuery().WithArgs("user1").WillReturnRows(rows)

	// Perform the request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])
	assert.NotEmpty(t, response["refresh_token"])
}
