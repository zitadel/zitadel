SELECT et.instance_id,
       et.execution_id,
       JSONB_AGG(
               JSON_OBJECT(
                       'position' : et.position,
                       'include' : et.include,
                       'target' : et.target_id
               )
       ) as targets
FROM projections.executions1_targets AS et
         INNER JOIN projections.targets2 AS t
                    ON et.instance_id = t.instance_id
                        AND et.target_id IS NOT NULL
                        AND et.target_id = t.id
GROUP BY et.instance_id, et.execution_id