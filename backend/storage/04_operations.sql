-- postgres
DROP TABLE IF EXISTS properties;
DROP TABLE IF EXISTS parents CASCADE;
DROP TABLE IF EXISTS objects CASCADE;
DROP TABLE IF EXISTS indexed_properties;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS models;

DROP TYPE IF EXISTS object CASCADe;
DROP TYPE IF EXISTS model CASCADE;

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

CREATE TYPE operation_type AS ENUM (
    -- inserts a new object, if the object already exists the operation will fail
    -- path is ignored
    'create'
    -- if path is null an upsert is performed and the payload is overwritten
    -- if path is not null the value is set at the given path
    , 'set'
    -- drops an object if path is null
    -- if path is set but no value, the field at the given path is dropped
    -- if path and value are set and the field is an array the value is removed from the array
    , 'delete'
    -- adds a value to an array
    -- or a field if it does not exist, if the field exists the operation will fail
    , 'add'
);

CREATE TYPE object_manipulation AS (
    path TEXT[]
    , operation operation_type
    , value JSONB
);

CREATE TABLE IF NOT EXISTS parents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    , parent_model TEXT NOT NULL
    , parent_id TEXT NOT NULL
    , child_model TEXT NOT NULL
    , child_id TEXT NOT NULL

    , FOREIGN KEY (parent_model, parent_id) REFERENCES objects(model, id) ON DELETE CASCADE
    , FOREIGN KEY (child_model, child_id) REFERENCES objects(model, id) ON DELETE CASCADE
);

CREATE TYPE parent_operation AS ENUM (
    'add'
    , 'remove'
);

CREATE TYPE parent_manipulation AS (
    model TEXT
    , id TEXT
    ,  operation parent_operation
);

CREATE OR REPLACE FUNCTION jsonb_set_ensure_path(_jsonb JSONB, _path TEXT[], _value JSONB)
RETURNS JSONB
LANGUAGE plpgsql
AS $$
DECLARE
    i INT;
BEGIN
    IF _jsonb #> _path IS NOT NULL THEN
        RETURN JSONB_SET(_jsonb, _path, _value);
    END IF;

    FOR i IN REVERSE ARRAY_LENGTH(_path, 1)..1 LOOP
        _value := JSONB_BUILD_OBJECT(_path[i], _value);
        IF _jsonb #> _path[:i] IS NOT NULL THEN
            EXIT;
        END IF;
    END LOOP;

    RETURN _jsonb || _value;
END;
$$;

-- current: {} 
-- '{"a": {"b": {"c": {"d": {"e": 1}}}}}'::JSONB #> '{a,b,c}' = 1

drop function manipulate_object;
drop function object_set;
CREATE OR REPLACE FUNCTION object_set(
    _model TEXT
    , _model_revision SMALLINT
    , _id TEXT

    , _manipulations object_manipulation[]
    , _parents parent_manipulation[]
)
RETURNS objects
LANGUAGE plpgsql
AS $$
DECLARE
    _manipulation object_manipulation;
