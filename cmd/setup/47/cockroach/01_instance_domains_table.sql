CREATE TABLE IF NOT EXISTS instance_domains (
    instance_id TEXT
    , domain TEXT
    , is_generated BOOLEAN NOT NULL
    , is_primary BOOLEAN NOT NULL DEFAULT FALSE
    
    , change_date TIMESTAMPTZ NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    
    , latest_position NUMERIC NOT NULL
    , latest_in_position_order INT2 NOT NULL

    , PRIMARY KEY (instance_id, domain)
    , CONSTRAINT fk_instance_id FOREIGN KEY (instance_id) REFERENCES instances (id)
);
