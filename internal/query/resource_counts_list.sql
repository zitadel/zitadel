SELECT id,
    instance_id,
    parent_type,
    parent_id,
    resource_name,
    updated_at,
    amount
FROM projections.resource_counts
WHERE id > $1
ORDER BY id
LIMIT $2;
