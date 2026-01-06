#!/bin/bash
# Build script for go-webapi-db
# This script ensures dependencies are downloaded before building

set -e

echo "Downloading dependencies..."
go mod download

echo "Tidying modules..."
go mod tidy

echo "Building application..."
go build -o bin/server ./cmd/server

echo "Build successful! Binary created at bin/server"

