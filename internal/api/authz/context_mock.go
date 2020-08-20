package authz

import "context"

func NewMockContext(orgID, userID string) context.Context {
	return context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
}
