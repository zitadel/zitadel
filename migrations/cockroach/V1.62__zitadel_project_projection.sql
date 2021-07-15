use zitadel;

CREATE TABLE projections.projects (
    id TEXT,
    name TEXT,
    creation_date TIMESTAMPTZ,
    owner_id TEXT,
    creator_id TEXT,
    state INT2,

    PRIMARY KEY (id)
);

