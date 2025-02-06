WITH target_rows AS (
	SELECT instance_id, aggregate_type, aggregate_id, sequence
	FROM eventstore.events2
	WHERE (instance_id, aggregate_type, aggregate_id, sequence) > ('', '', '', 0)
	ORDER BY instance_id, aggregate_type, aggregate_id, sequence
	LIMIT 1000
), u AS (
	UPDATE eventstore.events2
	SET in_tx_order_tmp = in_tx_order
	WHERE (instance_id,	aggregate_type,	aggregate_id, sequence) IN (
		SELECT instance_id,	aggregate_type,	aggregate_id, sequence
		FROM target_rows
	)
	RETURNING *
), n AS (
	SELECT count(*) FROM u
)
SELECT target_rows.*, n.* FROM target_rows, n
OFFSET 999;
