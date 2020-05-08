BEGIN;

CREATE TABLE management.users (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    user_state SMALLINT,
    last_login TIMESTAMPTZ,
    password_change TIMESTAMPTZ,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    nick_Name TEXT,
    display_name TEXT,
    preferred_language TEXT,
    gender SMALLINT,
    email TEXT,
    is_email_verified BOOLEAN,
    phone TEXT,
    is_phone_verified BOOLEAN,
    country TEXT,
    locality TEXT,
    postal_code TEXT,
    region TEXT,
    street_address TEXT,
    otp_state SMALLINT,
    sequence BIGINT,

    PRIMARY KEY (id)
);

CREATE TABLE management.user_grants (
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
