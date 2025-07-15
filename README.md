# Message Processor Service

This project is a demonstration of **Kafka-based asynchronous messaging** using **Golang**, the **Gin web framework**, and the **Segmentio Kafka client** (`github.com/segmentio/kafka-go`). It implements a stateless HTTP API for sending and consuming **message events** in a distributed, non-blocking manner.

Each message consists of sender and receiver information, content, timestamp, and status. When a new message is submitted via the API, it is wrapped into a `MessageEvent` and **published to a Kafka broker**. Meanwhile, **Kafka consumers run as worker goroutines**, fetching and processing events from the broker concurrently. The number of active workers is configured via an environment variable.

---


## ✨ Features

Below are the key features provided by this project:

- 📨 Message API  
  - Exposes an HTTP endpoint to accept new message requests via POST.  
  - Each message includes `SenderID`, `ReceiverID`, and the `Message` body.  
  - Validates input using struct tags and provides meaningful error responses.  

- 🧵 Kafka Event Producer
  - StWraps each message into a `MessageEvent` struct with an event type (e.g., `"sending-message"`).  
  - Publishes the event to a Kafka topic using the Segmentio Kafka writer.  

- 🛠 Concurrent Kafka Consumers
  - Starts multiple consumer workers in parallel based on the `.env` configuration (`KAFKA_CONSUMER_WORKERS`).  
  - Each worker listens to the Kafka topic, deserializes the `MessageEvent`, and invokes a custom handler.  
  - Includes worker indexing for tracking which worker processed which message.  

- 🔁 Event Status Tracking
  - Each message can have a status such as `pending`, `sent`, `delivered`, or `failed`.  
  - Demonstrates how status can be handled and updated in a streaming pipeline (future expansion possible).  

---

## 🧭 Business Process Flow

The following diagram illustrates the end-to-end flow of how a message is processed by the system — from HTTP **request submission**, **message validation**, **Kafka publishing**, to **asynchronous consumption** by worker consumers.

```pgsql
┌──────────────────────────────────────────────┐
│            [1] Client Sends Request          │
│----------------------------------------------│
│ - POST /send-message                         │
│ - Body: { sender_id, receiver_id, message }  │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│      [2] Gin Handler Validates Payload       │
│----------------------------------------------│
│ - Check required fields                      │
│ - Add UUID, timestamp, and default status    │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│     [3] Wrap Message into MessageEvent       │
│----------------------------------------------│
│ - EventType = "sending-message"              │
│ - Payload = Message struct                   │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│        [4] Publish to Kafka Topic            │
│----------------------------------------------│
│ - Uses segmentio Kafka writer                │
│ - Topic = "messaging"                        │
└──────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│    [5] Kafka Consumers Start (N Workers)        │
│-------------------------------------------------│
│ - Configurable via .env (KAFKA_CONSUMER_WORKERS)│
│ - Each worker listens on same group ID          │
└─────────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│   [6] MessageEvent is Processed by Worker    │
│----------------------------------------------│
│ - Deserialized from Kafka message            │
│ - Printed/logged by worker                   │
│ - Status updated (if needed)                 │
└──────────────────────────────────────────────┘

```
---


## 🤖 Tech Stack

This project utilizes a concise and modern tech stack focused on **asynchronous messaging**, **stateless API design**, and **high concurrency**. The selected tools ensure simplicity, maintainability, and scalability.

| **Component**             | **Description**                                                                             |
|---------------------------|---------------------------------------------------------------------------------------------|
| **Language**              | Go (Golang), a statically typed, compiled language known for concurrency and efficiency     |
| **Web Framework**         | Gin, a fast and minimalist HTTP web framework for Go                                        |
| **Validation**            | `go-playground/validator.v9` for input validation and data integrity enforcement            |
| **Message Broker**        | Kafka, used for publishing and consuming messaging events asynchronously                    |

---

## 🧱 Architecture Overview

This project follows a **modular** and **maintainable** architecture inspired by **Clean Architecture** principles. Each domain feature (e.g., **entity**, **handler**, **service**) is organized into self-contained modules with clear separation of concerns.

```bash
📁 go-kafka-messaging-demo/
├── 📂cmd/                                  # Contains the application's entry point.
├── 📂config/
│   └── 📂async/                            # Config for async-related components, like Kafka producer/consumer settings
├── 📂docker/                               # Docker-related configuration for building and running services
│   └── 📂app/                              # Contains Dockerfile to build the main Go application image
├── 📂internal/                             # Core domain logic and business use cases, organized by module
│   ├── 📂entity/                           # Data models/entities representing business concepts
│   ├── 📂handler/                          # HTTP handlers (controllers) that parse requests and return responses
│   └── 📂service/                          # Business logic layer orchestrating operations
├── 📂pkg/                                  # Reusable utility and middleware packages shared across modules
│   ├── 📂kafka/                            # Kafka-related abstractions and interfaces
│   │   └── 📂handler/                      # Kafka consumer handlers to process consumed messages
│   ├── 📂middleware/                       # Request processing middleware
│   │   └── 📂headers/                      # Manages request headers like CORS, security
│   └── 📂util/                             # General utility functions and helpers
│       ├── 📂kafka-util/                   # Kafka configuration and utility helpers
│       └── 📂validation-util/              # Common input validators (e.g., UUID, numeric range)
└── 📂routes/                               # Route definitions, groups APIs, and applies middleware per route scope
```

