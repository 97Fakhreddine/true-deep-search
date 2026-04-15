#!/usr/bin/env bash

set -e

echo "🧪 Running tests..."

go test ./... -v

echo "✅ Tests done"