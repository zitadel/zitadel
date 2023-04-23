WITH failures AS (
    SELECT 
        failure_count 
    FROM 
        projections.failed_events 
    WHERE 
        projection_name = $1 
        AND failed_sequence = $2 
        AND instance_id = $3
) SELECT COALESCE((SELECT failure_count FROM failures), 0) AS failure_count