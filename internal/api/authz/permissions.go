package authz

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func getUserMethodPermissions(ctx context.Context, t *TokenVerifier, requiredPerm string, authConfig Config) (_ context.Context, _ []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	ctxData := GetCtxData(ctx)
	if ctxData.IsZero() {
		return nil, nil, errors.ThrowUnauthenticated(nil, "AUTH-rKLWEH", "context missing")
	}
	grant, err := t.ResolveGrant(ctx)
	if err != nil {
		return nil, nil, err
	}
	if grant == nil {
		return context.WithValue(ctx, requestPermissionsKey, []string{}), []string{}, nil
	}
	requestPermissions, allPermissions := mapGrantToPermissions(requiredPerm, grant, authConfig)
	ctx = context.WithValue(ctx, allPermissionsKey, allPermissions)
	return context.WithValue(ctx, requestPermissionsKey, requestPermissions), requestPermissions, nil
}

func mapGrantToPermissions(requiredPerm string, grant *Grant, authConfig Config) ([]string, []string) {
	requestPermissions := make([]string, 0)
	allPermissions := make([]string, 0)
	for _, role := range grant.Roles {
		requestPermissions, allPermissions = mapRoleToPerm(requiredPerm, role, authConfig, requestPermissions, allPermissions)
	}

	return requestPermissions, allPermissions
}

func mapRoleToPerm(requiredPerm, actualRole string, authConfig Config, requestPermissions, allPermissions []string) ([]string, []string) {
	roleName, roleContextID := SplitPermission(actualRole)
	perms := authConfig.getPermissionsFromRole(roleName)

	for _, p := range perms {
		permWithCtx := addRoleContextIDToPerm(p, roleContextID)
		if !ExistsPerm(allPermissions, permWithCtx) {
			allPermissions = append(allPermissions, permWithCtx)
		}

		p, _ = SplitPermission(p)
		if p == requiredPerm {
			if !ExistsPerm(requestPermissions, permWithCtx) {
				requestPermissions = append(requestPermissions, permWithCtx)
			}
		}
	}
	return requestPermissions, allPermissions
}

func addRoleContextIDToPerm(perm, roleContextID string) string {
	if roleContextID != "" {
		perm = perm + ":" + roleContextID
	}
	return perm
}

func ExistsPerm(existingPermissions []string, perm string) bool {
	for _, existingPermission := range existingPermissions {
		if existingPermission == perm {
			return true
		}
	}
	return false
}
