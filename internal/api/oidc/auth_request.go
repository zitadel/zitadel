package oidc

import (
	"context"
	"strings"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/user/model"
)

const (
	loginClientHeader = "x-zitadel-login-client"
)

func (o *OPStorage) CreateAuthRequest(ctx context.Context, req *oidc.AuthRequest, userID string) (_ op.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	headers, _ := http_utils.HeadersFromCtx(ctx)
	if loginClient := headers.Get(loginClientHeader); loginClient != "" {
		return o.createAuthRequestLoginClient(ctx, req, userID, loginClient)
	}

	return o.createAuthRequest(ctx, req, userID)
}

func (o *OPStorage) createAuthRequestLoginClient(ctx context.Context, req *oidc.AuthRequest, hintUserID, loginClient string) (op.AuthRequest, error) {
	project, err := o.query.ProjectByClientID(ctx, req.ClientID, false)
	if err != nil {
		return nil, err
	}
	scope, err := o.assertProjectRoleScopesByProject(ctx, project, req.Scopes)
	if err != nil {
		return nil, err
	}
	audience, err := o.audienceFromProjectID(ctx, project.ID)
	if err != nil {
		return nil, err
	}
	authRequest := &command.AuthRequest{
		LoginClient:   loginClient,
		ClientID:      req.ClientID,
		RedirectURI:   req.RedirectURI,
		State:         req.State,
		Nonce:         req.Nonce,
		Scope:         scope,
		Audience:      audience,
		ResponseType:  ResponseTypeToBusiness(req.ResponseType),
		CodeChallenge: CodeChallengeToBusiness(req.CodeChallenge, req.CodeChallengeMethod),
		Prompt:        PromptToBusiness(req.Prompt),
		UILocales:     UILocalesToBusiness(req.UILocales),
		MaxAge:        MaxAgeToBusiness(req.MaxAge),
		LoginHint:     req.LoginHint,
		HintUserID:    hintUserID,
	}

	err = o.command.AddAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}
	return &AuthRequestV2{authRequest}, nil
}

func (o *OPStorage) createAuthRequest(ctx context.Context, req *oidc.AuthRequest, userID string) (_ op.AuthRequest, err error) {
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sd436", "no user agent id")
	}
	req.Scopes, err = o.assertProjectRoleScopes(ctx, req.ClientID, req.Scopes)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "OIDC-Gqrfg", "Errors.Internal")
	}
	authRequest := CreateAuthRequestToBusiness(ctx, req, userAgentID, userID)
	resp, err := o.repo.CreateAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) audienceFromProjectID(ctx context.Context, projectID string) ([]string, error) {
	projectIDQuery, err := query.NewAppProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}
	appIDs, err := o.query.SearchClientIDs(ctx, &query.AppSearchQueries{Queries: []query.SearchQuery{projectIDQuery}}, false)
	if err != nil {
		return nil, err
	}

	return append(appIDs, projectID), nil
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

	accessTokenLifetime, _, _, _, err := o.getOIDCSettings(ctx)
	if err != nil {
		return "", time.Time{}, err
	}

	resp, err := o.command.AddUserToken(setContextUserSystem(ctx), userOrgID, userAgentID, applicationID, req.GetSubject(), req.GetAudience(), req.GetScopes(), accessTokenLifetime) //PLANNED: lifetime from client
	if err != nil {
		return "", time.Time{}, err
	}
	return resp.TokenID, resp.Expiration, nil
}

func (o *OPStorage) CreateAccessAndRefreshTokens(ctx context.Context, req op.TokenRequest, refreshToken string) (_, _ string, _ time.Time, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, applicationID, userOrgID, authTime, authMethodsReferences := getInfoFromRequest(req)
	scopes, err := o.assertProjectRoleScopes(ctx, applicationID, req.GetScopes())
	if err != nil {
		return "", "", time.Time{}, errors.ThrowPreconditionFailed(err, "OIDC-Df2fq", "Errors.Internal")
	}
	if request, ok := req.(op.RefreshTokenRequest); ok {
		request.SetCurrentScopes(scopes)
	}

	accessTokenLifetime, _, refreshTokenIdleExpiration, refreshTokenExpiration, err := o.getOIDCSettings(ctx)
	if err != nil {
		return "", "", time.Time{}, err
	}

	resp, token, err := o.command.AddAccessAndRefreshToken(setContextUserSystem(ctx), userOrgID, userAgentID, applicationID, req.GetSubject(),
		refreshToken, req.GetAudience(), scopes, authMethodsReferences, accessTokenLifetime,
		refreshTokenIdleExpiration, refreshTokenExpiration, authTime) //PLANNED: lifetime from client
	if err != nil {
		if errors.IsErrorInvalidArgument(err) {
			err = oidc.ErrInvalidGrant().WithParent(err)
		}
		return "", "", time.Time{}, err
	}
	return resp.TokenID, token, resp.Expiration, nil
}

func getInfoFromRequest(req op.TokenRequest) (string, string, string, time.Time, []string) {
	authReq, ok := req.(*AuthRequest)
	if ok {
		return authReq.AgentID, authReq.ApplicationID, authReq.UserOrgID, authReq.AuthTime, authReq.GetAMR()
	}
	refreshReq, ok := req.(*RefreshTokenRequest)
	if ok {
		return refreshReq.UserAgentID, refreshReq.ClientID, "", refreshReq.AuthTime, refreshReq.AuthMethodsReferences
	}
	return "", "", "", time.Time{}, nil
}

