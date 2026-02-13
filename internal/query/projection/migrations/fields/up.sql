CREATE TABLE eventstore.fields (
    id TEXT NOT NULL DEFAULT gen_random_uuid()
    , instance_id TEXT NOT NULL
    , resource_owner TEXT NOT NULL

    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL

    , object_type TEXT NOT NULL
    , object_id TEXT NOT NULL
    , object_revision INT2
    
    , field_name TEXT NOT NULL

    -- all the values of fields are inserted into value column as jsonb and if we need to index something we store it to the type specific column additionally
    , "value" JSONB NOT NULL
    , number_value NUMERIC GENERATED ALWAYS AS (CASE WHEN should_index AND JSONB_TYPEOF("value") = 'number' THEN "value"::NUMERIC ELSE NULL END) STORED
    , text_value TEXT GENERATED ALWAYS AS (CASE WHEN should_index AND JSONB_TYPEOF("value") = 'string' THEN "value" #>> '{}' ELSE NULL END) STORED
    , bool_value BOOLEAN GENERATED ALWAYS AS (CASE WHEN should_index AND JSONB_TYPEOF("value") = 'boolean' THEN "value"::BOOLEAN ELSE NULL END) STORED
    
    -- if true the value must be unique within an instance
    , value_must_be_unique BOOLEAN
    -- if set to true the primitive value is indexed
    , should_index BOOLEAN

    , PRIMARY KEY (instance_id, id)
    -- TODO: create issue to enable the foreign key as soon as the objects table is implemented
    -- , CONSTRAINT f_objects_fk FOREIGN KEY (instance_id, resource_owner, object_type, object_id, object_revision) REFERENCES eventstore.objects (instance_id, resource_owner, object_type, object_id, object_revision) ON DELETE CASCADE

    -- the constraint ensures that a primitive value is set if the value must be unique
    , CONSTRAINT primitive_value_for_unique_check CHECK (
        CASE 
            WHEN value_must_be_unique THEN num_nonnulls(number_value, text_value, bool_value) = 1
            ELSE true
        END
    )
    -- the constraint ensures that a primitive value is set if the value must be indexed
    , CONSTRAINT primitive_value_for_index CHECK (
        CASE 
            WHEN should_index THEN num_nonnulls(number_value, text_value, bool_value) = 1
            ELSE true
        END
    )
);

-- unique constraints for primitive values
CREATE UNIQUE INDEX IF NOT EXISTS f_number_unique_idx ON eventstore.fields (instance_id, field_name, number_value) WHERE value_must_be_unique;
CREATE UNIQUE INDEX IF NOT EXISTS f_text_unique_idx ON eventstore.fields (instance_id, field_name, text_value) WHERE value_must_be_unique;
CREATE UNIQUE INDEX IF NOT EXISTS f_bool_unique_idx ON eventstore.fields (instance_id, field_name, bool_value) WHERE value_must_be_unique;

-- search index for primitive values
CREATE INDEX IF NOT EXISTS f_number_value_idx ON eventstore.fields (instance_id, object_type, field_name, number_value) 
    INCLUDE (resource_owner, object_id, object_revision, "value") 
    WHERE number_value IS NOT NULL ;
CREATE INDEX IF NOT EXISTS f_text_value_idx ON eventstore.fields (instance_id, object_type, field_name, text_value) 
    INCLUDE (resource_owner, object_id, object_revision, "value") 
    WHERE text_value IS NOT NULL ;
CREATE INDEX IF NOT EXISTS f_bool_value_idx ON eventstore.fields (instance_id, object_type, field_name, bool_value) 
    INCLUDE (resource_owner, object_id, object_revision, "value") 
    WHERE bool_value IS NOT NULL ;

-- search index for object by id
CREATE INDEX IF NOT EXISTS f_object_idx ON eventstore.fields (instance_id, object_type, object_id, object_revision)
    INCLUDE (resource_owner, field_name, "value") ;