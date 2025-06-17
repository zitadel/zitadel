CREATE TABLE IF NOT EXISTS projections.hosted_login_translations (
    instance_id TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_type TEXT NOT NULL CHECK (aggregate_type = 'instance' OR aggregate_type = 'org'),
    creation_date TIMESTAMPTZ NOT NULL,
    change_date TIMESTAMPTZ,
    sequence BIGINT NOT NULL,
    locale TEXT NOT NULL CHECK (LENGTH(TRIM(locale)) >= 2),
    file JSONB NOT NULL DEFAULT '{}',

    PRIMARY KEY (instance_id, aggregate_id, aggregate_type, locale)
);

