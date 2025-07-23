package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yoanesber/go-kafka-messaging-demo/internal/entity"
	"github.com/yoanesber/go-kafka-messaging-demo/internal/service"
	validation "github.com/yoanesber/go-kafka-messaging-demo/pkg/util/validation-util"
)

type MessageHandler struct {
	MessageService service.MessageService
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		MessageService: messageService,
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var message entity.Message

	// Bind JSON request to Message struct
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// Send the message using the MessageService
	// This will validate the message struct and publish it to Kafka
	if err := h.MessageService.SendMessage(c.Request.Context(), &message); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			// If validation errors, return 422 Unprocessable Entity
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   "Validation error",
				"details": validation.FormatValidationErrors(err),
			})
			return
		}
		c.JSON(500, gin.H{"error": "Internal server error", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Message sent successfully", "id": message.ID})
}
