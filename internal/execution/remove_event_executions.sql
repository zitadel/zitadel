DELETE
FROM projections.execution_handler
WHERE instance_id = $1
  AND resource_owner = $2
  AND aggregate_type = $3
  AND aggregate_version = $4
  AND aggregate_id = $5
  AND sequence = $6
  AND event_type = $7