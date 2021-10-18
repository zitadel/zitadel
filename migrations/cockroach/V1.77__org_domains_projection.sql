CREATE TABLE zitadel.projections.org_domains (
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    domain TEXT,
    org_id TEXT,
    is_verified BOOLEAN,
    is_primary BOOLEAN,
    validation_type SMALLINT,

    PRIMARY KEY (org_id, domain)
);