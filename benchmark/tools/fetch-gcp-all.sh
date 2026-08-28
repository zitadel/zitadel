#!/usr/bin/env bash
# Everything the v4.17.1 pages still need from Google Cloud, in one run.
#
#   bash tools/fetch-gcp-all.sh
#
# Run in Cloud Shell. Produces one tarball to send back. Individual steps that fail
# are reported and skipped rather than taking the whole run down with them.
set -uo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PROJECT="${PROJECT:-zitadel-cloud}"
SERVICE="${SERVICE:-zitadel-qa-us1-cr-us-central1}"
DBID="${DBID:-zitadel-cloud:us1}"

STAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUT="${PWD}/bench-gcp-${STAMP}"
mkdir -p "${OUT}/sweep" "${OUT}/bursts" "${OUT}/logs"
FAILED=()
step() { echo; echo "=== $* ==="; }
note() { echo "  !! $*"; FAILED+=("$*"); }

step "preflight"
if ! gcloud auth print-access-token >/dev/null 2>&1; then
  echo "  not authenticated. In Cloud Shell this should not happen; run: gcloud auth login" >&2
  exit 1
fi
if ! gcloud run services describe "$SERVICE" --project="$PROJECT" --region=us-central1 \
      --format='value(metadata.name)' >/dev/null 2>&1; then
  echo "  Cloud Run service '${SERVICE}' not found in ${PROJECT}. Available:" >&2
  gcloud run services list --project="$PROJECT" --format='value(metadata.name)' >&2 || true
  echo "  Re-run as: SERVICE=<name> bash tools/fetch-gcp-all.sh" >&2
  exit 1
fi
echo "  project=${PROJECT} service=${SERVICE} db=${DBID}"
echo "  output=${OUT}"

# ---------------------------------------------------------------- 1. the eleven windows
step "1/4  Cloud Run + Cloud SQL metrics, 11 test windows"
( cd "${OUT}/sweep" && SERVICE="$SERVICE" PROJECT="$PROJECT" DBID="$DBID" \
    bash "${HERE}/fetch-gcp-metrics.sh" ) || note "sweep metrics failed"

step "2/4  Query Insights ranking + autovacuum, 11 test windows"
( cd "${OUT}/sweep" && PROJECT="$PROJECT" DBID="$DBID" \
    bash "${HERE}/fetch-gcp-queries.sh" ) || note "sweep queries failed"

# ---------------------------------------------------------------- 2. the two 503 bursts
step "3/4  metrics across the two 503 bursts"
cat > "${OUT}/bursts/windows.tsv" <<'EOF'
burst1_manipulate_user|2026-08-27T13:38:00Z|2026-08-27T13:50:00Z
burst2_human_password_login|2026-08-27T19:30:00Z|2026-08-27T19:50:00Z
EOF
( cd "${OUT}/bursts" && WINDOWS_FILE="${OUT}/bursts/windows.tsv" \
    SERVICE="$SERVICE" PROJECT="$PROJECT" DBID="$DBID" \
    bash "${HERE}/fetch-gcp-metrics.sh" ) || note "burst metrics failed"

# ---------------------------------------------------------------- 3. the log side
step "4/4  Cloud Run logs around each burst"
RES="resource.type=\"cloud_run_revision\" AND resource.labels.service_name=\"${SERVICE}\""

logq() { # logq <outfile> <filter> [extra-format]
  local f="$1" filter="$2" fmt="${3:-value(timestamp,severity,textPayload,jsonPayload.message,jsonPayload.error)}"
  gcloud logging read "$filter" --project="$PROJECT" --limit="${4:-5000}" --order=asc \
    --format="$fmt" > "${OUT}/logs/${f}" 2>"${OUT}/logs/${f}.err" \
    || note "log query ${f} failed (see ${f}.err)"
  local n; n="$(wc -l < "${OUT}/logs/${f}")"
  printf '  %-34s %s lines%s\n' "$f" "$n" "$( [ "$n" -ge "${4:-5000}" ] && printf ' *** TRUNCATED AT LIMIT ***' )"
}

for b in "burst1 2026-08-27T13:38:00Z 2026-08-27T13:50:00Z" \
         "burst2 2026-08-27T19:30:00Z 2026-08-27T19:50:00Z"; do
  set -- $b; name="$1"; s="$2"; e="$3"
  W="AND timestamp>=\"${s}\" AND timestamp<=\"${e}\""
  logq "${name}_panics.txt" \
    "${RES} ${W} AND (textPayload:\"panic\" OR textPayload:\"fatal error\" OR textPayload:\"goroutine\" OR jsonPayload.message:\"panic\" OR jsonPayload.message:\"fatal\" OR jsonPayload.error:\"panic\")"
  logq "${name}_restarts.txt" \
    "${RES} ${W} AND (textPayload:\"STARTUP\" OR textPayload:\"Container called exit\" OR textPayload:\"Memory limit\" OR textPayload:\"terminated\" OR jsonPayload.message:\"STARTUP\" OR jsonPayload.message:\"terminated\" OR severity>=WARNING)" \
    'value(timestamp,severity,labels."run.googleapis.com/instance_id",textPayload,jsonPayload.message,jsonPayload.error)'
  logq "${name}_503s.txt" \
    "${RES} ${W} AND httpRequest.status=503" \
    'table(timestamp,httpRequest.status,httpRequest.requestUrl)'
  logq "${name}_deploys.txt" \
    "protoPayload.serviceName=\"run.googleapis.com\" AND resource.labels.service_name=\"${SERVICE}\" ${W}" \
    'table(timestamp,protoPayload.methodName,protoPayload.authenticationInfo.principalEmail)' 50
done

# The question that needs both windows at once: is the same container implicated twice?
step "instance ids across the whole sweep envelope"
gcloud logging read \
  "${RES} AND timestamp>=\"2026-08-27T13:30:00Z\" AND timestamp<=\"2026-08-27T20:40:00Z\" AND (textPayload:\"STARTUP\" OR textPayload:\"Container called exit\" OR textPayload:\"Memory limit\" OR textPayload:\"terminated\")" \
  --project="$PROJECT" --limit=1000 --order=asc \
  --format='value(timestamp,labels."run.googleapis.com/instance_id",textPayload)' \
  > "${OUT}/logs/instance_ids_raw.txt" 2>"${OUT}/logs/instance_ids_raw.err" \
  || note "instance id sweep failed"
awk -F'\t' '{print $2}' "${OUT}/logs/instance_ids_raw.txt" 2>/dev/null \
  | sort | uniq -c | sort -rn > "${OUT}/logs/instance_ids_tally.txt"
echo "  distinct instance ids: $(wc -l < "${OUT}/logs/instance_ids_tally.txt")  (baseline is 7)"

# ---------------------------------------------------------------- 4. wrap up
step "done"
TAR="bench-gcp-${STAMP}.tar.gz"
tar -czf "$TAR" -C "$(dirname "$OUT")" "$(basename "$OUT")"
echo "  ${TAR}  ($(du -h "$TAR" | cut -f1))"
if [ "${#FAILED[@]}" -gt 0 ]; then
  echo; echo "  ${#FAILED[@]} step(s) did not complete:"
  printf '   - %s\n' "${FAILED[@]}"
  echo "  The tarball still contains everything that did. Send it anyway."
fi
echo
echo "Send back: ${TAR}"
