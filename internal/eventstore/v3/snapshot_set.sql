INSERT INTO eventstore.snapshots(
	instance_id,
    snapshot_type,
    aggregate_id,
    "position",
    change_date,
    payload
)
SELECT $1, $2, $3, $4, $5, $6
WHERE  (
    -- only upsert snapshot with newer position
	SELECT "position"
	FROM eventstore.snapshots
	WHERE instance_id = $1
		AND snapshot_type = $2
		AND aggregate_id = $3
) < $4
ON CONFLICT (
    instance_id,
    snapshot_type,
    aggregate_id
)
DO UPDATE SET
    "position" = EXCLUDED."position" ,
    change_date = EXCLUDED.change_date,
    payload = EXCLUDED.payload
;