func (o *OPStorage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	tokenView, err := o.repo.RefreshTokenByToken(ctx, refreshToken)
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
		logging.Error("no user agent id")
		return errors.ThrowPreconditionFailed(nil, "OIDC-fso7F", "no user agent id")
	}
	userIDs, err := o.repo.UserSessionUserIDsByAgentID(ctx, userAgentID)
	if err != nil {
		logging.WithError(err).Error("error retrieving user sessions")
		return err
	}
	if len(userIDs) == 0 {
		return nil
	}
	data := authz.CtxData{
		UserID: userID,
	}
	err = o.command.HumansSignOut(authz.SetCtxData(ctx, data), userAgentID, userIDs)
	logging.OnError(err).Error("error signing out")
	return err
}

func (o *OPStorage) RevokeToken(ctx context.Context, token, userID, clientID string) *oidc.Error {
	refreshToken, err := o.repo.RefreshTokenByID(ctx, token, userID)
	if err == nil {
		if refreshToken.ClientID != clientID {
			return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
		}
		_, err = o.command.RevokeRefreshToken(ctx, refreshToken.UserID, refreshToken.ResourceOwner, refreshToken.ID)
		if err == nil || errors.IsNotFound(err) {
			return nil
		}
		return oidc.ErrServerError().WithParent(err)
	}
	accessToken, err := o.repo.TokenByIDs(ctx, userID, token)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return oidc.ErrServerError().WithParent(err)
	}
	if accessToken.ApplicationID != clientID {
		return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
	}
	_, err = o.command.RevokeAccessToken(ctx, userID, accessToken.ResourceOwner, accessToken.ID)
	if err == nil || errors.IsNotFound(err) {
		return nil
	}
	return oidc.ErrServerError().WithParent(err)
}

func (o *OPStorage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	refreshToken, err := o.repo.RefreshTokenByToken(ctx, token)
	if err != nil {
		return "", "", op.ErrInvalidRefreshToken
	}
	if refreshToken.ClientID != clientID {
		return "", "", oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
	}
	return refreshToken.UserID, refreshToken.ID, nil
}

func (o *OPStorage) assertProjectRoleScopes(ctx context.Context, clientID string, scopes []string) ([]string, error) {
	for _, scope := range scopes {
		if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
			return scopes, nil
		}
	}
	projectID, err := o.query.ProjectIDFromOIDCClientID(ctx, clientID, false)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-AEG4d", "Errors.Internal")
	}
	project, err := o.query.ProjectByID(ctx, false, projectID, false)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-w4wIn", "Errors.Internal")
	}
	if !project.ProjectRoleAssertion {
		return scopes, nil
	}
	projectIDQuery, err := query.NewProjectRoleProjectIDSearchQuery(project.ID)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-Cyc78", "Errors.Internal")
	}
	roles, err := o.query.SearchProjectRoles(ctx, true, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectIDQuery}}, false)
	if err != nil {
		return nil, err
	}
	for _, role := range roles.ProjectRoles {
		scopes = append(scopes, ScopeProjectRolePrefix+role.Key)
	}
	return scopes, nil
}

func (o *OPStorage) assertProjectRoleScopesByProject(ctx context.Context, project *query.Project, scopes []string) ([]string, error) {
	for _, scope := range scopes {
		if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
			return scopes, nil
		}
	}
	if !project.ProjectRoleAssertion {
		return scopes, nil
	}
	projectIDQuery, err := query.NewProjectRoleProjectIDSearchQuery(project.ID)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-Cyc78", "Errors.Internal")
	}
	roles, err := o.query.SearchProjectRoles(ctx, true, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectIDQuery}}, false)
	if err != nil {
		return nil, err
	}
	for _, role := range roles.ProjectRoles {
		scopes = append(scopes, ScopeProjectRolePrefix+role.Key)
	}
	return scopes, nil
}

func (o *OPStorage) assertClientScopesForPAT(ctx context.Context, token *model.TokenView, clientID, projectID string) error {
	token.Audience = append(token.Audience, clientID)
	projectIDQuery, err := query.NewProjectRoleProjectIDSearchQuery(projectID)
	if err != nil {
		return errors.ThrowInternal(err, "OIDC-Cyc78", "Errors.Internal")
	}
	roles, err := o.query.SearchProjectRoles(ctx, true, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectIDQuery}}, false)
	if err != nil {
		return err
	}
	for _, role := range roles.ProjectRoles {
		token.Scopes = append(token.Scopes, ScopeProjectRolePrefix+role.Key)
	}
	return nil
}

func setContextUserSystem(ctx context.Context) context.Context {
	data := authz.CtxData{
		UserID: "SYSTEM",
	}
	return authz.SetCtxData(ctx, data)
}

func (o *OPStorage) getOIDCSettings(ctx context.Context) (accessTokenLifetime, idTokenLifetime, refreshTokenIdleExpiration, refreshTokenExpiration time.Duration, _ error) {
	oidcSettings, err := o.query.OIDCSettingsByAggID(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil && !errors.IsNotFound(err) {
		return time.Duration(0), time.Duration(0), time.Duration(0), time.Duration(0), err
	}

	if oidcSettings != nil {
		return oidcSettings.AccessTokenLifetime, oidcSettings.IdTokenLifetime, oidcSettings.RefreshTokenIdleExpiration, oidcSettings.RefreshTokenExpiration, nil
	}
	return o.defaultAccessTokenLifetime, o.defaultIdTokenLifetime, o.defaultRefreshTokenIdleExpiration, o.defaultRefreshTokenExpiration, nil
}
