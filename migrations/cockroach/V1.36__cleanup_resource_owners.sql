WITH resource_owners AS (
	WITH first_sequences AS (
		WITH duplicates AS (
			select 
				aggregate_id,
				aggregate_type
			from (
				select 
					aggregate_type, 
					aggregate_id, 
					resource_owner 
				from events
				where aggregate_type not like '%.%'
				group by 
					aggregate_type, 
					aggregate_id, 
					resource_owner
			) group by 
				aggregate_id, 
				aggregate_type 
			having count(resource_owner) > 1
		)
		SELECT 
			MIN(event_sequence) AS seq,
			aggregate_type,
			aggregate_id
		FROM 
			eventstore.events 
		WHERE 
			aggregate_id IN (select aggregate_id from duplicates)
		GROUP BY 
			aggregate_type,
			aggregate_id
		ORDER BY aggregate_id, seq
	)
	SELECT 
		f.*, 
		e.resource_owner
	FROM 
		first_sequences f
	JOIN
		events e ON f.seq = e.event_sequence
)
UPDATE events e 
SET resource_owner = r.resource_owner
FROM resource_owners r where e.aggregate_id = r.aggregate_id;
