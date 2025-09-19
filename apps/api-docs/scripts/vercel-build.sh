#!/bin/bash

# Simplified version generation for Vercel builds
# This script generates only the current version and relies on config for version list

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARTIFACTS_DIR="$SCRIPT_DIR/../.artifacts"

echo "🚀 Building API docs with version support..."

# Always generate current artifacts first
echo "🔨 Generating current artifacts..."
cd "$SCRIPT_DIR/.."
pnpm run ensure-plugins
pnpm run generate

# Check if we're building from a git tag
GIT_TAG=""
if git describe --exact-match --tags HEAD 2>/dev/null; then
    GIT_TAG=$(git describe --exact-match --tags HEAD)
    echo "📦 Building from git tag: $GIT_TAG"
    
    # Create version snapshot for this tag
    VERSIONS_DIR="$ARTIFACTS_DIR/versions"
    TAG_DIR="$VERSIONS_DIR/$GIT_TAG"
    
    echo "📁 Creating version snapshot: $GIT_TAG"
    mkdir -p "$TAG_DIR"
    
    # Copy artifacts
    if [ -d "$ARTIFACTS_DIR/openapi3" ]; then
        cp -r "$ARTIFACTS_DIR/openapi3" "$TAG_DIR/"
    fi
    if [ -d "$ARTIFACTS_DIR/openapi" ]; then
        cp -r "$ARTIFACTS_DIR/openapi" "$TAG_DIR/"
    fi
    
    # Create metadata
    cat > "$TAG_DIR/metadata.json" << EOF
{
  "version": "$GIT_TAG",
  "gitRef": "$GIT_TAG",
  "generatedAt": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "gitCommit": "$(git rev-parse HEAD)",
  "gitBranch": "$(git rev-parse --abbrev-ref HEAD)",
  "isTagBuild": true
}
EOF
    
    echo "✅ Created version snapshot for $GIT_TAG"
else
    echo "📝 Building from branch (no tag detected)"
fi

echo "✅ Build preparation complete!"
