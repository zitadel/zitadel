ALTER INDEX IF EXISTS projections.user_metadata5_metadata_key_idx
    RENAME TO user_metadata5_key_idx;
DROP INDEX IF EXISTS projections.user_metadata5_metadata_value_idx;