INSERT INTO eventstore.unique_constraints (
    instance_id
    , unique_type
    , unique_field
)
SELECT instance_id, unique_type, unique_field
FROM unnest($1::TEXT[], $2::TEXT[], $3::TEXT[])
    AS requested(instance_id, unique_type, unique_field)
