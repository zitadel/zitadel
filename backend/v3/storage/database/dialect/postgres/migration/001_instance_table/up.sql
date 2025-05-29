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
-- CREATE UNIQUE INDEX instance_name_index ON zitadel.instances (name);
