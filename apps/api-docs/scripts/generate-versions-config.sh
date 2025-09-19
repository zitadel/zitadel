#!/bin/bash

# Script to auto-generate versions.config.json based on available versions and git tags
# Usage: ./scripts/generate-versions-config.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/../versions.config.json"
VERSIONS_DIR="$SCRIPT_DIR/../.artifacts/versions"

echo "🔧 Auto-generating versions.config.json..."

# Get available versions from filesystem
AVAILABLE_VERSIONS=()
if [ -d "$VERSIONS_DIR" ]; then
    for version_dir in "$VERSIONS_DIR"/*; do
        if [ -d "$version_dir" ]; then
            version=$(basename "$version_dir")
            AVAILABLE_VERSIONS+=("$version")
        fi
    done
fi

# Get recent git tags (last 10 stable versions)
GIT_TAGS=()
if git rev-parse --git-dir > /dev/null 2>&1; then
    while IFS= read -r tag; do
        GIT_TAGS+=("$tag")
    done < <(git tag --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -10)
fi

# Combine and deduplicate versions using a simpler approach
ALL_VERSIONS=()

# Function to check if version is already in array
version_exists() {
    local version="$1"
    for existing in "${ALL_VERSIONS[@]}"; do
        if [[ "$existing" == "$version" ]]; then
            return 0
        fi
    done
    return 1
}

# Add available versions first (they have artifacts)
for version in "${AVAILABLE_VERSIONS[@]}"; do
    if ! version_exists "$version"; then
        ALL_VERSIONS+=("$version")
    fi
done

# Add recent git tags (might not have artifacts yet)
for version in "${GIT_TAGS[@]}"; do
    if ! version_exists "$version"; then
        ALL_VERSIONS+=("$version")
    fi
done

# Always include main/latest
if ! version_exists "main"; then
    ALL_VERSIONS+=("main")
fi

echo "📋 Found versions: ${ALL_VERSIONS[*]}"

# Determine default version (latest stable tag with artifacts, or latest stable tag, or main)
DEFAULT_VERSION="main"
for version in "${AVAILABLE_VERSIONS[@]}"; do
    if [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        DEFAULT_VERSION="$version"
        break
    fi
done

if [[ "$DEFAULT_VERSION" == "main" ]]; then
    for version in "${GIT_TAGS[@]}"; do
        DEFAULT_VERSION="$version"
        break
    done
fi

echo "🎯 Default version: $DEFAULT_VERSION"

# Generate JSON config
cat > "$CONFIG_FILE" << EOF
{
  "versions": [
EOF

# Add each version
for i in "${!ALL_VERSIONS[@]}"; do
    version="${ALL_VERSIONS[$i]}"
    
    # Determine properties
    if [[ "$version" == "main" ]]; then
        name="Latest (Main Branch)"
        git_ref="main"
        is_stable="false"
    elif [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        name="$version"
        git_ref="$version"
        is_stable="true"
    else
        name="$version"
        git_ref="$version"
        is_stable="false"
    fi
    
    # Check if artifacts exist
    version_has_artifacts="false"
    for available_version in "${AVAILABLE_VERSIONS[@]}"; do
        if [[ "$available_version" == "$version" ]]; then
            version_has_artifacts="true"
            break
        fi
    done
    
    if [[ "$version_has_artifacts" == "true" ]]; then
        enabled="true"
    else
        enabled="true"  # Still enable it, will be marked as unavailable in UI
    fi
    
    # Add comma except for last item
    comma=","
    if [[ $i -eq $((${#ALL_VERSIONS[@]} - 1)) ]]; then
        comma=""
    fi
    
    cat >> "$CONFIG_FILE" << EOF
    {
      "id": "$version",
      "name": "$name",
      "gitRef": "$git_ref",
      "enabled": $enabled,
      "isStable": $is_stable
    }$comma
EOF
done

# Add settings section
cat >> "$CONFIG_FILE" << EOF
  ],
  "settings": {
    "defaultVersion": "$DEFAULT_VERSION",
    "autoGenerate": true,
    "maxVersions": 15,
    "includePrerelease": false
  },
  "metadata": {
    "generatedAt": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
    "generatedBy": "scripts/generate-versions-config.sh",
    "availableVersionsCount": ${#AVAILABLE_VERSIONS[@]},
    "totalVersionsCount": ${#ALL_VERSIONS[@]}
  }
}
EOF

echo "✅ Generated $CONFIG_FILE with ${#ALL_VERSIONS[@]} versions"
echo "📊 Breakdown:"
echo "   - Available (with artifacts): ${#AVAILABLE_VERSIONS[@]}"
echo "   - Total configured: ${#ALL_VERSIONS[@]}"
echo "   - Default: $DEFAULT_VERSION"

# Show the generated config
echo ""
echo "📄 Generated configuration preview:"
head -20 "$CONFIG_FILE"
echo "..."
tail -10 "$CONFIG_FILE"
