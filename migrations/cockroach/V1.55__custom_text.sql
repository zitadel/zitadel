CREATE TABLE notification.custom_texts (
   aggregate_id TEXT,

   creation_date TIMESTAMPTZ,
   change_date TIMESTAMPTZ,
   sequence BIGINT,

   template TEXT,
   language TEXT,
   key TEXT,
   text TEXT,

   PRIMARY KEY (aggregate_id, template, key, language)
);
