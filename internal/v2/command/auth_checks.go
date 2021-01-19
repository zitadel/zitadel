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
	containsID := false
	if grantID != "" {
		containsID = listContainsID(ids, grantID)
		if containsID {
			return nil
		}
	}
	containsID = listContainsID(ids, projectID)
	if !containsID {
		return caos_errors.ThrowPermissionDenied(nil, "EVENT-Shu7e", "Errors.UserGrant.NoPermissionForProject")
	}
	return nil
}

func listContainsID(ids []string, id string) bool {
	containsID := false
	for _, i := range ids {
		if i == id {
			containsID = true
			break
		}
	}
	return containsID
}
