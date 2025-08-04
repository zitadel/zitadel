#!/bin/bash
set -e

echo "Downloading protoc-gen-connect-openapi plugin..."
echo "Architecture: $(uname -m)"
echo "OS: $(uname)"

# Create directory if it doesn't exist
mkdir -p protoc-gen-connect-openapi
cd ./protoc-gen-connect-openapi/

# Skip download if plugin already exists and is executable
if [ -f "protoc-gen-connect-openapi" ] && [ -x "protoc-gen-connect-openapi" ]; then
  echo "Plugin already exists and is executable"
  ./protoc-gen-connect-openapi --version || echo "Plugin version check failed, but file exists"
  exit 0
fi

# Clean up any partial downloads
rm -f protoc-gen-connect-openapi.tar.gz protoc-gen-connect-openapi

# Determine download URL based on OS and architecture
if [ "$(uname)" = "Darwin" ]; then
  echo "Downloading for Darwin..."
  URL="https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v0.18.0/protoc-gen-connect-openapi_0.18.0_darwin_all.tar.gz"
else
  ARCH=$(uname -m)
  case $ARCH in
    x86_64)
      ARCH="amd64"
      ;;
    aarch64|arm64)
      ARCH="arm64"
      ;;
    *)
      echo "Unsupported architecture: $ARCH"
      exit 1
      ;;
  esac
  echo "Downloading for Linux ${ARCH}..."
  URL="https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v0.18.0/protoc-gen-connect-openapi_0.18.0_linux_${ARCH}.tar.gz"
fi

# Download with retries
echo "Downloading from: $URL"
curl -L -o protoc-gen-connect-openapi.tar.gz "$URL" || {
  echo "Download failed, trying with different curl options..."
  curl -L --fail --retry 3 --retry-delay 1 -o protoc-gen-connect-openapi.tar.gz "$URL"
}

echo "Extracting plugin..."
tar -xzf protoc-gen-connect-openapi.tar.gz

# Verify extraction
if [ ! -f "protoc-gen-connect-openapi" ]; then
  echo "ERROR: Plugin binary not found after extraction"
  ls -la
  exit 1
fi

# Make sure the plugin is executable
chmod +x protoc-gen-connect-openapi

# Verify plugin works
echo "Plugin installed successfully"
ls -la protoc-gen-connect-openapi
./protoc-gen-connect-openapi --version || echo "Plugin version check failed, but installation completed"

# Clean up
rm -f protoc-gen-connect-openapi.tar.gz