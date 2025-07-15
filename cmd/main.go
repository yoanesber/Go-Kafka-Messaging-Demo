package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/yoanesber/go-kafka-producer-consumer-demo/config/async"
	kafka "github.com/yoanesber/go-kafka-producer-consumer-demo/pkg/kafka"
	validation "github.com/yoanesber/go-kafka-producer-consumer-demo/pkg/util/validation-util"
	"github.com/yoanesber/go-kafka-producer-consumer-demo/routes"
)

var (
	kafkaInitialized     bool
	validatorInitialized bool
)

func main() {
	// Create base context with cancel for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get environment variables
	env := os.Getenv("ENV")
	port := os.Getenv("PORT")
	isSSL := os.Getenv("IS_SSL")
	apiVersion := os.Getenv("API_VERSION")

	if env == "" || port == "" || isSSL == "" || apiVersion == "" {
		fmt.Println("Environment variables ENV, PORT, IS_SSL, and API_VERSION must be set.")
		return
	}

	// Set Gin mode
	gin.SetMode(gin.DebugMode)
	if env == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	r := routes.SetupRouter()
	r.SetTrustedProxies(nil) // Set trusted proxies to nil to avoid issues with forwarded headers

	// Init all dependencies
	initializeDependencies()

	// Graceful shutdown
	gracefulShutdown(cancel)

	// Start the server
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server on port %s: %v\n", port, err)
		return
	}
}

func initializeDependencies() {
	if !validatorInitialized {
		if !validation.Init() {
			fmt.Println("Failed to initialize validator. Exiting...")
		} else {
			validatorInitialized = true
		}
	}

	if !kafkaInitialized {
		if !async.InitKafka() {
			fmt.Println("Failed to initialize Kafka. Exiting...")
		} else {
			kafkaInitialized = true

			// Start consuming messages from Kafka
			fmt.Println("Starting Kafka message consumption...")
			kafka.StartConsumer()
			fmt.Println("Kafka message consumption started.")
		}
	}
}

func gracefulShutdown(cancel context.CancelFunc) {
	// Handle graceful shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		fmt.Printf("Received signal: %s. Initiating graceful shutdown...\n", sig)

		// Cancel context
		cancel()

		// Clean up resources
		if kafkaInitialized {
			fmt.Println("Closing Kafka connections...")
			async.CloseKafka()
		}

		if validatorInitialized {
			fmt.Println("Clearing validator...")
			validation.ClearValidator()
		}

		os.Exit(0)
	}()
}
