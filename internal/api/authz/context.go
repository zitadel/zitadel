package authz

import (
	"context"
	"github.com/caos/zitadel/internal/errors"

	"github.com/caos/logging"
)

type key int

const (
	requestPermissionsKey key = 1
	dataKey               key = 2
	allPermissionsKey     key = 3
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

func VerifyTokenAndWriteCtxData(ctx context.Context, token, orgID string, t *TokenVerifier, method string) (_ context.Context, err error) {
	if orgID != "" {
		err := t.ExistsOrg(ctx, orgID)
		if err != nil {
			return nil, errors.ThrowPermissionDenied(nil, "AUTH-Bs7Ds", "Organisation doesn't exist")
		}
	}

	userID, clientID, agentID, err := verifyAccessToken(ctx, token, t, method)
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

func GetRequestPermissionsFromCtx(ctx context.Context) []string {
	ctxPermission, _ := ctx.Value(requestPermissionsKey).([]string)
	return ctxPermission
}

func GetAllPermissionsFromCtx(ctx context.Context) []string {
	ctxPermission, _ := ctx.Value(allPermissionsKey).([]string)
	return ctxPermission
}
