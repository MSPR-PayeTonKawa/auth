package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MSPR-PayeTonKawa/auth/tokens"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	// Setup the Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the handler
	handler := Handler{}

	// Define the endpoint and attach the handler
	router.POST("/refresh", handler.Refresh)

	// Create a valid refresh token
	td, _ := tokens.CreateToken("user1")
	validRefreshToken := td.RefreshToken

	// Create a request to the endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer([]byte(`{"refresh_token": "`+validRefreshToken+`"}`)))
	req.Header.Set("Content-Type", "application/json")

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
