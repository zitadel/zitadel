CREATE TABLE IF NOT EXISTS zitadel.instances(
  id TEXT NOT NULL CHECK (id <> '') PRIMARY KEY,
  name TEXT NOT NULL CHECK (name <> ''),
  default_org_id TEXT, -- NOT NULL,
  iam_project_id TEXT, -- NOT NULL,
  console_client_id TEXT, -- NOT NULL,
  console_app_id TEXT, -- NOT NULL,
  default_language TEXT, -- NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX instance_name_not_deleted_idx ON zitadel.instances (name)
    WHERE deleted_at IS NOT NULL;

CREATE OR REPLACE FUNCTION zitadel.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at := NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.instances
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();
