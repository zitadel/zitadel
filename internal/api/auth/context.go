package auth

import (
	"context"
	"time"

	"github.com/caos/logging"
)

type CtxKeyPermissions struct{}
type CtxKeyData struct{}

type CtxData struct {
	UserID    string
	OrgID     string
	ProjectID string
	AgentID   string
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
	ResolveGrants(ctx context.Context, sub, orgID string) ([]*Grant, error)
	GetProjectIDByClientID(ctx context.Context, clientID string) (string, error)
}

func VerifyTokenAndWriteCtxData(ctx context.Context, token, orgID string, t TokenVerifier) (_ context.Context, err error) {
	userID, clientID, agentID, err := verifyAccessToken(ctx, token, t)
	if err != nil {
		return nil, err
	}

	projectID, err := t.GetProjectIDByClientID(ctx, clientID)
	logging.LogWithFields("AUTH-GfAoV", "clientID", clientID).OnError(err).Warn("could not read projectid by clientid")

	return context.WithValue(ctx, CtxKeyData{}, &CtxData{UserID: userID, OrgID: orgID, ProjectID: projectID, AgentID: agentID}), nil
}

func GetCtxData(ctx context.Context) CtxData {
	if data := ctx.Value(CtxKeyData{}); data != nil {
		ctxData, ok := data.(*CtxData)
		if ok {
			return *ctxData
		}
		time.Now()
	}
	return CtxData{}
}

func GetPermissionsFromCtx(ctx context.Context) []string {
	if data := ctx.Value(CtxKeyPermissions{}); data != nil {
		ctxPermission, _ := data.([]string)
		return ctxPermission
	}
	return nil
}
