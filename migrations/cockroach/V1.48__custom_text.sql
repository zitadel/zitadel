ALTER TABLE adminapi.features ADD COLUMN custom_text BOOLEAN;
ALTER TABLE auth.features ADD COLUMN custom_text BOOLEAN;
ALTER TABLE authz.features ADD COLUMN custom_text BOOLEAN;
ALTER TABLE management.features ADD COLUMN custom_text BOOLEAN;

CREATE TABLE adminapi.message_texts (
     aggregate_id TEXT,

     creation_date TIMESTAMPTZ,
     change_date TIMESTAMPTZ,
     message_text_state SMALLINT,
     sequence BIGINT,

     message_text_type TEXT,
     language TEXT,
     title TEXT,
     pre_header TEXT,
     subject TEXT,
     greeting TEXT,
     text TEXT,
     button_text TEXT,
     footer_text TEXT,

     PRIMARY KEY (aggregate_id, message_text_type, language)
);


CREATE TABLE management.message_texts (
   aggregate_id TEXT,

   creation_date TIMESTAMPTZ,
   change_date TIMESTAMPTZ,
   message_text_state SMALLINT,
   sequence BIGINT,

   message_text_type TEXT,
   language TEXT,
   title TEXT,
   pre_header TEXT,
   subject TEXT,
   greeting TEXT,
   text TEXT,
   button_text TEXT,
   footer_text TEXT,

   PRIMARY KEY (aggregate_id, message_text_type, language)
);

GRANT SELECT ON TABLE adminapi.message_texts TO notification;
GRANT SELECT ON TABLE management.message_texts TO notification;
ALTER TABLE management.message_texts OWNER TO admin;
ALTER TABLE adminapi.message_texts OWNER TO admin;
