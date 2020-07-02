BEGIN;

CREATE TABLE admin_api.iam_members (
    user_id TEXT,

    iam_id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    user_name TEXT,
    email_address TEXT,
    first_name TEXT,
    last_name TEXT,
    roles TEXT ARRAY,
    sequence BIGINT,

    PRIMARY KEY (user_id)
);

COMMIT;
