#!/usr/bin/env bash
# Convert a k6 CSV run into the docs output.json format.
# Usage: tools/csv2json.sh <in.csv> <out.json>
#
# Needs DuckDB. Resolution order: $DUCKDB, duckdb on PATH, tools/.bin/duckdb.
# If none is present it downloads the CLI into tools/.bin (gitignored).
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

resolve_duckdb() {
  if [ -n "${DUCKDB:-}" ] && [ -x "${DUCKDB}" ]; then echo "${DUCKDB}"; return; fi
  if command -v duckdb >/dev/null 2>&1; then command -v duckdb; return; fi
  if [ -x "${HERE}/.bin/duckdb" ]; then echo "${HERE}/.bin/duckdb"; return; fi
  # Pinned rather than "latest" so a run is reproducible; override with DUCKDB_VERSION.
  # Platforms other than these must install duckdb themselves and point DUCKDB at it.
  version="${DUCKDB_VERSION:-v1.1.3}"
  case "$(uname -s)-$(uname -m)" in
    Linux-x86_64)   asset="duckdb_cli-linux-amd64.zip" ;;
    Linux-aarch64)  asset="duckdb_cli-linux-aarch64.zip" ;;
    Darwin-x86_64|Darwin-arm64) asset="duckdb_cli-osx-universal.zip" ;;
    *)
      echo "csv2json: no DuckDB CLI download for $(uname -s)-$(uname -m)." >&2
      echo "          Install duckdb and re-run with DUCKDB=/path/to/duckdb" >&2
      exit 1
      ;;
  esac
  echo "duckdb not found, downloading ${version} (${asset}) into ${HERE}/.bin" >&2
  mkdir -p "${HERE}/.bin"
  if ! curl -fsSL -o "${HERE}/.bin/duckdb.zip" \
    "https://github.com/duckdb/duckdb/releases/download/${version}/${asset}" >&2; then
    rm -f "${HERE}/.bin/duckdb.zip"
    echo "csv2json: failed to download DuckDB ${version} (${asset})." >&2
    echo "          Install duckdb and re-run with DUCKDB=/path/to/duckdb," >&2
    echo "          or pick another release with DUCKDB_VERSION=vX.Y.Z" >&2
    exit 1
  fi
  python3 -c "import zipfile,sys;zipfile.ZipFile(sys.argv[1]).extractall(sys.argv[2])" \
    "${HERE}/.bin/duckdb.zip" "${HERE}/.bin" >&2
  chmod +x "${HERE}/.bin/duckdb"
  rm -f "${HERE}/.bin/duckdb.zip"
  echo "${HERE}/.bin/duckdb"
}

DUCK="$(resolve_duckdb)"
mkdir -p "$(dirname "$2")" "${HERE}/.bin/tmp"
# `sed` reshapes DuckDB's json mode into the tab-indented layout the existing
# published output.json files use.
"${DUCK}" -c "SET temp_directory='${HERE}/.bin/tmp'; SET variable csv='$1';" -f "${HERE}/convert.sql" \
  | sed -e 's/^\[{/[\n\t{/' -e 's/^{/\t{/' -e 's/}\]$/}\n]/' > "$2"
