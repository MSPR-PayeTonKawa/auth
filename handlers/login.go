package handlers

import (
	"log"
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
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// Prepare the SQL statement
	stmt, err := h.db.PrepareContext(c, "SELECT user_id, email, password_hash FROM users WHERE email = $1")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query(creds.Email)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	// Check if a user was found
	if !rows.Next() {
		log.Println("User not found:", creds.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
		return
	}

	// Scan the result into a UserCredentials object
	var user types.UserCredentials
	err = rows.Scan(&user.UserID, &user.Email, &user.Password)
	if err != nil {
		log.Println("Error scanning rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Check if user exists and password is correct
	if user.UserID == "" || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		log.Println("Invalid credentials for user:", creds.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
		return
	}

	tokens, err := tokens.CreateToken(user.UserID, user.Email)
	if err != nil {
		log.Println("Error creating tokens:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken})
}
