package kafka

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func StartConsumer(groupId string, topics []string, processMessage func(kafka.Message)) {
	kafkaHost := os.Getenv("KAFKA_HOST")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	if kafkaHost == "" || username == "" || password == "" {
		log.Fatal("Kafka host, username or password not set")
	}

	log.Printf("Kafka Host: %s, Username: %s", kafkaHost, username)

	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		KeepAlive:     10 * time.Second,
		SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:                []string{kafkaHost},
		GroupID:                groupId,
		Topic:                  topics[0],
		MinBytes:               10e3, // 10KB
		MaxBytes:               10e6, // 10MB
		Dialer:                 dialer,
		CommitInterval:         time.Second,
		PartitionWatchInterval: time.Second,
	})

	ctx := context.Background()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Consumer error: %v\n", err)
			continue
		}
		processMessage(msg)
		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("Failed to commit message: %v\n", err)
		}
	}
}
