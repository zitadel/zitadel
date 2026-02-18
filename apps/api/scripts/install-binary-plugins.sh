#!/usr/bin/env bash
# Downloads platform-specific binary proto plugins that are not available via go install.
# Must be run from the workspace root (same as other generate-install commands).
# Supports Linux (including WSL2) and macOS only.
set -euo pipefail

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
BIN_DIR="${PWD}/.artifacts/bin/${GOOS}/${GOARCH}"

mkdir -p "$BIN_DIR"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

# ----- protoc-gen-grpc-web v1.5.0 (grpc/grpc-web) -----
case "$GOOS" in
  linux)  WEB_OS="linux"  ;;
  darwin) WEB_OS="darwin" ;;
  *) echo "Unsupported OS: $GOOS" >&2; exit 1 ;;
esac
case "$GOARCH" in
  amd64) WEB_ARCH="x86_64"  ;;
  arm64) WEB_ARCH="aarch64" ;;
  *) echo "Unsupported arch: $GOARCH" >&2; exit 1 ;;
esac

echo "Downloading protoc-gen-grpc-web v1.5.0 (${WEB_OS}/${WEB_ARCH})..."
curl -sL \
  "https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-${WEB_OS}-${WEB_ARCH}" \
  -o "${TMP}/protoc-gen-grpc-web"
install -m 755 "${TMP}/protoc-gen-grpc-web" "${BIN_DIR}/protoc-gen-grpc-web"

# ----- protoc-gen-js v3.21.4 (protocolbuffers/protobuf-javascript) -----
# Note: macOS asset uses "osx"; ARM uses "aarch_64" (underscore, not hyphen).
case "$GOOS" in
  linux)  JS_OS="linux" ;;
  darwin) JS_OS="osx"   ;;
  *) echo "Unsupported OS: $GOOS" >&2; exit 1 ;;
esac
case "$GOARCH" in
  amd64) JS_ARCH="x86_64"   ;;
  arm64) JS_ARCH="aarch_64" ;;
  *) echo "Unsupported arch: $GOARCH" >&2; exit 1 ;;
esac

echo "Downloading protoc-gen-js v3.21.4 (${JS_OS}/${JS_ARCH})..."
curl -sL \
  "https://github.com/protocolbuffers/protobuf-javascript/releases/download/v3.21.4/protobuf-javascript-3.21.4-${JS_OS}-${JS_ARCH}.tar.gz" \
  -o "${TMP}/p.tar.gz"
tar -xzf "${TMP}/p.tar.gz" -C "${TMP}"
install -m 755 "${TMP}/bin/protoc-gen-js" "${BIN_DIR}/protoc-gen-js"

echo "Binary proto plugins installed to ${BIN_DIR}"
