CREATE TABLE zitadel.projections.label_policies (
    id STRING NOT NULL, --TODO: pk
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    state INT2 NULL,
    resource_owner TEXT,
    
    is_default BOOLEAN,
    hide_login_name_suffix BOOLEAN,
	font_url STRING,
	watermark_disabled BOOLEAN,
	should_error_popup BOOLEAN,
	light_primary_color STRING,
	light_warn_color STRING,
	light_background_color STRING,
	light_font_color STRING,
	light_logo_url STRING,
	light_icon_url STRING,
	dark_primary_color STRING,
	dark_warn_color STRING,
	dark_background_color STRING,
	dark_font_color STRING,
	dark_logo_url STRING,
	dark_icon_url STRING
);