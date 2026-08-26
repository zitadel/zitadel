#!/usr/bin/env bash
# Usage: ./run-bench.sh <target> <vus> <duration>
set -euo pipefail
cd "$(dirname "$0")"
if [ -f .hold ]; then
  echo "HOLD: .hold present, skipping $1"
  exit 0
fi
set -a; . ./.env; set +a
TARGET="$1"; VUS="$2"; DURATION="$3"
STAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
mkdir -p summaries
LOG="summaries/${TARGET}_${STAMP}.log"

# Pre-flight: a previous failed run leaves an orphan org behind whose machine
# users (zitachine-N) collide instance-wide with the next run's setup, failing it
# with 409 User already exists. Safe here: no k6 is running yet at this point, so
# no live benchmark org can match. Deletion is asynchronous, so wait for it.
if [ -n "${ADMIN_PAT:-}" ]; then
  echo "preflight: checking for orphan benchmark orgs" | tee -a "${LOG}"
  ./cleanup-stale-orgs.sh --confirm 2>&1 | tee -a "${LOG}" || true
  for _ in $(seq 1 30); do
    left="$(curl -sf -m 20 -X POST "${ZITADEL_HOST}/v2/organizations/_search" \
      -H "Authorization: Bearer ${ADMIN_PAT}" -H "Content-Type: application/json" \
      -d '{"query":{"limit":1000}}' \
      | jq --arg skip "$(tr '\n' ',' < .skip-orgs 2>/dev/null || true)" \
           '[.result[]? | select(.name | test("^load-test-"))
                        | select(. as $o | $skip | contains($o.id) | not)] | length')" || left=0
    [ "${left:-0}" -eq 0 ] && break
    sleep 2
  done
  echo "preflight: ${left:-0} benchmark org(s) remaining" | tee -a "${LOG}"
fi

echo "start_utc=$(date -u '+%Y-%m-%d %H:%M:%S UTC') target=${TARGET} vus=${VUS} duration=${DURATION} host=${ZITADEL_HOST}" | tee -a "${LOG}"
make "${TARGET}" \
  ZITADEL_HOST="${ZITADEL_HOST}" \
  ADMIN_LOGIN_NAME="${ADMIN_LOGIN_NAME}" \
  ADMIN_PASSWORD="${ADMIN_PASSWORD}" \
  VUS="${VUS}" DURATION="${DURATION}" 2>&1 | tee -a "${LOG}"
echo "end_utc=$(date -u '+%Y-%m-%d %H:%M:%S UTC')" | tee -a "${LOG}"
