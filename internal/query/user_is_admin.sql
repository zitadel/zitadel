-- This query returns `True` when at least one membership is found.
-- No rows are returned in case no membership is found.
-- The UNION ALLs (instead of regular UNION) prevents a high cost sort.
-- The LIMIT 1 terminates the query as soon as the first membership is found.
SELECT true AS b
    FROM projections.instance_members4
    WHERE instance_id = $1
    AND user_id = $2
UNION ALL
SELECT true AS b
    FROM projections.org_members4
    WHERE instance_id = $1
    AND user_id = $2
UNION ALL
SELECT true AS b
    FROM projections.project_members4
    WHERE instance_id = $1
    AND user_id = $2
UNION ALL
SELECT true AS b
    FROM projections.project_grant_members4
    WHERE instance_id = $1
    AND user_id = $2
LIMIT 1;
