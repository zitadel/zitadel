CREATE TABLE zitadel.projections.password_complexity_policies (
    id STRING NOT NULL, --TODO: pk
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    state INT2 NULL,
    resource_owner TEXT,
    
    is_default BOOLEAN,
    min_length INT8 NULL,
    has_lowercase BOOL NULL,
    has_uppercase BOOL NULL,
    has_symbol BOOL NULL,
    has_number BOOL NULL,

    PRIMARY KEY (id)
);