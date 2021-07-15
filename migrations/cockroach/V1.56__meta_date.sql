CREATE TABLE auth.meta_data (
   aggregate_id TEXT,

   key TEXT,
   value TEXT,

   resource_owner TEXT,
   creation_date TIMESTAMPTZ,
   change_date TIMESTAMPTZ,
   sequence BIGINT,

   PRIMARY KEY (aggregate_id, key)
);

CREATE TABLE management.meta_data (
    aggregate_id TEXT,

    key TEXT,
    value TEXT,

    resource_owner TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (aggregate_id, key)
);
