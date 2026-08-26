-- k6 CSV -> per-second latency percentiles for the benchmark docs.
-- Only test-phase custom trend metrics: rows written during setup carry
-- group='::setup', so filtering on an empty group drops the whole setup phase,
-- and the k6 built-ins are excluded by name.
SET memory_limit='4GB';
SET threads=6;
.mode json
SELECT
    metric_name,
    to_timestamp(CAST(timestamp AS BIGINT)) AS timestamp,
    quantile_cont(metric_value, 0.5)  AS p50,
    quantile_cont(metric_value, 0.95) AS p95,
    quantile_cont(metric_value, 0.99) AS p99
FROM read_csv(getvariable('csv'), header=true, all_varchar=false)
WHERE "group" IS NULL
  AND metric_name NOT IN (
      'checks','data_received','data_sent','http_req_blocked','http_req_connecting',
      'http_req_duration','http_req_failed','http_req_receiving','http_req_sending',
      'http_req_tls_handshaking','http_req_waiting','http_reqs',
      'iteration_duration','iterations','vus','vus_max'
  )
GROUP BY metric_name, timestamp
ORDER BY timestamp, metric_name;
