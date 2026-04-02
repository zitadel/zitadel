"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import {
  getPermissionsForRoles,
  hasPermission,
  hasAnyPermission,
  type RoleKey,
} from "./roles";

interface PermissionContextType {
  /** All resolved permissions for the current user */
  permissions: Set<string>;
  /** The user's role keys */
  roles: string[];
  /** Check if user has a specific permission */
  can: (permission: string) => boolean;
  /** Check if user has ANY of the given permissions */
  canAny: (permissions: string[]) => boolean;
  /** Whether permissions have been loaded */
  isLoaded: boolean;
}

const PermissionContext = createContext<PermissionContextType | undefined>(
  undefined
);

/**
 * PermissionProvider loads the current user's roles and resolves their permissions.
 *
 * For now, roles are passed as a prop (can be loaded from API/token introspection later).
 * The provider resolves the role→permission mapping from defaults.yaml and exposes
 * `can()` / `canAny()` helpers for permission checks throughout the UI.
 */
export function PermissionProvider({
  children,
  initialRoles = [],
}: {
  children: React.ReactNode;
  initialRoles?: string[];
}) {
  const [roles, setRoles] = useState<string[]>(initialRoles);
  const [permissions, setPermissions] = useState<Set<string>>(new Set());
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    const resolved = getPermissionsForRoles(roles);
    setPermissions(resolved);
    setIsLoaded(true);
  }, [roles]);

  const can = (permission: string) => hasPermission(permissions, permission);
  const canAny = (perms: string[]) => hasAnyPermission(permissions, perms);

  return (
    <PermissionContext.Provider
      value={{ permissions, roles, can, canAny, isLoaded }}
    >
      {children}
    </PermissionContext.Provider>
  );
}

/**
 * Hook to access the current user's permissions.
 *
 * Usage:
 *   const { can, canAny } = usePermissions();
 *   if (can("user.write")) { ... }
 *   if (canAny(["org.read", "org.global.read"])) { ... }
 */
export function usePermissions() {
  const context = useContext(PermissionContext);
  if (!context) {
    throw new Error(
      "usePermissions must be used within a PermissionProvider"
    );
  }
  return context;
}
