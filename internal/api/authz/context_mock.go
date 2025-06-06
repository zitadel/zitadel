package authz

import (
	"context"

	"golang.org/x/text/language"
)

func NewMockContext(instanceID, orgID, userID string, language language.Tag) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
	return context.WithValue(ctx, instanceKey, &instance{id: instanceID, defaultLanguage: language})
}

func NewMockContextWithAgent(instanceID, orgID, userID, agentID string) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID, AgentID: agentID})
	return context.WithValue(ctx, instanceKey, &instance{id: instanceID})
}

func NewMockContextWithPermissions(instanceID, orgID, userID string, permissions []string) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
	ctx = context.WithValue(ctx, instanceKey, &instance{id: instanceID})
	return context.WithValue(ctx, requestPermissionsKey, permissions)
}
