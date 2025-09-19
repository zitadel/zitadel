#!/bin/bash

# Simple script to generate OpenAPI spec for current branch only
# No branch switching - just generates for whatever branch you're on

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DOCS_DIR="$SCRIPT_DIR/.."
ARTIFACTS_DIR="$API_DOCS_DIR/.artifacts"

echo "🚀 Generating OpenAPI spec for current branch..."

cd "$API_DOCS_DIR"

# Get current branch/version
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "📍 Current branch: $CURRENT_BRANCH"

# Create artifacts directory
mkdir -p "$ARTIFACTS_DIR/versions/$CURRENT_BRANCH"

# Generate OpenAPI spec for current branch only
echo "🔧 Generating OpenAPI spec..."
pnpm run generate:openapi

# Copy generated files to artifacts
if [ -d "public/openapi" ]; then
    echo "📦 Copying OpenAPI specs to artifacts..."
    cp -r public/openapi/* "$ARTIFACTS_DIR/versions/$CURRENT_BRANCH/"
    echo "✅ Generated OpenAPI spec for $CURRENT_BRANCH"
else
    echo "❌ No OpenAPI specs found in public/openapi"
    exit 1
fi

echo "🎉 Done! OpenAPI spec generated for $CURRENT_BRANCH"
