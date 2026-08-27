#!/usr/bin/env bash
# Pull Cloud Run + Cloud SQL metrics for the v4.17.1 benchmark windows.
# Run in Cloud Shell (already authenticated). Emits:
#   bench-metrics-raw.csv  -> target,metric,timestamp,value   (60s aligned)
#   bench-metrics.txt      -> per-target median/max summary   (paste this back)
set -uo pipefail

PROJECT="${PROJECT:-zitadel-cloud}"
SERVICE="${SERVICE:-}"                       # Cloud Run service name -- MUST be set
DBID="${DBID:-zitadel-cloud:us1}"            # Cloud SQL database_id = project:instance
ALIGN="${ALIGN:-60s}"

TOK="$(gcloud auth print-access-token)"
API="https://monitoring.googleapis.com/v3/projects/${PROJECT}/timeSeries"

if [ -z "$SERVICE" ]; then
  echo "SERVICE is unset. Cloud Run services in ${PROJECT}:" >&2
  gcloud run services list --project="$PROJECT" --format='value(metadata.name,status.url)' >&2 || true
  echo "Re-run as:  SERVICE=<name> bash $0" >&2
  exit 1
fi

WINDOWS="
add_session|2026-08-26T20:12:34Z|2026-08-26T20:42:35Z
introspect|2026-08-26T20:44:10Z|2026-08-26T21:14:10Z
machine_client_credentials_login|2026-08-26T21:15:46Z|2026-08-26T21:45:47Z
human_password_login|2026-08-26T22:44:22Z|2026-08-26T23:14:31Z
machine_jwt_profile_grant|2026-08-26T23:15:24Z|2026-08-26T23:45:25Z
machine_pat_login|2026-08-26T23:46:24Z|2026-08-27T00:16:24Z
oidc_session|2026-08-27T00:18:30Z|2026-08-27T00:48:32Z
otp_session|2026-08-27T00:49:39Z|2026-08-27T01:19:45Z
user_info|2026-08-27T01:21:07Z|2026-08-27T01:51:07Z
password_session|2026-08-27T01:56:56Z|2026-08-27T02:27:03Z
manipulate_user|2026-08-27T11:09:50Z|2026-08-27T11:40:02Z
"

RUN_RES='resource.type="cloud_run_revision" AND resource.labels.service_name="'"$SERVICE"'"'
SQL_RES='resource.type="cloudsql_database" AND resource.labels.database_id="'"$DBID"'"'

# label | filter | perSeriesAligner | crossSeriesReducer
METRICS="
zit_cpu_p50|metric.type=\"run.googleapis.com/container/cpu/utilizations\" AND ${RUN_RES}|ALIGN_PERCENTILE_50|REDUCE_MEAN
zit_cpu_p99|metric.type=\"run.googleapis.com/container/cpu/utilizations\" AND ${RUN_RES}|ALIGN_PERCENTILE_99|REDUCE_MAX
zit_mem_p50|metric.type=\"run.googleapis.com/container/memory/utilizations\" AND ${RUN_RES}|ALIGN_PERCENTILE_50|REDUCE_MEAN
zit_mem_p99|metric.type=\"run.googleapis.com/container/memory/utilizations\" AND ${RUN_RES}|ALIGN_PERCENTILE_99|REDUCE_MAX
zit_instances|metric.type=\"run.googleapis.com/container/instance_count\" AND ${RUN_RES}|ALIGN_MEAN|REDUCE_SUM
zit_req_per_s|metric.type=\"run.googleapis.com/request_count\" AND ${RUN_RES}|ALIGN_RATE|REDUCE_SUM
zit_lat_p50_ms|metric.type=\"run.googleapis.com/request_latencies\" AND ${RUN_RES}|ALIGN_PERCENTILE_50|REDUCE_MEAN
zit_lat_p95_ms|metric.type=\"run.googleapis.com/request_latencies\" AND ${RUN_RES}|ALIGN_PERCENTILE_95|REDUCE_MEAN
zit_lat_p99_ms|metric.type=\"run.googleapis.com/request_latencies\" AND ${RUN_RES}|ALIGN_PERCENTILE_99|REDUCE_MEAN
zit_startups|metric.type=\"run.googleapis.com/container/startup_latencies\" AND ${RUN_RES}|ALIGN_COUNT|REDUCE_SUM
db_cpu|metric.type=\"cloudsql.googleapis.com/database/cpu/utilization\" AND ${SQL_RES}|ALIGN_MEAN|REDUCE_MEAN
db_mem|metric.type=\"cloudsql.googleapis.com/database/memory/utilization\" AND ${SQL_RES}|ALIGN_MEAN|REDUCE_MEAN
db_connections|metric.type=\"cloudsql.googleapis.com/database/postgresql/num_backends\" AND ${SQL_RES}|ALIGN_MEAN|REDUCE_SUM
db_tx_per_s|metric.type=\"cloudsql.googleapis.com/database/postgresql/transaction_count\" AND ${SQL_RES}|ALIGN_RATE|REDUCE_SUM
db_read_iops|metric.type=\"cloudsql.googleapis.com/database/disk/read_ops_count\" AND ${SQL_RES}|ALIGN_RATE|REDUCE_SUM
db_write_iops|metric.type=\"cloudsql.googleapis.com/database/disk/write_ops_count\" AND ${SQL_RES}|ALIGN_RATE|REDUCE_SUM
"

RAW=bench-metrics-raw.csv
OUT=bench-metrics.txt
echo "target,metric,timestamp,value" > "$RAW"
: > "$OUT"

{
  echo "project=${PROJECT} service=${SERVICE} db=${DBID} align=${ALIGN}"
  echo "generated for zitadel v4.17.1 benchmark windows (all times UTC)"
  echo
} >> "$OUT"

while IFS='|' read -r target start end; do
  [ -z "${target:-}" ] && continue
  echo "== $target  $start -> $end" | tee -a "$OUT"
  while IFS='|' read -r label filter aligner reducer; do
    [ -z "${label:-}" ] && continue
    resp="$(curl -sS -G "$API" \
      -H "Authorization: Bearer ${TOK}" \
      --data-urlencode "filter=${filter}" \
      --data-urlencode "interval.startTime=${start}" \
      --data-urlencode "interval.endTime=${end}" \
      --data-urlencode "aggregation.alignmentPeriod=${ALIGN}" \
      --data-urlencode "aggregation.perSeriesAligner=${aligner}" \
      --data-urlencode "aggregation.crossSeriesReducer=${reducer}")"

    if echo "$resp" | jq -e '.error' >/dev/null 2>&1; then
      printf '   %-16s ERROR: %s\n' "$label" "$(echo "$resp" | jq -r '.error.message' | head -c 160)" | tee -a "$OUT"
      continue
    fi

    echo "$resp" | jq -r --arg t "$target" --arg m "$label" \
      '.timeSeries[]?.points[]? | [$t,$m,.interval.endTime,(.value.doubleValue // .value.int64Value // .value.distributionValue.mean // empty)] | @csv' \
      >> "$RAW"

    stats="$(echo "$resp" | jq -r '
      [ .timeSeries[]?.points[]? | (.value.doubleValue // (.value.int64Value|tonumber?) // empty) ] | sort
      | if length==0 then "no data"
        else "n=\(length) median=\(.[(length/2)|floor]) max=\(.[length-1]) min=\(.[0])" end')"
    printf '   %-16s %s\n' "$label" "$stats" | tee -a "$OUT"
  done <<< "$METRICS"
  echo | tee -a "$OUT"
done <<< "$WINDOWS"

echo "Wrote $OUT (summary) and $RAW (60s series)."
