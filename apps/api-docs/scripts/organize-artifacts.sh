#!/bin/bash

# Script to organize artifacts into version-specific folders
# This creates a cleaner structure for version management

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARTIFACTS_DIR="$SCRIPT_DIR/../.artifacts"

# Get current version from git or default to "latest"
CURRENT_VERSION=${1:-"latest"}

echo "Organizing artifacts for version: $CURRENT_VERSION"

# Create version-specific directory structure
VERSION_DIR="$ARTIFACTS_DIR/versions/$CURRENT_VERSION"
mkdir -p "$VERSION_DIR"

# Copy current artifacts to versioned directory
if [ -d "$ARTIFACTS_DIR/openapi3" ]; then
    echo "Copying OpenAPI 3.x specs..."
    cp -r "$ARTIFACTS_DIR/openapi3" "$VERSION_DIR/"
fi

if [ -d "$ARTIFACTS_DIR/openapi" ]; then
    echo "Copying OpenAPI 2.x specs..."
    cp -r "$ARTIFACTS_DIR/openapi" "$VERSION_DIR/"
fi

# Create a metadata file for this version
cat > "$VERSION_DIR/metadata.json" << EOF
{
  "version": "$CURRENT_VERSION",
  "generatedAt": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "gitCommit": "$(git rev-parse HEAD 2>/dev/null || echo 'unknown')",
  "gitBranch": "$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')"
}
EOF

echo "✓ Artifacts organized in: $VERSION_DIR"
echo "Available versions:"
ls -1 "$ARTIFACTS_DIR/versions/" 2>/dev/null || echo "No versions found"
