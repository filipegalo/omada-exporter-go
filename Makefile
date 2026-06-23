.PHONY: help build run test test-race cover fmt vet lint tidy clean docker-build docker-run

APP_NAME    := omada_exporter
APP_VERSION ?= 0.0.0
IMAGE       := $(APP_NAME):$(APP_VERSION)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	CGO_ENABLED=0 go build -o $(APP_NAME) .

run: ## Run the exporter locally
	go run .

test: ## Run tests
	go test ./...

test-race: ## Run tests with the race detector
	go test -race ./...

cover: ## Run tests with coverage report
	go test -cover ./...

fmt: ## Format source
	gofmt -w .

vet: ## Run go vet
	go vet ./...

lint: fmt vet ## Format and vet

tidy: ## Tidy go modules
	go mod tidy

clean: ## Remove build artifacts
	rm -f $(APP_NAME)

docker-build: ## Build the Docker image
	docker build --build-arg APP_VERSION=$(APP_VERSION) -t $(IMAGE) .

docker-run: ## Run the container from .env
	docker run -d --env-file .env -p 8080:8080 --name $(APP_NAME) $(IMAGE)
