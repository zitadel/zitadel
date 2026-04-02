/**
 * ZITADEL role keys and their permission mappings.
 * Sourced from cmd/defaults.yaml InternalAuthZ.RolePermissionMappings.
 *
 * This is used to gate UI elements: nav items, buttons, and pages
 * based on the authenticated user's roles.
 */

/** All known ZITADEL role keys */
export type RoleKey =
  | "SYSTEM_OWNER"
  | "SYSTEM_OWNER_VIEWER"
  | "IAM_OWNER"
  | "IAM_OWNER_VIEWER"
  | "IAM_ORG_MANAGER"
  | "IAM_USER_MANAGER"
  | "IAM_ADMIN_IMPERSONATOR"
  | "IAM_END_USER_IMPERSONATOR"
  | "IAM_LOGIN_CLIENT"
  | "ORG_OWNER"
  | "ORG_OWNER_VIEWER"
  | "ORG_SETTINGS_MANAGER"
  | "ORG_USER_MANAGER"
  | "ORG_USER_PERMISSION_EDITOR"
  | "ORG_PROJECT_PERMISSION_EDITOR"
  | "ORG_PROJECT_CREATOR"
  | "ORG_ADMIN_IMPERSONATOR"
  | "ORG_END_USER_IMPERSONATOR"
  | "PROJECT_OWNER"
  | "PROJECT_OWNER_VIEWER"
  | "PROJECT_OWNER_GLOBAL"
  | "PROJECT_OWNER_VIEWER_GLOBAL"
  | "PROJECT_GRANT_OWNER";

/**
 * Role → Permissions mapping from defaults.yaml.
 * Only the roles relevant for console navigation are included.
 * Full list in cmd/defaults.yaml lines 1415-1900+.
 */
