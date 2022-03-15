package authz

import "context"

func NewMockContext(tenantID, orgID, userID string) context.Context {
	return context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID, TenantID: tenantID})
}

func NewMockContextWithPermissions(tenantID, orgID, userID string, permissions []string) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID, TenantID: tenantID})
	return context.WithValue(ctx, requestPermissionsKey, permissions)
}
