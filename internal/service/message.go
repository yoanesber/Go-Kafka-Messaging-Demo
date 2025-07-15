package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/yoanesber/go-kafka-producer-consumer-demo/internal/entity"
	kafkautil "github.com/yoanesber/go-kafka-producer-consumer-demo/pkg/util/kafka-util"
	validator "github.com/yoanesber/go-kafka-producer-consumer-demo/pkg/util/validation-util"
)

const (
	TopicMessage = "messaging" // Kafka topic for messages
)

type MessageService interface {
	SendMessage(ctx context.Context, message *entity.Message) error
	ReadMessage(worker string, message *entity.Message) error
}

type messageService struct {
	// You can add fields here if needed, e.g., a repository or a logger
}

func NewMessageService() MessageService {
	return &messageService{}
}

func (s *messageService) SendMessage(ctx context.Context, message *entity.Message) error {
	// Default values
	message.ID = uuid.New().String() // Generate a new UUID for the Message
	message.Timestamp = time.Now()   // Set the current timestamp
	message.Status = entity.MessageStatusPending

	// Validate the message struct
	if err := validator.ValidateStruct(message); err != nil {
		return err
	}

	// Create the event with the message as payload
	messageEvent := entity.MessageEvent{
		EventType: entity.EventTypeSendingMessage,
		Payload:   *message, // Use the message struct as the payload
	}

	// Publish message to Kafka
	if err := kafkautil.PublishMessage(TopicMessage, message.ID, messageEvent); err != nil {
		message.Status = entity.MessageStatusFailed
	} else {
		message.Status = entity.MessageStatusSent
	}

	// Here you can add additional logic, such as saving the message to a database
	// or performing other operations after sending the message.

	// Log the Message sent
	fmt.Println()
	fmt.Printf("Sending message:\n"+
		"ID: %s\nSenderID: %s\nReceiverID: %s\nMessage: %s\n"+
		"Status: %s\nTimestamp: %s\n",
		message.ID, message.SenderID, message.ReceiverID, message.Message,
		message.Status, message.Timestamp.Format(time.RFC3339))

	return nil
}

func (s *messageService) ReadMessage(worker string, message *entity.Message) error {
	fmt.Println()
	fmt.Printf("Reading message by %s:\n"+
		"ID: %s\nSenderID: %s\nReceiverID: %s\nMessage: %s\n",
		worker, message.ID, message.SenderID, message.ReceiverID, message.Message)

	// Here you can add logic to process the message, such as updating its status
	message.Status = entity.MessageStatusDelivered

	// then you can save the updated message to a database or perform other actions

	return nil
}
