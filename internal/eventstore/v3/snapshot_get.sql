SELECT "position", change_date, payload
FROM eventstore.snapshots
WHERE instance_id = $1
    AND snapshot_type = $2
    AND aggregate_id = $3
;
