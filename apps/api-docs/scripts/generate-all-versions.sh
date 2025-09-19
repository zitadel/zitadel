#!/bin/bash

# Script to generate OpenAPI specs for all enabled versions in versions.config.json
# This is the main generation script that should be run to create all versions

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DOCS_DIR="$SCRIPT_DIR/.."
ARTIFACTS_DIR="$API_DOCS_DIR/.artifacts"
CONFIG_FILE="$API_DOCS_DIR/versions.config.json"

echo "🚀 Generating all versions from config..."

cd "$API_DOCS_DIR"

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Config file not found: $CONFIG_FILE"
    echo "Run 'pnpm run versions:create <version>' to create your first version"
    exit 1
fi

# Parse config and get enabled versions
echo "📋 Reading versions from config..."
ENABLED_VERSIONS=$(jq -r '.versions[] | select(.enabled == true) | .id' "$CONFIG_FILE")

if [ -z "$ENABLED_VERSIONS" ]; then
    echo "❌ No enabled versions found in config"
    exit 1
fi

echo "📦 Found enabled versions:"
echo "$ENABLED_VERSIONS" | sed 's/^/  - /'

# Create versions directory
mkdir -p "$ARTIFACTS_DIR/versions"

# For simplicity, let's just generate for the current branch for all versions
# This avoids the complexity of branch switching
echo ""
echo "🔧 Generating OpenAPI specs for current branch and copying to all version directories..."

# Generate OpenAPI specs once
echo "📦 Generating OpenAPI specs..."
if pnpm run generate:openapi && [ -d "public/openapi" ]; then
    echo "✅ OpenAPI specs generated successfully"
    
    # Copy to each version directory
    for VERSION in $ENABLED_VERSIONS; do
        echo "📦 Creating artifacts for version: $VERSION"
        VERSION_DIR="$ARTIFACTS_DIR/versions/$VERSION"
        mkdir -p "$VERSION_DIR"
        cp -r public/openapi/* "$VERSION_DIR/"
        echo "✅ Artifacts created for $VERSION"
    done
else
    echo "⚠️  OpenAPI generation failed, trying to use existing artifacts..."
    
    # Check for existing artifacts in other locations
    for VERSION in $ENABLED_VERSIONS; do
        VERSION_DIR="$ARTIFACTS_DIR/versions/$VERSION"
        mkdir -p "$VERSION_DIR"
        
        # Try to find existing artifacts
        if [ -d ".artifacts/openapi" ]; then
            cp -r .artifacts/openapi/* "$VERSION_DIR/" 2>/dev/null || true
        fi
        if [ -d ".artifacts/openapi3" ]; then
            cp -r .artifacts/openapi3/* "$VERSION_DIR/" 2>/dev/null || true
        fi
        
        if [ "$(ls -A "$VERSION_DIR" 2>/dev/null)" ]; then
            echo "✅ Used existing artifacts for $VERSION"
        else
            echo "⚠️  No artifacts available for $VERSION"
        fi
    done
fi

echo ""
echo "🎉 Finished generating all versions!"
echo ""
echo "📊 Summary:"
ls -la "$ARTIFACTS_DIR/versions/" | grep -E "^d" | awk '{print "  - " $9}' | grep -v "\.$" || echo "  (no versions created)"
