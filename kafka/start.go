package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func StartConsumer(brokers string, groupId string, topics []string, processMessage func(*kafka.Message)) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:29092", // Changed from `kafka:9092` to `localhost:29092`
		"group.id":           "yourGroupId",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %s", err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			processMessage(msg)
			c.CommitMessage(msg) // Commit the message after processing
		} else {
			// handle error
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}