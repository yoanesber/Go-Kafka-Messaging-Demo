package kafka

import (
	"os"
	"strconv"

	"github.com/yoanesber/go-kafka-messaging-demo/pkg/kafka/handler"
	kafkautil "github.com/yoanesber/go-kafka-messaging-demo/pkg/util/kafka-util"
)

const (
	topicMessage = "messaging"
)

func StartConsumer() {
	// Get the number of workers from the environment variable, default to 1 if not set
	numWorkersStr := os.Getenv("KAFKA_CONSUMER_WORKERS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil || numWorkers <= 0 {
		numWorkers = 1
	}

	for i := 0; i < numWorkers; i++ {
		go kafkautil.ConsumeMessages(i, topicMessage, handler.HandleMessaging)
	}
}
