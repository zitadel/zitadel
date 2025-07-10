#!/bin/bash

# Script to prepare standalone version of the login app
set -e

echo "Preparing standalone version..."

# Copy standalone configs
cp package.standalone.json package.json
cp tsconfig.standalone.json tsconfig.json
cp .eslintrc.standalone.cjs .eslintrc.cjs
cp prettier.config.standalone.mjs prettier.config.mjs
cp tailwind.config.standalone.mjs tailwind.config.mjs

# Install dependencies unless --no-install is passed
if [ "$1" != "--no-install" ]; then
    echo "Installing dependencies..."
    npm install
fi

echo "Standalone version prepared successfully!"
if [ "$1" != "--no-install" ]; then
    echo "You can now run:"
    echo "  npm run dev        - Start development server"
    echo "  npm run build      - Build for production"
    echo "  npm run start      - Start production server"
fi
