package main

import (
	"log"
	"net/http"
	"os"

	kafka "github.com/MSPR-PayeTonKawa/auth/kafka"

	"github.com/MSPR-PayeTonKawa/auth/database"
	"github.com/MSPR-PayeTonKawa/auth/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Connect to the database
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	go kafka.StartConsumer("user-group", []string{"user-topic"}, kafka.ProcessMessage)

	// Set up Gin router
	r := gin.Default()

	// Define a route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	h := handlers.NewHandler(db)

	r.POST("/login", h.Login)
	r.POST("/refresh", h.Refresh)
	r.POST("/verify", h.VerifyToken)

	port := os.Getenv("PORT")

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
