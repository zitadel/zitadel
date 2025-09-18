#!/bin/bash

# Utility script to manage version artifacts in organized folders
# Usage: ./manage-versions.sh [create|list|delete] [version]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARTIFACTS_DIR="$SCRIPT_DIR/../.artifacts"
VERSIONS_DIR="$ARTIFACTS_DIR/versions"

usage() {
    echo "Usage: $0 <command> [version]"
    echo ""
    echo "Commands:"
    echo "  create <version>    Create a new version from current artifacts"
    echo "  list                List all available versions"
    echo "  delete <version>    Delete a specific version"
    echo "  current             Show current artifacts info"
    echo ""
    echo "Examples:"
    echo "  $0 create v2.71.9"
    echo "  $0 list"
    echo "  $0 delete v2.70.0"
}

create_version() {
    local version="$1"
    
    if [ -z "$version" ]; then
        echo "Error: Version name required"
        usage
        exit 1
    fi
    
    echo "Creating version: $version"
    
    # Check if current artifacts exist
    if [ ! -d "$ARTIFACTS_DIR/openapi3" ]; then
        echo "Error: No current artifacts found. Run 'pnpm run generate' first."
        exit 1
    fi
    
    # Create version directory
    local version_dir="$VERSIONS_DIR/$version"
    mkdir -p "$version_dir"
    
    # Copy artifacts
    echo "Copying OpenAPI 3.x specs..."
    cp -r "$ARTIFACTS_DIR/openapi3" "$version_dir/"
    
    if [ -d "$ARTIFACTS_DIR/openapi" ]; then
        echo "Copying OpenAPI 2.x specs..."
        cp -r "$ARTIFACTS_DIR/openapi" "$version_dir/"
    fi
    
    # Create metadata
    cat > "$version_dir/metadata.json" << EOF
{
  "version": "$version",
  "generatedAt": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "gitCommit": "$(git rev-parse HEAD 2>/dev/null || echo 'unknown')",
  "gitBranch": "$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')",
  "gitTag": "$(git describe --tags --exact-match 2>/dev/null || echo 'unknown')"
}
EOF
    
    echo "✓ Version $version created successfully"
    echo "Location: $version_dir"
}

list_versions() {
    echo "Available versions:"
    echo ""
    
    # Current/latest
    if [ -d "$ARTIFACTS_DIR/openapi3" ]; then
        echo "📁 latest (current artifacts)"
    fi
    
    # Organized versions
    if [ -d "$VERSIONS_DIR" ]; then
        for version_dir in "$VERSIONS_DIR"/*; do
            if [ -d "$version_dir" ]; then
                local version=$(basename "$version_dir")
                local metadata_file="$version_dir/metadata.json"
                
                if [ -f "$metadata_file" ]; then
                    local generated_at=$(jq -r '.generatedAt // "unknown"' "$metadata_file" 2>/dev/null || echo "unknown")
                    local git_branch=$(jq -r '.gitBranch // "unknown"' "$metadata_file" 2>/dev/null || echo "unknown")
                    echo "📁 $version (generated: $generated_at, branch: $git_branch)"
                else
                    echo "📁 $version (no metadata)"
                fi
            fi
        done
    fi
    
    # Legacy versions
    local legacy_dir="$ARTIFACTS_DIR/../.artifacts-versioned"
    if [ -d "$legacy_dir" ]; then
        echo ""
        echo "Legacy versions (consider migrating):"
        for legacy_version in "$legacy_dir"/*; do
            if [ -d "$legacy_version" ]; then
                echo "📦 $(basename "$legacy_version") (legacy)"
            fi
        done
    fi
}

delete_version() {
    local version="$1"
    
    if [ -z "$version" ]; then
        echo "Error: Version name required"
        usage
        exit 1
    fi
    
    local version_dir="$VERSIONS_DIR/$version"
    
    if [ ! -d "$version_dir" ]; then
        echo "Error: Version $version not found"
        exit 1
    fi
    
    echo "Are you sure you want to delete version $version? (y/N)"
    read -r confirm
    
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        rm -rf "$version_dir"
        echo "✓ Version $version deleted"
    else
        echo "Cancelled"
    fi
}

show_current() {
    echo "Current artifacts information:"
    echo ""
    
    if [ -d "$ARTIFACTS_DIR/openapi3" ]; then
        local spec_count=$(find "$ARTIFACTS_DIR/openapi3" -name "*.yaml" | wc -l)
        echo "OpenAPI 3.x specs: $spec_count files"
        echo "Last modified: $(stat -f %Sm "$ARTIFACTS_DIR/openapi3" 2>/dev/null || echo 'unknown')"
    else
        echo "No current artifacts found"
    fi
    
    echo "Git info:"
    echo "  Branch: $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')"
    echo "  Commit: $(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
    echo "  Tag: $(git describe --tags --exact-match 2>/dev/null || echo 'none')"
}

# Main script logic
case "${1:-}" in
    create)
        create_version "$2"
        ;;
    list)
        list_versions
        ;;
    delete)
        delete_version "$2"
        ;;
    current)
        show_current
        ;;
    *)
        usage
        exit 1
        ;;
esac
