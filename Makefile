# Vietnam Administrative API Makefile

.PHONY: help setup build run test clean docker docker-run deploy

# Default target
help:
	@echo "🇻🇳 Vietnam Administrative API"
	@echo ""
	@echo "Available commands:"
	@echo "  setup       - Initialize project and download dependencies"
	@echo "  build       - Build the application binary"
	@echo "  run         - Run the application in development mode"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker      - Build Docker image"
	@echo "  docker-run  - Run application in Docker container"
	@echo "  deploy      - Deploy with docker-compose"
	@echo ""

# Setup project
setup:
	@echo "🚀 Setting up project..."
	@mkdir -p data
	@if [ -f "province.json" ]; then cp province.json data/; fi
	@if [ -f "ward.json" ]; then cp ward.json data/; fi
	@go mod tidy
	@echo "✅ Setup completed!"

# Build application
build: setup
	@echo "🔨 Building application..."
	@go build -ldflags="-s -w" -o vietnam-admin-api .
	@echo "✅ Build completed: vietnam-admin-api"

# Run in development mode
run: setup
	@echo "🏃 Running in development mode..."
	@GIN_MODE=debug go run main.go

# Run production build
run-prod: build
	@echo "🚀 Running production binary..."
	@GIN_MODE=release ./vietnam-admin-api

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -cover ./...

# Benchmark tests
benchmark:
	@echo "⚡ Running benchmarks..."
	@go test -bench=. ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f vietnam-admin-api
	@go clean

# Build Docker image
docker:
	@echo "🐳 Building Docker image..."
	@docker build -t vietnam-admin-api:latest .
	@echo "✅ Docker image built: vietnam-admin-api:latest"

# Run Docker container
docker-run: docker
	@echo "🐳 Running Docker container..."
	@docker run -d \
		--name vietnam-api \
		-p 8100:8100 \
		-v $(PWD)/province.json:/root/data/province.json:ro \
		-v $(PWD)/ward.json:/root/data/ward.json:ro \
		vietnam-admin-api:latest
	@echo "✅ Container started: http://localhost:8100"

# Stop Docker container
docker-stop:
	@echo "🛑 Stopping Docker container..."
	@docker stop vietnam-api || true
	@docker rm vietnam-api || true

# Deploy with docker-compose
deploy:
	@echo "🚀 Deploying with docker-compose..."
	@docker-compose up -d --build
	@echo "✅ Deployed! Check: http://localhost:8100"

# Stop deployment
deploy-stop:
	@echo "🛑 Stopping deployment..."
	@docker-compose down

# View logs
logs:
	@docker-compose logs -f vietnam-api

# Check API health
health:
	@echo "🏥 Checking API health..."
	@curl -s http://localhost:8100/health | jq '.' || echo "API not responding"

# Load test (requires hey: go install github.com/rakyll/hey@latest)
load-test:
	@echo "⚡ Running load test..."
	@hey -n 1000 -c 10 http://localhost:8100/api/v1/provinces

# Development workflow
dev: clean setup run

# Production workflow
prod: clean build run-prod

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

# Security scan (requires gosec)
security:
	@echo "🔒 Running security scan..."
	@gosec ./... || echo "Install gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"

# Generate API documentation (requires swag)
docs:
	@echo "📚 Generating API documentation..."
	@swag init || echo "Install swag: go install github.com/swaggo/swag/cmd/swag@latest"

# Full CI pipeline
ci: fmt lint test build

# Install development tools
install-tools:
	@echo "🛠️ Installing development tools..."
	@go install github.com/rakyll/hey@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Tools installed!" 