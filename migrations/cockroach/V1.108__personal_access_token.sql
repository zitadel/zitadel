ALTER TABLE auth.tokens ADD COLUMN is_pat BOOLEAN DEFAULT false NOT NULL;

CREATE TABLE zitadel.projections.personal_access_tokens (
    id STRING
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , resource_owner STRING NOT NULL
    , sequence INT8 NOT NULL
    , user_id STRING NOT NULL
    , expiration TIMESTAMPTZ NOT NULL
    , scopes STRING[]

    , PRIMARY KEY (id)
);
