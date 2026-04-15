#!/usr/bin/env bash

set -e

echo "🚀 Starting HybridSearch (dev mode)..."

# Ensure go modules are ready
echo "📦 Tidying modules..."
go mod tidy

# Build
echo "🔨 Building..."
go build -o hybridsearch ./cmd/hybridsearch

# Run
echo "▶️ Running..."
./hybridsearch