#!/usr/bin/env bash
# Re-run every target the sweep logged as failed, once. Usage: retry-failed.sh <sweep.log>
cd "$(dirname "$0")"
LOG="$1"
mapfile -t failed < <(grep -oE "QUEUE (run|postprocess) FAILED [a-z_]+" "$LOG" | awk '{print $NF}' | sort -u)
if [ "${#failed[@]}" -eq 0 ]; then echo "=== RETRY: nothing failed ==="; exit 0; fi
echo "=== RETRY: ${failed[*]} ==="
for t in "${failed[@]}"; do
  echo "=== QUEUE start ${t} (retry) $(date -u '+%Y-%m-%d %H:%M:%S UTC') ==="
  if ./run-bench.sh "$t" 600 1800s; then
    ./postprocess.sh "$t" v4.17.1 || echo "=== QUEUE postprocess FAILED ${t} ==="
  else
    echo "=== QUEUE run FAILED ${t} ==="
  fi
  echo "=== QUEUE done ${t} (retry) $(date -u '+%Y-%m-%d %H:%M:%S UTC') ==="
done
echo "=== RETRY DONE ==="
