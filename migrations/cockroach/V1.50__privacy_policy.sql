CREATE TABLE auth.privacy_policies (
   aggregate_id TEXT,

   creation_date TIMESTAMPTZ,
   change_date TIMESTAMPTZ,
   state SMALLINT,
   sequence BIGINT,

   tos_link STRING,
   privacy_link STRING,

   PRIMARY KEY (aggregate_id)
);

CREATE TABLE adminapi.privacy_policies (
   aggregate_id TEXT,

   creation_date TIMESTAMPTZ,
   change_date TIMESTAMPTZ,
   state SMALLINT,
   sequence BIGINT,

   tos_link STRING,
   privacy_link STRING,

   PRIMARY KEY (aggregate_id)
);

CREATE TABLE management.privacy_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    state SMALLINT,
    sequence BIGINT,

    tos_link STRING,
    privacy_link STRING,

    PRIMARY KEY (aggregate_id)
);