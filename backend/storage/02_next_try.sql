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
    
    , path TEXT NOT NULL

    , value JSONB
    , text_value TEXT
    , number_value NUMERIC
    , boolean_value BOOLEAN

    , PRIMARY KEY (model, object_id, path)
    , FOREIGN KEY (model, object_id) REFERENCES objects(model, id) ON DELETE CASCADE
    , FOREIGN KEY (model, model_revision) REFERENCES models(name, revision) ON DELETE RESTRICT
);

CREATE OR REPLACE FUNCTION ip_value_converter()
RETURNS TRIGGER AS $$
BEGIN
    CASE jsonb_typeof(NEW.value)
        WHEN 'boolean' THEN
            NEW.boolean_value := NEW.value::BOOLEAN;
            NEW.value := NULL;
        WHEN 'number' THEN
            NEW.number_value := NEW.value::NUMERIC;
            NEW.value := NULL;
        WHEN 'string' THEN
            NEW.text_value := (NEW.value#>>'{}')::TEXT;
            NEW.value := NULL;
        ELSE
            -- do nothing
    END CASE;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ip_value_converter_before_insert
BEFORE INSERT
ON indexed_properties
FOR EACH ROW
EXECUTE FUNCTION ip_value_converter();

CREATE TRIGGER ip_value_converter_before_update
BEFORE UPDATE
ON indexed_properties
FOR EACH ROW
EXECUTE FUNCTION ip_value_converter();

CREATE INDEX ip_search_model_object ON indexed_properties (model, path, value) WHERE value IS NOT NULL;
CREATE INDEX ip_search_model_rev_object ON indexed_properties (model, model_revision, path, value) WHERE value IS NOT NULL;
CREATE INDEX ip_search_model_text ON indexed_properties (model, path, text_value) WHERE text_value IS NOT NULL;
CREATE INDEX ip_search_model_rev_text ON indexed_properties (model, model_revision, path, text_value) WHERE text_value IS NOT NULL;
CREATE INDEX ip_search_model_number ON indexed_properties (model, path, number_value) WHERE number_value IS NOT NULL;
CREATE INDEX ip_search_model_rev_number ON indexed_properties (model, model_revision, path, number_value) WHERE number_value IS NOT NULL;
CREATE INDEX ip_search_model_boolean ON indexed_properties (model, path, boolean_value) WHERE boolean_value IS NOT NULL;
CREATE INDEX ip_search_model_rev_boolean ON indexed_properties (model, model_revision, path, boolean_value) WHERE boolean_value IS NOT NULL;

CREATE TABLE IF NOT EXISTS parents (
    parent_model TEXT NOT NULL
    , parent_id TEXT NOT NULL
    , child_model TEXT NOT NULL
    , child_id TEXT NOT NULL

    , PRIMARY KEY (parent_model, parent_id, child_model, child_id)
    , FOREIGN KEY (parent_model, parent_id) REFERENCES objects(model, id) ON DELETE CASCADE
    , FOREIGN KEY (child_model, child_id) REFERENCES objects(model, id) ON DELETE CASCADE
);

INSERT INTO models VALUES
    ('instance', 1, ARRAY['name', 'domain.name'])
    , ('organization', 1, ARRAY['name'])
    , ('user', 1, ARRAY['username', 'email', 'firstname', 'lastname'])
;

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
            RETURN QUERY SELECT * FROM jsonb_to_rows(v, _path || k)
            UNION VALUES (_path, '{}'::JSONB);
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

CREATE OR REPLACE FUNCTION merge_payload(_old JSONB, _new JSONB)
RETURNS JSONB
LANGUAGE plpgsql
AS $$
DECLARE
    _fields CURSOR FOR SELECT DISTINCT ON (path) 
        path
        , last_value(value) over (partition by path) as value
    FROM (
        SELECT path, value FROM jsonb_to_rows(_old)
    UNION ALL
        SELECT path, value FROM jsonb_to_rows(_new)
    );
    _path text[];
    _value jsonb;
BEGIN
    OPEN _fields;
    LOOP
        FETCH _fields INTO _path, _value;
        EXIT WHEN NOT FOUND;
        IF jsonb_typeof(_value) = 'object' THEN
            IF _old #> _path IS NOT NULL THEN
                CONTINUE;
            END IF;
            _old = jsonb_set_lax(_old, _path, '{}'::jsonb, TRUE);
            CONTINUE;
        END IF;
            
        _old = jsonb_set_lax(_old, _path, _value, TRUE, 'delete_key');
    END LOOP;

    RETURN _old;
END;
$$;

CREATE OR REPLACE FUNCTION set_object(_object object)
RETURNS VOID AS $$
DECLARE
    _parent model;
BEGIN
    INSERT INTO objects (model, model_revision, id, payload)
    VALUES (_object.model, _object.model_revision, _object.id, _object.payload)
    ON CONFLICT (model, id) DO UPDATE
    SET
        payload = merge_payload(objects.payload, EXCLUDED.payload)
        , model_revision = EXCLUDED.model_revision;

    INSERT INTO indexed_properties (model, model_revision, object_id, path, value)
    SELECT
        *
    FROM (
        SELECT
            _object.model
            , _object.model_revision
            , _object.id
            , UNNEST(m.indexed_paths) AS "path"
            , _object.payload #> string_to_array(UNNEST(m.indexed_paths), '.') AS "value"
        FROM
            models m
        WHERE 
            m.name = _object.model 
            AND m.revision = _object.model_revision
        GROUP BY
            m.name
            , m.revision
    )
    WHERE
        "value" IS NOT NULL
    ON CONFLICT (model, object_id, path) DO UPDATE
    SET 
        value = EXCLUDED.value
        , text_value = EXCLUDED.text_value
        , number_value = EXCLUDED.number_value
        , boolean_value = EXCLUDED.boolean_value
    ;

    INSERT INTO parents (parent_model, parent_id, child_model, child_id)
    VALUES 
        (_object.model, _object.id, _object.model, _object.id)
    ON CONFLICT (parent_model, parent_id, child_model, child_id) DO NOTHING;

    IF _object.parents IS NULL THEN
        RETURN;
    END IF;

    FOREACH _parent IN ARRAY _object.parents
    LOOP
        INSERT INTO parents (parent_model, parent_id, child_model, child_id)
        SELECT 
            p.parent_model
            , p.parent_id
            , _object.model
            , _object.id
        FROM parents p
        WHERE 
            p.child_model = _parent.name
            AND p.child_id = _parent.id
        ON CONFLICT (parent_model, parent_id, child_model, child_id) DO NOTHING
        ;

        INSERT INTO parents (parent_model, parent_id, child_model, child_id)
        VALUES 
            (_parent.name, _parent.id, _object.model, _object.id)
        ON CONFLICT (parent_model, parent_id, child_model, child_id) DO NOTHING;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_objects(_objects object[])
RETURNS VOID AS $$
DECLARE
    _object object;
BEGIN
    FOREACH _object IN ARRAY _objects
    LOOP
        PERFORM set_object(_object);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- CREATE OR REPLACE FUNCTION set_objects(VARIADIC _objects object[])
-- RETURNS VOID AS $$
-- BEGIN
--     PERFORM set_objectS(_objects);
-- END;
-- $$ LANGUAGE plpgsql;

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

SELECT set_objects(
    ARRAY[
        ROW('instance', 1::smallint, 'i1', '{"domain": {"isVerified": true}}', NULL)::object
    ]
);


SELECT 
    o.*
FROM 
    indexed_properties ip 
JOIN 
    objects o
ON
    ip.model = o.model
    AND ip.object_id = o.id
WHERE
    ip.model = 'instance'
    AND ip.path = 'name'
    AND ip.text_value = 'i1';
;

select * from merge_payload(
    '{"a": "asdf", "b": {"c":{"d": 1, "g": {"h": [4,5,6]}}}, "f": [1,2,3]}'::jsonb
    , '{"b": {"c":{"d": 1, "g": {"i": [4,5,6]}}}, "a": null}'::jsonb
);