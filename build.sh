#!/bin/bash

# Create build directory if it doesn't exist
mkdir -p build

# Build for Windows
echo "Building for Windows..."
cd src
GOOS=windows GOARCH=amd64 go build -o ../build/tet-simulation-windows-amd64.exe

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o ../build/tet-simulation-linux-amd64

# Build for macOS (Intel)
echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o ../build/tet-simulation-darwin-amd64

# Build for macOS (Apple Silicon)
echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o ../build/tet-simulation-darwin-arm64
cd ..

# Make Unix executables executable
chmod +x build/tet-simulation-linux-amd64
chmod +x build/tet-simulation-darwin-amd64
chmod +x build/tet-simulation-darwin-arm64

echo "Build complete! Executables are in the build directory:"
ls -l build/ 