WITH failures AS (
    SELECT 
        failure_count 
    FROM 
        projections.failed_events2
    WHERE 
        projection_name = $1
        AND instance_id = $2
        AND aggregate_type = $3
        AND aggregate_id = $4
        AND failed_sequence = $5
) SELECT COALESCE((SELECT failure_count FROM failures), 0) AS failure_count