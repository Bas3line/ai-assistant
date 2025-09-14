#!/bin/bash

# Build script for AI Assistant API
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="ai-assistant"
BUILD_DIR="./bin"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

echo -e "${YELLOW}Building AI Assistant API...${NC}"
echo "Version: ${VERSION}"
echo "Build Time: ${BUILD_TIME}"
echo "Git Commit: ${GIT_COMMIT}"

# Create build directory
mkdir -p ${BUILD_DIR}

# Build for different platforms
build_for_platform() {
    local GOOS=$1
    local GOARCH=$2
    local output_name="${BINARY_NAME}"
    
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo -e "${YELLOW}Building for ${GOOS}/${GOARCH}...${NC}"
    
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags "${LDFLAGS}" \
        -o "${BUILD_DIR}/${GOOS}-${GOARCH}/${output_name}" \
        ./cmd/api
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built ${GOOS}/${GOARCH} successfully${NC}"
    else
        echo -e "${RED}✗ Failed to build ${GOOS}/${GOARCH}${NC}"
        exit 1
    fi
}

# Default build (current platform)
echo -e "${YELLOW}Building for current platform...${NC}"
go build -ldflags "${LDFLAGS}" -o "${BUILD_DIR}/${BINARY_NAME}" ./cmd/api

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Built successfully${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Cross-platform builds (optional)
if [ "$1" = "all" ]; then
    echo -e "${YELLOW}Building for all platforms...${NC}"
    build_for_platform "linux" "amd64"
    build_for_platform "linux" "arm64"
    build_for_platform "darwin" "amd64"
    build_for_platform "darwin" "arm64"
    build_for_platform "windows" "amd64"
fi

echo -e "${GREEN}Build completed!${NC}"
echo "Binary location: ${BUILD_DIR}/${BINARY_NAME}"