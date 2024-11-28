WITH RECURSIVE
    matched AS ((SELECT *
                 FROM projections.executions1
                 WHERE instance_id = $1
                   AND id = ANY($2)
                 ORDER BY id DESC
                 LIMIT 1)
                UNION ALL
                (SELECT *
                 FROM projections.executions1
                 WHERE instance_id = $1
                   AND id = ANY($3)
                 ORDER BY id DESC
                 LIMIT 1)),
    matched_targets_and_includes AS (SELECT pos.*
                                     FROM matched m
                                              JOIN
                                          projections.executions1_targets pos
                                          ON m.id = pos.execution_id
                                              AND m.instance_id = pos.instance_id
                                     ORDER BY execution_id,
                                              position),
    dissolved_execution_targets(execution_id, instance_id, position, "include", "target_id")
        AS (SELECT execution_id
                 , instance_id
                 , ARRAY [position]
                 , "include"
                 , "target_id"
            FROM matched_targets_and_includes
            UNION ALL
            SELECT e.execution_id
                 , p.instance_id
                 , e.position || p.position
                 , p."include"
                 , p."target_id"
            FROM dissolved_execution_targets e
                     JOIN projections.executions1_targets p
                          ON e.instance_id = p.instance_id
                              AND e.include IS NOT NULL
                              AND e.include = p.execution_id)
select e.execution_id, e.instance_id, e.target_id, t.target_type, t.endpoint, t.timeout, t.interrupt_on_error, t.signing_key
FROM dissolved_execution_targets e
         JOIN projections.targets2 t
              ON e.instance_id = t.instance_id
                  AND e.target_id = t.id
WHERE "include" = ''
ORDER BY position DESC;
