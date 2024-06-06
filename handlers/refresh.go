package handlers

import (
	"net/http"
	"time"

	"github.com/MSPR-PayeTonKawa/auth/tokens"
	"github.com/MSPR-PayeTonKawa/auth/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h Handler) Refresh(c *gin.Context) {
	var creds struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	claims := &types.Claims{}
	_, _ = jwt.ParseWithClaims(creds.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return tokens.JwtKey, nil
	})

	expirationTime, err := claims.Claims.GetExpirationTime()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get expiration time"})
		return
	}

	if time.Until(time.Unix(expirationTime.Unix(), 0)) > 30*time.Second {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	tokens, err := tokens.CreateToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken})
}
