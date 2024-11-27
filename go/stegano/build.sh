#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o ./linux/stegano main.go
GOOS=windows GOARCH=amd64 go build -o ./windows/stegano.exe main.go
GOOS=darwin GOARCH=amd64 go build -o ./macos/amd64/stegano main.go
GOOS=darwin GOARCH=arm64 go build -o ./macos/arm64/stegano main.go

