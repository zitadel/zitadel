use auth;
CREATE INDEX IF NOT EXISTS auth_code on auth.auth_requests (code);

CREATE INDEX IF NOT EXISTS default_event_query on eventstore.events (aggregate_type, aggregate_id, event_type, resource_owner);
DROP INDEX IF EXISTS event_type;
DROP INDEX IF EXISTS resource_owner;