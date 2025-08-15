-- Add composite indexes to improve introspection endpoint performance
-- These indexes optimize the queries in introspection_client_by_id.sql

-- Composite index for apps7_api_configs: (instance_id, client_id) 
-- This speeds up the WHERE clause in the first part of UNION
CREATE INDEX CONCURRENTLY IF NOT EXISTS apps7_api_configs_instance_client_idx 
ON projections.apps7_api_configs (instance_id, client_id);

-- Composite index for apps7_oidc_configs: (instance_id, client_id)
-- This speeds up the WHERE clause in the second part of UNION
CREATE INDEX CONCURRENTLY IF NOT EXISTS apps7_oidc_configs_instance_client_idx 
ON projections.apps7_oidc_configs (instance_id, client_id);

-- Composite index for apps7: (instance_id, id, state)
-- This speeds up the JOIN condition from config to apps7
CREATE INDEX CONCURRENTLY IF NOT EXISTS apps7_instance_id_state_idx 
ON projections.apps7 (instance_id, id, state);

-- Composite index for projects4: (instance_id, id, state)  
-- This speeds up the JOIN condition from apps7 to projects4
CREATE INDEX CONCURRENTLY IF NOT EXISTS projects4_instance_id_state_idx 
ON projections.projects4 (instance_id, id, state);

-- Composite index for orgs1: (instance_id, id, org_state)
-- This speeds up the JOIN condition from projects4 to orgs1
CREATE INDEX CONCURRENTLY IF NOT EXISTS orgs1_instance_id_state_idx 
ON projections.orgs1 (instance_id, id, org_state);

-- Composite index for authn_keys2: (instance_id, identifier, expiration)
-- This speeds up the keys lookup with all filtering conditions
CREATE INDEX CONCURRENTLY IF NOT EXISTS authn_keys2_instance_identifier_expiration_idx 
ON projections.authn_keys2 (instance_id, identifier, expiration) 
WHERE expiration > current_timestamp;