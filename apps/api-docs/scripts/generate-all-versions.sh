#!/bin/bash

# Simple MVP script to generate fixed versions from config
# Usage: ./scripts/generate-all-versions.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DOCS_DIR="$SCRIPT_DIR/.."
CONFIG_FILE="$API_DOCS_DIR/versions.config.simple.json"
ARTIFACTS_DIR="$API_DOCS_DIR/.artifacts"

echo "🚀 Generating fixed versions from config..."

# Check if config exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Config file not found: $CONFIG_FILE"
    exit 1
fi

# Clean up old versioned artifacts but keep current ones
echo "🧹 Cleaning old versioned artifacts..."
rm -rf "$ARTIFACTS_DIR/versions"
mkdir -p "$ARTIFACTS_DIR/versions"

cd "$API_DOCS_DIR"

# Save current branch
ORIGINAL_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "📍 Current branch: $ORIGINAL_BRANCH"

# Generate current OpenAPI specs first
echo "📝 Generating current OpenAPI specs..."
pnpm run ensure-plugins
buf generate ../../proto

# Read versions from config and generate each one
echo "📖 Reading versions from config..."

# Use jq to parse the JSON config
if ! command -v jq &> /dev/null; then
    echo "❌ jq is required but not installed. Please install jq."
    exit 1
fi

# Get all enabled versions
jq -r '.versions[] | select(.enabled == true) | .gitRef' "$CONFIG_FILE" | while read -r git_ref; do
    echo ""
    echo "🔄 Processing version: $git_ref"
    
    # Create version directory
    VERSION_DIR="$ARTIFACTS_DIR/versions/$git_ref"
    mkdir -p "$VERSION_DIR"
    
    if [ "$git_ref" = "main" ]; then
        # For main, just copy current artifacts
        echo "   📋 Using current artifacts for main..."
        if [ -d "$ARTIFACTS_DIR/openapi3" ]; then
            cp -r "$ARTIFACTS_DIR/openapi3" "$VERSION_DIR/"
        fi
        if [ -d "$ARTIFACTS_DIR/openapi" ]; then
            cp -r "$ARTIFACTS_DIR/openapi" "$VERSION_DIR/"
        fi
    else
        # For git tags, checkout, generate, and copy
        echo "   🔀 Checking out $git_ref..."
        cd ../../
        
        # Stash any local changes
        git stash push -m "Auto-stash before version generation" 2>/dev/null || true
        
        # Checkout the specific version
        if git checkout "$git_ref" 2>/dev/null; then
            echo "   ✅ Checked out $git_ref"
            
            # Generate artifacts for this version
            echo "   📝 Generating artifacts for $git_ref..."
            
            # Check if api-docs exists in this version
            if [ -d "apps/api-docs" ]; then
                cd apps/api-docs
                # Generate OpenAPI specs for this version
                pnpm run ensure-plugins
                buf generate ../../proto
                
                # Copy generated artifacts to version folder
                if [ -d "$ARTIFACTS_DIR/openapi3" ]; then
                    cp -r "$ARTIFACTS_DIR/openapi3" "$VERSION_DIR/"
                fi
                if [ -d "$ARTIFACTS_DIR/openapi" ]; then
                    cp -r "$ARTIFACTS_DIR/openapi" "$VERSION_DIR/"
                fi
            else
                echo "   📝 API docs not available in $git_ref, generating proto specs directly..."
                # Generate directly from proto without the api-docs setup
                if [ -f "proto/buf.yaml" ]; then
                    # Create a temporary directory for generation
                    TEMP_DIR="/tmp/zitadel-proto-gen-$git_ref"
                    mkdir -p "$TEMP_DIR"
                    
                    # Copy proto files
                    cp -r proto "$TEMP_DIR/"
                    cd "$TEMP_DIR"
                    
                    # Download the plugin if needed
                    if [ ! -f "protoc-gen-connect-openapi" ]; then
                        echo "   📥 Downloading protoc plugin..."
                        curl -L "https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v0.18.0/protoc-gen-connect-openapi_0.18.0_darwin_all.tar.gz" | tar -xz
                        chmod +x protoc-gen-connect-openapi
                    fi
                    
                    # Generate OpenAPI specs
                    export PATH="$PWD:$PATH"
                    buf generate proto
                    
                    # Copy results to version folder
                    if [ -d "openapi3" ]; then
                        cp -r openapi3 "$VERSION_DIR/"
                    fi
                    if [ -d "openapi" ]; then
                        cp -r openapi "$VERSION_DIR/"
                    fi
                    
                    # Clean up
                    rm -rf "$TEMP_DIR"
                else
                    echo "   ⚠️  No proto directory found in $git_ref"
                fi
                cd ../../apps/api-docs 2>/dev/null || cd ../../
            fi
            
            echo "   ✅ Generated artifacts for $git_ref"
        else
            echo "   ⚠️  Failed to checkout $git_ref, skipping..."
        fi
        
        # Go back to original branch
        cd ../../
        git checkout "$ORIGINAL_BRANCH" 2>/dev/null || git checkout apidocs
        
        # Restore stashed changes if any
        git stash pop 2>/dev/null || true
        
        cd apps/api-docs
    fi
    
    # Create metadata
    cat > "$VERSION_DIR/metadata.json" << EOF
{
  "version": "$git_ref",
  "gitRef": "$git_ref",
  "generatedAt": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "generatedBy": "scripts/generate-all-versions.sh"
}
EOF
    
    echo "   ✅ Created version: $git_ref"
done

# Generate current artifacts again (in case we're not on main)
echo ""
echo "📋 Regenerating current artifacts..."
pnpm run ensure-plugins
buf generate ../../proto

echo ""
echo "✅ All versions generated successfully!"
echo "📁 Available versions:"
ls -1 "$ARTIFACTS_DIR/versions/" 2>/dev/null || echo "No versions found"
echo ""
echo "🎯 Ready to use with: pnpm dev"
