package oidc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/errors"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
)

func (o *OPStorage) CreateAuthRequest(ctx context.Context, req *oidc.AuthRequest, userID string) (_ op.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sd436", "no user agent id")
	}
	app, err := o.repo.ApplicationByClientID(ctx, req.ClientID)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-AEG4d", "Errors.Internal")
	}
	req.Scopes, err = o.assertProjectRoleScopes(app, req.Scopes)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-Gqrfg", "Errors.Internal")
	}
	authRequest := CreateAuthRequestToBusiness(ctx, req, userAgentID, userID)
	resp, err := o.repo.CreateAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) AuthRequestByID(ctx context.Context, id string) (_ op.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
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

func (o *OPStorage) AuthRequestByCode(ctx context.Context, code string) (_ op.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	resp, err := o.repo.AuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) SaveAuthCode(ctx context.Context, id, code string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return errors.ThrowPreconditionFailed(nil, "OIDC-Dgus2", "no user agent id")
	}
	return o.repo.SaveAuthCode(ctx, id, code, userAgentID)
}

func (o *OPStorage) DeleteAuthRequest(ctx context.Context, id string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return o.repo.DeleteAuthRequest(ctx, id)
}

func (o *OPStorage) CreateAccessToken(ctx context.Context, req op.TokenRequest) (_ string, _ time.Time, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	var userAgentID, applicationID, userOrgID string
	authReq, ok := req.(*AuthRequest)
	if ok {
		userAgentID = authReq.AgentID
		applicationID = authReq.ApplicationID
		userOrgID = authReq.UserOrgID
	}
	resp, err := o.command.AddUserToken(ctx, userOrgID, userAgentID, applicationID, req.GetSubject(), req.GetAudience(), req.GetScopes(), o.defaultAccessTokenLifetime) //PLANNED: lifetime from client
	if err != nil {
		return "", time.Time{}, err
	}
	return resp.TokenID, resp.Expiration, nil
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

func (o *OPStorage) CreateAccessAndRefreshTokens(ctx context.Context, req op.TokenRequest, refreshToken string) (_, _ string, _ time.Time, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	var userAgentID, applicationID, userOrgID string
	var authTime time.Time
	var authMethodsReferences []string
	authReq, ok := req.(*AuthRequest)
	if ok {
		userAgentID = authReq.AgentID
		applicationID = authReq.ApplicationID
		userOrgID = authReq.UserOrgID
		authTime = authReq.AuthTime
		authMethodsReferences = authReq.GetAMR()
	}
	resp, token, err := o.command.AddUserAndRefreshToken(ctx, userOrgID, userAgentID, applicationID, req.GetSubject(),
		refreshToken, req.GetAudience(), req.GetScopes(), authMethodsReferences, o.defaultAccessTokenLifetime,
		o.defaultRefreshTokenIdleExpiration, o.defaultRefreshTokenExpiration, authTime) //PLANNED: lifetime from client
	if err != nil {
		return "", "", time.Time{}, err
	}
	return resp.TokenID, token, resp.Expiration, nil
}

func (o *OPStorage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	tokenView, err := o.repo.RefreshTokenByID(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return RefreshTokenRequestFromBusiness(tokenView), nil
}

func (o *OPStorage) TerminateSession(ctx context.Context, userID, clientID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return errors.ThrowPreconditionFailed(nil, "OIDC-fso7F", "no user agent id")
	}
	userIDs, err := o.repo.UserSessionUserIDsByAgentID(ctx, userAgentID)
	if err != nil {
		return err
	}
	return o.command.HumansSignOut(ctx, userAgentID, userIDs)
}

func (o *OPStorage) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey) {
	o.repo.GetSigningKey(ctx, keyCh, o.signingKeyAlgorithm)
}

func (o *OPStorage) GetKeySet(ctx context.Context) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return o.repo.GetKeySet(ctx)
}

func (o *OPStorage) assertProjectRoleScopes(app *proj_model.ApplicationView, scopes []string) ([]string, error) {
	if !app.ProjectRoleAssertion {
		return scopes, nil
	}
	for _, scope := range scopes {
		if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
			return scopes, nil
		}
	}
	roles, err := o.repo.ProjectRolesByProjectID(app.ProjectID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		scopes = append(scopes, ScopeProjectRolePrefix+role.Key)
	}
	return scopes, nil
}
