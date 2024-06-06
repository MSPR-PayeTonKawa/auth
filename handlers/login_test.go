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
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := NewHandler(db)
	router.POST("/login", handler.Login)

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"user_id", "password_hash"}).AddRow("user1", string(passwordHash))

	mock.ExpectPrepare("SELECT id, password_hash FROM users WHERE user_id = \\$1").
		ExpectQuery().
		WithArgs("user1").
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"user_id": "user1", "password": "pass123"}`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])
	assert.NotEmpty(t, response["refresh_token"])
}
