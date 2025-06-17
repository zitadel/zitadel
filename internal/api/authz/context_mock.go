package authz

import (
	"context"

	"golang.org/x/text/language"
)

type MockContextInstanceOpts func(i *instance)

func WithMockDefaultLanguage(lang language.Tag) MockContextInstanceOpts {
	return func(i *instance) {
		i.defaultLanguage = lang
	}
}

func NewMockContext(instanceID, orgID, userID string, opts ...MockContextInstanceOpts) context.Context {
	ctx := context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})

	i := &instance{id: instanceID}
	for _, o := range opts {
		o(i)
	}

	return context.WithValue(ctx, instanceKey, i)
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
