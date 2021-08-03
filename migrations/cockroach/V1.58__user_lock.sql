CREATE TABLE management.lockout_policies (
      aggregate_id TEXT,

      creation_date TIMESTAMPTZ,
      change_date TIMESTAMPTZ,
      lockout_policy_state SMALLINT,
      sequence BIGINT,

      max_password_attempts BIGINT,
      show_lockout_failures BOOLEAN,

      PRIMARY KEY (aggregate_id)
);

CREATE TABLE adminapi.lockout_policies (
     aggregate_id TEXT,

     creation_date TIMESTAMPTZ,
     change_date TIMESTAMPTZ,
     lockout_policy_state SMALLINT,
     sequence BIGINT,

     max_password_attempts BIGINT,
     show_lockout_failures BOOLEAN,

     PRIMARY KEY (aggregate_id)
);

CREATE TABLE auth.lockout_policies (
       aggregate_id TEXT,

       creation_date TIMESTAMPTZ,
       change_date TIMESTAMPTZ,
       lockout_policy_state SMALLINT,
       sequence BIGINT,

       max_password_attempts BIGINT,
       show_lockout_failures BOOLEAN,

       PRIMARY KEY (aggregate_id)
);

CREATE TABLE adminapi.user_lock (
   user_id TEXT,

   change_date TIMESTAMPTZ,
   sequence BIGINT,
   resourceowner TEXT,
   state SMALLINT,
   password_check_failed_count BIGINT,

   PRIMARY KEY (user_id)
);

DROP TABLE management.password_lockout_policies;
DROP TABLE adminapi.password_lockout_policies;
DROP TABLE auth.password_lockout_policies;