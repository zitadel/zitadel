#!/bin/bash

# Auto-version script for Vercel builds
# This script runs during Vercel build and creates a version if building from a git tag

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "🔍 Checking if we're building from a git tag..."

# Check if we're building from a tag
if [ -n "$VERCEL_GIT_COMMIT_REF" ]; then
    # Vercel provides the git ref (branch/tag name)
    GIT_REF="$VERCEL_GIT_COMMIT_REF"
    echo "📍 Vercel Git Ref: $GIT_REF"
    
    # Check if it's a version tag (starts with 'v' followed by numbers)
    if [[ "$GIT_REF" =~ ^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
        echo "🏷️  Detected version tag: $GIT_REF"
        
        # Generate artifacts first
        echo "🔨 Generating API artifacts..."
        pnpm run generate
        
        # Check if version already exists
        if [ -d ".artifacts/versions/$GIT_REF" ]; then
            echo "⚠️  Version $GIT_REF already exists, skipping creation"
        else
            # Create version
            echo "📦 Creating version: $GIT_REF"
            ./scripts/manage-versions.sh create "$GIT_REF"
            echo "✅ Version $GIT_REF created successfully"
        fi
    else
        echo "📝 Not a version tag, skipping auto-versioning"
    fi
else
    echo "🤷 No git ref available, skipping auto-versioning"
fi

echo "🎯 Auto-versioning complete"
