
CREATE TABLE adminapi.styling2 (
    aggregate_id TEXT NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    label_policy_state INT2 NOT NULL DEFAULT 0::INT2,
    sequence INT8 NULL,
    primary_color TEXT NULL,
    background_color TEXT NULL,
    warn_color TEXT NULL,
    font_color TEXT NULL,
    primary_color_dark TEXT NULL,
    background_color_dark TEXT NULL,
    warn_color_dark TEXT NULL,
    font_color_dark TEXT NULL,
    logo_url TEXT NULL,
    icon_url TEXT NULL,
    logo_dark_url TEXT NULL,
    icon_dark_url TEXT NULL,
    font_url TEXT NULL,
    err_msg_popup BOOL NULL,
    disable_watermark BOOL NULL,
    hide_login_name_suffix BOOL NULL,
    instance_id TEXT NOT NULL,
    owner_removed BOOL DEFAULT false,

    PRIMARY KEY (instance_id, aggregate_id, label_policy_state)
);

CREATE INDEX  st2_owner_removed_idx ON adminapi.styling2 (owner_removed);