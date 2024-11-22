AIR_DIR = ../.

AIR_BINARY = ../air-serivce

.PHONY: air-build
air-build:
	@echo "Building air server..."
	cd ${AIR_DIR} && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./airApp ./cmd/api

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: air-build
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"