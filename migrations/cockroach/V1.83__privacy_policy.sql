CREATE TABLE zitadel.projections.privacy_policies (
    id STRING NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    state INT2 NULL,
    resource_owner TEXT,
    
    is_default BOOLEAN,
    privacy_link TEXT,
    tos_link TEXT,

    PRIMARY KEY (id)
);