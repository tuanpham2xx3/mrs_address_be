#!/bin/bash

# Vietnam Admin API Setup Script
echo "🚀 Setting up Vietnam Administrative API..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Initialize Go module if not exists
if [ ! -f "go.mod" ]; then
    echo "📦 Initializing Go module..."
    go mod init vietnam-admin-api
fi

# Create data directory
echo "📁 Creating data directory..."
mkdir -p data

# Check if JSON files exist
if [ ! -f "province.json" ]; then
    echo "⚠️  Warning: province.json not found. Please copy it to the project root."
fi

if [ ! -f "ward.json" ]; then
    echo "⚠️  Warning: ward.json not found. Please copy it to the project root."
fi

# Copy JSON files to data directory if they exist
if [ -f "province.json" ]; then
    echo "📋 Copying province.json to data directory..."
    cp province.json data/
fi

if [ -f "ward.json" ]; then
    echo "📋 Copying ward.json to data directory..."
    cp ward.json data/
fi

# Download dependencies
echo "📥 Downloading Go dependencies..."
go mod tidy

# Build the application
echo "🔨 Building application..."
if go build -o vietnam-admin-api .; then
    echo "✅ Build successful!"
    echo ""
    echo "🎉 Setup completed!"
    echo ""
    echo "To run the application:"
    echo "  ./vietnam-admin-api"
    echo ""
    echo "Or run in development mode:"
    echo "  go run main.go"
    echo ""
    echo "API will be available at: http://localhost:8080"
    echo "Health check: http://localhost:8080/health"
else
    echo "❌ Build failed. Please check the errors above."
    exit 1
fi 