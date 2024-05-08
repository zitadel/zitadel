SELECT instance_id,
       execution_id,
       JSONB_AGG(
               JSON_OBJECT(
                       'position' : position,
                       'include' : include,
                       'target' : target_id
                   )
           ) as targets
FROM projections.executions1_targets
GROUP BY instance_id, execution_id