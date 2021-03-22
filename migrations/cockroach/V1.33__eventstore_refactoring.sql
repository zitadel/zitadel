CREATE TABLE eventstore.unique_constraints (
    unique_type TEXT,
	unique_field TEXT,
	PRIMARY KEY (unique_type, unique_field)
);

GRANT DELETE ON TABLE eventstore.unique_constraints to eventstore;

ALTER TABLE management.login_policies ADD COLUMN default_policy BOOLEAN;
ALTER TABLE adminapi.login_policies ADD COLUMN default_policy BOOLEAN;
ALTER TABLE auth.login_policies ADD COLUMN default_policy BOOLEAN;

CREATE INDEX event_type ON eventstore.events (event_type);
CREATE INDEX resource_owner ON eventstore.events (resource_owner);

CREATE USER queries WITH PASSWORD ${queriespassword};
GRANT SELECT ON TABLE eventstore.events TO queries;

ALTER TABLE management.org_members ADD COLUMN preferred_login_name TEXT;
ALTER TABLE management.project_members ADD COLUMN preferred_login_name TEXT;
ALTER TABLE management.project_grant_members ADD COLUMN preferred_login_name TEXT;
ALTER TABLE adminapi.iam_members ADD COLUMN preferred_login_name TEXT;

