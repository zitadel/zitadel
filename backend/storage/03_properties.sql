-- postgres
DROP TABLE IF EXISTS properties;
DROP TABLE IF EXISTS parents CASCADE;
DROP TABLE IF EXISTS objects CASCADE;
DROP TABLE IF EXISTS indexed_properties;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS models;

DROP TYPE IF EXISTS object CASCADe;
DROP TYPE IF EXISTS model CASCADE;

CREATE TYPE model AS (
    name TEXT
    , id TEXT
);

CREATE TYPE object AS (
    model TEXT
    , model_revision SMALLINT
    , id TEXT
    , payload JSONB
    , parents model[]
);

CREATE TABLE models (
    name TEXT
    , revision SMALLINT NOT NULL CONSTRAINT positive_revision CHECK (revision > 0)
    , indexed_paths TEXT[]

    , PRIMARY KEY (name, revision)
);

CREATE TABLE objects (
    model TEXT NOT NULL
    , model_revision SMALLINT NOT NULL

    , id TEXT NOT NULL
    , payload JSONB

    , PRIMARY KEY (model, id)
    , FOREIGN KEY (model, model_revision) REFERENCES models(name, revision) ON DELETE RESTRICT
);

CREATE TABLE indexed_properties (
    model TEXT NOT NULL
    , model_revision SMALLINT NOT NULL
    , object_id TEXT NOT NULL
    
    , path TEXT[] NOT NULL

    , value JSONB
    , text_value TEXT
    , number_value NUMERIC
    , boolean_value BOOLEAN

    , PRIMARY KEY (model, object_id, path)
    , FOREIGN KEY (model, object_id) REFERENCES objects(model, id) ON DELETE CASCADE
    , FOREIGN KEY (model, model_revision) REFERENCES models(name, revision) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS parents (
    parent_model TEXT NOT NULL
    , parent_id TEXT NOT NULL
    , child_model TEXT NOT NULL
    , child_id TEXT NOT NULL
    , PRIMARY KEY (parent_model, parent_id, child_model, child_id)
    , FOREIGN KEY (parent_model, parent_id) REFERENCES objects(model, id) ON DELETE CASCADE
    , FOREIGN KEY (child_model, child_id) REFERENCES objects(model, id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION jsonb_to_rows(j jsonb, _path text[] DEFAULT ARRAY[]::text[])
RETURNS TABLE (path text[], value jsonb)
LANGUAGE plpgsql
AS $$
DECLARE
    k text;
    v jsonb;
BEGIN
    FOR k, v IN SELECT * FROM jsonb_each(j) LOOP
        IF jsonb_typeof(v) = 'object' THEN
            -- Recursive call for nested objects, appending the key to the path
            RETURN QUERY SELECT * FROM jsonb_to_rows(v, _path || k);
        ELSE
            -- Base case: return the key path and value
            CASE WHEN jsonb_typeof(v) = 'null' THEN
                RETURN QUERY SELECT _path || k, NULL::jsonb;
            ELSE
                RETURN QUERY SELECT _path || k, v;
            END CASE;
        END IF;
    END LOOP;
END;
$$;

-- after insert trigger which is called after the object was inserted and then inserts the properties
CREATE OR REPLACE FUNCTION set_ip_from_object_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    _property RECORD;
    _model models;
BEGIN
    SELECT * INTO _model FROM models WHERE name = NEW.model AND revision = NEW.model_revision;

    FOR _property IN SELECT * FROM jsonb_to_rows(NEW.payload) LOOP
        IF ARRAY_TO_STRING(_property.path, '.') = ANY(_model.indexed_paths) THEN
            INSERT INTO indexed_properties (model, model_revision, object_id, path, value)
            VALUES (NEW.model, NEW.model_revision, NEW.id, _property.path, _property.value);
        END IF;
    END LOOP;
    RETURN NULL;
END;
$$;

CREATE TRIGGER set_ip_from_object_insert
AFTER INSERT ON objects
FOR EACH ROW
EXECUTE FUNCTION set_ip_from_object_insert();

-- before update trigger with is called before an object is updated
-- it updates the properties table first
-- and computes the correct payload for the object
-- partial update of the object is allowed
-- if the value of a property is set to null the properties and all its children are deleted
CREATE OR REPLACE FUNCTION set_ip_from_object_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    _property RECORD;
    _payload JSONB;
    _model models;
    _path_index INT;
BEGIN
    _payload := OLD.payload;
    SELECT * INTO _model FROM models WHERE name = NEW.model AND revision = NEW.model_revision;

    FOR _property IN SELECT * FROM jsonb_to_rows(NEW.payload) ORDER BY array_length(path, 1) LOOP
        -- set the properties
        CASE WHEN _property.value IS NULL THEN
            RAISE NOTICE 'DELETE PROPERTY: %', _property;
            DELETE FROM indexed_properties 
            WHERE model = NEW.model
                AND model_revision = NEW.model_revision
                AND object_id = NEW.id 
                AND path[:ARRAY_LENGTH(_property.path, 1)] = _property.path;
        ELSE
            RAISE NOTICE 'UPSERT PROPERTY: %', _property;
            DELETE FROM indexed_properties
            WHERE 
                model = NEW.model
                AND model_revision = NEW.model_revision
                AND object_id = NEW.id
                AND (
                    _property.path[:array_length(path, 1)] = path
                    OR path[:array_length(_property.path, 1)] = _property.path
                )
                AND array_length(path, 1) <> array_length(_property.path, 1);

            -- insert property if should be indexed
            IF ARRAY_TO_STRING(_property.path, '.') = ANY(_model.indexed_paths) THEN
                RAISE NOTICE 'path should be indexed: %', _property.path;
                INSERT INTO indexed_properties (model, model_revision, object_id, path, value)
                VALUES (NEW.model, NEW.model_revision, NEW.id, _property.path, _property.value)
                ON CONFLICT (model, object_id, path) DO UPDATE
                SET value = EXCLUDED.value;
            END IF;
        END CASE;

        -- if property is updated we can set it directly
        IF _payload #> _property.path IS NOT NULL THEN
            _payload = jsonb_set_lax(COALESCE(_payload, '{}'::JSONB), _property.path, _property.value, TRUE);
            EXIT;
        END IF;
        -- ensure parent object exists exists
        FOR _path_index IN 1..(array_length(_property.path, 1)-1) LOOP
            IF _payload #> _property.path[:_path_index] IS NOT NULL AND jsonb_typeof(_payload #> _property.path[:_path_index]) = 'object' THEN
                CONTINUE;
            END IF;

            _payload = jsonb_set(_payload, _property.path[:_path_index], '{}'::JSONB, TRUE);
            EXIT;
        END LOOP;
        _payload = jsonb_set_lax(_payload, _property.path, _property.value, TRUE, 'delete_key');

    END LOOP;

    -- update the payload
    NEW.payload = _payload;

    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER set_ip_from_object_update
BEFORE UPDATE ON objects
FOR EACH ROW
EXECUTE FUNCTION set_ip_from_object_update();


CREATE OR REPLACE FUNCTION set_object(_object object)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO objects (model, model_revision, id, payload)
    VALUES (_object.model, _object.model_revision, _object.id, _object.payload)
    ON CONFLICT (model, id) DO UPDATE
    SET 
        payload = EXCLUDED.payload
        , model_revision = EXCLUDED.model_revision
    ;

    INSERT INTO parents (parent_model, parent_id, child_model, child_id)
    SELECT
        p.name
        , p.id
        , _object.model
        , _object.id
    FROM UNNEST(_object.parents) AS p
    ON CONFLICT DO NOTHING;
END;
$$;

CREATE OR REPLACE FUNCTION set_objects(_objects object[])
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    _object object;
BEGIN
    FOREACH _object IN ARRAY _objects LOOP
        PERFORM set_object(_object);
    END LOOP;
END;
$$;





INSERT INTO models VALUES
    ('instance', 1, ARRAY['name', 'domain.name'])
    , ('organization', 1, ARRAY['name'])
    , ('user', 1, ARRAY['username', 'email', 'firstname', 'lastname'])
;

INSERT INTO objects VALUES
    ('instance', 1, 'i2', '{"name": "i2", "domain": {"name": "example2.com", "isVerified": false}}')
    , ('instance', 1, 'i3', '{"name": "i3", "domain": {"name": "example3.com", "isVerified": false}}')
    , ('instance', 1, 'i4', '{"name": "i4", "domain": {"name": "example4.com", "isVerified": false}}')
;


begin;
UPDATE objects SET payload = '{"domain": {"isVerified": true}}' WHERE model = 'instance';
rollback;


SELECT set_objects(
    ARRAY[
        ROW('instance', 1::smallint, 'i1', '{"name": "i1", "domain": {"name": "example2.com", "isVerified": false}}', NULL)::object
        , ROW('organization', 1::smallint, 'o1', '{"name": "o1", "description": "something useful"}', ARRAY[
            ROW('instance', 'i1')::model
        ])::object
        , ROW('user', 1::smallint, 'u1', '{"username": "u1", "description": "something useful", "firstname": "Silvan"}', ARRAY[
            ROW('instance', 'i1')::model
            , ROW('organization', 'o1')::model
        ])::object
    ]
);