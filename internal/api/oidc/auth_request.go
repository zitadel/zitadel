package oidc

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/handler"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/domain/federatedlogout"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	LoginClientHeader            = "x-zitadel-login-client"
	LoginPostLogoutRedirectParam = "post_logout_redirect"
	LoginLogoutHintParam         = "logout_hint"
	LoginUILocalesParam          = "ui_locales"
	LoginPath                    = "/login"
	LogoutPath                   = "/logout"
	LogoutDonePath               = "/logout/done"
)

func (o *OPStorage) CreateAuthRequest(ctx context.Context, req *oidc.AuthRequest, userID string) (_ op.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	// for backwards compatibility we pass the login client if set
	headers, _ := http_utils.HeadersFromCtx(ctx)
	loginClient := headers.Get(LoginClientHeader)

	// for backwards compatibility we'll use the new login if the header is set (no matter the other configs)
	if loginClient != "" {
		return o.createAuthRequestLoginClient(ctx, req, userID, loginClient)
	}

	// if the instance requires the v2 login, use it no matter what the application configured
	if authz.GetFeatures(ctx).LoginV2.Required {
		return o.createAuthRequestLoginClient(ctx, req, userID, loginClient)
	}

	version, err := o.query.OIDCClientLoginVersion(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	switch version {
	case domain.LoginVersion1:
		return o.createAuthRequest(ctx, req, userID)
	case domain.LoginVersion2:
		return o.createAuthRequestLoginClient(ctx, req, userID, loginClient)
	case domain.LoginVersionUnspecified:
		fallthrough
	default:
		// since we already checked for a login header, we can fall back to the v1 login
		return o.createAuthRequest(ctx, req, userID)
	}
}

func (o *OPStorage) createAuthRequestScopeAndAudience(ctx context.Context, clientID string, reqScope []string) (scope, audience []string, err error) {
	project, err := o.query.ProjectByClientID(ctx, clientID)
	if err != nil {
		return nil, nil, err
	}
	scope, err = o.assertProjectRoleScopesByProject(ctx, project, reqScope)
	if err != nil {
		return nil, nil, err
	}
	audience, err = o.audienceFromProjectID(ctx, project.ID)
	audience = domain.AddAudScopeToAudience(ctx, audience, scope)
	if err != nil {
		return nil, nil, err
	}
	return scope, audience, nil
}

func (o *OPStorage) createAuthRequestLoginClient(ctx context.Context, req *oidc.AuthRequest, hintUserID, loginClient string) (op.AuthRequest, error) {
	scope, audience, err := o.createAuthRequestScopeAndAudience(ctx, req.ClientID, req.Scopes)
	if err != nil {
		return nil, err
	}
	authRequest := &command.AuthRequest{
		LoginClient:      loginClient,
		ClientID:         req.ClientID,
		RedirectURI:      req.RedirectURI,
		State:            req.State,
		Nonce:            req.Nonce,
		Scope:            scope,
		Audience:         audience,
		NeedRefreshToken: slices.Contains(scope, oidc.ScopeOfflineAccess),
		ResponseType:     ResponseTypeToBusiness(req.ResponseType),
		ResponseMode:     ResponseModeToBusiness(req.ResponseMode),
		CodeChallenge:    CodeChallengeToBusiness(req.CodeChallenge, req.CodeChallengeMethod),
		Prompt:           PromptToBusiness(req.Prompt),
		UILocales:        UILocalesToBusiness(req.UILocales),
		MaxAge:           MaxAgeToBusiness(req.MaxAge),
		Issuer:           o.contextToIssuer(ctx),
	}
	if req.LoginHint != "" {
		authRequest.LoginHint = &req.LoginHint
	}
	if hintUserID != "" {
		authRequest.HintUserID = &hintUserID
	}

	aar, err := o.command.AddAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}
	return &AuthRequestV2{aar}, nil
}

