CREATE TYPE zitadel.organization_state AS ENUM (
	'ACTIVE',
	'INACTIVE'
);

CREATE TABLE zitadel.organizations(
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  instance_id TEXT NOT NULL,
  state zitadel.organization_state NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);
