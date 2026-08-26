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
  echo "duckdb not found, downloading CLI into ${HERE}/.bin" >&2
  mkdir -p "${HERE}/.bin"
  curl -sSL -o "${HERE}/.bin/duckdb.zip" \
    "https://github.com/duckdb/duckdb/releases/latest/download/duckdb_cli-linux-amd64.zip" >&2
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
