package authz

import (
	"context"

	"github.com/caos/logging"
)

type key int

const (
	permissionsKey key = 1
	dataKey        key = 2
)

type CtxData struct {
	UserID            string
	OrgID             string
	ProjectID         string
	AgentID           string
	PreferredLanguage string
}

func (ctxData CtxData) IsZero() bool {
	return ctxData.UserID == "" || ctxData.OrgID == ""
}

type Grants []*Grant

type Grant struct {
	OrgID string
	Roles []string
}

type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (string, string, string, error)
	ResolveGrant(ctx context.Context) (*Grant, error)
	GetProjectIDByClientID(ctx context.Context, clientID string) (string, error)
}

func VerifyTokenAndWriteCtxData(ctx context.Context, token, orgID string, t TokenVerifier) (_ context.Context, err error) {
	userID, clientID, agentID, err := verifyAccessToken(ctx, token, t)
	if err != nil {
		return nil, err
	}
	projectID, err := t.GetProjectIDByClientID(ctx, clientID)
	logging.LogWithFields("AUTH-GfAoV", "clientID", clientID).OnError(err).Warn("could not read projectid by clientid")
	return context.WithValue(ctx, dataKey, CtxData{UserID: userID, OrgID: orgID, ProjectID: projectID, AgentID: agentID}), nil
}

func SetCtxData(ctx context.Context, ctxData CtxData) context.Context {
	return context.WithValue(ctx, dataKey, ctxData)
}

func GetCtxData(ctx context.Context) CtxData {
	ctxData, _ := ctx.Value(dataKey).(CtxData)
	return ctxData
}

func GetPermissionsFromCtx(ctx context.Context) []string {
	ctxPermission, _ := ctx.Value(permissionsKey).([]string)
	return ctxPermission
}
