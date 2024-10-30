#!/bin/bash

# Create output directory
#mkdir -p bin

# Build for each platform
#echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o ./linux/docker-compose-converter main.go

#echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o ./windows/docker-compose-converter.exe main.go

#echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o ./macos/docker-compose-converter main.go

#echo "Builds completed. Check the 'bin' directory for executables."
