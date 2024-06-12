package handlers

import (
	"log"
	"net/http"

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
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	log.Printf("Received refresh token: %s", creds.RefreshToken)

	claims := &types.Claims{}
	token, err := jwt.ParseWithClaims(creds.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return tokens.JwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Printf("Error parsing token with claims: %v", err)
		log.Printf("Token valid: %v", token.Valid)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	log.Printf("Parsed claims: %+v", claims)

	tokens, err := tokens.CreateToken(claims.UserID, claims.Email)
	if err != nil {
		log.Printf("Error creating tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create tokens"})
		return
	}

	log.Printf("Created tokens: %+v", tokens)

	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken})
}
