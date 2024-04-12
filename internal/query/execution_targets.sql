WITH RECURSIVE
    dissolved_execution_targets(execution_id, resource_owner, instance_id, position, "include", "target")
        AS (SELECT execution_id
                 , resource_owner
                 , instance_id
                 , ARRAY [position]
                 , "include"
                 , "target"
            FROM matched_targets_and_includes
            UNION ALL
            SELECT e.execution_id
                 , resource_owner
                 , p.instance_id
                 , e.position || p.position
                 , p."include"
                 , p."target"
            FROM dissolved_execution_targets e
                     JOIN projections.executions1_targets p
                          ON e.instance_id = p.instance_id
                              AND e.resource_owner = p.resource_owner
                              AND e.include IS NOT NULL
                              AND e.include = p.execution_id),
    matched AS (SELECT *
                FROM projections.executions1
                WHERE instance_id = $1
                  AND resource_owner = $2
                  AND id @> $3
                ORDER BY id DESC
                LIMIT 1),
    matched_targets_and_includes AS (SELECT pos.*
                          FROM matched m
                                   JOIN
                               projections.executions1_targets pos
                               ON m.id = pos.execution_id
                                   AND m.resource_owner = pos.resource_owner
                                   AND m.instance_id = pos.instance_id
                          ORDER BY execution_id,
                                   position)
select *
FROM dissolved_execution_targets e
         JOIN projections.targets t
              ON e.instance_id = t.instance_id
                  AND e.resource_owner = t.resource_owner
                  AND e.target = t.id
WHERE "include" IS NULL
ORDER BY position DESC;
