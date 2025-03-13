package authz

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	authenticated = "authenticated"
)

// CheckUserAuthorization verifies that:
// - the token is active,
// - the organisation (**either** provided by ID or verified domain) exists
// - the user is permitted to call the requested endpoint (permission option in proto)
// it will pass the [CtxData] and permission of the user into the ctx [context.Context]
func CheckUserAuthorization(ctx context.Context, req interface{}, token, orgID, orgDomain string, verifier APITokenVerifier, SystemAuthConfig Config, authConfig Config, requiredAuthOption Option, method string) (ctxSetter func(context.Context) context.Context, err error) {
	ctx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	ctxData, err := VerifyTokenAndCreateCtxData(ctx, token, orgID, orgDomain, verifier)
	if err != nil {
		return nil, err
	}

	if requiredAuthOption.Permission == authenticated {
		return func(parent context.Context) context.Context {
			parent = addGetSystemUserRolesFuncToCtx(parent, SystemAuthConfig.RolePermissionMappings, nil, ctxData)
			return context.WithValue(parent, dataKey, ctxData)
		}, nil
	}

	requestedPermissions, allPermissions, err := getUserPermissions(ctx, verifier, requiredAuthOption.Permission, SystemAuthConfig.RolePermissionMappings, authConfig.RolePermissionMappings, ctxData, ctxData.OrgID)
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
		parent = addGetSystemUserRolesFuncToCtx(parent, SystemAuthConfig.RolePermissionMappings, requestedPermissions, ctxData)
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

type SystemUserAuthParams struct {
	MemberType        []int32
	InstanceID        []string
	AggregateID       []string
	Permissions       [][]string
	PermissionsLength []int32
}

func addGetSystemUserRolesFuncToCtx(ctx context.Context, systemUserRoleMap []RoleMapping, requestedPermissions []string, ctxData CtxData) context.Context {
	if len(ctxData.SystemMemberships) == 0 {
		return ctx
		// } else if ctxData.SystemMemberships[0].MemberType == MemberTypeSystem {
	} else {
		ctx = context.WithValue(ctx, systemUserRolesFuncKey, func() func(ctx context.Context) *SystemUserAuthParams {
			var systemUserAuthParams *SystemUserAuthParams
			chann := make(chan struct{}, 1)
			return func(ctx context.Context) *SystemUserAuthParams {
				chann <- struct{}{}
				if systemUserAuthParams != nil {
					return systemUserAuthParams
				}
				defer func() {
					<-chann
					close(chann)
				}()

				systemUserAuthParams = &SystemUserAuthParams{
					MemberType:        make([]int32, len(ctxData.SystemMemberships)),
					InstanceID:        make([]string, len(ctxData.SystemMemberships)),
					AggregateID:       make([]string, len(ctxData.SystemMemberships)),
					Permissions:       make([][]string, len(ctxData.SystemMemberships)),
					PermissionsLength: make([]int32, len(ctxData.SystemMemberships)),
				}

				for i, systemPerm := range ctxData.SystemMemberships {
					permissions := []string{}
					for _, role := range systemPerm.Roles {
						permissions = append(permissions, getPermissionsFromRole(systemUserRoleMap, role)...)
					}
					slices.Sort(permissions)
					permissions = slices.Compact(permissions)

					systemUserAuthParams.MemberType[i] = MemberTypeServerToMemberTypeDBMap[systemPerm.MemberType]
					systemUserAuthParams.InstanceID[i] = systemPerm.InstanceID
					systemUserAuthParams.AggregateID[i] = systemPerm.AggregateID
					systemUserAuthParams.Permissions[i] = permissions
					systemUserAuthParams.PermissionsLength[i] = int32(len(permissions))
				}
				return systemUserAuthParams
			}
		}())
	}
	return ctx
}

func GetSystemUserAuthParams(ctx context.Context) (*SystemUserAuthParams, error) {
	getSystemUserRolesFuncValue := ctx.Value(systemUserRolesFuncKey)
	if getSystemUserRolesFuncValue == nil {
		return &SystemUserAuthParams{}, nil
	}
	getSystemUserRolesFunc, ok := getSystemUserRolesFuncValue.(func(context.Context) *SystemUserAuthParams)
	if !ok {
		return nil, errors.New("unable to obtain systems role func")
	}
	return getSystemUserRolesFunc(ctx), nil
}
