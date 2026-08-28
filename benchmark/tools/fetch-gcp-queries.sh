#!/usr/bin/env bash
# Query Insights ranking + autovacuum log stats for the v4.17.1 benchmark windows.
# Run in Cloud Shell. Emits bench-queries.txt and bench-autovacuum.txt
set -uo pipefail

PROJECT="${PROJECT:-zitadel-cloud}"
DBID="${DBID:-zitadel-cloud:us1}"
INSTANCE="${INSTANCE:-us1}"
TOPN="${TOPN:-5}"

TOK="$(gcloud auth print-access-token)"
API="https://monitoring.googleapis.com/v3/projects/${PROJECT}"

# Test windows. Baked in so the script runs standalone in Cloud Shell, where this repo
# is not checked out. Regenerate after every republish with tools/sync-windows.sh.
# Override for a one-off window (e.g. a 503 post-mortem) with WINDOWS_FILE=/path.
WINDOWS="${WINDOWS_DEFAULT:-
add_session|2026-08-27T14:58:50Z|2026-08-27T15:28:51Z
human_password_login|2026-08-27T19:24:36Z|2026-08-27T19:54:52Z
introspect|2026-08-27T15:31:22Z|2026-08-27T16:01:22Z
machine_client_credentials_login|2026-08-27T16:03:44Z|2026-08-27T16:33:45Z
machine_jwt_profile_grant|2026-08-27T16:34:57Z|2026-08-27T17:04:57Z
machine_pat_login|2026-08-27T17:06:10Z|2026-08-27T17:36:10Z
manipulate_user|2026-08-27T13:41:19Z|2026-08-27T14:11:30Z
oidc_session|2026-08-27T17:38:43Z|2026-08-27T18:08:45Z
otp_session|2026-08-27T18:49:09Z|2026-08-27T19:19:16Z
password_session|2026-08-27T20:02:39Z|2026-08-27T20:32:49Z
user_info|2026-08-27T18:11:39Z|2026-08-27T18:41:40Z
}"
[ -n "${WINDOWS_FILE:-}" ] && WINDOWS="$(cat "${WINDOWS_FILE}")"

##############################################################################
# PART 0 -- is Query Insights on, and what labels does the metric carry?
##############################################################################
# Probe with the first configured window, not a frozen one: a stale probe window whose
# Insights data has aged out reports "not enabled" and skips the whole of Part 1.
# Sweep envelope for the by-table tally: earliest start to latest end across all windows.
SWEEP_START="$(printf '%s\n' "$WINDOWS" | sed '/^$/d' | cut -d"|" -f2 | sort | head -1)"
SWEEP_END="$(printf '%s\n' "$WINDOWS" | sed '/^$/d' | cut -d"|" -f3 | sort | tail -1)"
PROBE_START="$(printf '%s\n' "$WINDOWS" | sed '/^$/d' | head -1 | cut -d"|" -f2)"
PROBE_END="$(printf '%s\n' "$WINDOWS" | sed '/^$/d' | head -1 | cut -d"|" -f3)"
QOUT=bench-queries.txt
: > "$QOUT"
echo "### Part 0: metric descriptor discovery" | tee -a "$QOUT"
for MT in \
  "cloudsql.googleapis.com/database/postgresql/insights/perquery/execution_time" \
  "cloudsql.googleapis.com/database/postgresql/insights/perquery/latencies" \
  "cloudsql.googleapis.com/database/postgresql/insights/aggregate/execution_time" ; do
  desc="$(curl -sS -H "Authorization: Bearer ${TOK}" \
      "${API}/metricDescriptors/${MT}")"
  if echo "$desc" | jq -e '.error' >/dev/null 2>&1; then
    echo "  $MT -> NOT AVAILABLE: $(echo "$desc" | jq -r '.error.message'|head -c 120)" | tee -a "$QOUT"
  else
    echo "  $MT" | tee -a "$QOUT"
    echo "     kind=$(echo "$desc"|jq -r '.metricKind') type=$(echo "$desc"|jq -r '.valueType') unit=$(echo "$desc"|jq -r '.unit // "-"')" | tee -a "$QOUT"
    echo "     labels: $(echo "$desc"|jq -r '[.labels[]?.key]|join(", ")')" | tee -a "$QOUT"
  fi
done
echo | tee -a "$QOUT"

MT="cloudsql.googleapis.com/database/postgresql/insights/perquery/execution_time"

##############################################################################
# PART 1 -- top N query hashes per window, ranked by TOTAL execution time
##############################################################################
echo "### Part 1: top ${TOPN} queries per window, by total execution time" | tee -a "$QOUT"
for RTYPE in cloudsql_instance_database cloudsql_database; do
  probe="$(curl -sS -G "${API}/timeSeries" -H "Authorization: Bearer ${TOK}" \
    --data-urlencode "filter=metric.type=\"${MT}\" AND resource.type=\"${RTYPE}\"" \
    --data-urlencode "interval.startTime=${PROBE_START}" \
    --data-urlencode "interval.endTime=${PROBE_END}" \
    --data-urlencode "aggregation.alignmentPeriod=1800s" \
    --data-urlencode "aggregation.perSeriesAligner=ALIGN_DELTA")"
  if echo "$probe" | jq -e '.timeSeries|length>0' >/dev/null 2>&1; then
    echo "  using resource.type=${RTYPE}" | tee -a "$QOUT"
    echo "  available metric labels on returned series: $(echo "$probe"|jq -r '[.timeSeries[0].metric.labels|keys[]]|join(", ")')" | tee -a "$QOUT"
    break
  fi
  RTYPE=""
