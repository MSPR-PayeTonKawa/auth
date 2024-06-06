package handlers

import (
	"net/http"

	"github.com/MSPR-PayeTonKawa/auth/tokens"
	"github.com/MSPR-PayeTonKawa/auth/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login handles the login request and generates access and refresh tokens for the user.
// It takes a gin.Context object as a parameter and expects the user credentials to be provided in JSON format.
// If the request is valid and the user credentials are correct, it returns the access and refresh tokens in the response.
// If there is an error during the login process, it returns an appropriate error response.
func (h Handler) Login(c *gin.Context) {
	var creds types.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// Prepare the SQL statement
	stmt, err := h.db.Prepare("SELECT id, password_hash FROM users WHERE user_id = $1")
	if err != nil {
		// handle error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query(creds.UserID)
	if err != nil {
		// handle error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	// Check if a user was found
	if !rows.Next() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
		return
	}

	// Scan the result into a UserCredentials object
	var user types.UserCredentials
	err = rows.Scan(&user.UserID, &user.Password)
	if err != nil {
		// handle error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Check if user exists and password is correct
	if user.UserID == "" || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
		return
	}

	tokens, err := tokens.CreateToken(creds.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken})
}
