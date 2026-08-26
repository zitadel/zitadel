#!/usr/bin/env bash
# Delete exactly ONE benchmark org, by id AND expected name. Both must match, and
# the name must match the load-test-<ISO timestamp> convention. Used when a run
# fails mid-sweep and leaves an orphan while another run is legitimately in flight.
set -euo pipefail
cd "$(dirname "$0")"
set -a; . ./.env; set +a
PATTERN='^load-test-[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3}Z$'
ID="$1"; EXPECTED="$2"
resp="$(curl -sf -m 30 -X POST "${ZITADEL_HOST}/v2/organizations/_search" \
  -H "Authorization: Bearer ${ADMIN_PAT}" -H "Content-Type: application/json" \
  -d '{"query":{"limit":1000}}')"
name="$(jq -r --arg id "$ID" '.result[]? | select(.id==$id) | .name' <<<"$resp")"
[ -n "$name" ]              || { echo "ABORT: org $ID not found"; exit 1; }
[ "$name" = "$EXPECTED" ]   || { echo "ABORT: name mismatch: got '$name', expected '$EXPECTED'"; exit 1; }
[[ "$name" =~ $PATTERN ]]   || { echo "ABORT: '$name' does not match the benchmark convention"; exit 1; }
code="$(curl -s -m 30 -o /dev/null -w '%{http_code}' -X DELETE \
  "${ZITADEL_HOST}/v2/organizations/${ID}" -H "Authorization: Bearer ${ADMIN_PAT}")"
echo "deleted ${ID} ${name} -> HTTP ${code}"
