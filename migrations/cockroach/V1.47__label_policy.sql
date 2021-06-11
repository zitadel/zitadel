ALTER TABLE management.label_policies ALTER COLUMN label_policy_state SET DEFAULT 0;
ALTER TABLE management.label_policies ALTER COLUMN label_policy_state SET NOT NULL;
ALTER TABLE management.label_policies RENAME COLUMN secondary_color TO background_color;
ALTER TABLE management.label_policies ADD COLUMN warn_color TEXT;
ALTER TABLE management.label_policies ADD COLUMN font_color TEXT;
ALTER TABLE management.label_policies ADD COLUMN primary_color_dark TEXT;
ALTER TABLE management.label_policies ADD COLUMN background_color_dark TEXT;
ALTER TABLE management.label_policies ADD COLUMN warn_color_dark TEXT;
ALTER TABLE management.label_policies ADD COLUMN font_color_dark TEXT;
ALTER TABLE management.label_policies ADD COLUMN logo_url TEXT;
ALTER TABLE management.label_policies ADD COLUMN icon_url TEXT;
ALTER TABLE management.label_policies ADD COLUMN logo_dark_url TEXT;
ALTER TABLE management.label_policies ADD COLUMN icon_dark_url TEXT;
ALTER TABLE management.label_policies ADD COLUMN font_url TEXT;
ALTER TABLE management.label_policies ADD COLUMN err_msg_popup BOOLEAN;
ALTER TABLE management.label_policies ADD COLUMN disable_watermark BOOLEAN;

ALTER TABLE adminapi.label_policies ALTER COLUMN label_policy_state SET DEFAULT 0;
ALTER TABLE adminapi.label_policies ALTER COLUMN label_policy_state SET NOT NULL;
ALTER TABLE adminapi.label_policies RENAME COLUMN secondary_color TO background_color;
ALTER TABLE adminapi.label_policies ADD COLUMN warn_color TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN font_color TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN primary_color_dark TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN background_color_dark TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN warn_color_dark TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN font_color_dark TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN logo_url TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN icon_url TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN logo_dark_url TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN icon_dark_url TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN font_url TEXT;
ALTER TABLE adminapi.label_policies ADD COLUMN err_msg_popup BOOLEAN;
ALTER TABLE adminapi.label_policies ADD COLUMN disable_watermark BOOLEAN;

ALTER TABLE auth.label_policies ALTER COLUMN label_policy_state SET DEFAULT 0;
ALTER TABLE auth.label_policies ALTER COLUMN label_policy_state SET NOT NULL;
ALTER TABLE auth.label_policies RENAME COLUMN secondary_color TO background_color;
ALTER TABLE auth.label_policies ADD COLUMN warn_color TEXT;
ALTER TABLE auth.label_policies ADD COLUMN font_color TEXT;
ALTER TABLE auth.label_policies ADD COLUMN primary_color_dark TEXT;
ALTER TABLE auth.label_policies ADD COLUMN background_color_dark TEXT;
ALTER TABLE auth.label_policies ADD COLUMN warn_color_dark TEXT;
ALTER TABLE auth.label_policies ADD COLUMN font_color_dark TEXT;
ALTER TABLE auth.label_policies ADD COLUMN logo_url TEXT;
ALTER TABLE auth.label_policies ADD COLUMN icon_url TEXT;
ALTER TABLE auth.label_policies ADD COLUMN logo_dark_url TEXT;
ALTER TABLE auth.label_policies ADD COLUMN icon_dark_url TEXT;
ALTER TABLE auth.label_policies ADD COLUMN font_url TEXT;
ALTER TABLE auth.label_policies ADD COLUMN err_msg_popup BOOLEAN;
ALTER TABLE auth.label_policies ADD COLUMN disable_watermark BOOLEAN;


BEGIN;
ALTER TABLE management.label_policies DROP CONSTRAINT "primary";
ALTER TABLE management.label_policies ADD CONSTRAINT "primary" PRIMARY KEY (aggregate_id, label_policy_state);
ALTER TABLE adminapi.label_policies DROP CONSTRAINT "primary";
ALTER TABLE adminapi.label_policies ADD CONSTRAINT "primary" PRIMARY KEY (aggregate_id, label_policy_state);
ALTER TABLE auth.label_policies DROP CONSTRAINT "primary";
ALTER TABLE auth.label_policies ADD CONSTRAINT "primary" PRIMARY KEY (aggregate_id, label_policy_state);
COMMIT;


ALTER TABLE management.users ADD COLUMN avatar_key TEXT;
ALTER TABLE auth.users ADD COLUMN avatar_key TEXT;
ALTER TABLE adminapi.users ADD COLUMN avatar_key TEXT;

ALTER TABLE auth.user_sessions ADD COLUMN avatar_key TEXT;

CREATE TABLE adminapi.styling (
     aggregate_id TEXT,

     creation_date TIMESTAMPTZ,
     change_date TIMESTAMPTZ,
     label_policy_state SMALLINT NOT NULL DEFAULT 0,
     sequence BIGINT,

     primary_color TEXT,
     background_color TEXT,
     warn_color TEXT,
     font_color TEXT,
     primary_color_dark TEXT,
     background_color_dark TEXT,
     warn_color_dark TEXT,
     font_color_dark TEXT,
     logo_url TEXT,
     icon_url TEXT,
     logo_dark_url TEXT,
     icon_dark_url TEXT,
     font_url TEXT,

     err_msg_popup BOOLEAN,
     disable_watermark BOOLEAN,
     hide_login_name_suffix BOOLEAN,

     PRIMARY KEY (aggregate_id, label_policy_state)
);

ALTER TABLE adminapi.features RENAME COLUMN label_policy TO label_policy_private_label;
ALTER TABLE adminapi.features ADD COLUMN label_policy_watermark BOOLEAN;
ALTER TABLE auth.features RENAME COLUMN label_policy TO label_policy_private_label;
ALTER TABLE auth.features ADD COLUMN label_policy_watermark BOOLEAN;
ALTER TABLE authz.features RENAME COLUMN label_policy TO label_policy_private_label;
ALTER TABLE authz.features ADD COLUMN label_policy_watermark BOOLEAN;
ALTER TABLE management.features RENAME COLUMN label_policy TO label_policy_private_label;
ALTER TABLE management.features ADD COLUMN label_policy_watermark BOOLEAN;