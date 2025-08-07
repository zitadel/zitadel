ALTER TABLE IF EXISTS projections.login_names3_users ADD COLUMN IF NOT EXISTS user_name_lower TEXT GENERATED ALWAYS AS (lower(user_name)) STORED;
CREATE INDEX IF NOT EXISTS login_names3_users_search ON projections.login_names3_users (instance_id, user_name_lower) INCLUDE (resource_owner);

ALTER TABLE IF EXISTS projections.login_names3_domains ADD COLUMN IF NOT EXISTS name_lower TEXT GENERATED ALWAYS AS (lower(name)) STORED;
CREATE INDEX IF NOT EXISTS login_names3_domain_search ON projections.login_names3_domains (instance_id, resource_owner, name_lower);
CREATE INDEX IF NOT EXISTS login_names3_domain_search_result ON projections.login_names3_domains (instance_id, resource_owner) INCLUDE (is_primary);