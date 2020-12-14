package authz

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func getUserMethodPermissions(ctx context.Context, t *TokenVerifier, requiredPerm string, authConfig Config, ctxData CtxData) (requestedPermissions, allPermissions []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if ctxData.IsZero() {
		return nil, nil, errors.ThrowUnauthenticated(nil, "AUTH-rKLWEH", "context missing")
	}

	ctx = context.WithValue(ctx, dataKey, ctxData)
	grant, err := t.ResolveGrant(ctx)
	if err != nil {
		return nil, nil, err
	}
	if grant == nil {
		return requestedPermissions, nil, nil
	}
	requestedPermissions, allPermissions = mapGrantToPermissions(requiredPerm, grant, authConfig)
	return requestedPermissions, allPermissions, nil
}

func mapGrantToPermissions(requiredPerm string, grant *Grant, authConfig Config) (requestPermissions, allPermissions []string) {
	requestPermissions = make([]string, 0)
	allPermissions = make([]string, 0)
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
