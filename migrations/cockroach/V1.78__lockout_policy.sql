CREATE TABLE zitadel.projections.lockout_policies (
    id STRING NOT NULL, --TODO: pk
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    state INT2 NULL,
    resource_owner TEXT,
    
    is_default BOOLEAN,
    max_password_attempts INT8 NULL,
    show_failure BOOLEAN NULL
);