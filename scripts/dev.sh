#!/usr/bin/env bash

set -e

echo "🚀 Starting HybridSearch (dev mode)..."

cd "$(dirname "$0")/.."

echo "📦 Tidying modules..."
go mod tidy

echo "🏃 Running app..."
go run ./cmd/hybridsearch