-- check_permission answers "has user p_user_id of instance p_instance_id permission p_permission",
-- optionally scoped to organization, project, or project grant.
-- Checks are evaluated top-to-bottom; the first matching scope short-circuits via EXISTS.
-- Organization scope also authorizes descendant project and project-grant checks when
-- p_organization_id is provided.
--
-- p_organization_id   – NULL means no organization context.
-- p_project_id        – NULL means no project context.
-- p_project_grant_id  – NULL means no project grant context.
CREATE OR REPLACE FUNCTION zitadel.check_permission(
	p_instance_id TEXT
	, p_user_id TEXT
	, p_permission TEXT

	, p_organization_id TEXT DEFAULT NULL
	, p_project_id TEXT DEFAULT NULL
	, p_project_grant_id TEXT DEFAULT NULL
	, p_raise_if_denied BOOLEAN DEFAULT FALSE
)
RETURNS BOOLEAN
LANGUAGE plpgsql
STABLE
AS $$
DECLARE
	has_permission BOOLEAN;
BEGIN
	SELECT EXISTS (
		SELECT 1
		FROM zitadel.administrators administrator
		JOIN zitadel.administrator_roles administrator_role
			ON administrator_role.instance_id = administrator.instance_id
			AND administrator_role.administrator_id = administrator.id
		JOIN zitadel.administrator_role_permissions role_permission
			ON role_permission.instance_id = administrator_role.instance_id
			AND role_permission.role_name = administrator_role.role_name
		WHERE administrator.instance_id = p_instance_id
			AND administrator.user_id = p_user_id
			AND role_permission.permission = p_permission
			AND (
				-- 1. instance-level: administrator carries the permission across the whole instance
				administrator.scope = 'instance'
				-- 2. organization-level: when supplied, this also covers descendant project and
				--    project-grant checks within the same organization
				OR (
					p_organization_id IS NOT NULL
					AND administrator.scope = 'organization'
					AND administrator.organization_id = p_organization_id
				)
				-- 3. project-level: only when a project context was supplied
				OR (
					p_project_id IS NOT NULL
					AND administrator.scope = 'project'
					AND administrator.project_id = p_project_id
				)
				-- 4. project-grant-level: only when a project grant context was supplied
				OR (
					p_project_grant_id IS NOT NULL
					AND administrator.scope = 'project_grant'
					AND administrator.project_grant_id = p_project_grant_id
				)
			)
	) INTO has_permission;

	IF NOT has_permission AND p_raise_if_denied THEN
		PERFORM zitadel.raise_exception('ZIT01', 'Permission denied');
	END IF;

	RETURN has_permission;
END;
$$;

-- used to raise an error using a condition, 
-- p_id must be a valid SQLSTATE code, e.g. 'ZIT01' for permission denied.
-- e.g. multiple OR conditions and last one calls this function to raise the error if all previous conditions failed.
CREATE OR REPLACE FUNCTION zitadel.raise_exception(
	p_id TEXT
	, p_text TEXT
) RETURNS VOID
LANGUAGE plpgsql
STABLE
AS $$
BEGIN
	RAISE EXCEPTION USING MESSAGE = p_text, ERRCODE = p_id;
END;
$$;