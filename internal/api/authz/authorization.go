package authz

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	authenticated = "authenticated"
)

// CheckUserAuthorization verifies that:
// - the token is active,
// - the organization (**either** provided by ID or verified domain) exists
// - the user is permitted to call the requested endpoint (permission option in proto)
// it will pass the [CtxData] and permission of the user into the ctx [context.Context]
func CheckUserAuthorization(ctx context.Context, req interface{}, token, orgID, orgDomain string, verifier APITokenVerifier, systemRolePermissionMapping []RoleMapping, rolePermissionMapping []RoleMapping, requiredAuthOption Option, method string) (ctxSetter func(context.Context) context.Context, err error) {
	ctx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	ctxData, err := VerifyTokenAndCreateCtxData(ctx, token, orgID, orgDomain, verifier)
	if err != nil {
		return nil, err
	}

	if requiredAuthOption.Permission == authenticated {
		return func(parent context.Context) context.Context {
			parent = addGetSystemUserRolesToCtx(parent, systemRolePermissionMapping, ctxData)
			return context.WithValue(parent, dataKey, ctxData)
		}, nil
	}

	requestedPermissions, allPermissions, err := getUserPermissions(ctx, verifier, requiredAuthOption.Permission, systemRolePermissionMapping, rolePermissionMapping, ctxData, ctxData.OrgID)
	if err != nil {
		return nil, err
	}

	ctx, userPermissionSpan := tracing.NewNamedSpan(ctx, "checkUserPermissions")
	err = checkUserPermissions(req, requestedPermissions, requiredAuthOption)
	userPermissionSpan.EndWithError(err)
	if err != nil {
		return nil, err
	}

	return func(parent context.Context) context.Context {
		parent = context.WithValue(parent, dataKey, ctxData)
		parent = context.WithValue(parent, allPermissionsKey, allPermissions)
		parent = context.WithValue(parent, requestPermissionsKey, requestedPermissions)
		parent = addGetSystemUserRolesToCtx(parent, systemRolePermissionMapping, ctxData)
		return parent
	}, nil
}

func checkUserPermissions(req interface{}, userPerms []string, authOpt Option) error {
	if len(userPerms) == 0 {
		return zerrors.ThrowPermissionDenied(nil, "AUTH-5mWD2", "No matching permissions found")
	}

	if authOpt.CheckParam == "" {
		return nil
	}

	if HasGlobalPermission(userPerms) {
		return nil
	}

	if hasContextPermission(req, authOpt.CheckParam, userPerms) {
		return nil
	}

	return zerrors.ThrowPermissionDenied(nil, "AUTH-3jknH", "No matching permissions found")
}

func SplitPermission(perm string) (string, string) {
	splittedPerm := strings.Split(perm, ":")
	if len(splittedPerm) == 1 {
		return splittedPerm[0], ""
	}
	return splittedPerm[0], splittedPerm[1]
}

func hasContextPermission(req interface{}, fieldName string, permissions []string) bool {
	for _, perm := range permissions {
		_, ctxID := SplitPermission(perm)
		if checkPermissionContext(req, fieldName, ctxID) {
			return true
		}
	}
	return false
}

func checkPermissionContext(req interface{}, fieldName, roleContextID string) bool {
	field := getFieldFromReq(req, fieldName)
	return field != "" && field == roleContextID
}

func getFieldFromReq(req interface{}, field string) string {
	v := reflect.Indirect(reflect.ValueOf(req)).FieldByName(field)
	if reflect.ValueOf(v).IsZero() {
		return ""
	}
	return fmt.Sprintf("%v", v.Interface())
}

func HasGlobalPermission(perms []string) bool {
	for _, perm := range perms {
		_, ctxID := SplitPermission(perm)
		if ctxID == "" {
			return true
		}
	}
	return false
}

func GetAllPermissionCtxIDs(perms []string) []string {
	ctxIDs := make([]string, 0)
	for _, perm := range perms {
		_, ctxID := SplitPermission(perm)
		if ctxID != "" {
			ctxIDs = append(ctxIDs, ctxID)
		}
	}
	return ctxIDs
}

type SystemUserPermissionsDBQuery struct {
	MemberType  string   `json:"member_type"`
	AggregateID string   `json:"aggregate_id"`
	ObjectID    string   `json:"object_id"`
	Permissions []string `json:"permissions"`
}

func addGetSystemUserRolesToCtx(ctx context.Context, systemUserRoleMap []RoleMapping, ctxData CtxData) context.Context {
	if len(ctxData.SystemMemberships) == 0 {
		return ctx
	}
	systemUserPermissions := make([]SystemUserPermissionsDBQuery, len(ctxData.SystemMemberships))
	for i, systemPerm := range ctxData.SystemMemberships {
		permissions := make([]string, 0, len(systemPerm.Roles))
		for _, role := range systemPerm.Roles {
			permissions = append(permissions, getPermissionsFromRole(systemUserRoleMap, role)...)
		}
		slices.Sort(permissions)
		permissions = slices.Compact(permissions)

		systemUserPermissions[i].MemberType = systemPerm.MemberType.String()
		systemUserPermissions[i].AggregateID = systemPerm.AggregateID
		systemUserPermissions[i].Permissions = permissions
	}
	return context.WithValue(ctx, systemUserRolesFuncKey, systemUserPermissions)
}

func GetSystemUserPermissions(ctx context.Context) []SystemUserPermissionsDBQuery {
	getSystemUserRolesFuncValue := ctx.Value(systemUserRolesFuncKey)
	if getSystemUserRolesFuncValue == nil {
		return nil
	}
	systemUserRoles, ok := getSystemUserRolesFuncValue.([]SystemUserPermissionsDBQuery)
	if !ok {
		logging.WithFields("Authz").Error("unable to cast []SystemUserPermissionsDBQuery")
		return nil
	}
	return systemUserRoles
}
