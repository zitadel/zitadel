#!/usr/bin/env bash
# Downloads protoc-gen-grpc-web, protoc-gen-js, and protoc-gen-openapiv2
# for console protobuf generation.
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

# ── VERSIONS & CHECKSUMS ────────────────────────────────────────────────────
# To upgrade a plugin: update the VERSION and all four SHA256_<os>_<arch> vars.

GRPC_WEB_VERSION="1.5.0"
# checksums from upstream .sha256 sidecar files
GRPC_WEB_SHA256_linux_amd64="2e6e074497b221045a14d5a54e9fc910945bfdd1198b12b9fc23686a95671d64"
GRPC_WEB_SHA256_linux_arm64="522e958568cdeabdd20ef3c97394fc067ff8e374a728c08b7410bf5de8f57783"
GRPC_WEB_SHA256_darwin_amd64="1fa3ef92194d06c03448a5cba82759e9773e43d8b188866a1f1d4fc23bb1ecb7"
GRPC_WEB_SHA256_darwin_arm64="a12b759629b943ebac5528f58fac5039d9aa2fb7abb9e9684d1b481b35afbfc6"

PROTOC_GEN_JS_VERSION="4.0.2"
# checksums of release zip assets (this release ships no upstream checksum file)
PROTOC_GEN_JS_SHA256_linux_amd64="e4b0bc2c9fe32a21167c6d84a50df21c2f405552f9a6ed7d5e858d92cac46d22"
PROTOC_GEN_JS_SHA256_linux_arm64="5c0242143a99d2ca1f84453ada97bfb156fda2b7bdad9f8ca01aa8c40b681499"
PROTOC_GEN_JS_SHA256_darwin_amd64="9d61a636f392d020f1224aa6e533e65bd7b65da71219742a0b5a84de60d6f59d"
PROTOC_GEN_JS_SHA256_darwin_arm64="73092321c6655a61991d57b5d9a9c54640517f58bcb379d6b8db0220905cef5d"

OPENAPIV2_VERSION="2.22.0"
# checksums from grpc-gateway_${OPENAPIV2_VERSION}_checksums.txt
OPENAPIV2_SHA256_linux_amd64="72a6fc6a6d130189c549a6d88cbdef407d3ed1c95ab101ffb3d223d8b3778c90"
OPENAPIV2_SHA256_linux_arm64="4921799b8d80dde5f8cb89817d3ae04dee1e2560e141fd0fc79a2e544cc63178"
OPENAPIV2_SHA256_darwin_amd64="14c95d1305a81822cd17ef750cfe71e8471728eba19068e9142a70a6cbaf5847"
OPENAPIV2_SHA256_darwin_arm64="dc215925f49912d53a443107879d91898b56699e5b8bc4ed0d7b9ba94939dd86"

# ── HELPERS ──────────────────────────────────────────────────────────────────

verify_sha256() {
  local file="$1" expected="$2"
  local actual
  if command -v sha256sum &>/dev/null; then
    actual=$(sha256sum "$file" | cut -d' ' -f1)
  else
    actual=$(shasum -a 256 "$file" | cut -d' ' -f1)
  fi
  if [ "$actual" != "$expected" ]; then
    echo "ERROR: SHA256 mismatch for $(basename "$file")" >&2
    echo "  expected: $expected" >&2
    echo "  actual:   $actual" >&2
    exit 1
  fi
}

# ── PLATFORM DETECTION ───────────────────────────────────────────────────────

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

# ── INSTALL ───────────────────────────────────────────────────────────────────

# ----- protoc-gen-grpc-web (grpc/grpc-web) -----
case "$GOARCH" in
  amd64) WEB_ARCH="x86_64"  ;;
  arm64) WEB_ARCH="aarch64" ;;
esac
sha256_var="GRPC_WEB_SHA256_${GOOS}_${GOARCH}"
echo "Downloading protoc-gen-grpc-web v${GRPC_WEB_VERSION} (${GOOS}/${WEB_ARCH})..."
curl -fsSL \
  "https://github.com/grpc/grpc-web/releases/download/${GRPC_WEB_VERSION}/protoc-gen-grpc-web-${GRPC_WEB_VERSION}-${GOOS}-${WEB_ARCH}" \
  -o "${TMP}/protoc-gen-grpc-web"
verify_sha256 "${TMP}/protoc-gen-grpc-web" "${!sha256_var}"
install -m 755 "${TMP}/protoc-gen-grpc-web" "${BIN_DIR}/protoc-gen-grpc-web"

# ----- protoc-gen-js (protocolbuffers/protobuf-javascript) -----
# Note: macOS asset uses "osx"; ARM uses "aarch_64" (underscore, not hyphen).
case "$GOOS" in
  linux)  JS_OS="linux" ;;
  darwin) JS_OS="osx"   ;;
esac
case "$GOARCH" in
  amd64) JS_ARCH="x86_64"   ;;
  arm64) JS_ARCH="aarch_64" ;;
esac
sha256_var="PROTOC_GEN_JS_SHA256_${GOOS}_${GOARCH}"
echo "Downloading protoc-gen-js v${PROTOC_GEN_JS_VERSION} (${JS_OS}/${JS_ARCH})..."
curl -fsSL \
  "https://github.com/protocolbuffers/protobuf-javascript/releases/download/v${PROTOC_GEN_JS_VERSION}/protobuf-javascript-${PROTOC_GEN_JS_VERSION}-${JS_OS}-${JS_ARCH}.zip" \
  -o "${TMP}/p.zip"
verify_sha256 "${TMP}/p.zip" "${!sha256_var}"
unzip -q "${TMP}/p.zip" -d "${TMP}"
install -m 755 "${TMP}/bin/protoc-gen-js" "${BIN_DIR}/protoc-gen-js"

# ----- protoc-gen-openapiv2 (grpc-ecosystem/grpc-gateway) -----
case "$GOARCH" in
  amd64) OAI2_ARCH="x86_64" ;;
  arm64) OAI2_ARCH="arm64"  ;;
esac
sha256_var="OPENAPIV2_SHA256_${GOOS}_${GOARCH}"
echo "Downloading protoc-gen-openapiv2 v${OPENAPIV2_VERSION} (${GOOS}/${OAI2_ARCH})..."
curl -fsSL \
  "https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v${OPENAPIV2_VERSION}/protoc-gen-openapiv2-v${OPENAPIV2_VERSION}-${GOOS}-${OAI2_ARCH}" \
  -o "${TMP}/protoc-gen-openapiv2"
verify_sha256 "${TMP}/protoc-gen-openapiv2" "${!sha256_var}"
install -m 755 "${TMP}/protoc-gen-openapiv2" "${BIN_DIR}/protoc-gen-openapiv2"

echo "Console proto plugins installed to ${BIN_DIR}"
