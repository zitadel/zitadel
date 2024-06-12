-- TODO: whats the correct name of this table?
CREATE TABLE eventstore.lookup_fields (
    instance_id TEXT NOT NULL
    , resource_owner TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , field_name TEXT NOT NULL
    , number_value NUMERIC
    , text_value TEXT
    , CONSTRAINT one_of_values CHECK ((number_value IS NULL) <> (text_value IS NULL))
);

CREATE INDEX IF NOT EXISTS lf_field_number_idx ON eventstore.lookup_fields (instance_id, field_name, number_value) INCLUDE (resource_owner, aggregate_type, aggregate_id);
CREATE INDEX IF NOT EXISTS lf_field_text_idx ON eventstore.lookup_fields (instance_id, field_name, text_value) INCLUDE (resource_owner, aggregate_type, aggregate_id);

select 
    instance_id
    , resource_owner
    , aggregate_type
    , aggregate_id
    , field_name
    , text_value
from 
    eventstore.lookup_fields 
where 
    instance_id = '271204370027177212'
    and aggregate_type = 'project' 
    -- and field_name = 'project:app:oidc:client_id' 
    and field_name = 'project:app:id' 
    and text_value like '271204370027504892';