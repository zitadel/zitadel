#!/usr/bin/env bash
# Convert the newest k6 CSV into the docs page for <target>, then drop the CSV.
# Usage: postprocess.sh <target> [version]
set -euo pipefail
cd "$(dirname "$0")"
TARGET="$1"; VERSION="${2:-v4.17.1}"
CSV="$(ls -t output/*.csv | head -1)"
LOG="$(ls -t summaries/${TARGET}_*.log | head -1)"
OUT="output/${TARGET}-output.json"
echo "postprocess: target=${TARGET} csv=${CSV} log=${LOG}"
./tools/csv2json.sh "$(realpath "$CSV")" "${OUT}"
ROWS="$(jq 'length' "${OUT}")"
if [ "$ROWS" -lt 2 ]; then
  echo "postprocess: ERROR only $ROWS rows from $CSV -- keeping CSV for inspection" >&2
  exit 1
fi
python3 tools/gen_report.py "$TARGET" "$LOG" "${OUT}" "$VERSION"
rm -f "$CSV"
echo "postprocess: ${TARGET} done (${ROWS} rows), csv removed, df: $(df -h /home/tim | awk 'NR==2{print $4}') free"
