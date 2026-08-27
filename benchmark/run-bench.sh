#!/usr/bin/env bash
# Usage: ./run-bench.sh <target> <vus> <duration>
set -euo pipefail
cd "$(dirname "$0")"
if [ -f .hold ]; then
  echo "HOLD: .hold present, skipping $1"
  exit 0
fi
set -a; . ./.env; set +a

# Requests whose path carries a user/session/project id now set a grouped `name`
# tag at the call site (e.g. /v2beta/users/{userId}), so `name` stays low
# cardinality and is worth keeping. The raw `url` tag is still one time series per
# request - over 800k in a 30 minute run against k6's suggested limit of 100k,
# which drove k6 to 6.8GB and an OOM - so it stays out.
export K6_SYSTEM_TAGS="${K6_SYSTEM_TAGS:-proto,status,method,name,group,check,error,error_code,scenario,expected_response}"

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
