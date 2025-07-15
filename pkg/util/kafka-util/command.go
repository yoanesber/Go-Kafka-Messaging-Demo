package kafka_util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/yoanesber/go-kafka-producer-consumer-demo/config/async"
)

const (
	retries      = 3
	maxWaitTime  = 10 * time.Second       // Maximum time to wait for the message to be written
	maxSleepTime = 250 * time.Millisecond // Maximum time to wait before retrying
)

func PublishMessage(topic string, key string, value interface{}) error {
	// Get the Kafka writer for the specified topic
	writer, err := async.GetKafkaWriter(topic)
	if err != nil {
		fmt.Printf("failed to get Kafka writer for topic %s: %v\n", topic, err)
		return err
	}

	// Marshal the value to JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Create a new message
	msg := kafka.Message{
		Key:   []byte(key),
		Value: valueBytes,
		Time:  time.Now(),
	}

	// Try to write the message to the topic
	// We will retry a few times in case of transient errors
	for range retries {
		ctx, cancel := context.WithTimeout(context.Background(), maxWaitTime)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = writer.WriteMessages(ctx, msg)
		if err != nil {
			if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
				time.Sleep(maxSleepTime)
				continue
			}

			fmt.Printf("failed to write message to topic %s: %v\n", topic, err)
		}
		break
	}

	return nil
}

func ConsumeMessages(workerID int, topic string, handler func(string, kafka.Message) error) {
	// Get the Kafka reader for the specified topic
	reader, err := async.GetKafkaReader(topic)
	if err != nil {
		fmt.Printf("failed to get Kafka reader for topic %s: %v\n", topic, err)
		return
	}

	for {
		worker := fmt.Sprintf("Worker-%d", workerID)

		// Read messages from the topic
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("failed to read message from topic %s: %v\n", topic, err)
			continue
		}

		// Call the handler function with the received message
		if err := handler(worker, msg); err != nil {
			fmt.Printf("failed to handle message from topic %s: %v\n", topic, err)
			continue
		}
	}
}
