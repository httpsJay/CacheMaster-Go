# Variables
APP_NAME=backend-take-home-ovxzsw
GO_FILES=$(wildcard *.go)

# Commands
.PHONY: all build run test clean install-deps

all: install-deps run

build:
	@echo "Building the application..."
	@go build -o $(APP_NAME) $(GO_FILES)

run: build
	@echo "Running the application..."
	./$(APP_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)

install-deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Docker Commands
.PHONY: docker-build docker-run docker-clean

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --name $(APP_NAME) -d $(APP_NAME)

docker-clean:
	@echo "Stopping and removing Docker container..."
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true
	@docker rmi $(APP_NAME) || true
