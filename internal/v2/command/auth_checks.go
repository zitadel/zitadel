package command

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	caos_errors "github.com/caos/zitadel/internal/errors"
)

func checkExplicitProjectPermission(ctx context.Context, grantID, projectID string) error {
	permissions := authz.GetRequestPermissionsFromCtx(ctx)
	if authz.HasGlobalPermission(permissions) {
		return nil
	}
	ids := authz.GetAllPermissionCtxIDs(permissions)
	if grantID != "" && listContainsID(ids, grantID) {
		return nil
	}
	if listContainsID(ids, projectID) {
		return nil
	}
	return caos_errors.ThrowPermissionDenied(nil, "EVENT-Shu7e", "Errors.UserGrant.NoPermissionForProject")
}

func listContainsID(ids []string, id string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}
