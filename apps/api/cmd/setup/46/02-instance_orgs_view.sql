CREATE OR REPLACE VIEW eventstore.instance_orgs AS
SELECT instance_id, aggregate_id as org_id
FROM eventstore.fields
WHERE aggregate_type = 'org'
AND object_type = 'org'
AND field_name = 'state';
