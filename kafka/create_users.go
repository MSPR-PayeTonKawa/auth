package kafka

import (
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func ListenToCreateUserEvents() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"user_creation_topic"}, nil)

	defer c.Close()

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			// Here you can add the logic to extract the user details from the message
			// and store them in your service's database.
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
