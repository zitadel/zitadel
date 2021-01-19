CREATE TABLE eventstore.unique_usernames (
	unique_field TEXT,
	PRIMARY KEY (unique_field)
);

CREATE TABLE eventstore.unique_external_idps (
	unique_field TEXT,
	PRIMARY KEY (unique_field)
);

CREATE TABLE eventstore.unique_org_names (
	unique_field TEXT,
	PRIMARY KEY (unique_field)
);

CREATE TABLE eventstore.unique_project_names (
	unique_field TEXT,
	PRIMARY KEY (unique_field)
);

CREATE TABLE eventstore.unique_idp_config_names (
	unique_field TEXT,
	PRIMARY KEY (unique_field)
);