func (o *OPStorage) createAuthRequest(ctx context.Context, req *oidc.AuthRequest, userID string) (_ op.AuthRequest, err error) {
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, zerrors.ThrowPreconditionFailed(nil, "OIDC-sd436", "no user agent id")
	}
	scope, audience, err := o.createAuthRequestScopeAndAudience(ctx, req.ClientID, req.Scopes)
	if err != nil {
		return nil, err
	}
	req.Scopes = scope
	authRequest := CreateAuthRequestToBusiness(ctx, req, userAgentID, userID, audience)
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
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	if strings.HasPrefix(id, command.IDPrefixV2) {
		req, err := o.command.GetCurrentAuthRequest(ctx, id)
		if err != nil {
			return nil, err
		}
		return &AuthRequestV2{req}, nil
	}

	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, zerrors.ThrowPreconditionFailed(nil, "OIDC-D3g21", "no user agent id")
	}
	resp, err := o.repo.AuthRequestByIDCheckLoggedIn(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (o *OPStorage) AuthRequestByCode(ctx context.Context, code string) (_ op.AuthRequest, err error) {
	panic(o.panicErr("AuthRequestByCode"))
}

// decryptGrant decrypts a code or refresh_token
func (o *OPStorage) decryptGrant(grant string) (string, error) {
	decodedGrant, err := base64.RawURLEncoding.DecodeString(grant)
	if err != nil {
		return "", err
	}
	return o.encAlg.DecryptString(decodedGrant, o.encAlg.EncryptionKeyID())
}

func (o *OPStorage) SaveAuthCode(ctx context.Context, id, code string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	if strings.HasPrefix(id, command.IDPrefixV2) {
		return o.command.AddAuthRequestCode(ctx, id, code)
	}

	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return zerrors.ThrowPreconditionFailed(nil, "OIDC-Dgus2", "no user agent id")
	}
	return o.repo.SaveAuthCode(ctx, id, code, userAgentID)
}

func (o *OPStorage) DeleteAuthRequest(context.Context, string) error {
	panic(o.panicErr("DeleteAuthRequest"))
}

func (o *OPStorage) CreateAccessToken(context.Context, op.TokenRequest) (string, time.Time, error) {
	panic(o.panicErr("CreateAccessToken"))
}

func (o *OPStorage) CreateAccessAndRefreshTokens(context.Context, op.TokenRequest, string) (string, string, time.Time, error) {
	panic(o.panicErr("CreateAccessAndRefreshTokens"))
}

func (*OPStorage) panicErr(method string) error {
	return fmt.Errorf("OPStorage.%s should not be called anymore. This is a bug. Please report https://github.com/zitadel/zitadel/issues", method)
}

func (o *OPStorage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (_ op.RefreshTokenRequest, err error) {
	panic("TokenRequestByRefreshToken should not be called anymore. This is a bug. Please report https://github.com/zitadel/zitadel/issues")
}

func (o *OPStorage) TerminateSession(ctx context.Context, userID, clientID string) (err error) {
	_, err = o.terminateSession(ctx, userID)
	return err
}

func (o *OPStorage) terminateSession(ctx context.Context, userID string) (sessions []command.HumanSignOutSession, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		logging.Error("no user agent id")
		return nil, zerrors.ThrowPreconditionFailed(nil, "OIDC-fso7F", "no user agent id")
	}
	sessions, err = o.repo.UserSessionsByAgentID(ctx, userAgentID)
	if err != nil {
		logging.WithError(err).Error("error retrieving user sessions")
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}
	data := authz.CtxData{
		UserID: userID,
	}
	err = o.command.HumansSignOut(authz.SetCtxData(ctx, data), userAgentID, sessions)
	logging.OnError(err).Error("error signing out")
	return sessions, err
}

