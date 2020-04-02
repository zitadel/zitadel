package auth

import "context"

func NewMockContext(orgID, userID string) context.Context {
	return context.WithValue(nil, dataKey, CtxData{UserID: userID, OrgID: orgID})
}
