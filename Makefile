.PHONY: build clean install test run help

# Build the application
build:
	go build -o yeetrap.exe

# Build for multiple platforms
build-all:
	GOOS=windows GOARCH=amd64 go build -o bin/yeetrap-windows-amd64.exe
	GOOS=linux GOARCH=amd64 go build -o bin/yeetrap-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/yeetrap-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/yeetrap-darwin-arm64

# Clean build artifacts
clean:
	rm -f yeetrap yeetrap.exe
	rm -rf bin/
	rm -rf downloads/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Run the application
run:
	go run main.go

# Display help
help:
	@echo "YeeTrap - Build Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  build      - Build the application for current platform"
	@echo "  build-all  - Build for Windows, Linux, and macOS"
	@echo "  clean      - Remove build artifacts"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  test       - Run tests"
	@echo "  run        - Run the application"
	@echo "  help       - Show this help message"


