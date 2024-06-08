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

func TestVerifyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := Handler{}
	router.POST("/verify", handler.VerifyToken)

	td, _ := tokens.CreateToken("user1")
	validAccessToken := td.AccessToken

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer([]byte(`{"access_token": "`+validAccessToken+`"}`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user1", response["user_id"])
}
