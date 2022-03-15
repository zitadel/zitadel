package authz

import "context"

func NewMockContext(instanceID, orgID, userID string) context.Context {
	return context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID, InstanceID: instanceID})
}

func NewMockContextWithPermissions(instanceID, orgID, userID string, permissions []string) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID, InstanceID: instanceID})
	return context.WithValue(ctx, requestPermissionsKey, permissions)
}
