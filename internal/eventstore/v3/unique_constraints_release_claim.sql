WITH requested AS MATERIALIZED (
    SELECT instance_id, unique_type, LOWER(unique_field) AS unique_field, claimant, ord
    FROM unnest($1::TEXT[], $2::TEXT[], $3::TEXT[], $4::TEXT[]) WITH ORDINALITY
        AS r(instance_id, unique_type, unique_field, claimant, ord)
), released AS (
    DELETE FROM eventstore.unique_constraints AS stored
    USING requested
    WHERE (stored.instance_id, stored.unique_type, stored.unique_field) =
          (requested.instance_id, requested.unique_type, requested.unique_field)
        AND stored.owners @> ARRAY[requested.claimant]
)
-- The select sees the snapshot before the delete. Like releasing one by one in
-- command order, a release fails if another claimant holds the key and no
-- release up to it in the same push names that claimant. A missing row counts
-- as released.
SELECT NOT EXISTS (
    SELECT 1
    FROM requested
    JOIN eventstore.unique_constraints AS stored
        ON (stored.instance_id, stored.unique_type, stored.unique_field) =
           (requested.instance_id, requested.unique_type, requested.unique_field)
    WHERE EXISTS (SELECT 1 FROM unnest(stored.owners) AS owner WHERE owner LIKE 'claimant:%')
        AND NOT EXISTS (
            SELECT 1 FROM requested AS earlier
            WHERE earlier.ord <= requested.ord
                AND (earlier.instance_id, earlier.unique_type, earlier.unique_field) =
                    (stored.instance_id, stored.unique_type, stored.unique_field)
                AND stored.owners @> ARRAY[earlier.claimant]
        )
)
