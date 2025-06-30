CREATE TYPE zitadel.organization_state AS ENUM (
	'active',
	'inactive'
);

CREATE TABLE zitadel.organizations(
  id TEXT NOT NULL CHECK (id <> ''),
  name TEXT NOT NULL CHECK (name <> ''),
  instance_id TEXT NOT NULL CHECK (instance_id <> '') REFERENCES zitadel.instances (id),
  state zitadel.organization_state NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,

  PRIMARY KEY (instance_id, id),
  UNIQUE      (instance_id, name)
);

-- users are able to set the id for organizations
CREATE INDEX org_id_not_deleted_idx ON zitadel.organizations (id)
    WHERE deleted_at IS NULL;

CREATE INDEX org_name_not_deleted_idx ON zitadel.organizations (name)
    WHERE deleted_at IS NULL;

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.organizations
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();
