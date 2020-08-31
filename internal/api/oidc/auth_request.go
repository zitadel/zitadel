package oidc

import (
	"context"
	"fmt"
	"time"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/errors"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
)

func (o *OPStorage) CreateAuthRequest(ctx context.Context, req *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sd436", "no user agent id")
	}
	authRequest := CreateAuthRequestToBusiness(ctx, req, userAgentID, userID)
	resp, err := o.repo.CreateAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-D3g21", "no user agent id")
	}
	resp, err := o.repo.AuthRequestByIDCheckLoggedIn(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	resp, err := o.repo.AuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) SaveAuthCode(ctx context.Context, id, code string) error {
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return errors.ThrowPreconditionFailed(nil, "OIDC-Dgus2", "no user agent id")
	}
	return o.repo.SaveAuthCode(ctx, id, code, userAgentID)
}

func (o *OPStorage) DeleteAuthRequest(ctx context.Context, id string) error {
	return o.repo.DeleteAuthRequest(ctx, id)
}

func (o *OPStorage) CreateToken(ctx context.Context, authReq op.AuthRequest) (string, time.Time, error) {
	app, err := o.repo.ApplicationByClientID(ctx, authReq.GetClientID())
	if err != nil {
		return "", time.Time{}, err
	}
	grants, err := o.repo.UserGrantsByProjectAndUserID(app.ProjectID, authReq.GetSubject())
	scopes := append(authReq.GetScopes(), grantsToScopes(grants)...)
	req, _ := authReq.(*AuthRequest)
	resp, err := o.repo.CreateToken(ctx, req.AgentID, req.ApplicationID, req.UserID, req.Audience, scopes, o.defaultAccessTokenLifetime) //PLANNED: lifetime from client
	if err != nil {
		return "", time.Time{}, err
	}
	return resp.ID, resp.Expiration, nil
}

func grantsToScopes(grants []*grant_model.UserGrantView) []string {
	scopes := make([]string, 0)
	for _, grant := range grants {
		for _, role := range grant.RoleKeys {
			scopes = append(scopes, fmt.Sprintf("%v:%v", grant.ResourceOwner, role))
		}
	}
	return scopes
}

func (o *OPStorage) TerminateSession(ctx context.Context, userID, clientID string) error {
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return errors.ThrowPreconditionFailed(nil, "OIDC-fso7F", "no user agent id")
	}
	return o.repo.SignOut(ctx, userAgentID)
}

func (o *OPStorage) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, errCh chan<- error, timer <-chan time.Time) {
	o.repo.GetSigningKey(ctx, keyCh, errCh, timer)
}

func (o *OPStorage) GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error) {
	return o.repo.GetKeySet(ctx)
}

func (o *OPStorage) SaveNewKeyPair(ctx context.Context) error {
	return o.repo.GenerateSigningKeyPair(ctx, o.signingKeyAlgorithm)
}
