#!/bin/bash

# Ensure output directory exists
mkdir -p pkg/pb

# Clean up old files
find . -name "*.pb.go" -delete

# Generate Go code and gRPC service code
# Using module mapping to ensure it lands in pkg/pb/ as per go_package option
protoc --go_out=. --go_opt=module=github.com/user/mangahub \
    --go-grpc_out=. --go-grpc_opt=module=github.com/user/mangahub \
    api/proto/mangahub.proto

echo "✅ gRPC code generated in pkg/pb/"
