WITH requested AS MATERIALIZED (
    SELECT aggregate_type, aggregate_id, instance_id, resource_owner, object_type, object_id,
           object_revision::SMALLINT AS object_revision, field_name, value::JSONB AS value, should_index
    FROM unnest($1::TEXT[], $2::TEXT[], $3::TEXT[], $4::TEXT[], $5::TEXT[], $6::TEXT[],
                $7::TEXT[], $8::TEXT[], $9::TEXT[], $10::BOOLEAN[])
        AS fields(aggregate_type, aggregate_id, instance_id, resource_owner, object_type, object_id,
                  object_revision, field_name, value, should_index)
), updated AS (
    UPDATE eventstore.fields AS stored
    SET value = requested.value, value_must_be_unique = FALSE, should_index = requested.should_index
    FROM requested
    WHERE (stored.instance_id, stored.resource_owner, stored.aggregate_type, stored.aggregate_id,
           stored.object_type, stored.object_id, stored.object_revision, stored.field_name) =
          (requested.instance_id, requested.resource_owner, requested.aggregate_type, requested.aggregate_id,
           requested.object_type, requested.object_id, requested.object_revision, requested.field_name)
    RETURNING stored.instance_id, stored.resource_owner, stored.aggregate_type, stored.aggregate_id,
              stored.object_type, stored.object_id, stored.object_revision, stored.field_name
)
INSERT INTO eventstore.fields (instance_id, resource_owner, aggregate_type, aggregate_id,
                              object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index)
SELECT instance_id, resource_owner, aggregate_type, aggregate_id,
       object_type, object_id, object_revision, field_name, value, FALSE, should_index
FROM requested
WHERE NOT EXISTS (
    SELECT 1 FROM updated
    WHERE (updated.instance_id, updated.resource_owner, updated.aggregate_type, updated.aggregate_id,
           updated.object_type, updated.object_id, updated.object_revision, updated.field_name) =
          (requested.instance_id, requested.resource_owner, requested.aggregate_type, requested.aggregate_id,
           requested.object_type, requested.object_id, requested.object_revision, requested.field_name)
);
