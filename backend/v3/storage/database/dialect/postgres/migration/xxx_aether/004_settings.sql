CREATE TABLE aether.settings(
    instance_id TEXT NOT NULL
    , collection_id TEXT

    , id SERIAL NOT NULL
    , type TEXT NOT NULL

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES aether.instances(id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, collection_id) REFERENCES aether.collections(instance_id, id) ON DELETE CASCADE
);

CREATE TABLE aether.setting_properties(
    instance_id TEXT NOT NULL
    , setting_id INT
    , path LTREE NOT NULL
    , value JSONB NOT NULL

    , PRIMARY KEY (instance_id, setting_id, path)
    , FOREIGN KEY (instance_id, setting_id) REFERENCES aether.settings(instance_id, id) ON DELETE CASCADE
);
