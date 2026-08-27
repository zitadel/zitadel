#!/usr/bin/env bash
# Usage: tools/status_breakdown.sh <in.csv> <out.txt>
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DUCK="${DUCKDB:-}"
[ -x "$DUCK" ] || DUCK="$(command -v duckdb 2>/dev/null || echo "${HERE}/.bin/duckdb")"
mkdir -p "${HERE}/.bin/tmp"
"${DUCK}" -c "SET temp_directory='${HERE}/.bin/tmp'; SET variable csv='$1';" \
  -f "${HERE}/status_breakdown.sql" > "$2"
