DROP TABLE IF EXISTS properties;
DROP TABLE IF EXISTS parents;
DROP TABLE IF EXISTS objects;

CREATE TABLE IF NOT EXISTS objects (
    type TEXT NOT NULL
    , id TEXT NOT NULL

    , PRIMARY KEY (type, id)
);

TRUNCATE objects CASCADE;
INSERT INTO objects VALUES
    ('instance', 'i1')
    , ('organization', 'o1')
    , ('user', 'u1')
    , ('user', 'u2')
    , ('organization', 'o2')
    , ('user', 'u3')
    , ('project', 'p3')
    
    , ('instance', 'i2')
    , ('organization', 'o3')
    , ('user', 'u4')
    , ('project', 'p1')
    , ('project', 'p2')
    , ('application', 'a1')
    , ('application', 'a2')
    , ('org_domain', 'od1')
    , ('org_domain', 'od2')
;

CREATE TABLE IF NOT EXISTS parents (
    parent_type TEXT NOT NULL
    , parent_id TEXT NOT NULL
    , child_type TEXT NOT NULL
    , child_id TEXT NOT NULL
    , PRIMARY KEY (parent_type, parent_id, child_type, child_id)
    , FOREIGN KEY (parent_type, parent_id) REFERENCES objects(type, id) ON DELETE CASCADE
    , FOREIGN KEY (child_type, child_id) REFERENCES objects(type, id) ON DELETE CASCADE
);

INSERT INTO parents VALUES
    ('instance', 'i1', 'organization', 'o1')
    , ('organization', 'o1', 'user', 'u1')
    , ('organization', 'o1', 'user', 'u2')
    , ('instance', 'i1', 'organization', 'o2')
    , ('organization', 'o2', 'user', 'u3')
    , ('organization', 'o2', 'project', 'p3')

    , ('instance', 'i2', 'organization', 'o3')
    , ('organization', 'o3', 'user', 'u4')
    , ('organization', 'o3', 'project', 'p1')
    , ('organization', 'o3', 'project', 'p2')
    , ('project', 'p1', 'application', 'a1')
    , ('project', 'p2', 'application', 'a2')
    , ('organization', 'o3', 'org_domain', 'od1')
    , ('organization', 'o3', 'org_domain', 'od2')
;

CREATE TABLE properties (
    object_type TEXT NOT NULL
    , object_id TEXT NOT NULL
    , key TEXT NOT NULL
    , value JSONB NOT NULL
    , should_index BOOLEAN NOT NULL DEFAULT FALSE

    , PRIMARY KEY (object_type, object_id, key)
    , FOREIGN KEY (object_type, object_id) REFERENCES objects(type, id) ON DELETE CASCADE
);

CREATE INDEX properties_object_indexed ON properties (object_type, object_id) INCLUDE (value) WHERE should_index;
CREATE INDEX properties_value_indexed ON properties (object_type, key, value) WHERE should_index;

TRUNCATE properties;
INSERT INTO properties VALUES
    ('instance', 'i1', 'name', '"Instance 1"', TRUE)
    , ('instance', 'i1', 'description', '"Instance 1 description"', FALSE)
    , ('instance', 'i2', 'name', '"Instance 2"', TRUE)
    , ('organization', 'o1', 'name', '"Organization 1"', TRUE)
    , ('org_domain', 'od1', 'domain', '"example.com"', TRUE)
    , ('org_domain', 'od1', 'is_primary', 'true', TRUE)
    , ('org_domain', 'od1', 'is_verified', 'true', FALSE)
    , ('org_domain', 'od2', 'domain', '"example.org"', TRUE)
    , ('org_domain', 'od2', 'is_primary', 'false', TRUE)
    , ('org_domain', 'od2', 'is_verified', 'false', FALSE)
;

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    , type TEXT NOT NULL
    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , revision SMALLINT NOT NULL
    , creator TEXT NOT NULL
    , payload JSONB

    , global_sequence NUMERIC NOT NULL DEFAULT pg_current_xact_id()::TEXT::NUMERIC
    , sequence_order SMALLINT NOT NULL CHECK (sequence_order >= 0)

    -- , object_type TEXT NOT NULL
    -- , object_id TEXT NOT NULL

    -- , FOREIGN KEY (object_type, object_id) REFERENCES objects(type, id)
);

CREATE TYPE property (
    -- key must be a json path
    key TEXT
    -- value should be a primitive type
    , value JSONB
    -- indicates wheter the property should be indexed
    , should_index BOOLEAN
);

CREATE TYPE parent (
    parent_type TEXT
    , parent_id TEXT
);

CREATE TYPE object (
    type TEXT
    , id TEXT
    , properties property[]
    -- an object automatically inherits the parents of its parent
    , parents parent[]
);

CREATE TYPE command (
    type TEXT
    , revision SMALLINT
    , creator TEXT
    , payload JSONB

    -- if properties is null the objects and all its child objects get deleted
    -- if the value of a property is null the property and all sub fields get deleted
    -- for example if the key is 'a.b' and the value is null the property 'a.b.c' will be deleted as well
    , objects object[]
);

CREATE OR REPLACE PROCEDURE update_object(_object object)
AS $$
DECLARE
    _property property;
BEGIN
    FOR _property IN ARRAY _object.properties LOOP
        IF _property.value IS NULL THEN
            DELETE FROM properties
            WHERE object_type = _object.type
            AND object_id = _object.id
            AND key LIKE CONCAT(_property.key, '%');
        ELSE
            INSERT INTO properties (object_type, object_id, key, value, should_index)
            VALUES (_object.type, _object.id, _property.key, _property.value, _property.should_index)
            ON CONFLICT (object_type, object_id, key) DO UPDATE SET (value, should_index) = (_property.value, _property.should_index);
        END IF;
    END LOOP;
END;

CREATE OR REPLACE PROCEDURE delete_object(_type, _id) 
AS $$
BEGIN
    WITH RECURSIVE objects_to_delete (_type, _id) AS (
        SELECT $1, $2

        UNION

        SELECT p.child_type, p.child_id
        FROM parents p
        JOIN objects_to_delete o ON p.parent_type = o.type AND p.parent_id = o.id
    )
    DELETE FROM objects
    WHERE (type, id) IN (SELECT * FROM objects_to_delete)
END;

CREATE OR REPLACE FUNCTION push(_commands command[])
RETURNS NUMMERIC AS $$
DECLARE
    _command command;
    _index INT;

    _object object;
BEGIN
    FOR _index IN 1..array_length(_commands, 1) LOOP
        _command := _commands[_index];
        INSERT INTO events (type, revision, creator, payload)
        VALUES (_command.type, _command.revision, _command.creator, _command.payload);

        FOREACH _object IN ARRAY _command.objects LOOP
            IF _object.properties IS NULL THEN
                PERFORM delete_object(_object.type, _object.id);
            ELSE
                PERFORM update_object(_object);
            END IF;
    END LOOP;
    RETURN pg_current_xact_id()::TEXT::NUMERIC;
END;
$$ LANGUAGE plpgsql;


BEGIN;


RETURNING *
;

rollback;

SELECT
    *
FROM
    properties
WHERE
    (object_type, object_id) IN (
        SELECT 
            object_type
            , object_id
        FROM 
            properties 
        where 
            object_type = 'instance' 
            and key = 'name' 
            and value = '"Instance 1"' 
            and should_index 
    )
;