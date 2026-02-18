#!/usr/bin/env bash
# Downloads protoc-gen-grpc-web and protoc-gen-js for console protobuf generation.
# Must be run from the workspace root.
# Supports Linux (including WSL2) and macOS only.
#
# WHY binary downloads instead of BSR remote plugins:
#   Using locally installed binaries avoids a network call to buf.build on every
#   cache miss and removes the dependency on BSR availability during CI builds.
#   Nx caches the outputs, so the downloads only happen when versions change.
#
# WHY uname instead of "go env GOOS/GOARCH":
#   Keeps the script dependency-free (no Go toolchain required).
set -euo pipefail

_uname_os=$(uname -s | tr '[:upper:]' '[:lower:]')
_uname_arch=$(uname -m)
case "$_uname_os" in
  linux)  GOOS="linux"  ;;
  darwin) GOOS="darwin" ;;
  *) echo "Unsupported OS: $_uname_os" >&2; exit 1 ;;
esac
case "$_uname_arch" in
  x86_64)          GOARCH="amd64" ;;
  aarch64 | arm64) GOARCH="arm64" ;;
  *) echo "Unsupported arch: $_uname_arch" >&2; exit 1 ;;
esac

BIN_DIR="${PWD}/.artifacts/bin/${GOOS}/${GOARCH}"
mkdir -p "$BIN_DIR"

TMP=$(mktemp -d "${TMPDIR:-/tmp}/zitadel-console-proto-plugins.XXXXXX")
trap 'rm -rf "$TMP"' EXIT

# ----- protoc-gen-grpc-web v1.5.0 (grpc/grpc-web) -----
case "$GOARCH" in
  amd64) WEB_ARCH="x86_64"  ;;
  arm64) WEB_ARCH="aarch64" ;;
esac

echo "Downloading protoc-gen-grpc-web v1.5.0 (${GOOS}/${WEB_ARCH})..."
curl -fsSL \
  "https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-${GOOS}-${WEB_ARCH}" \
  -o "${TMP}/protoc-gen-grpc-web"
install -m 755 "${TMP}/protoc-gen-grpc-web" "${BIN_DIR}/protoc-gen-grpc-web"

# ----- protoc-gen-js v3.21.4 (protocolbuffers/protobuf-javascript) -----
# Note: macOS asset uses "osx"; ARM uses "aarch_64" (underscore, not hyphen).
case "$GOOS" in
  linux)  JS_OS="linux" ;;
  darwin) JS_OS="osx"   ;;
esac
case "$GOARCH" in
  amd64) JS_ARCH="x86_64"   ;;
  arm64) JS_ARCH="aarch_64" ;;
esac

echo "Downloading protoc-gen-js v3.21.4 (${JS_OS}/${JS_ARCH})..."
curl -fsSL \
  "https://github.com/protocolbuffers/protobuf-javascript/releases/download/v3.21.4/protobuf-javascript-3.21.4-${JS_OS}-${JS_ARCH}.tar.gz" \
  -o "${TMP}/p.tar.gz"
tar -xzf "${TMP}/p.tar.gz" -C "${TMP}"
install -m 755 "${TMP}/bin/protoc-gen-js" "${BIN_DIR}/protoc-gen-js"

echo "Console proto plugins installed to ${BIN_DIR}"
