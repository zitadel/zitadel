SELECT instance_id,
       resource_owner,
       aggregate_type,
       aggregate_version,
       aggregate_id,
       sequence,
       event_type,
       created_at,
       user_id,
       event_data,
       targets_data
FROM projections.execution_handler
WHERE instance_id = $1
LIMIT $2
FOR UPDATE SKIP LOCKED