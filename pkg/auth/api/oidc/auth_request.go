package oidc

import (
	"context"
	"strings"
	"time"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/errors"
)

func (o *OPStorage) CreateAuthRequest(ctx context.Context, req *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	//userAgentCtx := ctx.Value(UserAgentContext)
	var userAgentID string
	//var err error
	//if userAgentCtx != nil {
	//	userAgent, err = o.processor.GetUserAgent(ctx, userAgentCtx.(string))
	//}
	//if userAgentCtx == nil || err != nil {
	//	agent := CreateAgentFromContext(ctx)
	//	userAgent, err = o.processor.CreateUserAgent(ctx, agent)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	authRequest := CreateAuthRequestToBusiness(ctx, req, userAgentID, userID)
	resp, err := o.repo.CreateAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	//ids := strings.Split(id, ":")
	//if len(ids) != 2 {
	//	return nil, errors.ThrowInvalidArgument(nil, "OIDC-seM5E6", "invalid id")
	//}
	resp, err := o.repo.AuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) DeleteAuthRequest(ctx context.Context, id string) error {
	//TODO: What to do?
	// return o.processor.TerminateAuthSession(ctx, id)
	return nil
}

func (o *OPStorage) CreateToken(ctx context.Context, authReq op.AuthRequest) (string, time.Time, error) {
	ids := strings.Split(authReq.GetID(), ":")
	if len(ids) != 2 {
		return "", time.Time{}, errors.ThrowInvalidArgument(nil, "OIDC-seM5E6", "invalid id")
	}
	resp, err := o.repo.CreateToken(ctx, &model.CreateToken{AgentID: ids[0], AuthSessionID: ids[1]})
	if err != nil {
		return "", time.Time{}, err
	}
	return resp.ID, resp.Expiration, nil
}

func (o *OPStorage) TerminateSession(ctx context.Context, userID, clientID string) error {

	userAgentID := ctx.Value(UserAgentContext).(string)
	userAgent, err := o.processor.GetUserAgent(ctx, userAgentID)
	if err != nil {
		return err
	}
	return o.repo.SignOut(ctx, "", userID)
}

func (o *OPStorage) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, errCh chan<- error, timer <-chan time.Time) {
	o.processor.GetSigningKey(ctx, keyCh, errCh, timer)
}

func (o *OPStorage) GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error) {
	return o.processor.GetKeySet(ctx)
}

func (o *OPStorage) SaveNewKeyPair(ctx context.Context) error {
	return o.processor.SaveKeyPair(ctx)
}
