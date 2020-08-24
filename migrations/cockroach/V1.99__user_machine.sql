ALTER TABLE management.users ADD COLUMN machine_name STRING, ADD COLUMN machine_description STRING, ADD COLUMN user_type STRING;
ALTER TABLE adminapi.users ADD COLUMN machine_name STRING, ADD COLUMN machine_description STRING, ADD COLUMN user_type STRING;
ALTER TABLE auth.users ADD COLUMN machine_name STRING, ADD COLUMN machine_description STRING, ADD COLUMN user_type STRING;

--TODO: (adminapi)iam_members
--TODO: (auth,authz,management)user_grants
--TODO: (auth)user_sessions
--tODO: (management)project_grant_members