BEGIN
    FOREACH _manipulation IN ARRAY _manipulations LOOP
        CASE _manipulation.operation
            WHEN 'create' THEN
                INSERT INTO objects (model, model_revision, id, payload)
                VALUES (_model, _model_revision, _id, _manipulation.value);
            WHEN 'delete' THEN
                CASE
                    WHEN _manipulation.path IS NULL THEN
                        DELETE FROM objects
                        WHERE 
                            model = _model
                            AND model_revision = _model_revision
                            AND id = _id;
                    WHEN _manipulation.value IS NULL THEN
                        UPDATE 
                            objects
                        SET 
                            payload = payload #- _manipulation.path
                        WHERE 
                            model = _model
                            AND model_revision = _model_revision
                            AND id = _id;
                    ELSE
                        UPDATE 
                            objects
                        SET 
                            payload = jsonb_set(payload, _manipulation.path, (SELECT JSONB_AGG(v) FROM JSONB_PATH_QUERY(payload, ('$.' || ARRAY_TO_STRING(_manipulation.path, '.') || '[*]')::jsonpath) AS v WHERE v <> _manipulation.value)) 
                        WHERE 
                            model = _model
                            AND model_revision = _model_revision
                            AND id = _id;
                    END CASE;
            WHEN 'set' THEN
                IF _manipulation.path IS NULL THEN
                    INSERT INTO objects (model, model_revision, id, payload)
                    VALUES (_model, _model_revision, _id, _manipulation.value)
                    ON CONFLICT (model, model_revision, id) 
                    DO UPDATE SET payload = _manipulation.value;
                ELSE
                    UPDATE 
                        objects
                    SET 
                        payload = jsonb_set_ensure_path(payload, _manipulation.path, _manipulation.value)
                    WHERE 
                        model = _model
                        AND model_revision = _model_revision
                        AND id = _id;
                END IF;
            WHEN 'add' THEN
                UPDATE 
                    objects
                SET 
                    -- TODO: parent field must exist
                    payload = CASE
                        WHEN jsonb_typeof(payload #> _manipulation.path) IS NULL THEN
                            jsonb_set_ensure_path(payload, _manipulation.path, _manipulation.value)
                        WHEN jsonb_typeof(payload #> _manipulation.path) = 'array' THEN
                            jsonb_set(payload, _manipulation.path, COALESCE(payload #> _manipulation.path, '[]'::JSONB) || _manipulation.value)
                        -- ELSE
                        --     RAISE EXCEPTION 'Field at path % is not an array', _manipulation.path;
                    END
                WHERE
                    model = _model
                    AND model_revision = _model_revision
                    AND id = _id;
            --     TODO: RAISE EXCEPTION 'Field at path % is not an array', _manipulation.path;
        END CASE;
    END LOOP;

    FOREACH _parent IN ARRAY _parents LOOP
        CASE _parent.operation
            WHEN 'add' THEN
                -- insert the new parent and all its parents
                INSERT INTO parents (parent_model, parent_id, child_model, child_id)
                (
                    SELECT 
                        id
                    FROM parents p
                    WHERE 
                        p.child_model = _parent.model
                        AND p.child_id = _parent.id
                    UNION
                    SELECT 
                        _parent.model
                        , _parent.id
                        , _model
                        , _id
                )
                ON CONFLICT (parent_model, parent_id, child_model, child_id) DO NOTHING;
            WHEN 'remove' THEN
                -- remove the parent including the objects childs parent
                DELETE FROM parents
                WHERE id IN (
                    SELECT 
                        id
                    FROM 
                        parents p
                    WHERE 
                        p.child_model = _model
                        AND p.child_id = _id
                        AND p.parent_model = _parent.model
                        AND p.parent_id = _parent.id
                    UNION
                    SELECT
                        id
                    FROM (
                        SELECT 
                            id
                        FROM
                            parents p
                        WHERE 
                            p.parent_model = _model
                            AND p.parent_id = _id
                    )
                    WHERE
                        
                );
        END CASE;
    END LOOP;
    RETURN NULL;
END;
$$;

INSERT INTO models VALUES
    ('instance', 1, ARRAY['name', 'domain.name'])
    , ('organization', 1, ARRAY['name'])
    , ('user', 1, ARRAY['username', 'email', 'firstname', 'lastname'])
;

rollback;
BEGIN;
SELECT * FROM manipulate_object(
    'instance'
    , 1::SMALLINT
    , 'i1'
    , ARRAY[
        ROW(NULL, 'create', '{"name": "i1"}'::JSONB)::object_manipulation
        , ROW(ARRAY['domain'], 'set', '{"name": "example.com", "isVerified": false}'::JSONB)::object_manipulation
        , ROW(ARRAY['domain', 'isVerified'], 'set', 'true'::JSONB)::object_manipulation
        , ROW(ARRAY['domain', 'name'], 'delete', NULL)::object_manipulation
        , ROW(ARRAY['domain', 'name'], 'add', '"i1.com"')::object_manipulation
        , ROW(ARRAY['managers'], 'set', '[]'::JSONB)::object_manipulation
        , ROW(ARRAY['managers', 'objects'], 'add', '[{"a": "asdf"}, {"a": "qewr"}]'::JSONB)::object_manipulation
        , ROW(ARRAY['managers', 'objects'], 'delete', '{"a": "asdf"}'::JSONB)::object_manipulation
        , ROW(ARRAY['some', 'objects'], 'set', '{"a": "asdf"}'::JSONB)::object_manipulation
        -- , ROW(NULL, 'delete', NULL)::object_manipulation
    ]
);
select * from objects;
ROLLBACK;

BEGIN;
SELECT * FROM manipulate_object(
    'instance'
    , 1::SMALLINT
    , 'i1'
    , ARRAY[
        ROW(NULL, 'create', '{"name": "i1"}'::JSONB)::object_manipulation
        , ROW(ARRAY['domain', 'name'], 'set', '"example.com"'::JSONB)::object_manipulation
    ]
);
select * from objects;
ROLLBACK;

select jsonb_path_query_array('{"a": [12, 13, 14, 15]}'::JSONB, ('$.a ? (@ != $val)')::jsonpath, jsonb_build_object('val', '12'));