CREATE TABLE adminapi.user_lock (
   user_id TEXT,

   change_date TIMESTAMPTZ,
   sequence BIGINT,
   resourceowner TEXT,
   state SMALLINT,
   password_check_failed_count BIGINT,

   PRIMARY KEY (user_id)
);