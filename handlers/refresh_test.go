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
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := Handler{}
	router.POST("/refresh", handler.Refresh)

	td, _ := tokens.CreateToken("user1")
	validRefreshToken := td.RefreshToken

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer([]byte(`{"refresh_token": "`+validRefreshToken+`"}`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])
	assert.NotEmpty(t, response["refresh_token"])
}
