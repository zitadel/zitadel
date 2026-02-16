import React, { useMemo } from 'react';
import yaml from 'js-yaml';
import { Check, X } from 'lucide-react';
import { cn } from '@/utils/cn'; // Assuming utils/cn exists, or I'll use clsx/tailwind-merge directly if not sure

// Fallback utility if generic cn doesn't exist, checking imports first is safer, but I'll inline a simple merge if needed
// Checking package.json I saw clsx and tailwind-merge are dependencies.

function classNames(...classes) {
  return classes.filter(Boolean).join(' ');
}

export default function PermissionTable({ data }) {
  const { roles, permissions, matrix } = useMemo(() => {
    if (!data) return { roles: [], permissions: [], matrix: {} };

    let parsedData = data;
    if (typeof data === 'string') {
      try {
        parsedData = yaml.load(data);
      } catch (e) {
        console.error('Failed to parse YAML data:', e);
        return { roles: [], permissions: [], matrix: {} };
      }
    }

    // 1. Identify all sources.
    const sources = [
      parsedData?.InternalAuthZ?.RolePermissionMappings,
      parsedData?.SystemAuthZ?.RolePermissionMappings,
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

  if (!roles.length) {
    return <div className="p-4 text-sm text-fd-muted-foreground bg-fd-muted rounded-md border border-fd-border">No permission data available.</div>;
  }

  return (
    <div className="overflow-x-auto my-6 rounded-lg border border-fd-border bg-fd-card shadow-sm">
      <table className="w-full text-sm text-left">
        <thead className="text-xs font-semibold text-fd-muted-foreground uppercase bg-fd-muted/50 border-b border-fd-border">
          <tr>
            <th scope="col" className="px-4 py-3 sticky left-0 z-10 bg-fd-muted/95 backdrop-blur shadow-[1px_0_0_0_rgba(0,0,0,0.1)] dark:shadow-[1px_0_0_0_rgba(255,255,255,0.1)]">
              Permission
            </th>
            {roles.map((role) => (
              <th key={role} scope="col" className="px-4 py-3 whitespace-nowrap text-center min-w-[100px]">
                {role.replace(/_/g, ' ')}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-fd-border">
          {permissions.map((perm) => (
            <tr key={perm} className="hover:bg-fd-accent/50 transition-colors">
              <th scope="row" className="px-4 py-2 font-medium text-fd-foreground whitespace-nowrap sticky left-0 z-10 bg-fd-card shadow-[1px_0_0_0_rgba(0,0,0,0.1)] dark:shadow-[1px_0_0_0_rgba(255,255,255,0.1)]">
                {perm}
              </th>
              {roles.map((role) => {
                const hasPermission = matrix.get(role).has(perm);
                return (
                  <td
                    key={`${role}-${perm}`}
                    className={classNames(
                      "px-4 py-2 text-center",
                      hasPermission ? "bg-green-50/50 dark:bg-green-900/10" : ""
                    )}
                  >
                    {hasPermission ? (
                      <Check className="w-4 h-4 text-green-600 dark:text-green-400 mx-auto" strokeWidth={3} />
                    ) : (
                      <span className="sr-only">No</span>
                    )}
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