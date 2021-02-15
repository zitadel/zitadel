CREATE TABLE authz.user_memberships (
   user_id TEXT,
   member_type SMALLINT,
   aggregate_id TEXT,
   object_id TEXT,

   roles TEXT ARRAY,
   display_name TEXT,
   resource_owner TEXT,
   resource_owner_name TEXT,
   creation_date TIMESTAMPTZ,
   change_date TIMESTAMPTZ,
   sequence BIGINT,

   PRIMARY KEY (user_id, member_type, aggregate_id, object_id)
);

ALTER TABLE authz.user_memberships OWNER TO admin;