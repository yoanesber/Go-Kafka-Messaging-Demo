package async

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Writers map[string]*kafka.Writer
	Readers map[string]*kafka.Reader
}

var (
	kafkaClient *KafkaClient
	once        sync.Once

	kafkaBrokers      []string
	kafkaTopics       []string
	kafkaGroupID      string
	kafkaReadTimeout  time.Duration
	kafkaWriteTimeout time.Duration
)

const (
	defaultKafkaGroupID        = "default-group"
	defaultKafkaReadTimeout    = 10 * time.Second
	defaultKafkaWriteTimeout   = 10 * time.Second
	defaultKafkaReaderMinBytes = int(10e3) // 10KB
	defaultKafkaReaderMaxBytes = int(10e6) // 10MB
)

func InitKafka() bool {
	isSuccess := true
	once.Do(func() {
		if !loadKafkaEnv() {
			isSuccess = false
			return
		}

		client := &KafkaClient{
			Writers: make(map[string]*kafka.Writer),
			Readers: make(map[string]*kafka.Reader),
		}

		for _, topic := range kafkaTopics {
			if topic == "" {
				fmt.Printf("Kafka topic cannot be empty\n")
				continue
			}

			// Initialize writer
			writer := initKafkaWriter(topic)
			client.Writers[topic] = writer

			// Initialize reader
			reader := initKafkaReader(topic)
			client.Readers[topic] = reader
		}

		kafkaClient = client
		fmt.Printf("Kafka client initialized with topics: %v\n", kafkaTopics)
	})

	return isSuccess
}

func GetKafkaWriter(topic string) (*kafka.Writer, error) {
	if kafkaClient == nil {
		return nil, fmt.Errorf("kafka client is not initialized")
	}

	writer, exists := kafkaClient.Writers[topic]
	if !exists {
		return nil, fmt.Errorf("kafka writer for topic %s does not exist", topic)
	}

	return writer, nil
}

func GetKafkaReader(topic string) (*kafka.Reader, error) {
	if kafkaClient == nil {
		return nil, fmt.Errorf("kafka client is not initialized")
	}

	reader, exists := kafkaClient.Readers[topic]
	if !exists {
		return nil, fmt.Errorf("kafka reader for topic %s does not exist", topic)
	}

	return reader, nil
}

func CloseKafka() {
	if kafkaClient != nil {
		for topic, writer := range kafkaClient.Writers {
			if err := writer.Close(); err != nil {
				fmt.Printf("Failed to close Kafka writer for topic %s: %v\n", topic, err)
				continue
			}

			fmt.Printf("Kafka writer for topic %s closed successfully\n", topic)
		}

		for topic, reader := range kafkaClient.Readers {
			if err := reader.Close(); err != nil {
				fmt.Printf("Failed to close Kafka reader for topic %s: %v\n", topic, err)
				continue
			}

			fmt.Printf("Kafka reader for topic %s closed successfully\n", topic)
		}

		kafkaClient = nil
		fmt.Println("Kafka client closed successfully")
		return
	}

	once = sync.Once{} // Reset the once to allow re-initialization
	kafkaClient = nil  // Clear the kafkaClient variable to prevent further use
}

func loadKafkaEnv() bool {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		fmt.Println("KAFKA_BROKERS must be set")
		return false
	}
	kafkaBrokers = strings.Split(brokers, ",")

	topics := os.Getenv("KAFKA_TOPICS")
	if topics == "" {
		fmt.Println("KAFKA_TOPICS must be set")
		return false
	}
	kafkaTopics = strings.Split(topics, ",")

	kafkaGroupID = os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = defaultKafkaGroupID
	}

	timeoutStr := os.Getenv("KAFKA_READ_TIMEOUT_MS")
	if timeoutStr == "" {
		kafkaReadTimeout = defaultKafkaReadTimeout
	} else {
		ms, err := strconv.Atoi(timeoutStr)
		if err != nil {
			fmt.Printf("Invalid KAFKA_READ_TIMEOUT_MS value: %s\n", timeoutStr)
			return false
		}
		kafkaReadTimeout = time.Duration(ms) * time.Millisecond
	}

	timeoutStr = os.Getenv("KAFKA_WRITE_TIMEOUT_MS")
	if timeoutStr == "" {
		kafkaWriteTimeout = defaultKafkaWriteTimeout
	} else {
		ms, err := strconv.Atoi(timeoutStr)
		if err != nil {
			fmt.Printf("Invalid KAFKA_WRITE_TIMEOUT_MS value: %s\n", timeoutStr)
			return false
		}
		kafkaWriteTimeout = time.Duration(ms) * time.Millisecond
	}

	return true
}

func initKafkaWriter(topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      kafkaBrokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: kafkaWriteTimeout,
	})
}

func initKafkaReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:         kafkaBrokers,
		Topic:           topic,
		GroupID:         kafkaGroupID,
		MinBytes:        defaultKafkaReaderMinBytes,
		MaxBytes:        defaultKafkaReaderMaxBytes,
		MaxWait:         kafkaReadTimeout,
		CommitInterval:  time.Second, // auto-commit offset tiap detik
		StartOffset:     kafka.FirstOffset,
		GroupBalancers:  []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}},
		ReadLagInterval: -1,
	})
}
