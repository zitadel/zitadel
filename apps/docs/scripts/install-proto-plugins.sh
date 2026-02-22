#!/usr/bin/env bash
# Downloads protoc-gen-connect-openapi for docs OpenAPI generation.
# Must be run from the workspace root.
# Supports Linux (including WSL2) and macOS only.
#
# WHY a binary download instead of "go install":
#   This script runs as an Nx target that docs/generate-proto-docs depends on.
#   Docs OpenAPI generation happens both in CI (has Go) and on Vercel
#   (Node.js-only, no Go). Using a pre-built binary means Vercel never needs a
#   Go toolchain. Nx remote cache means the download is skipped on repeated runs.
#
# WHY uname instead of "go env GOOS/GOARCH":
#   Same reason — this script must work without Go present.
set -euo pipefail

# ── VERSIONS & CHECKSUMS ────────────────────────────────────────────────────
# To upgrade: update VERSION and the SHA256_<os>_<arch> vars below.
# Note: macOS ships a universal binary, so darwin_amd64 and darwin_arm64 share
#       the same tarball and checksum.

CONNECT_OPENAPI_VERSION="0.25.2"
# checksums from upstream checksums.txt
CONNECT_OPENAPI_SHA256_linux_amd64="90821ab96f16747dc5a4ab93900a6bc6b968d63f47553a08274dc594301bced3"
CONNECT_OPENAPI_SHA256_linux_arm64="9726418d7af55f87e3bbb2a4fa6233016b200caf5a285ce1c6e6428ef225536d"
CONNECT_OPENAPI_SHA256_darwin_amd64="a6bc3d87b4caaa7048d22fe0f45e3b5e778b4b209f84a836b3335fb2fd14b588"
CONNECT_OPENAPI_SHA256_darwin_arm64="a6bc3d87b4caaa7048d22fe0f45e3b5e778b4b209f84a836b3335fb2fd14b588"

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

TMP=$(mktemp -d "${TMPDIR:-/tmp}/zitadel-docs-proto-plugins.XXXXXX")
trap 'rm -rf "$TMP"' EXIT

# ── INSTALL ───────────────────────────────────────────────────────────────────

# ----- protoc-gen-connect-openapi (sudorandom/protoc-gen-connect-openapi) -----
# macOS ships a single universal ("all") tarball for both amd64 and arm64.
case "$GOOS" in
  linux)  OAI_ARCH="$GOARCH" ;;
  darwin) OAI_ARCH="all"     ;;
esac
sha256_var="CONNECT_OPENAPI_SHA256_${GOOS}_${GOARCH}"
echo "Downloading protoc-gen-connect-openapi v${CONNECT_OPENAPI_VERSION} (${GOOS}/${OAI_ARCH})..."
curl -fsSL \
  "https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v${CONNECT_OPENAPI_VERSION}/protoc-gen-connect-openapi_${CONNECT_OPENAPI_VERSION}_${GOOS}_${OAI_ARCH}.tar.gz" \
  -o "${TMP}/oai.tar.gz"
verify_sha256 "${TMP}/oai.tar.gz" "${!sha256_var}"
tar -xzf "${TMP}/oai.tar.gz" -C "${TMP}" protoc-gen-connect-openapi
install -m 755 "${TMP}/protoc-gen-connect-openapi" "${BIN_DIR}/protoc-gen-connect-openapi"

echo "protoc-gen-connect-openapi installed to ${BIN_DIR}"
