DELETE FROM eventstore.unique_constraints AS target
USING unnest($1::TEXT[], $2::TEXT[]) AS requested(instance_id, owner)
WHERE target.instance_id = requested.instance_id
    AND target.owners @> ARRAY[requested.owner]
