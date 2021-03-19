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
	memberships, err := t.SearchMyMemberships(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(memberships) == 0 {
		return requestedPermissions, nil, nil
	}
	requestedPermissions, allPermissions = mapMembershipsToPermissions(requiredPerm, memberships, authConfig)
	return requestedPermissions, allPermissions, nil
}

func mapMembershipsToPermissions(requiredPerm string, memberships []*Membership, authConfig Config) (requestPermissions, allPermissions []string) {
	requestPermissions = make([]string, 0)
	allPermissions = make([]string, 0)
	for _, membership := range memberships {
		requestPermissions, allPermissions = mapMembershipToPerm(requiredPerm, membership, authConfig, requestPermissions, allPermissions)
	}

	return requestPermissions, allPermissions
}

func mapMembershipToPerm(requiredPerm string, membership *Membership, authConfig Config, requestPermissions, allPermissions []string) ([]string, []string) {
	roleNames, roleContextID := roleWithContext(membership)
	for _, roleName := range roleNames {
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

func roleWithContext(membership *Membership) (roles []string, ctxID string) {
	if membership.MemberType == MemberTypeProject || membership.MemberType == MemberTypeProjectGrant {
		return membership.Roles, membership.ObjectID
	}
	return membership.Roles, ""
}