func (o *OPStorage) TerminateSessionFromRequest(ctx context.Context, endSessionRequest *op.EndSessionRequest) (redirectURI string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	// check for the login client header
	headers, _ := http_utils.HeadersFromCtx(ctx)

	// V2:
	// In case there is no id_token_hint and login V2 is either required by feature
	// or requested via header (backwards compatibility),
	// we'll redirect to the UI (V2) and let it decide which session to terminate
	//
	// If there's no id_token_hint and for v1 logins, we handle them separately
	if endSessionRequest.IDTokenHintClaims == nil && (authz.GetFeatures(ctx).LoginV2.Required || headers.Get(LoginClientHeader) != "") {
		redirectURI := v2PostLogoutRedirectURI(endSessionRequest.RedirectURI)
		logoutURI := authz.GetFeatures(ctx).LoginV2.BaseURI
		// if no logout uri is set, fallback to the default configured in the runtime config
		if logoutURI == nil || logoutURI.String() == "" {
			logoutURI, err = url.Parse(o.defaultLogoutURLV2)
			if err != nil {
				return "", err
			}
		} else {
			logoutURI = logoutURI.JoinPath(LogoutPath)
		}
		return buildLoginV2LogoutURL(logoutURI, redirectURI, endSessionRequest.LogoutHint, endSessionRequest.UILocales), nil
	}

	// V1:
	// We check again for the id_token_hint param and if a session is set in it.
	// All explicit V2 sessions with empty id_token_hint are handled above and all V2 session contain a sessionID
	// So if any condition is not met, we handle the request as a V1 request and do a (v1) TerminateSession,
	// which terminates all sessions of the user agent, identified by cookie.
	if endSessionRequest.IDTokenHintClaims == nil || endSessionRequest.IDTokenHintClaims.SessionID == "" {
		sessions, err := o.terminateSession(ctx, endSessionRequest.UserID)
		if err != nil {
			return "", err
		}
		if len(sessions) == 1 {
			if path := o.federatedLogout(ctx, sessions[0].ID, endSessionRequest.RedirectURI); path != "" {
				return path, nil
			}
		}
		return endSessionRequest.RedirectURI, nil
	}

	// V1:
	// If the sessionID is prefixed by V1, we also terminate a v1 session, but based on the SingleV1SessionTermination feature flag,
	// we either terminate all sessions of the user agent or only the specific session
	if strings.HasPrefix(endSessionRequest.IDTokenHintClaims.SessionID, handler.IDPrefixV1) {
		err = o.terminateV1Session(ctx, endSessionRequest.UserID, endSessionRequest.IDTokenHintClaims.SessionID)
		if err != nil {
			return "", err
		}
		if path := o.federatedLogout(ctx, endSessionRequest.IDTokenHintClaims.SessionID, endSessionRequest.RedirectURI); path != "" {
			return path, nil
		}
		return endSessionRequest.RedirectURI, nil
	}

	// V2:
	// Terminate the v2 session of the id_token_hint
	_, err = o.command.TerminateSessionWithoutTokenCheck(ctx, endSessionRequest.IDTokenHintClaims.SessionID)
	if err != nil {
		return "", err
	}
	return v2PostLogoutRedirectURI(endSessionRequest.RedirectURI), nil
}

// federatedLogout checks whether the session has an idp session linked and the IDP template is configured for federated logout.
// If so, it creates a federated logout request and stores it in the cache and returns the logout path.
func (o *OPStorage) federatedLogout(ctx context.Context, sessionID string, postLogoutRedirectURI string) string {
	session, err := o.repo.UserSessionByID(ctx, sessionID)
	if err != nil {
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "sessionID", sessionID).
			WithError(err).Error("error retrieving user session")
		return ""
	}
	if session.SelectedIDPConfigID.String == "" {
		return ""
	}
	identityProvider, err := o.query.IDPTemplateByID(ctx, false, session.SelectedIDPConfigID.String, false, nil)
	if err != nil {
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "idpID", session.SelectedIDPConfigID.String, "sessionID", sessionID).
			WithError(err).Error("error retrieving idp template")
		return ""
	}
	if identityProvider.SAMLIDPTemplate == nil || !identityProvider.FederatedLogoutEnabled {
		return ""
	}
	o.federateLogoutCache.Set(ctx, &federatedlogout.FederatedLogout{
		InstanceID:            authz.GetInstance(ctx).InstanceID(),
		FingerPrintID:         authz.GetCtxData(ctx).AgentID,
		SessionID:             sessionID,
		IDPID:                 session.SelectedIDPConfigID.String,
		UserID:                session.UserID,
		PostLogoutRedirectURI: postLogoutRedirectURI,
		State:                 federatedlogout.StateCreated,
	})
	return login.ExternalLogoutPath(sessionID)
}

