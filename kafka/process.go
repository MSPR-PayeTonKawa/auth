package kafka

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/MSPR-PayeTonKawa/auth/database"
	"github.com/segmentio/kafka-go"
)

// ProcessMessage processes the Kafka message
func ProcessMessage(msg kafka.Message) {
	var user struct {
		UserID   int    `json:"user_id"`
		Email    string `json:"email"`
		Password string `json:"hashed_password"`
	}

	if err := json.Unmarshal(msg.Value, &user); err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return
	}

	db, err := database.ConnectDatabase()
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
		return
	}
	defer db.Close()

	if err := storeUser(db, user.UserID, user.Email, user.Password); err != nil {
		log.Printf("Error storing user data: %s", err)
	}
}

// storeUser stores the user data in the database
func storeUser(db *sql.DB, userID int, userEmail string, password string) error {
	_, err := db.Exec("INSERT INTO users (user_id, email, password_hash) VALUES ($1, $2, $3)", userID, userEmail, password)
	if err == nil {
		log.Printf("User with ID %d and email %s created successfully", userID, userEmail)
	}
	return err
}
