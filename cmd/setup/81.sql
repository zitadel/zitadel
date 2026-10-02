ALTER TABLE IF EXISTS projections.security_policies3
ADD COLUMN IF NOT EXISTS enable_client_id_metadata_document BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE IF EXISTS projections.security_policies3
ADD COLUMN IF NOT EXISTS client_id_metadata_document_allowed_urls TEXT[];
ALTER TABLE IF EXISTS projections.security_policies3
ADD COLUMN IF NOT EXISTS client_id_metadata_document_allow_any_url BOOLEAN NOT NULL DEFAULT FALSE;
