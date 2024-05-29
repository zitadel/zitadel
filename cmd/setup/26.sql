CREATE TABLE IF NOT EXISTS auth.users3 (
    instance_id TEXT NOT NULL,
    id TEXT NOT NULL,
    resource_owner TEXT NOT NULL,
    change_date TIMESTAMPTZ NULL,
    password_set BOOL NULL,
    password_change TIMESTAMPTZ NULL,
    last_login TIMESTAMPTZ NULL,
    init_required BOOL NULL,
    mfa_init_skipped TIMESTAMPTZ NULL,
    username_change_required BOOL NULL,
    passwordless_init_required BOOL NULL,
    password_init_required BOOL NULL,

    PRIMARY KEY (instance_id, id)
)
