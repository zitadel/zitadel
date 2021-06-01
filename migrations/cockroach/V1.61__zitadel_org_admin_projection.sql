use zitadel;

CREATE TABLE projections.org_admins (
    org_id TEXT,
    org_name TEXT,
    org_creation_date TIMESTAMPTZ,
    owner_id TEXT,
    owner_language VARCHAR(10),
    owner_email TEXT,
    owner_first_name TEXT,
	owner_last_name TEXT,
    owner_gender INT2,

    PRIMARY KEY (org_id, owner_id)
);
