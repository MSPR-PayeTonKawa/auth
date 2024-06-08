package handlers

import (
	"net/http"

	"github.com/MSPR-PayeTonKawa/auth/tokens"
	"github.com/MSPR-PayeTonKawa/auth/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h Handler) VerifyToken(c *gin.Context) {
	var creds struct {
		AccessToken string `json:"access_token"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	claims := &types.Claims{}
	token, err := jwt.ParseWithClaims(creds.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return tokens.JwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": claims.UserID})
}
