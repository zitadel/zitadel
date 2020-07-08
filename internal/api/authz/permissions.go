package authz

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
)

func getUserMethodPermissions(ctx context.Context, t *TokenVerifier, requiredPerm string, authConfig Config) (context.Context, []string, error) {
	ctxData := GetCtxData(ctx)
	if ctxData.IsZero() {
		return nil, nil, errors.ThrowUnauthenticated(nil, "AUTH-rKLWEH", "context missing")
	}
	grant, err := t.ResolveGrant(ctx)
	if err != nil {
		return nil, nil, err
	}
	if grant == nil {
		return context.WithValue(ctx, permissionsKey, []string{}), []string{}, nil
	}
	permissions := mapGrantToPermissions(requiredPerm, grant, authConfig)
	return context.WithValue(ctx, permissionsKey, permissions), permissions, nil
}

func mapGrantToPermissions(requiredPerm string, grant *Grant, authConfig Config) []string {
	resolvedPermissions := make([]string, 0)
	for _, role := range grant.Roles {
		resolvedPermissions = mapRoleToPerm(requiredPerm, role, authConfig, resolvedPermissions)
	}

	return resolvedPermissions
}

func mapRoleToPerm(requiredPerm, actualRole string, authConfig Config, resolvedPermissions []string) []string {
	roleName, roleContextID := SplitPermission(actualRole)
	perms := authConfig.getPermissionsFromRole(roleName)

	for _, p := range perms {
		if p == requiredPerm {
			p = addRoleContextIDToPerm(p, roleContextID)
			if !ExistsPerm(resolvedPermissions, p) {
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

func ExistsPerm(existing []string, perm string) bool {
	for _, e := range existing {
		if e == perm {
			return true
		}
	}
	return false
}
