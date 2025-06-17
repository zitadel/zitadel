CREATE TYPE zitadel.organization_state AS ENUM (
	'ACTIVE',
	'INACTIVE'
);

CREATE TABLE zitadel.organizations(
  id TEXT NOT NULL CHECK (id <> '') PRIMARY KEY,
  name TEXT NOT NULL CHECK (name <> ''),
  instance_id TEXT NOT NULL CHECK (instance_id <> ''),
  state zitadel.organization_state NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.organizations
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();
