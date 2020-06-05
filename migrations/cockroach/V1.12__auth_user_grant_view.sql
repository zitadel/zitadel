BEGIN;

CREATE TABLE auth.user_grants (
    id TEXT,
    resource_owner TEXT,
    project_id TEXT,
    user_id TEXT,
    org_name TEXT,
    org_domain TEXT,
    project_name TEXT,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    role_keys TEXT Array,

    grant_state SMALLINT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (id)
);

COMMIT;
