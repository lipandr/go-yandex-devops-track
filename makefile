AGENT_BINARY=agentApp
SERVER_BINARY=serverApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
build_up: build_agent build_server
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

restart: down build_up

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	rm ./cmd/agent/agentApp || true
	rm ./cmd/server/serverApp || true
	@echo "Done!"

## build_agent: builds the broker binary as a linux executable
build_agent:
	@echo "Building agent binary..."
	cd ./cmd/agent && env GOOS=linux CGO_ENABLED=0 go build -o ${AGENT_BINARY} ./
	@echo "Done!"

## build_server: builds the logger binary as a linux executable
build_server:
	@echo "Building server binary..."
	cd ./cmd/server && env GOOS=linux CGO_ENABLED=0 go build -o ${SERVER_BINARY} ./
	@echo "Done!"
