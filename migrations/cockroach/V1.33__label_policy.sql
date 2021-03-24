
ALTER TABLE management.label_policies ADD COLUMN hide_login_name_suffix BOOLEAN;
ALTER TABLE adminapi.label_policies ADD COLUMN hide_login_name_suffix BOOLEAN;


CREATE TABLE auth.label_policies (
   aggregate_id TEXT,

   creation_date TIMESTAMPTZ,
   change_date TIMESTAMPTZ,
   label_policy_state SMALLINT,
   sequence BIGINT,

   primary_color TEXT,
   secondary_color TEXT,
   hide_login_name_suffix BOOLEAN

   PRIMARY KEY (aggregate_id)
);

GRANT SELECT ON TABLE auth.label_policies TO notification;