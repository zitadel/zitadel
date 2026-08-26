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
