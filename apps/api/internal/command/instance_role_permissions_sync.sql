/*
This query creates a change set of permissions that need to be added or removed.
It compares the current state in the fields table (thru the role_permissions view)
against a passed role permission mapping as JSON, created from Zitadel's config:

{
	"IAM_ADMIN_IMPERSONATOR": ["admin.impersonation", "impersonation"],
	"IAM_END_USER_IMPERSONATOR": ["impersonation"],
	"FOO_BAR": ["foo.bar", "bar.foo"]
 }

It uses an aggregate_id as first argument which may be an instance_id or 'SYSTEM'
for system level permissions.
*/
WITH target AS (
	-- unmarshal JSON representation into flattened tabular data
	SELECT
		key AS role,
		jsonb_array_elements_text(value) AS permission
	FROM jsonb_each($2::jsonb)
), add AS (
    -- find all role permissions that exist in `target` and not in `role_permissions`
	SELECT t.role, t.permission
	FROM eventstore.role_permissions p
	RIGHT JOIN target t
		ON p.aggregate_id = $1::text
		AND p.role = t.role
		AND p.permission = t.permission
	WHERE p.role IS NULL
), remove AS (
    -- find all role permissions that exist `role_permissions` and not in `target`
	SELECT p.role, p.permission
	FROM eventstore.role_permissions p
	LEFT JOIN target t
		ON p.role = t.role
		AND p.permission = t.permission
	WHERE p.aggregate_id = $1::text
	AND t.role IS NULL
)
-- return the required operations
SELECT
	'add' AS operation,
	role,
	permission
FROM add
UNION ALL
SELECT
	'remove' AS operation,
	role,
	permission
FROM remove
;
