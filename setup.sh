#!/bin/bash

# YeeTrap Setup Script for macOS/Linux
# This script helps set up YeeTrap for first-time use

echo "========================================"
echo "   YeeTrap Setup Script"
echo "========================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Check if Go is installed
echo -e "${YELLOW}Checking for Go installation...${NC}"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo -e "${GREEN}✓ Go is installed: $GO_VERSION${NC}"
else
    echo -e "${RED}❌ Go is not installed!${NC}"
    echo -e "${RED}Please install Go from https://golang.org/dl/${NC}"
    exit 1
fi

# Check if yt-dlp is installed
echo -e "${YELLOW}Checking for yt-dlp installation...${NC}"
if command -v yt-dlp &> /dev/null; then
    YTDLP_VERSION=$(yt-dlp --version)
    echo -e "${GREEN}✓ yt-dlp is installed: v$YTDLP_VERSION${NC}"
else
    echo -e "${RED}❌ yt-dlp is not installed!${NC}"
    echo -e "${YELLOW}Please install yt-dlp:${NC}"
    
    # Detect OS and suggest installation method
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "  macOS: brew install yt-dlp"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "  Linux: sudo apt install yt-dlp"
        echo "     or: pip install yt-dlp"
    fi
    echo ""
fi

# Create config directory
echo ""
echo -e "${YELLOW}Creating configuration directory...${NC}"
CONFIG_DIR="$HOME/.yeetrap"
if [ ! -d "$CONFIG_DIR" ]; then
    mkdir -p "$CONFIG_DIR"
    chmod 700 "$CONFIG_DIR"
    echo -e "${GREEN}✓ Created directory: $CONFIG_DIR${NC}"
else
    echo -e "${GREEN}✓ Configuration directory already exists${NC}"
fi

# Check for credentials file
echo ""
echo -e "${YELLOW}Checking for credentials file...${NC}"
CREDENTIALS_PATH="$CONFIG_DIR/credentials.json"
if [ ! -f "$CREDENTIALS_PATH" ]; then
    echo -e "${RED}❌ credentials.json not found!${NC}"
    echo ""
    echo -e "${YELLOW}Please follow these steps:${NC}"
    echo "1. Go to https://console.cloud.google.com/"
    echo "2. Create a new project or select existing one"
    echo "3. Enable YouTube Data API v3"
    echo "4. Create OAuth2 credentials (Desktop app)"
    echo "5. Download the credentials JSON file"
    echo "6. Save it to: $CREDENTIALS_PATH"
    echo ""
    echo -e "${CYAN}See credentials_example.txt for the expected file format${NC}"
else
    echo -e "${GREEN}✓ credentials.json found${NC}"
fi

# Build the application
echo ""
echo -e "${YELLOW}Building YeeTrap...${NC}"
go build -o yeetrap
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ YeeTrap built successfully${NC}"
    chmod +x yeetrap
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Test the executable
echo ""
echo -e "${YELLOW}Testing YeeTrap...${NC}"
./yeetrap version
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ YeeTrap is working correctly${NC}"
else
    echo -e "${RED}❌ YeeTrap test failed${NC}"
    exit 1
fi

# Summary
echo ""
echo "========================================"
echo "   Setup Complete!"
echo "========================================"
echo ""

if [ -f "$CREDENTIALS_PATH" ]; then
    echo -e "${GREEN}Next steps:${NC}"
    echo "1. Run: ./yeetrap auth"
    echo "2. Run: ./yeetrap list"
    echo "3. Run: ./yeetrap download"
else
    echo -e "${YELLOW}Next steps:${NC}"
    echo "1. Add your credentials.json to: $CREDENTIALS_PATH"
    echo "2. Run: ./yeetrap auth"
    echo "3. Run: ./yeetrap list"
    echo "4. Run: ./yeetrap download"
fi

echo ""
echo -e "${CYAN}For help: ./yeetrap --help${NC}"
echo -e "${CYAN}Documentation: README.md and QUICKSTART.md${NC}"
echo ""


