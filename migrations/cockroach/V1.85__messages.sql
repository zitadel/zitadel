CREATE TABLE zitadel.projections.mail_templates (
    aggregate_id TEXT NOT NULL,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    state SMALLINT,
    sequence BIGINT,
    is_default BOOLEAN,

    template BYTES,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE zitadel.projections.message_texts (
   aggregate_id TEXT,

   creation_date TIMESTAMPTZ,
   change_date TIMESTAMPTZ,
   mail_text_state SMALLINT,
   sequence BIGINT,

   mail_text_type TEXT,
   language TEXT,
   title TEXT,
   pre_header TEXT,
   subject TEXT,
   greeting TEXT,
   text TEXT,
   button_text TEXT,

   PRIMARY KEY (aggregate_id, mail_text_type, language)
);