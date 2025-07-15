package entity

import (
	"time"
)

const (
	// MessageStatusPending indicates the Message is pending to be sent
	MessageStatusPending = "pending"
	// MessageStatusSent indicates the Message has been sent successfully
	MessageStatusSent = "sent"
	// MessageStatusDelivered indicates the Message has been delivered
	MessageStatusDelivered = "delivered"
	// MessageStatusFailed indicates the Message failed to send
	MessageStatusFailed = "failed"
)

type Message struct {
	ID         string    `json:"id"`                              // UUID, unique identifier for each Message
	SenderID   string    `json:"sender_id" validate:"required"`   // ID of the sender (could be a user ID or system ID)
	ReceiverID string    `json:"receiver_id" validate:"required"` // ID of the receiver (could be a user ID or system ID)
	Message    string    `json:"message" validate:"required"`     // Message content
	Timestamp  time.Time `json:"timestamp"`                       // When it was created/sent
	Status     string    `json:"status"`                          // Status of the message (e.g., "sent", "failed", "delivered")
}