export const ROLE_PERMISSIONS: Record<string, readonly string[]> = {
  IAM_OWNER: [
    "iam.read", "iam.write", "iam.policy.read", "iam.policy.write", "iam.policy.delete",
    "iam.member.read", "iam.member.write", "iam.member.delete",
    "iam.idp.read", "iam.idp.write", "iam.idp.delete",
    "iam.action.read", "iam.action.write", "iam.action.delete",
    "iam.flow.read", "iam.flow.write", "iam.flow.delete",
    "iam.feature.read", "iam.feature.write", "iam.feature.delete",
    "iam.restrictions.read", "iam.restrictions.write",
    "iam.web_key.read", "iam.web_key.write", "iam.web_key.delete",
    "org.read", "org.global.read", "org.create", "org.write", "org.delete",
    "org.member.read", "org.member.write", "org.member.delete",
    "org.idp.read", "org.idp.write", "org.idp.delete",
    "org.action.read", "org.action.write", "org.action.delete",
    "org.flow.read", "org.flow.write", "org.flow.delete",
    "org.feature.read", "org.feature.write", "org.feature.delete",
    "user.read", "user.global.read", "user.write", "user.delete",
    "user.grant.read", "user.grant.write", "user.grant.delete",
    "user.membership.read", "user.credential.write", "user.passkey.write",
    "user.feature.read", "user.feature.write", "user.feature.delete",
    "policy.read", "policy.write", "policy.delete",
    "project.read", "project.create", "project.write", "project.delete",
    "project.member.read", "project.member.write", "project.member.delete",
    "project.role.read", "project.role.write", "project.role.delete",
    "project.app.read", "project.app.write", "project.app.delete",
    "project.grant.read", "project.grant.write", "project.grant.delete",
    "project.grant.member.read", "project.grant.member.write", "project.grant.member.delete",
    "events.read", "milestones.read",
    "session.read", "session.write", "session.delete",
    "action.target.read", "action.target.write", "action.target.delete",
    "action.execution.read", "action.execution.write",
    "userschema.read", "userschema.write", "userschema.delete",
  ],
  IAM_OWNER_VIEWER: [
    "iam.read", "iam.policy.read", "iam.member.read", "iam.idp.read",
    "iam.action.read", "iam.flow.read", "iam.restrictions.read", "iam.feature.read",
    "iam.web_key.read",
    "org.read", "org.member.read", "org.idp.read", "org.action.read", "org.flow.read", "org.feature.read",
    "user.read", "user.global.read", "user.grant.read", "user.membership.read", "user.feature.read",
    "policy.read",
    "project.read", "project.member.read", "project.role.read", "project.app.read",
    "project.grant.read", "project.grant.member.read",
    "events.read", "milestones.read",
    "action.target.read", "action.execution.read",
    "userschema.read", "session.read",
  ],
  IAM_ORG_MANAGER: [
    "org.read", "org.global.read", "org.create", "org.write", "org.delete",
    "org.member.read", "org.member.write", "org.member.delete",
    "org.idp.read", "org.idp.write", "org.idp.delete",
    "org.action.read", "org.action.write", "org.action.delete",
    "org.flow.read", "org.flow.write", "org.flow.delete",
    "org.feature.read", "org.feature.write", "org.feature.delete",
    "user.read", "user.global.read", "user.write", "user.delete",
    "user.grant.read", "user.grant.write", "user.grant.delete",
    "user.membership.read", "user.credential.write", "user.passkey.write",
    "user.feature.read", "user.feature.write", "user.feature.delete",
    "policy.read", "policy.write", "policy.delete",
    "project.read", "project.create", "project.write", "project.delete",
    "project.member.read", "project.member.write", "project.member.delete",
    "project.role.read", "project.role.write", "project.role.delete",
    "project.app.read", "project.app.write", "project.app.delete",
    "project.grant.read", "project.grant.write", "project.grant.delete",
    "project.grant.member.read", "project.grant.member.write", "project.grant.member.delete",
    "session.read", "session.delete",
  ],
  IAM_USER_MANAGER: [
    "org.read", "org.global.read", "org.member.read", "org.member.delete",
    "user.read", "user.global.read", "user.write", "user.delete",
    "user.grant.read", "user.grant.write", "user.grant.delete",
    "user.membership.read", "user.passkey.write",
    "user.feature.read", "user.feature.write", "user.feature.delete",
    "project.read", "project.member.read", "project.role.read", "project.app.read",
    "project.grant.read", "project.grant.write", "project.grant.delete",
    "project.grant.member.read",
    "session.read", "session.delete",
  ],
  ORG_OWNER: [
    "org.read", "org.global.read", "org.write", "org.delete",
    "org.member.read", "org.member.write", "org.member.delete",
    "org.idp.read", "org.idp.write", "org.idp.delete",
    "org.action.read", "org.action.write", "org.action.delete",
    "org.flow.read", "org.flow.write", "org.flow.delete",
    "org.feature.read", "org.feature.write", "org.feature.delete",
    "user.read", "user.global.read", "user.write", "user.delete",
    "user.grant.read", "user.grant.write", "user.grant.delete",
    "user.membership.read", "user.credential.write", "user.passkey.write",
    "user.feature.read", "user.feature.write", "user.feature.delete",
    "policy.read", "policy.write", "policy.delete",
    "project.read", "project.create", "project.write", "project.delete",
    "project.member.read", "project.member.write", "project.member.delete",
    "project.role.read", "project.role.write", "project.role.delete",
    "project.app.read", "project.app.write",
    "project.grant.read", "project.grant.write", "project.grant.delete",
    "project.grant.member.read", "project.grant.member.write", "project.grant.member.delete",
    "session.read", "session.delete",
  ],
  ORG_OWNER_VIEWER: [
    "org.read", "org.member.read", "org.idp.read", "org.action.read", "org.flow.read", "org.feature.read",
    "user.read", "user.global.read", "user.grant.read", "user.membership.read", "user.feature.read",
    "policy.read",
    "project.read", "project.member.read", "project.role.read", "project.app.read",
    "project.grant.read", "project.grant.member.read",
  ],
  ORG_USER_MANAGER: [
    "org.read",
    "user.read", "user.global.read", "user.write", "user.delete",
    "user.grant.read", "user.grant.write", "user.grant.delete",
    "user.membership.read",
    "user.feature.read", "user.feature.write", "user.feature.delete",
    "policy.read",
    "project.read", "project.role.read",
    "session.read", "session.delete",
  ],
  ORG_SETTINGS_MANAGER: [
    "org.read", "org.write", "org.member.read",
    "org.idp.read", "org.idp.write", "org.idp.delete",
    "org.feature.read", "org.feature.write", "org.feature.delete",
    "policy.read", "policy.write", "policy.delete",
  ],
  ORG_USER_PERMISSION_EDITOR: [
    "org.read", "org.member.read",
    "user.read", "user.global.read", "user.grant.read", "user.grant.write", "user.grant.delete",
    "policy.read",
    "project.read", "project.member.read", "project.role.read", "project.app.read",
    "project.grant.read", "project.grant.member.read",
  ],
  ORG_PROJECT_PERMISSION_EDITOR: [
    "org.read", "org.member.read",
    "user.read", "user.global.read", "user.grant.read", "user.grant.write", "user.grant.delete",
    "policy.read",
    "project.read", "project.member.read", "project.role.read", "project.app.read",
    "project.grant.read", "project.grant.write", "project.grant.delete",
    "project.grant.member.read", "project.grant.member.write", "project.grant.member.delete",
  ],
} as const;

/**
 * Resolve permissions from a list of role keys.
 * Returns a Set of all permissions the user has.
 */
export function getPermissionsForRoles(roles: string[]): Set<string> {
  const permissions = new Set<string>();
  for (const role of roles) {
    const rolePerms = ROLE_PERMISSIONS[role];
    if (rolePerms) {
      for (const perm of rolePerms) {
        permissions.add(perm);
      }
    }
  }
  return permissions;
}

/**
 * Check if a set of permissions includes a required permission.
 */
export function hasPermission(
  userPermissions: Set<string>,
  required: string
): boolean {
  return userPermissions.has(required);
}

/**
 * Check if a set of permissions includes ANY of the required permissions.
 */
export function hasAnyPermission(
  userPermissions: Set<string>,
  required: string[]
): boolean {
  return required.some((p) => userPermissions.has(p));
}