func buildLoginV2LogoutURL(logoutURI *url.URL, redirectURI, logoutHint string, uiLocales []language.Tag) string {
	if strings.HasSuffix(logoutURI.Path, "/") && len(logoutURI.Path) > 1 {
		logoutURI.Path = strings.TrimSuffix(logoutURI.Path, "/")
	}

	q := logoutURI.Query()
	q.Set(LoginPostLogoutRedirectParam, redirectURI)
	if logoutHint != "" {
		q.Set(LoginLogoutHintParam, logoutHint)
	}
	if len(uiLocales) > 0 {
		locales := make([]string, len(uiLocales))
		for i, locale := range uiLocales {
			locales[i] = locale.String()
		}
		q.Set(LoginUILocalesParam, strings.Join(locales, " "))
	}
	logoutURI.RawQuery = q.Encode()
	return logoutURI.String()
}

// v2PostLogoutRedirectURI will take care that the post_logout_redirect_uri is correctly set for v2 logins.
// The default value set by the [op.SessionEnder] only handles V1 logins. In case the redirect_uri is set to the default
// we'll return the path for the v2 login.
func v2PostLogoutRedirectURI(redirectURI string) string {
	if redirectURI != login.DefaultLoggedOutPath {
		return redirectURI
	}
	return LogoutDonePath
}

// terminateV1Session terminates "v1" sessions created through the login UI.
// Depending on the flag, we either terminate a single session or all of the user agent
func (o *OPStorage) terminateV1Session(ctx context.Context, userID, sessionID string) error {
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: userID})
	// if the flag is active we only terminate the specific session
	if authz.GetFeatures(ctx).OIDCSingleV1SessionTermination {
		userAgentID, err := o.repo.UserAgentIDBySessionID(ctx, sessionID)
		if err != nil {
			return err
		}
		return o.command.HumansSignOut(ctx, userAgentID, []command.HumanSignOutSession{{ID: sessionID, UserID: userID}})
	}
	// otherwise we search for all active sessions within the same user agent of the current session id
	userAgentID, sessions, err := o.repo.ActiveUserSessionsBySessionID(ctx, sessionID)
	if err != nil {
		logging.WithError(err).Error("error retrieving user sessions")
		return err
	}
	if len(sessions) == 0 {
		return nil
	}
	return o.command.HumansSignOut(ctx, userAgentID, sessions)
}

func (o *OPStorage) RevokeToken(ctx context.Context, token, userID, clientID string) (err *oidc.Error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		// check for nil, because `err` is not an error and EndWithError would panic
		if err == nil {
			span.End()
			return
		}
		span.EndWithError(err)
	}()

	if strings.HasPrefix(token, command.IDPrefixV2) {
		err := o.command.RevokeOIDCSessionToken(ctx, token, clientID)
		if err == nil {
			return nil
		}
		if zerrors.IsPreconditionFailed(err) {
			return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
		}
		return oidc.ErrServerError().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError)
	}

	return o.revokeTokenV1(ctx, token, userID, clientID)
}

func (o *OPStorage) revokeTokenV1(ctx context.Context, token, userID, clientID string) *oidc.Error {
	refreshToken, err := o.repo.RefreshTokenByID(ctx, token, userID)
	if err == nil {
		if refreshToken.ClientID != clientID {
			return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
		}
		_, err = o.command.RevokeRefreshToken(ctx, refreshToken.UserID, refreshToken.ResourceOwner, refreshToken.ID)
		if err == nil || zerrors.IsNotFound(err) {
			return nil
		}
		return oidc.ErrServerError().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError)
	}
	accessToken, err := o.repo.TokenByIDs(ctx, userID, token)
	if err != nil {
		if zerrors.IsNotFound(err) {
			return nil
		}
		return oidc.ErrServerError().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError)
	}
	if accessToken.ApplicationID != clientID {
		return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
	}
	_, err = o.command.RevokeAccessToken(ctx, userID, accessToken.ResourceOwner, accessToken.ID)
	if err == nil || zerrors.IsNotFound(err) {
		return nil
	}
	return oidc.ErrServerError().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError)
}

