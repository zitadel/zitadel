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
echo "start_utc=$(date -u '+%Y-%m-%d %H:%M:%S UTC') target=${TARGET} vus=${VUS} duration=${DURATION} host=${ZITADEL_HOST}" | tee "${LOG}"
make "${TARGET}" \
  ZITADEL_HOST="${ZITADEL_HOST}" \
  ADMIN_LOGIN_NAME="${ADMIN_LOGIN_NAME}" \
  ADMIN_PASSWORD="${ADMIN_PASSWORD}" \
  VUS="${VUS}" DURATION="${DURATION}" 2>&1 | tee -a "${LOG}"
echo "end_utc=$(date -u '+%Y-%m-%d %H:%M:%S UTC')" | tee -a "${LOG}"
