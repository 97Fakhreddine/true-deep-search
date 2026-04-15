#!/usr/bin/env bash

set -e

echo "🔍 Running lint checks..."

# Format
echo "🎨 Formatting..."
go fmt ./...

# Vet
echo "🧠 Running go vet..."
go vet ./...

# Optional: staticcheck (if installed)
if command -v staticcheck &> /dev/null
then
    echo "⚡ Running staticcheck..."
    staticcheck ./...
else
    echo "⚠️ staticcheck not installed (skip)"
fi

echo "✅ Lint passed"