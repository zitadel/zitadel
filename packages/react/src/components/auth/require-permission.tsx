"use client";

import { usePermissions } from "../../context/permissions";

/**
 * Wrapper component that only renders children if the user has the required permission.
 *
 * Usage:
 *   <RequirePermission permission="user.read">
 *     <UserTable />
 *   </RequirePermission>
 *
 *   <RequirePermission permission="org.write" fallback={<p>No access</p>}>
 *     <OrgEditForm />
 *   </RequirePermission>
 */
export function RequirePermission({
  children,
  permission,
  anyOf,
  fallback = null,
}: {
  children: React.ReactNode;
  /** Single required permission */
  permission?: string;
  /** Alternative: user needs ANY of these permissions */
  anyOf?: string[];
  /** What to render if user lacks permission */
  fallback?: React.ReactNode;
}) {
  const { can, canAny, isLoaded } = usePermissions();

  if (!isLoaded) {
    return null;
  }

  if (permission && !can(permission)) {
    return <>{fallback}</>;
  }

  if (anyOf && !canAny(anyOf)) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}