done
if [ -z "${RTYPE:-}" ]; then
  echo "  !! No perquery time series returned. Query Insights is probably NOT enabled" | tee -a "$QOUT"
  echo "     on ${DBID}. This data cannot be recovered retroactively." | tee -a "$QOUT"
else
  while IFS='|' read -r target start end; do
    [ -z "${target:-}" ] && continue
    echo | tee -a "$QOUT"
    echo "== $target  $start -> $end" | tee -a "$QOUT"
    curl -sS -G "${API}/timeSeries" -H "Authorization: Bearer ${TOK}" \
      --data-urlencode "filter=metric.type=\"${MT}\" AND resource.type=\"${RTYPE}\"" \
      --data-urlencode "interval.startTime=${start}" \
      --data-urlencode "interval.endTime=${end}" \
      --data-urlencode "aggregation.alignmentPeriod=1800s" \
      --data-urlencode "aggregation.perSeriesAligner=ALIGN_DELTA" \
      --data-urlencode "aggregation.crossSeriesReducer=REDUCE_SUM" \
      --data-urlencode "aggregation.groupByFields=metric.label.query_hash" \
      --data-urlencode "aggregation.groupByFields=metric.label.querystring" \
      --data-urlencode "aggregation.groupByFields=metric.label.user" \
    | tee "/tmp/qi_${target}.json" > /dev/null
    jq -r --argjson n "$TOPN" '
        [ .timeSeries[]? | {
            hash: (.metric.labels.query_hash // "?"),
            q:    (.metric.labels.querystring // .metric.labels.query_string // null),
            tot:  ([.points[]?.value.doubleValue // (.points[]?.value.int64Value|tonumber?) // 0]|add)
          } ]
        | sort_by(-.tot) | .[:$n]
        | to_entries[] | "   \(.key+1). total=\(.value.tot)  hash=\(.value.hash)\(if .value.q then "\n      SQL: \(.value.q|gsub("\\s+";" ")|.[0:400])" else "\n      SQL: <no querystring label returned>" end)"
      ' "/tmp/qi_${target}.json" | tee -a "$QOUT"
  done <<< "$WINDOWS"
fi

##############################################################################
# PART 2 -- autovacuum / autoanalyze on eventstore.events2 (v4.17 feature)
##############################################################################
VOUT=bench-autovacuum.txt
: > "$VOUT"
echo "### autovacuum / autoanalyze of eventstore.events2 per window" | tee -a "$VOUT"
BASE='resource.type="cloudsql_database" AND resource.labels.database_id="'"$DBID"'"'

while IFS='|' read -r target start end; do
  [ -z "${target:-}" ] && continue
  echo | tee -a "$VOUT"
  echo "== $target  $start -> $end" | tee -a "$VOUT"
  gcloud logging read \
    "${BASE} AND textPayload:(\"automatic vacuum of table\" OR \"automatic analyze of table\") AND textPayload:\"eventstore.events2\" AND timestamp>=\"${start}\" AND timestamp<=\"${end}\"" \
    --project="$PROJECT" --limit=1000 --format='value(timestamp,textPayload)' 2>/dev/null \
  | tee "/tmp/av_${target}.raw" \
  | awk -v OFS='' '
      /automatic vacuum/  {v++; kind="vacuum"}
      /automatic analyze/ {a++; kind="analyze"}
      match($0, /elapsed: [0-9.]+ s/) {
        e=substr($0, RSTART+9, RLENGTH-11)+0
        if (kind=="vacuum") { ve+=e; if(e>vm) vm=e } else { ae+=e; if(e>am) am=e }
      }
      END {
        printf "   vacuum : runs=%d  total_elapsed=%.1fs  max=%.1fs\n", v+0, ve+0, vm+0
        printf "   analyze: runs=%d  total_elapsed=%.1fs  max=%.1fs\n", a+0, ae+0, am+0
      }' | tee -a "$VOUT"
done <<< "$WINDOWS"

echo | tee -a "$VOUT"
echo "### sanity check: ALL autovacuum log lines in the sweep, by table" | tee -a "$VOUT"
gcloud logging read \
  "${BASE} AND textPayload:(\"automatic vacuum of table\" OR \"automatic analyze of table\") AND timestamp>=\"${SWEEP_START}\" AND timestamp<=\"${SWEEP_END}\"" \
  --project="$PROJECT" --limit=5000 --format='value(textPayload)' 2>/dev/null \
| grep -oE '(vacuum|analyze) of table "[^"]+"' | sort | uniq -c | sort -rn | head -20 | tee -a "$VOUT"

echo
# raw per-event CSV for correlation against the db_cpu series
echo "target,timestamp,kind,elapsed_s" > bench-autovacuum-events.csv
for f in /tmp/av_*.raw; do
  t="$(basename "$f" .raw)"; t="${t#av_}"
  awk -v t="$t" -F'\t' '{
      kind = /automatic vacuum/ ? "vacuum" : (/automatic analyze/ ? "analyze" : "")
      if (kind=="" ) next
      e=""; if (match($0, /elapsed: [0-9.]+ s/)) e=substr($0, RSTART+9, RLENGTH-11)
      print t "," $1 "," kind "," e
    }' "$f" >> bench-autovacuum-events.csv
done
jq -s '.' /tmp/qi_*.json > bench-queries.json 2>/dev/null || true
echo "Wrote $QOUT, $VOUT, bench-queries.json and bench-autovacuum-events.csv"
echo "NOTE: if Part 2 is empty everywhere, check log_autovacuum_min_duration on ${INSTANCE}."
echo "      At the default of -1 Postgres logs nothing and no filter can recover it."
