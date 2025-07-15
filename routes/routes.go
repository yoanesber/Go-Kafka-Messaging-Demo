package routes

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/yoanesber/go-kafka-producer-consumer-demo/internal/handler"
	"github.com/yoanesber/go-kafka-producer-consumer-demo/internal/service"
	"github.com/yoanesber/go-kafka-producer-consumer-demo/pkg/middleware/headers"
)

func SetupRouter() *gin.Engine {
	// Create a new Gin router instance
	r := gin.Default()

	// Set up middleware for the router
	r.Use(
		headers.SecurityHeaders(),
		headers.CorsHeaders(),
		headers.ContentType(),
		gzip.Gzip(gzip.DefaultCompression),
	)

	// Set up the API group
	api := r.Group("/api")
	{
		// Set the service and handler for the message API
		s := service.NewMessageService()
		h := handler.NewMessageHandler(s)

		// Define the routes for the API
		api.POST("/send-message", h.SendMessage)
	}

	// This handler will be called when no other route matches the request
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Not Found",
			"message": "The requested resource could not be found",
		})
	})

	// This handler will be called when a request method is not allowed for the requested resource
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"error":   "Method Not Allowed",
			"message": "The requested method is not allowed for this resource",
		})
	})

	return r
}
