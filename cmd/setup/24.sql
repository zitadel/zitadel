ALTER TABLE auth.tokens ADD COLUMN actor jsonb;
ALTER TABLE auth.refresh_tokens ADD COLUMN actor jsonb;