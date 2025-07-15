package entity

const (
	EventTypeSendingMessage = "sending-message" // Event type for sending messages
)

type MessageEvent struct {
	EventType string  `json:"event_type"` // Type of the event, e.g., "sending-message"
	Payload   Message `json:"payload"`    // The message payload
}
