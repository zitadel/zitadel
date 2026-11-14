DELETE FROM eventstore.fields AS stored
USING unnest($1::TEXT[], $2::TEXT[], $3::TEXT[], $4::TEXT[], $5::TEXT[], $6::TEXT[], $7::TEXT[], $8::TEXT[])
    AS requested(aggregate_type, aggregate_id, instance_id, resource_owner, object_type, object_id, object_revision, field_name)
WHERE (stored.instance_id, stored.resource_owner, stored.aggregate_type, stored.aggregate_id,
       stored.object_type, stored.object_id, stored.object_revision, stored.field_name) =
      (requested.instance_id, requested.resource_owner, requested.aggregate_type, requested.aggregate_id,
       requested.object_type, requested.object_id, requested.object_revision::SMALLINT, requested.field_name);
