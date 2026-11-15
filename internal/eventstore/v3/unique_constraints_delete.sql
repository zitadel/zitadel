-- Resolve keys before joining targets to avoid a tenant-wide cross product.
WITH matches AS MATERIALIZED (
    SELECT requested.instance_id, requested.unique_type, matched.unique_field
    FROM unnest($1::text[], $2::text[], $3::text[]) AS requested(instance_id, unique_type, unique_field)
    CROSS JOIN LATERAL (
        -- Prefer legacy case-sensitive rows before checking claims, even if the exact match is claimed.
        SELECT stored.unique_field
        FROM eventstore.unique_constraints AS stored
        WHERE stored.instance_id = requested.instance_id
            AND stored.unique_type = requested.unique_type
            AND stored.unique_field IN (requested.unique_field, LOWER(requested.unique_field))
        ORDER BY (stored.unique_field = requested.unique_field) DESC
        LIMIT 1
    ) AS matched
)
DELETE FROM eventstore.unique_constraints AS target
USING matches
WHERE target.instance_id = matches.instance_id
    AND target.unique_type = matches.unique_type
    AND target.unique_field = matches.unique_field
    -- only a claimed remove can release a claimed key
    AND NOT EXISTS (SELECT 1 FROM unnest(target.owners) AS owner WHERE owner LIKE 'claimant:%')
