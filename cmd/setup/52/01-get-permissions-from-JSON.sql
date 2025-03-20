DROP FUNCTION IF EXISTS eventstore.get_system_permissions;

CREATE OR REPLACE FUNCTION eventstore.get_system_permissions(
    permissions_json JSONB
    , permm TEXT
)
RETURNS TABLE (
    member_type TEXT,
    aggregate_id TEXT,
    object_id TEXT
)
  LANGUAGE 'plpgsql'
AS $$
BEGIN
    RETURN QUERY
    SELECT res.member_type, res.aggregate_id, res.object_id  FROM (
    SELECT 
        (perm)->>'member_type' AS member_type,
        (perm)->>'aggregate_id' AS aggregate_id,
        (perm)->>'object_id' AS object_id,
         permission-- AS permission
        FROM jsonb_array_elements(permissions_json) AS perm
        CROSS JOIN jsonb_array_elements_text(perm->'permissions') AS permission) AS res
        WHERE res. permission= permm;
END;
$$;

