ALTER TABLE adminapi.features RENAME COLUMN custom_text TO custom_text_message;
ALTER TABLE auth.features RENAME COLUMN custom_text TO custom_text_message;
ALTER TABLE authz.features RENAME COLUMN custom_text TO custom_text_message;
ALTER TABLE management.features RENAME COLUMN custom_text TO custom_text_message;

ALTER TABLE adminapi.features ADD COLUMN custom_text_login BOOLEAN;
ALTER TABLE auth.features ADD COLUMN custom_text_login BOOLEAN;
ALTER TABLE authz.features ADD COLUMN custom_text_login BOOLEAN;
ALTER TABLE management.features ADD COLUMN custom_text_login BOOLEAN;