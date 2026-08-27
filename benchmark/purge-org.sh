#!/usr/bin/env bash
# Delete a benchmark organization that is too large to remove in one call.
#
# Org deletion sends every child object to the pusher in a single call, which
# exceeds the maximum argument size once an org has churned enough users, and
# fails with a 500. Delete the users first, a page at a time, then the org.
#
# Usage: purge-org.sh <org-id> <exact-name> [--confirm]
# SAFETY: the org must exist, its name must match the argument exactly, and that
# name must match load-test-<ISO timestamp>. Every user is re-checked to belong to
# this org immediately before deletion. Dry-run unless --confirm is passed.
set -euo pipefail
cd "$(dirname "$0")"
set -a; . ./.env; set +a
PATTERN='^load-test-[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3}Z$'
PAGE=1000
PARALLEL=${PARALLEL:-25}
ORG="$1"; EXPECTED="$2"; CONFIRM="${3:-}"

if pgrep -x k6 >/dev/null; then
  echo "REFUSING: a k6 run is in flight" >&2; exit 1
fi

name="$(curl -sf -m 30 -X POST "${ZITADEL_HOST}/v2/organizations/_search" \
  -H "Authorization: Bearer ${ADMIN_PAT}" -H "Content-Type: application/json" \
  -d '{"query":{"limit":1000}}' | jq -r --arg id "$ORG" '.result[]? | select(.id==$id) | .name')"
[ -n "$name" ]            || { echo "ABORT: org $ORG not found"; exit 1; }
[ "$name" = "$EXPECTED" ] || { echo "ABORT: name mismatch: '$name' != '$EXPECTED'"; exit 1; }
[[ "$name" =~ $PATTERN ]] || { echo "ABORT: '$name' is not a benchmark org"; exit 1; }

count() {
  curl -sf -m 60 -X POST "${ZITADEL_HOST}/v2/users" \
    -H "Authorization: Bearer ${ADMIN_PAT}" -H "Content-Type: application/json" \
    -d "{\"query\":{\"limit\":1},\"queries\":[{\"organizationIdQuery\":{\"organizationId\":\"${ORG}\"}}]}" \
    | jq -r '.details.totalResult // 0'
}

total="$(count)"
echo "org ${ORG} (${name}) holds ${total} user(s)"
if [ "$CONFIRM" != "--confirm" ]; then
  echo "dry run -- re-run with --confirm to delete these users and then the org"; exit 0
fi

export ZITADEL_HOST ADMIN_PAT ORG
del_one() {
  id="$1"
  # re-check ownership immediately before deleting
  owner="$(curl -sf -m 30 "${ZITADEL_HOST}/v2/users/${id}" \
    -H "Authorization: Bearer ${ADMIN_PAT}" | jq -r '.user.details.resourceOwner // ""')" || return 0
  [ "$owner" = "${ORG}" ] || { echo "SKIP ${id} (owner ${owner:-unknown})"; return 0; }
  curl -sf -m 60 -o /dev/null -X DELETE "${ZITADEL_HOST}/v2/users/${id}" \
    -H "Authorization: Bearer ${ADMIN_PAT}" || echo "FAILED ${id}"
}
export -f del_one

page=0
while :; do
  ids="$(curl -sf -m 60 -X POST "${ZITADEL_HOST}/v2/users" \
    -H "Authorization: Bearer ${ADMIN_PAT}" -H "Content-Type: application/json" \
    -d "{\"query\":{\"limit\":${PAGE}},\"queries\":[{\"organizationIdQuery\":{\"organizationId\":\"${ORG}\"}}]}" \
    | jq -r '.result[]?.userId')"
  [ -z "$ids" ] && break
  n=$(wc -l <<<"$ids")
  page=$((page+1))
  echo "page ${page}: deleting ${n} user(s) with ${PARALLEL} in flight ..."
  xargs -P "${PARALLEL}" -I{} bash -c 'del_one "$@"' _ {} <<<"$ids"
  left="$(count)"
  echo "page ${page} done, ${left} user(s) left"
  [ "$left" -eq 0 ] && break
done

echo "deleting org ${ORG} ..."
code="$(curl -s -m 120 -o /tmp/purge-org.out -w '%{http_code}' -X DELETE \
  "${ZITADEL_HOST}/v2/organizations/${ORG}" -H "Authorization: Bearer ${ADMIN_PAT}")"
echo "org delete -> HTTP ${code}"
[ "$code" = "200" ] || { head -c 300 /tmp/purge-org.out; echo; exit 1; }
