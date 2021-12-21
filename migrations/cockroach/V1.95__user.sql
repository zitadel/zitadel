CREATE TABLE zitadel.projections.users(
    id STRING
    , creation_date TIMESTAMPTZ
    , change_date TIMESTAMPTZ
    , resource_owner STRING NOT NULL
    , state INT2
    , sequence INT8

    , username STRING

    , PRIMARY KEY (id)
    , INDEX idx_username (username)
);

CREATE TABLE zitadel.projections.users_machines(
	user_id STRING REFERENCES zitadel.projections.users (id) ON DELETE CASCADE
	
    , name STRING NOT NULL
    , description STRING

    , PRIMARY KEY (user_id)
);

CREATE TABLE zitadel.projections.users_humans(
    user_id STRING REFERENCES zitadel.projections.users (id) ON DELETE CASCADE
    
    --profile
    , first_name STRING NOT NULL
    , last_name STRING NOT NULL
    , nick_name STRING
    , display_name STRING
    , preferred_language VARCHAR(10)
    , gender INT2 
    , avater_key STRING

    --email
    , email STRING NOT NULL
    , is_email_verified BOOLEAN NOT NULL DEFAULT false
    
    --phone
    , phone STRING
    , is_phone_verified BOOLEAN

    , PRIMARY KEY (user_id)
);
