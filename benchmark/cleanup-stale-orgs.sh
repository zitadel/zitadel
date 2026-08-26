#!/usr/bin/env bash
# List (and optionally delete) stale benchmark organizations.
#
# SAFETY: only organizations whose name matches EXACTLY
#   load-test-<ISO8601 timestamp>      e.g. load-test-2026-08-26T20:11:56.440Z
# are ever considered. Anything else is skipped, always.
# Dry-run by default; pass --confirm to actually delete.
# Never run this while a benchmark is in flight: the live run's own org matches
# the pattern too. The script refuses if a k6 process is running.
set -euo pipefail
cd "$(dirname "$0")"
set -a; . ./.env; set +a
PATTERN='^load-test-[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3}Z$'
CONFIRM="${1:-}"

if pgrep -f "xk6-modules/k6 run" >/dev/null; then
  echo "REFUSING: a k6 run is in flight; its org matches the pattern. Stop it first." >&2
  exit 1
fi

resp="$(curl -sf -m 30 -X POST "${ZITADEL_HOST}/v2/organizations/_search" \
  -H "Authorization: Bearer ${ADMIN_PAT}" -H "Content-Type: application/json" \
  -d '{"query":{"limit":1000}}')"

echo "--- all organizations ---"
jq -r '.result[]? | "\(.id)\t\(.name)"' <<<"$resp"

# orgs known to be undeletable (server returns 500); retrying them every run
# would hammer the deployment under measurement
SKIP_FILE=".skip-orgs"
mapfile -t stale < <(jq -r --arg re "$PATTERN" \
  '.result[]? | select(.name | test($re)) | "\(.id)\t\(.name)"' <<<"$resp" \
  | { if [ -f "$SKIP_FILE" ]; then grep -vFf "$SKIP_FILE"; else cat; fi; })

echo "--- matching the benchmark convention (${#stale[@]}) ---"
printf '%s\n' "${stale[@]:-}" | sed '/^$/d'

if [ "${#stale[@]}" -eq 0 ]; then echo "nothing to clean"; exit 0; fi
if [ "$CONFIRM" != "--confirm" ]; then
  echo "dry run -- re-run with --confirm to delete the ${#stale[@]} org(s) above"; exit 0
fi

for row in "${stale[@]}"; do
  id="${row%%$'\t'*}"; name="${row#*$'\t'}"
  # belt and braces: re-check the name right before deleting
  if ! [[ "$name" =~ $PATTERN ]]; then echo "SKIP (pattern) $id $name"; continue; fi
  code="$(curl -s -m 30 -o /dev/null -w '%{http_code}' -X DELETE \
    "${ZITADEL_HOST}/v2/organizations/${id}" -H "Authorization: Bearer ${ADMIN_PAT}")"
  echo "deleted ${id} ${name} -> HTTP ${code}"
done
