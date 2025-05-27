-- the projection for instances happens before md/setup/54.go is run,
-- hence why the zitadel schema is added below
CREATE SCHEMA IF NOT EXISTS zitadel;

CREATE TABLE IF NOT EXISTS zitadel.instances(
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  default_org_id TEXT, -- NOT NULL,
  iam_project_id TEXT, -- NOT NULL,
  console_client_id TEXT, -- NOT NULL,
  console_app_id TEXT, -- NOT NULL,
  default_language TEXT, -- NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE UNIQUE INDEX instance_name_index ON zitadel.instances (name);
