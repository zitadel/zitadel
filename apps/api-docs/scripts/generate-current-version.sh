#!/bin/bash

# Enhanced script to manage versions and optionally generate OpenAPI specs
# Usage: ./scripts/generate-current-version.sh [version_id]
# If version_id is provided, it will be added to versions.config.json

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DOCS_DIR="$SCRIPT_DIR/.."
ARTIFACTS_DIR="$API_DOCS_DIR/.artifacts"
CONFIG_FILE="$API_DOCS_DIR/versions.config.json"

# Get version parameter (optional)
VERSION_ID="$1"

echo "🚀 Managing API versions..."

cd "$API_DOCS_DIR"

# Get current branch/version
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "📍 Current branch: $CURRENT_BRANCH"

# Use provided version ID or default to current branch
FINAL_VERSION="${VERSION_ID:-$CURRENT_BRANCH}"
echo "📦 Version ID: $FINAL_VERSION"

# Create artifacts directory
mkdir -p "$ARTIFACTS_DIR/versions/$FINAL_VERSION"

# Try to generate OpenAPI spec for current branch
echo "🔧 Attempting to generate OpenAPI spec..."
if pnpm run generate:openapi && [ -d "public/openapi" ]; then
    echo "📦 Copying OpenAPI specs to artifacts..."
    cp -r public/openapi/* "$ARTIFACTS_DIR/versions/$FINAL_VERSION/"
    echo "✅ Generated OpenAPI spec for $FINAL_VERSION"
else
    echo "⚠️  OpenAPI generation failed or no specs found, using existing artifacts if available"
    
    # Check if we have existing artifacts for this version
    if [ -d "$ARTIFACTS_DIR/versions/$FINAL_VERSION" ] && [ "$(ls -A "$ARTIFACTS_DIR/versions/$FINAL_VERSION" 2>/dev/null)" ]; then
        echo "📦 Using existing artifacts for $FINAL_VERSION"
    else
        echo "❌ No existing artifacts found for $FINAL_VERSION"
        # Don't exit - still allow config updates
    fi
fi

# If version ID was provided, add it to config
if [ -n "$VERSION_ID" ]; then
    echo "🔧 Adding version $VERSION_ID to config..."
    
    # Check if config file exists
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "❌ Config file not found: $CONFIG_FILE"
        exit 1
    fi
    
    # Check if version already exists in config
    if jq -e ".versions[] | select(.id == \"$VERSION_ID\")" "$CONFIG_FILE" > /dev/null 2>&1; then
        echo "⚠️  Version $VERSION_ID already exists in config"
    else
        # Add new version to config
        echo "➕ Adding new version $VERSION_ID to config..."
        
        # Determine if it's a stable version (starts with 'v' and contains dots)
        IS_STABLE="true"
        if [[ "$VERSION_ID" == "main" ]] || [[ "$VERSION_ID" == "develop" ]] || [[ "$VERSION_ID" == *"alpha"* ]] || [[ "$VERSION_ID" == *"beta"* ]]; then
            IS_STABLE="false"
        fi
        
        # Create new version object
        NEW_VERSION=$(jq -n \
            --arg id "$VERSION_ID" \
            --arg name "$VERSION_ID" \
            --arg gitRef "$VERSION_ID" \
            --argjson enabled true \
            --argjson isStable "$IS_STABLE" \
            '{
                id: $id,
                name: $name,
                gitRef: $gitRef,
                enabled: $enabled,
                isStable: $isStable
            }')
        
        # Add to versions array and update config
        jq ".versions += [$NEW_VERSION]" "$CONFIG_FILE" > "${CONFIG_FILE}.tmp" && mv "${CONFIG_FILE}.tmp" "$CONFIG_FILE"
        
        echo "✅ Added version $VERSION_ID to config"
    fi
fi

echo "🎉 Done! Version management completed for $FINAL_VERSION"
