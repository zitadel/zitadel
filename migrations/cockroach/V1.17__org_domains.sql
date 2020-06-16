BEGIN;

CREATE TABLE management.org_domains (
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    domain TEXT,
    org_id TEXT,
    verified BOOLEAN,
    primary_domain BOOLEAN,

    PRIMARY KEY (org_id, domain)
);

COMMIT;