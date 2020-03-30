package auth

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
)

func getUserMethodPermissions(ctx context.Context, t TokenVerifier, requiredPerm string, authConfig *Config) (context.Context, []string, error) {
	ctxData := GetCtxData(ctx)
	if ctxData.IsZero() {
		return nil, nil, errors.ThrowUnauthenticated(nil, "AUTH-rKLWEH", "context missing")
	}
	grants, err := t.ResolveGrants(ctx, ctxData.UserID, ctxData.OrgID)
	if err != nil {
		return nil, nil, err
	}
	permissions := mapGrantsToPermissions(requiredPerm, grants, authConfig)
	return context.WithValue(ctx, permissionsKey, permissions), permissions, nil
}

func mapGrantsToPermissions(requiredPerm string, grants []*Grant, authConfig *Config) []string {
	resolvedPermissions := make([]string, 0)
	for _, grant := range grants {
		for _, role := range grant.Roles {
			resolvedPermissions = mapRoleToPerm(requiredPerm, role, authConfig, resolvedPermissions)
		}
	}
	return resolvedPermissions
}

func mapRoleToPerm(requiredPerm, actualRole string, authConfig *Config, resolvedPermissions []string) []string {
	roleName, roleContextID := SplitPermission(actualRole)
	perms := authConfig.getPermissionsFromRole(roleName)

	for _, p := range perms {
		if p == requiredPerm {
			p = addRoleContextIDToPerm(p, roleContextID)
			if !existsPerm(resolvedPermissions, p) {
				resolvedPermissions = append(resolvedPermissions, p)
			}
		}
	}
	return resolvedPermissions
}

func addRoleContextIDToPerm(perm, roleContextID string) string {
	if roleContextID != "" {
		perm = perm + ":" + roleContextID
	}
	return perm
}

func existsPerm(existing []string, perm string) bool {
	for _, e := range existing {
		if e == perm {
			return true
		}
	}
	return false
}
