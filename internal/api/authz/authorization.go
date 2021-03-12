package authz

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

const (
	authenticated = "authenticated"
)

func CheckUserAuthorization(ctx context.Context, req interface{}, token, orgID string, verifier *TokenVerifier, authConfig Config, requiredAuthOption Option, method string) (ctxSetter func(context.Context) context.Context, err error) {
	ctx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	ctxData, err := VerifyTokenAndCreateCtxData(ctx, token, orgID, verifier, method)
	if err != nil {
		return nil, err
	}

	if requiredAuthOption.Feature != "" {
		err = checkOrgFeatures(ctx, verifier, requiredAuthOption.Feature, ctxData)
		if err != nil {
			return nil, err
		}
	}

	if requiredAuthOption.Permission == authenticated {
		return func(parent context.Context) context.Context {
			return context.WithValue(parent, dataKey, ctxData)
		}, nil
	}

	requestedPermissions, allPermissions, err := getUserMethodPermissions(ctx, verifier, requiredAuthOption.Permission, authConfig, ctxData)
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
		return parent
	}, nil
}

func checkOrgFeatures(ctx context.Context, t *TokenVerifier, requiredFeature string, ctxData CtxData) error {
	//TODO: impl
	fmt.Println(requiredFeature, ctxData.OrgID)
	return nil
}

func checkUserPermissions(req interface{}, userPerms []string, authOpt Option) error {
	if len(userPerms) == 0 {
		return errors.ThrowPermissionDenied(nil, "AUTH-5mWD2", "No matching permissions found")
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

	return errors.ThrowPermissionDenied(nil, "AUTH-3jknH", "No matching permissions found")
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

func HasGlobalExplicitPermission(perms []string, permToCheck string) bool {
	for _, perm := range perms {
		p, ctxID := SplitPermission(perm)
		if p == permToCheck {
			if ctxID == "" {
				return true
			}
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

func GetExplicitPermissionCtxIDs(perms []string, searchPerm string) []string {
	ctxIDs := make([]string, 0)
	for _, perm := range perms {
		p, ctxID := SplitPermission(perm)
		if p == searchPerm {
			if ctxID != "" {
				ctxIDs = append(ctxIDs, ctxID)
			}
		}
	}
	return ctxIDs
}
