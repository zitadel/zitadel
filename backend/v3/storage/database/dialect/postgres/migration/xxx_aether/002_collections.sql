CREATE TABLE aether.collections(
    instance_id TEXT NOT NULL
    , id TEXT NOT NULL
    , parent_id TEXT
    
    , type TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , last_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES aether.instances(id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, parent_id) REFERENCES aether.collections(instance_id, id) ON DELETE CASCADE
);

CREATE TABLE aether.collection_properties(
    instance_id TEXT NOT NULL
    , collection_id TEXT NOT NULL
    
    , path LTREE NOT NULL
    , value JSONB NOT NULL

    , is_linking BOOLEAN NOT NULL DEFAULT FALSE

    , PRIMARY KEY (instance_id, collection_id, path) --TODO: currently a path can only exist once per user, even across collections
    , FOREIGN KEY (instance_id, collection_id) REFERENCES aether.collections(instance_id, id) ON DELETE CASCADE

    , UNIQUE NULLS NOT DISTINCT (instance_id, collection_id, path)
);

CREATE FUNCTION aether.update_collection_last_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE aether.collections
    SET last_updated_at = NOW()
    WHERE instance_id = NEW.instance_id AND id = NEW.collection_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_collection_last_updated_at
AFTER INSERT OR UPDATE OR DELETE ON aether.collection_properties
FOR EACH ROW EXECUTE FUNCTION aether.update_collection_last_updated_at();

CREATE TABLE aether.collection_links(
    instance_id TEXT NOT NULL
    , collection_id TEXT NOT NULL
    , linked_collection_id TEXT NOT NULL

    , PRIMARY KEY (instance_id, collection_id, linked_collection_id)
    , FOREIGN KEY (instance_id, collection_id) REFERENCES aether.collections(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, linked_collection_id) REFERENCES aether.collections(instance_id, id) ON DELETE CASCADE
);

CREATE FUNCTION aether.update_collection_links()
RETURNS TRIGGER AS $$
BEGIN
-- TODO: handle json arrays for multiple links
    IF NEW.is_linking THEN
        INSERT INTO aether.collection_links (instance_id, collection_id, linked_collection_id)
        VALUES (NEW.instance_id, NEW.collection_id, (NEW.value #>> '{}')::TEXT)
        ON CONFLICT DO UPDATE SET linked_collection_id = EXCLUDED.linked_collection_id;
    ELSIF OLD.is_linking THEN
        DELETE FROM aether.collection_links
        WHERE instance_id = OLD.instance_id 
            AND collection_id = OLD.collection_id 
            AND linked_collection_id = (OLD.value #>> '{}')::TEXT;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_collection_links
AFTER INSERT OR UPDATE OR DELETE ON aether.collection_properties
FOR EACH ROW EXECUTE FUNCTION aether.update_collection_links();