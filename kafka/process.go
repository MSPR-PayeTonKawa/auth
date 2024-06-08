package kafka

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/MSPR-PayeTonKawa/auth/database"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"golang.org/x/crypto/bcrypt"
)

func ProcessMessage(msg *kafka.Message) {
	var user struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
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

	if err := storeUser(db, user.UserID, user.Password); err != nil {
		log.Printf("Error storing user data: %s", err)
	}
}

func storeUser(db *sql.DB, userID, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (user_id, password_hash) VALUES ($1, $2)", userID, string(hashedPassword))
	return err
}
