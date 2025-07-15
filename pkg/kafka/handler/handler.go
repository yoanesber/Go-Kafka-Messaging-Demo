package handler

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/yoanesber/go-kafka-producer-consumer-demo/internal/entity"
	"github.com/yoanesber/go-kafka-producer-consumer-demo/internal/service"
)

func HandleMessaging(worker string, msg kafka.Message) error {
	// Unmarshal the message value into a MessageEvent struct
	var msgEvent entity.MessageEvent
	if err := json.Unmarshal(msg.Value, &msgEvent); err != nil {
		return fmt.Errorf("failed to unmarshal message value: %w", err)
	}

	// Map the event type to transaction handler
	eventType := msgEvent.EventType
	switch eventType {
	case entity.EventTypeSendingMessage:
		s := service.NewMessageService()
		if err := s.ReadMessage(worker, &msgEvent.Payload); err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}
	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}

	return nil
}
