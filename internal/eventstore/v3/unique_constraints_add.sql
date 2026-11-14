INSERT INTO eventstore.unique_constraints (
    instance_id
    , unique_type
    , unique_field
    , claimant_type
    , claimant_id
)
SELECT instance_id, unique_type, unique_field, NULLIF(claimant_type, ''), NULLIF(claimant_id, '')
FROM unnest($1::text[], $2::text[], $3::text[], $4::text[], $5::text[])
    AS requested(instance_id, unique_type, unique_field, claimant_type, claimant_id)
