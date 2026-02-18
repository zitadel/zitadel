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

# ----- protoc-gen-connect-openapi v0.25.2 (sudorandom/protoc-gen-connect-openapi) -----
# Note: macOS ships a universal binary ("darwin_all") for both amd64 and arm64.
case "$GOOS" in
  linux)  OAI_ARCH="$GOARCH" ;;
  darwin) OAI_ARCH="all"      ;;
esac

echo "Downloading protoc-gen-connect-openapi v0.25.2 (${GOOS}/${OAI_ARCH})..."
curl -fsSL \
  "https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v0.25.2/protoc-gen-connect-openapi_0.25.2_${GOOS}_${OAI_ARCH}.tar.gz" \
  -o "${TMP}/oai.tar.gz"
tar -xzf "${TMP}/oai.tar.gz" -C "${TMP}" protoc-gen-connect-openapi
install -m 755 "${TMP}/protoc-gen-connect-openapi" "${BIN_DIR}/protoc-gen-connect-openapi"

echo "protoc-gen-connect-openapi installed to ${BIN_DIR}"