---

## 🛠️ Installation & Setup  

Follow the instructions below to get the project up and running in your local development environment. You may run it natively or via Docker depending on your preference.  

### ✅ Prerequisites

Make sure the following tools are installed on your system:

| **Tool**                                                      | **Description**                           |
|---------------------------------------------------------------|-------------------------------------------|
| [Go](https://go.dev/dl/)                                      | Go programming language (v1.20+)          |
| [Make](https://www.gnu.org/software/make/)                    | Build automation tool (`make`)            |
| [Apache Kafka](https://kafka.apache.org/)                     | Distributed event streaming platform for async processing |
| [Docker](https://www.docker.com/)                             | Containerization platform (optional)      |

### 🔁 Clone the Project  

Clone the repository:  

```bash
git clone https://github.com/yoanesber/Go-Kafka-Messaging-Demo.git
cd Go-Kafka-Messaging-Demo
```

### ⚙️ Configure `.env` File  

Set up your **kafka** and **application** in `.env` file. Create a `.env` file at the project root directory:  

```properties
# Application configuration
ENV=PRODUCTION
API_VERSION=1.0
PORT=1000
IS_SSL=FALSE

# Kafka configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPICS=messaging
KAFKA_GROUP_ID=messaging-group
KAFKA_READ_TIMEOUT_MS=5000
KAFKA_WRITE_TIMEOUT_MS=5000
KAFKA_CONSUMER_WORKERS=3
```
---


## 🚀 Running the Application  

This section provides step-by-step instructions to run the application either **locally** or via **Docker containers**.

- **Notes**:  
  - All commands are defined in the `Makefile`.
  - To run using `make`, ensure that `make` is installed on your system.
  - To run the application in containers, make sure `Docker` is installed and running.
  - Ensure you have `Go` installed on your system

### 📦 Install Dependencies

Make sure all Go modules are properly installed:  

```bash
make tidy
```

### 🧪 Run Unit Tests

```bash
make test
```

### 🔧 Run Locally (Non-containerized)

Ensure Kafka are running locally, then:

```bash
make run
```

### 🐳 Run Using Docker

To build and run all services (Kafka, Go app):

```bash
make docker-up
```

To stop and remove all containers:

```bash
make docker-down
```

- **Notes**:  
  - Before running the application inside Docker, make sure to update your environment variables `.env`
    - Change `KAFKA_BROKERS=localhost:9092` to `KAFKA_BROKERS=kafka-server:9092`.

### 🟢 Application is Running

Now your application is accessible at:
```bash
http://localhost:1000
```

---

## 🧪 Testing Scenarios  

### 🔐 Sending Message

**Endpoint**: `POST http://localhost:1000/api/send-message`

**Request**:

```json
{
  "sender_id": "a2f3cbe1-0e4e-4b3b-bb7e-8ff9b6d4a124",
  "receiver_id": "f4a1e8d7-22d7-4b3a-b6d1-c9ea2ff6a9b3",
  "message": "Hello, how are you doing today?"
}
```

**Response**:

```json
{
    "id": "f38d7d4d-5da0-4188-a314-9b94f85c090c",
    "message": "Message sent successfully"
}
```

**On Producing the Message**:

Once the API successfully receives the request, the message will be packaged and published to a Kafka topic. The log in the terminal will look like this:

```bash
Sending message:
ID: f38d7d4d-5da0-4188-a314-9b94f85c090c
SenderID: a2f3cbe1-0e4e-4b3b-bb7e-8ff9b6d4a124
ReceiverID: f4a1e8d7-22d7-4b3a-b6d1-c9ea2ff6a9b3
Message: Hello, how are you doing today?
Status: sent
Timestamp: 2025-06-22T16:02:12+07:00
```

**Note**:
- **ID**: Unique UUID generated for the message
- **SenderID / ReceiverID**: Sender and receiver IDs
- **Status**: Will be set to sent when published
- **Timestamp**: The sending time is printed in the local time zone

**On Consuming the Message**:

Each published message will be read by one of the Kafka workers. Example log from a worker:

```bash
Reading message by Worker-0:
ID: f38d7d4d-5da0-4188-a314-9b94f85c090c
SenderID: a2f3cbe1-0e4e-4b3b-bb7e-8ff9b6d4a124
ReceiverID: f4a1e8d7-22d7-4b3a-b6d1-c9ea2ff6a9b3
Message: Hello, how are you doing today?
```

**Note**:
- Indicates that `Worker-0` successfully received and read the message
- This helps verify the consumption process is running according to the Kafka worker count configuration.