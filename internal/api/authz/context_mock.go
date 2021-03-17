package authz

import "context"

func NewMockContext(orgID, userID string) context.Context {
	return context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
}

func NewMockContextWithPermissions(orgID, userID string, permissions []string) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
	return context.WithValue(ctx, requestPermissionsKey, permissions)
}
