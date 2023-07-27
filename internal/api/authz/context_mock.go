package authz

import "context"

func NewMockContext(instanceID, orgID, userID string) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
	return context.WithValue(ctx, instanceKey, &instance{id: instanceID})
}

func NewMockContextWithPermissions(instanceID, orgID, userID string, permissions []string) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
	ctx = context.WithValue(ctx, instanceKey, &instance{id: instanceID})
	return context.WithValue(ctx, requestPermissionsKey, permissions)
}
