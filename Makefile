# Variables for the application container
APP_CONTAINER_IMAGE=my-go-app
APP_CONTAINER_NAME=go-app
APP_DOCKER_CONTEXT=.
APP_DOCKERFILE=./docker/app/Dockerfile
APP_ENV_FILE=.env
APP_PORT=1000

# Variables for the Kafka container
KAFKA_CONTAINER_IMAGE=bitnami/kafka:latest
KAFKA_CONTAINER_NAME=kafka-server
KAFKA_PORT=9092
KAFKA_CLUSTER_ID=CcsfQr-zTweZlDPJZf-4EQ

# Network for the application and RabbitMQ containers
NETWORK=app-network

## ---- Development Commands ----
# Install dependencies
tidy:
	@echo -e "Running go mod tidy..."
	@go mod tidy

# Run the application in development mode
run:
	@echo -e "Running the application..."
	@dotenv -e .env -- go run ./cmd/main.go

# Test the application
test:
	@echo -e "Running tests..."
	@dotenv -e .env -- go test -v ./tests/...




## ---- Docker related targets ----
# Create a Docker network if it doesn't exist
docker-create-network:
	docker network inspect $(NETWORK) >NUL 2>&1 || docker network create $(NETWORK)

docker-remove-network:
	docker network rm $(NETWORK)




## --- Kafka related targets ---
# Run Kafka with KRaft mode (Kafka without Zookeeper)
docker-run-kafka:
	docker run --name $(KAFKA_CONTAINER_NAME) --network $(NETWORK) -p $(KAFKA_PORT):$(KAFKA_PORT) \
	-e KAFKA_CFG_NODE_ID=1 \
	-e KAFKA_CFG_PROCESS_ROLES=broker,controller \
	-e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@$(KAFKA_CONTAINER_NAME):9093 \
	-e KAFKA_CFG_LISTENERS=PLAINTEXT://:$(KAFKA_PORT),CONTROLLER://:9093 \
	-e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://$(KAFKA_CONTAINER_NAME):$(KAFKA_PORT) \
	-e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
	-e KAFKA_CFG_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
	-e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER \
	-e KAFKA_CLUSTER_ID=$(KAFKA_CLUSTER_ID) \
	-v kafka-data:/bitnami/kafka \
	-d $(KAFKA_CONTAINER_IMAGE)

# Remove Kafka container
docker-remove-kafka:
	docker stop $(KAFKA_CONTAINER_NAME)
	docker rm $(KAFKA_CONTAINER_NAME)




## --- Application related targets ---
docker-build-app:
	docker build -f $(APP_DOCKERFILE) -t $(APP_CONTAINER_IMAGE) $(APP_DOCKER_CONTEXT)

# Run the application container
docker-run-app:
	docker run --name $(APP_CONTAINER_NAME) --network $(NETWORK) -p $(APP_PORT):$(APP_PORT) \
	--env-file $(APP_ENV_FILE) \
	-d $(APP_CONTAINER_IMAGE)

# Build and run the application container
docker-build-run-app: docker-build-app docker-run-app

docker-remove-app:
	docker stop $(APP_CONTAINER_NAME)
	docker rm $(APP_CONTAINER_NAME)

docker-up: docker-create-network \
	docker-run-kafka \
	docker-build-run-app

docker-down: docker-remove-app \
	docker-remove-kafka \
	docker-remove-network

.PHONY: tidy run test \
	docker-create-network docker-remove-network \
	docker-run-kafka docker-remove-kafka \
	docker-build-app docker-run-app docker-build-run-app docker-remove-app \
	docker-up docker-down