func (o *OPStorage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	plainToken, err := o.decryptGrant(token)
	if err != nil {
		return "", "", op.ErrInvalidRefreshToken
	}
	if strings.HasPrefix(plainToken, command.IDPrefixV2) {
		oidcSession, err := o.command.OIDCSessionByRefreshToken(ctx, plainToken)
		if err != nil {
			return "", "", op.ErrInvalidRefreshToken
		}
		return oidcSession.UserID, oidcSession.OIDCRefreshTokenID(oidcSession.RefreshTokenID), nil
	}
	refreshToken, err := o.repo.RefreshTokenByToken(ctx, token)
	if err != nil {
		return "", "", op.ErrInvalidRefreshToken
	}
	if refreshToken.ClientID != clientID {
		return "", "", oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
	}
	return refreshToken.UserID, refreshToken.ID, nil
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
		return nil, zerrors.ThrowInternal(err, "OIDC-Cyc78", "Errors.Internal")
	}
	roles, err := o.query.SearchProjectRoles(ctx, true, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectIDQuery}}, nil)
	if err != nil {
		return nil, err
	}
	for _, role := range roles.ProjectRoles {
		scopes = append(scopes, ScopeProjectRolePrefix+role.Key)
	}
	return scopes, nil
}

func setContextUserSystem(ctx context.Context) context.Context {
	data := authz.CtxData{
		UserID: "SYSTEM",
	}
	return authz.SetCtxData(ctx, data)
}

func CreateErrorCallbackURL(authReq op.AuthRequest, reason, description, uri string, authorizer op.Authorizer) (string, error) {
	e := struct {
		Error       string `schema:"error"`
		Description string `schema:"error_description,omitempty"`
		URI         string `schema:"error_uri,omitempty"`
		State       string `schema:"state,omitempty"`
	}{
		Error:       reason,
		Description: description,
		URI:         uri,
		State:       authReq.GetState(),
	}
	callback, err := op.AuthResponseURL(authReq.GetRedirectURI(), authReq.GetResponseType(), authReq.GetResponseMode(), e, authorizer.Encoder())
	if err != nil {
		return "", err
	}
	return callback, nil
}

func CreateCodeCallbackURL(ctx context.Context, authReq op.AuthRequest, authorizer op.Authorizer) (string, error) {
	code, err := op.CreateAuthRequestCode(ctx, authReq, authorizer.Storage(), authorizer.Crypto())
	if err != nil {
		return "", err
	}
	codeResponse := struct {
		code  string
		state string
	}{
		code:  code,
		state: authReq.GetState(),
	}
	return op.AuthResponseURL(authReq.GetRedirectURI(), authReq.GetResponseType(), authReq.GetResponseMode(), &codeResponse, authorizer.Encoder())
}

func (s *Server) CreateTokenCallbackURL(ctx context.Context, req op.AuthRequest) (string, error) {
	provider := s.Provider()
	opClient, err := provider.Storage().GetClientByClientID(ctx, req.GetClientID())
	if err != nil {
		return "", err
	}
	client, ok := opClient.(*Client)
	if !ok {
		return "", zerrors.ThrowInternal(nil, "OIDC-waeN6", "Error.Internal")
	}

	session, state, err := s.command.CreateOIDCSessionFromAuthRequest(
		setContextUserSystem(ctx),
		req.GetID(),
		implicitFlowComplianceChecker(),
		slices.Contains(client.GrantTypes(), oidc.GrantTypeRefreshToken),
		client.client.BackChannelLogoutURI,
	)
	if err != nil {
		return "", err
	}
	resp, err := s.accessTokenResponseFromSession(ctx, client, session, state, client.client.ProjectID, client.client.ProjectRoleAssertion, client.client.AccessTokenRoleAssertion, client.client.IDTokenRoleAssertion, client.client.IDTokenUserinfoAssertion)
	if err != nil {
		return "", err
	}
	callback, err := op.AuthResponseURL(req.GetRedirectURI(), req.GetResponseType(), req.GetResponseMode(), resp, provider.Encoder())
	if err != nil {
		return "", err
	}
	return callback, err
}

func implicitFlowComplianceChecker() command.AuthRequestComplianceChecker {
	return func(_ context.Context, authReq *command.AuthRequestWriteModel) error {
		if err := authReq.CheckAuthenticated(); err != nil {
			return err
		}
		return nil
	}
}

