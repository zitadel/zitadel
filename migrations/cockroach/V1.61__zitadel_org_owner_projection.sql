use zitadel;

CREATE TABLE projections.org_owners_orgs (
    id TEXT,
    name TEXT,
    creation_date TIMESTAMPTZ,

    PRIMARY KEY (id)
);


CREATE TABLE projections.org_owners_users (
    org_id TEXT,
    owner_id TEXT,
    language VARCHAR(10),
    email TEXT,
    first_name TEXT,
	last_name TEXT,
    gender INT2,

    PRIMARY KEY (owner_id, org_id),
    CONSTRAINT fk_org FOREIGN KEY (org_id) REFERENCES projections.org_owners_orgs (id) ON DELETE CASCADE
);

CREATE VIEW projections.org_owners AS (
    SELECT o.id AS org_id, 
        o.name AS org_name, 
        o.creation_date,
        u.owner_id,
        u.language,
        u.email,
        u.first_name,
        u.last_name,
        u.gender
    FROM projections.org_owners_orgs o
    JOIN projections.org_owners_users u ON o.id = u.org_id
);