-- Make sure no other transaction is writing to fields table.
-- Will block transactions with events that modify the fields table.
-- SELECT is still possible during this time.
LOCK TABLE eventstore.fields IN SHARE ROW EXCLUSIVE MODE;

DELETE FROM eventstore.fields
  WHERE object_type IN (
    'instance_member_role',
    'org_member_role',
    'project_member_role'
  );

INSERT INTO eventstore.fields(
    instance_id,
    resource_owner,
    aggregate_type,
    aggregate_id,
    object_type,
    object_id,
    object_revision,
    field_name,
    value,
    value_must_be_unique,
    should_index
)

SELECT
  instance_id, 
  resource_owner,
  'instance' as aggregate_type,
  id as aggregate_id,
  'instance_member_role' as object_type,
  user_id as object_id,
  1::smallint as object_revision,
  'instance_role' as field_name,
  to_jsonb(unnest(roles)) as value,
  false as value_must_be_unique,
  true as should_index
FROM projections.instance_members4

UNION ALL

SELECT
  instance_id, 
  resource_owner,
  'org' as aggregate_type,
  org_id as aggregate_id,
  'org_member_role' as object_type,
  user_id as object_id,
  1::smallint as object_revision,
  'org_role' as field_name,
  to_jsonb(unnest(roles)) as value,
  false as value_must_be_unique,
  true as should_index
FROM projections.org_members4

UNION ALL

SELECT
  instance_id, 
  resource_owner,
  'project' as aggregate_type,
  project_id as aggregate_id,
  'project_member_role' as object_type,
  user_id as object_id,
  1::smallint as object_revision,
  'project_role' as field_name,
  to_jsonb(unnest(roles)) as value,
  false as value_must_be_unique,
  true as should_index
FROM projections.project_members4;
