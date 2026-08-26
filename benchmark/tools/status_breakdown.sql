-- Test-phase HTTP status / error-code breakdown, so a failure rate in the k6
-- summary can be diagnosed after the raw CSV is deleted.
SET memory_limit='4GB';
SET threads=6;
.mode box
SELECT status,
       coalesce(nullif(error_code,''),'-') AS error_code,
       coalesce(nullif(error,''),'-')      AS error,
       count(*)                            AS requests
FROM read_csv(getvariable('csv'), header=true, all_varchar=false)
WHERE metric_name = 'http_reqs' AND "group" IS NULL
GROUP BY ALL
ORDER BY requests DESC
LIMIT 30;

-- Non-2xx over time: distinguishes a steady failure rate (e.g. expiring
-- credentials) from bursts (e.g. a container being recycled).
SELECT strftime(to_timestamp(CAST(timestamp AS BIGINT)), '%H:%M') AS minute,
       status,
       count(*) AS requests
FROM read_csv(getvariable('csv'), header=true, all_varchar=false)
WHERE metric_name = 'http_reqs' AND "group" IS NULL
  AND (CAST(status AS INTEGER) >= 400 OR CAST(status AS INTEGER) = 0)
GROUP BY ALL
ORDER BY minute, status
LIMIT 60;
