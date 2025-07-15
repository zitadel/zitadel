#!/bin/bash

# Script to prepare standalone version of the login app
set -e

echo "🔧 Preparing standalone version..."

# Parse arguments
INSTALL_DEPS=true
USE_LATEST=false

for arg in "$@"; do
  case $arg in
    --no-install)
      INSTALL_DEPS=false
      shift
      ;;
    --latest)
      USE_LATEST=true
      shift
      ;;
    *)
      # Unknown option
      ;;
  esac
done

# Build arguments for Node.js script
NODE_ARGS=""
if [ "$INSTALL_DEPS" = true ]; then
  NODE_ARGS="$NODE_ARGS --install"
fi
if [ "$USE_LATEST" = true ]; then
  NODE_ARGS="$NODE_ARGS --latest"
fi

# Check if Node.js scripts exist
if [ ! -f "scripts/prepare-standalone.js" ]; then
  echo "❌ scripts/prepare-standalone.js not found!"
  echo "   Make sure you're in the correct directory"
  exit 1
fi

# Run the enhanced Node.js prepare script
node scripts/prepare-standalone.js $NODE_ARGS

echo ""
echo "✅ Standalone version prepared successfully!"

if [ "$INSTALL_DEPS" = false ]; then
  echo ""
  echo "📝 Next steps:"
  echo "  npm install        - Install dependencies"
  echo "  npm run dev        - Start development server"
  echo "  npm run build      - Build for production"
  echo "  npm run start      - Start production server"
fi
