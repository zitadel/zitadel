package authz

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func CheckPermission(ctx context.Context, resolver MembershipsResolver, roleMappings []RoleMapping, permission, orgID, resourceID string) (err error) {
	requestedPermissions, _, err := getUserPermissions(ctx, resolver, permission, roleMappings, GetCtxData(ctx), orgID)
	if err != nil {
		return err
	}

	_, userPermissionSpan := tracing.NewNamedSpan(ctx, "checkUserPermissions")
	err = checkUserResourcePermissions(requestedPermissions, resourceID)
	userPermissionSpan.EndWithError(err)

	return err
}

// getUserPermissions retrieves the memberships of the authenticated user (on instance and provided organisation level),
// and maps them to permissions. It will return the requested permission(s) and all other granted permissions separately.
func getUserPermissions(ctx context.Context, resolver MembershipsResolver, requiredPerm string, roleMappings []RoleMapping, ctxData CtxData, orgID string) (requestedPermissions, allPermissions []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if ctxData.IsZero() {
		return nil, nil, errors.ThrowUnauthenticated(nil, "AUTH-rKLWEH", "context missing")
	}

	if ctxData.SystemMemberships != nil {
		requestedPermissions, allPermissions = mapMembershipsToPermissions(requiredPerm, ctxData.SystemMemberships, roleMappings)
		return requestedPermissions, allPermissions, nil
	}

	ctx = context.WithValue(ctx, dataKey, ctxData)
	memberships, err := resolver.SearchMyMemberships(ctx, orgID, false)
	if err != nil {
		return nil, nil, err
	}
	if len(memberships) == 0 {
		memberships, err = resolver.SearchMyMemberships(ctx, orgID, true)
		if len(memberships) == 0 {
			return nil, nil, errors.ThrowNotFound(nil, "AUTHZ-cdgFk", "membership not found")
		}
		if err != nil {
			return nil, nil, err
		}
	}
	requestedPermissions, allPermissions = mapMembershipsToPermissions(requiredPerm, memberships, roleMappings)
	return requestedPermissions, allPermissions, nil
}

// checkUserResourcePermissions checks that if a user i granted either the requested permission globally (project.write)
// or the specific resource (project.write:123)
func checkUserResourcePermissions(userPerms []string, resourceID string) error {
	if len(userPerms) == 0 {
		return errors.ThrowPermissionDenied(nil, "AUTH-AWfge", "No matching permissions found")
	}

	if resourceID == "" {
		return nil
	}

	if HasGlobalPermission(userPerms) {
		return nil
	}

	if hasContextResourcePermission(userPerms, resourceID) {
		return nil
	}

	return errors.ThrowPermissionDenied(nil, "AUTH-Swrgg2", "No matching permissions found")
}

func hasContextResourcePermission(permissions []string, resourceID string) bool {
	for _, perm := range permissions {
		_, ctxID := SplitPermission(perm)
		if resourceID == ctxID {
			return true
		}
	}
	return false
}

func mapMembershipsToPermissions(requiredPerm string, memberships []*Membership, roleMappings []RoleMapping) (requestPermissions, allPermissions []string) {
	requestPermissions = make([]string, 0)
	allPermissions = make([]string, 0)
	for _, membership := range memberships {
		requestPermissions, allPermissions = mapMembershipToPerm(requiredPerm, membership, roleMappings, requestPermissions, allPermissions)
	}

	return requestPermissions, allPermissions
}

func mapMembershipToPerm(requiredPerm string, membership *Membership, roleMappings []RoleMapping, requestPermissions, allPermissions []string) ([]string, []string) {
	roleNames, roleContextID := roleWithContext(membership)
	for _, roleName := range roleNames {
		perms := getPermissionsFromRole(roleMappings, roleName)

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
