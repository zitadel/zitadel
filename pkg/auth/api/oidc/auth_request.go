package oidc

import (
	"context"
	"strings"
	"time"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
)

func (o *OPStorage) CreateAuthRequest(ctx context.Context, req *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	//userAgentCtx := ctx.Value(UserAgentContext)
	userAgentID, ok := UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sd436", "no user agent id")
	}
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
	req, err := o.repo.AuthRequestByID(ctx, authReq.GetID())
	if err != nil {
		return "", time.Time{}, err
	}
	resp, err := o.repo.CreateToken(ctx, req.AgentID, req.ApplicationID, req.UserID, req.Request.(*model.AuthRequestOIDC).Scopes, lifetime)
	if err != nil {
		return "", time.Time{}, err
	}
	return resp.ID, resp.Expiration, nil
}

func (o *OPStorage) TerminateSession(ctx context.Context, userID, clientID string) error {
	userAgentID, ok := UserAgentIDFromCtx(ctx)
	if !ok {
		return errors.ThrowPreconditionFailed(nil, "OIDC-fso7F", "no user agent id")
	}
	return o.repo.SignOut(ctx, userAgentID, userID)
}

func (o *OPStorage) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, errCh chan<- error, timer <-chan time.Time) {
	o.repo.GetSigningKey(ctx, keyCh, errCh, timer)
}

func (o *OPStorage) GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error) {
	return o.repo.GetKeySet(ctx)
}

func (o *OPStorage) SaveNewKeyPair(ctx context.Context) error {
	return o.repo.SaveKeyPair(ctx)
}
