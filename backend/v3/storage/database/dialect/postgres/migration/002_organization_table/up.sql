CREATE TYPE zitadel.organization_state AS ENUM (
	'active',
	'inactive'
);

CREATE TABLE zitadel.organizations(
  id TEXT NOT NULL CHECK (id <> ''),
  name TEXT NOT NULL CHECK (name <> ''),
  instance_id TEXT NOT NULL REFERENCES zitadel.instances (id) ON DELETE CASCADE,
  state zitadel.organization_state NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

  PRIMARY KEY (instance_id, id)
);

CREATE UNIQUE INDEX org_unique_instance_id_name_idx
  ON zitadel.organizations (instance_id, name);

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.organizations
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();
