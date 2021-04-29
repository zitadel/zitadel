CREATE TABLE auth.refresh_tokens (
     id TEXT,

     creation_date TIMESTAMPTZ,
     change_date TIMESTAMPTZ,

     resource_owner TEXT,
     application_id TEXT,
     user_agent_id TEXT,
     user_id TEXT,
     idle_expiration TIMESTAMPTZ,
     expiration TIMESTAMPTZ,
     sequence BIGINT,
     scopes TEXT ARRAY,
     audience TEXT ARRAY,
     token TEXT,

     PRIMARY KEY (id)
);
