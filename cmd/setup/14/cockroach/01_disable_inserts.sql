CREATE USER temp_events_owner;
ALTER TABLE eventstore.events OWNER TO temp_events_owner;
REVOKE INSERT ON eventstore.events FROM {{.username}};