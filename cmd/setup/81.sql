ALTER TABLE IF EXISTS projections.security_policies3
ADD COLUMN IF NOT EXISTS enable_client_id_metadata_document BOOLEAN NOT NULL DEFAULT FALSE;
