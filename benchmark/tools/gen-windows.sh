#!/usr/bin/env bash
# Emit the benchmark test windows as `target|startISO|endISO`, read from the published
# pages' `Test start` / `Test end` rows -- the only absolute record of when a test ran.
#
#   tools/gen-windows.sh [version]        # default: newest vN directory
#
# Pipe into a file and pass it to the fetch scripts as WINDOWS_FILE, or use
# tools/sync-windows.sh to rewrite their embedded defaults in place.
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOCS="${DOCS:-${HERE}/../../apps/docs/content/apis/benchmarks}"

VERSION="${1:-}"
if [ -z "$VERSION" ]; then
  VERSION="$(ls -1 "$DOCS" | grep -E '^v[0-9]' | sort -V | tail -1)"
fi
DIR="${DOCS}/${VERSION}"
[ -d "$DIR" ] || { echo "gen-windows: no such version: ${VERSION}" >&2; exit 1; }

found=0
for f in "$DIR"/*/index.mdx; do
  target="$(basename "$(dirname "$f")")"
  s="$(sed -nE 's/^\| Test start[[:space:]]*\|[[:space:]]*([0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9:]{8}) UTC.*/\1T\2Z/p' "$f" | head -1)"
  e="$(sed -nE 's/^\| Test end[[:space:]]*\|[[:space:]]*([0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9:]{8}) UTC.*/\1T\2Z/p' "$f" | head -1)"
  if [ -z "$s" ] || [ -z "$e" ]; then
    echo "gen-windows: ${target} has no usable Test start/end -- pages before v4.17.1 carry a bare HH:MM" >&2
    continue
  fi
  echo "${target}|${s}|${e}"
  found=$((found+1))
done
[ "$found" -gt 0 ] || { echo "gen-windows: no windows found in ${VERSION}" >&2; exit 1; }
echo "gen-windows: ${found} window(s) from ${VERSION}" >&2
