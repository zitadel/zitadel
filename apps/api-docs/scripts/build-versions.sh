#!/bin/bash

# Script to generate API artifacts for specified versions during build
# This runs during Vercel build and generates versions that don't exist yet

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/../versions.config.json"
ARTIFACTS_DIR="$SCRIPT_DIR/../.artifacts"
VERSIONS_DIR="$ARTIFACTS_DIR/versions"

echo "🚀 Starting version generation..."

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ No versions.config.json found. Creating default config..."
    exit 1
fi

# Parse config and get enabled versions
ENABLED_VERSIONS=$(node -e "
const config = require('$CONFIG_FILE');
const enabled = config.versions.filter(v => v.enabled).map(v => ({ id: v.id, gitRef: v.gitRef }));
console.log(JSON.stringify(enabled));
")

echo "📋 Enabled versions from config:"
echo "$ENABLED_VERSIONS" | node -e "
const versions = JSON.parse(require('fs').readFileSync(0, 'utf8'));
versions.forEach(v => console.log(\`  - \${v.id} (git: \${v.gitRef})\`));
"

# Create versions directory
mkdir -p "$VERSIONS_DIR"

# Generate each version
echo "$ENABLED_VERSIONS" | node -e "
const versions = JSON.parse(require('fs').readFileSync(0, 'utf8'));
versions.forEach(version => {
    console.log(\`\nGenerating version: \${version.id}\`);
    
    const { execSync } = require('child_process');
    const fs = require('fs');
    const path = require('path');
    
    const versionDir = path.join('$VERSIONS_DIR', version.id);
    
    // Skip if version already exists (unless it's 'main' which should always be fresh)
    if (fs.existsSync(versionDir) && version.id !== 'main') {
        console.log(\`  ✓ Version \${version.id} already exists, skipping\`);
        return;
    }
    
    try {
        // Create version directory
        fs.mkdirSync(versionDir, { recursive: true });
        
        // Store current branch/commit
        const currentBranch = execSync('git rev-parse --abbrev-ref HEAD', { encoding: 'utf8' }).trim();
        const currentCommit = execSync('git rev-parse HEAD', { encoding: 'utf8' }).trim();
        
        // Checkout the target ref if it's different from current
        if (version.gitRef !== currentBranch && version.gitRef !== 'HEAD') {
            console.log(\`  📦 Checking out \${version.gitRef}\`);
            execSync(\`git fetch origin \${version.gitRef} || git fetch origin\`, { stdio: 'inherit' });
            execSync(\`git checkout \${version.gitRef}\`, { stdio: 'inherit' });
        }
        
        // Generate artifacts for this version
        console.log(\`  🔨 Generating artifacts for \${version.id}\`);
        execSync('pnpm run generate', { 
            stdio: 'inherit',
            cwd: '$SCRIPT_DIR/..'
        });
        
        // Copy artifacts to version directory
        console.log(\`  📁 Copying artifacts to \${version.id}\`);
        if (fs.existsSync('$ARTIFACTS_DIR/openapi3')) {
            execSync(\`cp -r $ARTIFACTS_DIR/openapi3 \${versionDir}/\`, { stdio: 'inherit' });
        }
        if (fs.existsSync('$ARTIFACTS_DIR/openapi')) {
            execSync(\`cp -r $ARTIFACTS_DIR/openapi \${versionDir}/\`, { stdio: 'inherit' });
        }
        
        // Create metadata
        const metadata = {
            version: version.id,
            gitRef: version.gitRef,
            generatedAt: new Date().toISOString(),
            gitCommit: execSync('git rev-parse HEAD', { encoding: 'utf8' }).trim(),
            gitBranch: execSync('git rev-parse --abbrev-ref HEAD', { encoding: 'utf8' }).trim()
        };
        
        fs.writeFileSync(
            path.join(versionDir, 'metadata.json'),
            JSON.stringify(metadata, null, 2)
        );
        
        console.log(\`  ✅ Generated version \${version.id}\`);
        
        // Restore original branch/commit if we changed it
        if (version.gitRef !== currentBranch && version.gitRef !== 'HEAD') {
            console.log(\`  🔙 Restoring \${currentBranch}\`);
            execSync(\`git checkout \${currentBranch}\`, { stdio: 'inherit' });
        }
        
    } catch (error) {
        console.error(\`  ❌ Failed to generate version \${version.id}:\`, error.message);
        
        // Try to restore original branch on error
        try {
            if (version.gitRef !== currentBranch && version.gitRef !== 'HEAD') {
                execSync(\`git checkout \${currentBranch}\`, { stdio: 'pipe' });
            }
        } catch (restoreError) {
            console.error('Failed to restore original branch:', restoreError.message);
        }
    }
});
"

# Generate current/latest artifacts
echo -e "\n🔨 Generating current artifacts..."
cd "$SCRIPT_DIR/.."
pnpm run generate

echo -e "\n✅ Version generation complete!"
echo "📊 Available versions:"
ls -la "$VERSIONS_DIR" 2>/dev/null || echo "No versions directory yet"
