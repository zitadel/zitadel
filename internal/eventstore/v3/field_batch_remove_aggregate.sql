DELETE FROM eventstore.fields AS stored
USING unnest($1::TEXT[], $2::TEXT[], $3::TEXT[], $4::TEXT[])
    AS requested(aggregate_type, aggregate_id, instance_id, resource_owner)
WHERE (stored.instance_id, stored.resource_owner, stored.aggregate_type, stored.aggregate_id) =
      (requested.instance_id, requested.resource_owner, requested.aggregate_type, requested.aggregate_id);
