BEGIN;

CREATE TABLE management.orgs (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    org_state SMALLINT,
    sequence BIGINT,

    domain TEXT,
    name TEXT,

    PRIMARY KEY (id)
);

CREATE TABLE management.org_members (
    user_id TEXT,
    org_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    user_name TEXT,
    email_address TEXT,
    first_name TEXT,
    last_name TEXT,
    roles TEXT ARRAY,
    sequence BIGINT,

    PRIMARY KEY (org_id, user_id)
);

COMMIT;
