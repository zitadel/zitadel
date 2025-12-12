import React, { useMemo } from 'react';
import styles from './styles.module.css';

// Note: We no longer need 'js-yaml' here because the loader parses it before the component sees it.

export default function PermissionTable({ data }) {
  const { roles, permissions, matrix } = useMemo(() => {
    if (!data) return { roles: [], permissions: [], matrix: {} };

    // 1. Identify all sources. The 'data' prop is already a JS Object.
    const sources = [
      data?.InternalAuthZ?.RolePermissionMappings,
      data?.SystemAuthZ?.RolePermissionMappings,
    ].filter(Boolean).flat();

    // 2. Aggregate Data
    const roleMap = new Map(); 
    const allPermissions = new Set();

    sources.forEach((mapping) => {
      const roleName = mapping.Role;
      const rolePerms = mapping.Permissions || [];

      if (!roleMap.has(roleName)) {
        roleMap.set(roleName, new Set());
      }

      rolePerms.forEach((perm) => {
        roleMap.get(roleName).add(perm);
        allPermissions.add(perm);
      });
    });

    // 3. Sort for display
    const sortedRoles = Array.from(roleMap.keys()).sort();
    const sortedPermissions = Array.from(allPermissions).sort();

    return {
      roles: sortedRoles,
      permissions: sortedPermissions,
      matrix: roleMap,
    };
  }, [data]);

  return (
    <div className={styles.container}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>Permission</th>
            {roles.map((role) => (
              <th key={role}>{role}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {permissions.map((perm) => (
            <tr key={perm}>
              <th>{perm}</th>
              {roles.map((role) => {
                const hasPermission = matrix.get(role).has(perm);
                return (
                  <td 
                    key={`${role}-${perm}`} 
                    className={hasPermission ? styles.yes : styles.no}
                  >
                    {hasPermission ? 'yes' : 'no'}
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}