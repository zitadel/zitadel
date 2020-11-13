ALTER TABLE management.users ADD COLUMN u2f_verified_ids TEXT ARRAY;
ALTER TABLE auth.users ADD COLUMN u2f_verified_ids TEXT ARRAY;
ALTER TABLE adminapi.users ADD COLUMN u2f_verified_ids TEXT ARRAY;