func (s *Server) authorizeCallbackHandler(w http.ResponseWriter, r *http.Request) {
	authorizer := s.Provider()
	authReq, err := func(ctx context.Context) (authReq *AuthRequest, err error) {
		ctx, span := tracing.NewSpan(ctx)
		r = r.WithContext(ctx)
		defer func() { span.EndWithError(err) }()

		id, err := op.ParseAuthorizeCallbackRequest(r)
		if err != nil {
			return nil, err
		}
		authReq, err = s.getAuthRequestV1ByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if !authReq.Done() {
			return authReq, oidc.ErrInteractionRequired().WithDescription("Unfortunately, the user may be not logged in and/or additional interaction is required.")
		}
		return authReq, s.authResponse(authReq, authorizer, w, r)
	}(r.Context())
	if err != nil {
		// we need to make sure there's no empty interface passed
		if authReq == nil {
			op.AuthRequestError(w, r, nil, err, authorizer)
			return
		}
		op.AuthRequestError(w, r, authReq, err, authorizer)
	}
}

func (s *Server) authResponse(authReq *AuthRequest, authorizer op.Authorizer, w http.ResponseWriter, r *http.Request) (err error) {
	ctx, span := tracing.NewSpan(r.Context())
	r = r.WithContext(ctx)
	defer func() { span.EndWithError(err) }()

	client, err := authorizer.Storage().GetClientByClientID(ctx, authReq.GetClientID())
	if err != nil {
		op.AuthRequestError(w, r, authReq, err, authorizer)
		return err
	}
	if authReq.GetResponseType() == oidc.ResponseTypeCode {
		op.AuthResponseCode(w, r, authReq, authorizer)
		return nil
	}
	return s.authResponseToken(authReq, authorizer, client, w, r)
}

func (s *Server) authResponseToken(authReq *AuthRequest, authorizer op.Authorizer, opClient op.Client, w http.ResponseWriter, r *http.Request) (err error) {
	ctx, span := tracing.NewSpan(r.Context())
	r = r.WithContext(ctx)
	defer func() { span.EndWithError(err) }()

	client, ok := opClient.(*Client)
	if !ok {
		return zerrors.ThrowInternal(nil, "OIDC-waeN6", "Error.Internal")
	}

	scope := authReq.GetScopes()
	session, err := s.command.CreateOIDCSession(ctx,
		authReq.UserID,
		authReq.UserOrgID,
		client.client.ClientID,
		client.client.BackChannelLogoutURI,
		scope,
		authReq.Audience,
		authReq.AuthMethods(),
		authReq.AuthTime,
		authReq.GetNonce(),
		authReq.PreferredLanguage,
		authReq.ToUserAgent(),
		domain.TokenReasonAuthRequest,
		nil,
		slices.Contains(scope, oidc.ScopeOfflineAccess),
		authReq.SessionID,
		authReq.oidc().ResponseType,
	)
	if err != nil {
		op.AuthRequestError(w, r, authReq, err, authorizer)
		return err
	}
	resp, err := s.accessTokenResponseFromSession(ctx, client, session, authReq.GetState(), client.client.ProjectID, client.client.ProjectRoleAssertion, client.client.AccessTokenRoleAssertion, client.client.IDTokenRoleAssertion, client.client.IDTokenUserinfoAssertion)
	if err != nil {
		op.AuthRequestError(w, r, authReq, err, authorizer)
		return err
	}

	if authReq.GetResponseMode() == oidc.ResponseModeFormPost {
		if err = op.AuthResponseFormPost(w, authReq.GetRedirectURI(), resp, authorizer.Encoder()); err != nil {
			op.AuthRequestError(w, r, authReq, err, authorizer)
			return err
		}
		return nil
	}

	callback, err := op.AuthResponseURL(authReq.GetRedirectURI(), authReq.GetResponseType(), authReq.GetResponseMode(), resp, authorizer.Encoder())
	if err != nil {
		op.AuthRequestError(w, r, authReq, err, authorizer)
		return err
	}
	http.Redirect(w, r, callback, http.StatusFound)
	return nil
}
