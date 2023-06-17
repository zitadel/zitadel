BEGIN;
ALTER TABLE eventstore.events ALTER COLUMN editor_service DROP NOT NULL;

-- not used anymore
DROP INDEX IF EXISTS eventstore.events@previous_sequence_unique CASCADE;
DROP INDEX IF EXISTS eventstore.events@prev_agg_type_seq_unique CASCADE;
DROP INDEX IF EXISTS eventstore.events@events_event_sequence_instance_id_key CASCADE;
-- dropped because equal to the PK new
DROP INDEX IF EXISTS eventstore.events@agg_type CASCADE;
DROP INDEX IF EXISTS eventstore.events@agg_type_agg_id CASCADE;
DROP INDEX IF EXISTS eventstore.events@agg_type_seq CASCADE;
DROP INDEX IF EXISTS eventstore.events@max_sequence CASCADE;
COMMIT;

BEGIN;
ALTER TABLE eventstore.events ALTER PRIMARY KEY USING COLUMNS (instance_id, aggregate_type, aggregate_id, event_sequence DESC);
COMMIT;

BEGIN;
CREATE INDEX IF NOT EXISTS es_handler_idx ON eventstore.events (instance_id, aggregate_type, event_type, created_at) STORING (aggregate_version, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_user, editor_service, resource_owner);
CREATE INDEX IF NOT EXISTS es_agg_id_event_idx ON eventstore.events (aggregate_type, aggregate_id, event_type) STORING (aggregate_version, previous_aggregate_sequence, previous_aggregate_type_sequence, created_at, event_data, editor_user, editor_service, resource_owner);
CREATE INDEX IF NOT EXISTS es_active_instances_new ON eventstore.events (created_at DESC, instance_id) USING HASH;
COMMIT;

-- TODO: create in projection package
-- BEGIN;
-- CREATE INDEX IF NOT EXISTS im_instance ON projections.instance_members3 (instance_id) STORING (creation_date, change_date, user_owner_removed, roles, sequence, resource_owner, owner_removed); 
-- CREATE INDEX IF NOT EXISTS om_instance ON projections.org_members3 (instance_id) STORING (creation_date, change_date, user_owner_removed, roles, sequence, resource_owner, owner_removed); 
-- CREATE INDEX IF NOT EXISTS pgm_instance ON projections.project_grant_members3 (instance_id) STORING (creation_date, change_date, user_owner_removed, roles, sequence, resource_owner, owner_removed, granted_org_removed); 
-- CREATE INDEX IF NOT EXISTS p_instance ON projections.project_members3 (instance_id) STORING (creation_date, change_date, user_owner_removed, roles, sequence, resource_owner, owner_removed); 
-- CREATE INDEX IF NOT EXISTS lnu_instance_ro_id ON projections.login_names2_users (instance_id, resource_owner, id) STORING (user_name, owner_removed);
-- COMMIT;