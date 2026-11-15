INSERT INTO eventstore.unique_constraints (
    instance_id
    , unique_type
    , unique_field
    , owners
)
SELECT instance_id, unique_type, unique_field, ARRAY(SELECT jsonb_array_elements_text(owners::JSONB))
FROM unnest($1::TEXT[], $2::TEXT[], $3::TEXT[], $4::TEXT[])
    AS requested(instance_id, unique_type, unique_field, owners)
