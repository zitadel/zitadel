CREATE TABLE aether.identities(
    instance_id TEXT NOT NULL
    , collection_id TEXT
    , id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , last_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES aether.instances(id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, collection_id) REFERENCES aether.collections(instance_id, id) ON DELETE CASCADE
);

CREATE TABLE aether.identity_properties(
    instance_id TEXT NOT NULL
    -- if the property is linked to a collection
    , collection_id TEXT
    , identity_id TEXT NOT NULL
    
    , path LTREE NOT NULL
    , value JSONB NOT NULL
    , is_identifier BOOLEAN NOT NULL DEFAULT FALSE

    , PRIMARY KEY (instance_id, identity_id, path) --TODO: currently a path can only exist once per user, even across collections
    , FOREIGN KEY (instance_id, identity_id) REFERENCES aether.identities(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, collection_id) REFERENCES aether.collections(instance_id, id) ON DELETE CASCADE

    , UNIQUE NULLS NOT DISTINCT (instance_id, collection_id, identity_id, path)
);

CREATE TABLE aether.identity_identifiers(
    instance_id TEXT NOT NULL
    , collection_id TEXT 
    , identity_id TEXT NOT NULL
    
    , path LTREE NOT NULL
    , value JSONB NOT NULL

    , PRIMARY KEY (instance_id, identity_id, path)
    , FOREIGN KEY (instance_id, identity_id) REFERENCES aether.identities(instance_id, id) ON DELETE CASCADE ON UPDATE CASCADE
    , FOREIGN KEY (instance_id, identity_id, path) REFERENCES aether.identity_properties(instance_id, identity_id, path) ON DELETE CASCADE ON UPDATE CASCADE
    , UNIQUE NULLS NOT DISTINCT (instance_id, collection_id, value)
);

CREATE INDEX idx_identity_identifiers_value ON aether.identity_identifiers (instance_id, value);

CREATE FUNCTION aether.ensure_identity_identifier()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_identifier AND NEW.value IS DISTINCT FROM OLD.value THEN
        INSERT INTO aether.identity_identifiers(instance_id, collection_id, identity_id, path, value)
        VALUES (NEW.instance_id, NEW.collection_id, NEW.identity_id, NEW.path, NEW.value)
        ON CONFLICT (instance_id, identity_id, path) 
          DO UPDATE SET value = EXCLUDED.value;
    ELSIF NOT NEW.is_identifier AND OLD.is_identifier THEN
        DELETE FROM aether.identity_identifiers
        WHERE instance_id = NEW.instance_id
          AND identity_id = NEW.identity_id
          AND path = NEW.path;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_ensure_identity_identifier
AFTER INSERT OR UPDATE ON aether.identity_properties
FOR EACH ROW EXECUTE FUNCTION aether.ensure_identity_identifier();

CREATE FUNCTION aether.update_identity_last_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE aether.identities
       SET last_updated_at = NOW()
     WHERE instance_id = NEW.instance_id
       AND id = NEW.identity_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_identity_last_updated_at
AFTER INSERT OR UPDATE ON aether.identity_properties
FOR EACH ROW EXECUTE FUNCTION aether.update_identity_last_updated_at();