package oidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
	"gopkg.in/square/go-jose.v2"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/actions/object"
	"github.com/zitadel/zitadel/internal/api/authz"
	api_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	ScopeProjectRolePrefix  = "urn:zitadel:iam:org:project:role:"
	ScopeProjectsRoles      = "urn:zitadel:iam:org:projects:roles"
	ClaimProjectRoles       = "urn:zitadel:iam:org:project:roles"
	ClaimProjectRolesFormat = "urn:zitadel:iam:org:project:%s:roles"
	ScopeUserMetaData       = "urn:zitadel:iam:user:metadata"
	ClaimUserMetaData       = ScopeUserMetaData
	ScopeResourceOwner      = "urn:zitadel:iam:user:resourceowner"
	ClaimResourceOwner      = ScopeResourceOwner + ":"
	ClaimActionLogFormat    = "urn:zitadel:iam:action:%s:log"

	oidcCtx = "oidc"
)

func (o *OPStorage) GetClientByClientID(ctx context.Context, id string) (_ op.Client, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	client, err := o.query.AppByOIDCClientID(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if client.State != domain.AppStateActive {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sdaGg", "client is not active")
	}
	projectIDQuery, err := query.NewProjectRoleProjectIDSearchQuery(client.ProjectID)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-mPxqP", "Errors.Internal")
	}
	projectRoles, err := o.query.SearchProjectRoles(ctx, true, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectIDQuery}}, false)
	if err != nil {
		return nil, err
	}
	allowedScopes := make([]string, len(projectRoles.ProjectRoles))
	for i, role := range projectRoles.ProjectRoles {
		allowedScopes[i] = ScopeProjectRolePrefix + role.Key
	}

	accessTokenLifetime, idTokenLifetime, _, _, err := o.getOIDCSettings(ctx)
	if err != nil {
		return nil, err
	}

	return ClientFromBusiness(client, o.defaultLoginURL, o.defaultLoginURLV2, accessTokenLifetime, idTokenLifetime, allowedScopes)
}

func (o *OPStorage) GetKeyByIDAndClientID(ctx context.Context, keyID, userID string) (_ *jose.JSONWebKey, err error) {
	return o.GetKeyByIDAndIssuer(ctx, keyID, userID)
}

func (o *OPStorage) GetKeyByIDAndIssuer(ctx context.Context, keyID, issuer string) (_ *jose.JSONWebKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	publicKeyData, err := o.query.GetAuthNKeyPublicKeyByIDAndIdentifier(ctx, keyID, issuer, false)
	if err != nil {
		return nil, err
	}
	publicKey, err := crypto.BytesToPublicKey(publicKeyData)
	if err != nil {
		return nil, err
	}
	return &jose.JSONWebKey{
		KeyID: keyID,
		Use:   "sig",
		Key:   publicKey,
	}, nil
}

func (o *OPStorage) ValidateJWTProfileScopes(ctx context.Context, subject string, scopes []string) ([]string, error) {
	user, err := o.query.GetUserByID(ctx, true, subject, false)
	if err != nil {
		return nil, err
	}
	return o.checkOrgScopes(ctx, user, scopes)
}

func (o *OPStorage) AuthorizeClientIDSecret(ctx context.Context, id string, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctx = authz.SetCtxData(ctx, authz.CtxData{
		UserID: oidcCtx,
		OrgID:  oidcCtx,
	})
	app, err := o.query.AppByClientID(ctx, id, false)
	if err != nil {
		return err
	}
	if app.OIDCConfig != nil {
		return o.command.VerifyOIDCClientSecret(ctx, app.ProjectID, app.ID, secret)
	}
	return o.command.VerifyAPIClientSecret(ctx, app.ProjectID, app.ID, secret)
}

func (o *OPStorage) SetUserinfoFromToken(ctx context.Context, userInfo *oidc.UserInfo, tokenID, subject, origin string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if strings.HasPrefix(tokenID, command.IDPrefixV2) {
		token, err := o.query.ActiveAccessTokenByToken(ctx, tokenID)
		if err != nil {
			return err
		}
		if err = o.isOriginAllowed(ctx, token.ClientID, origin); err != nil {
			return err
		}
		return o.setUserinfo(ctx, userInfo, token.UserID, token.ClientID, token.Scope, nil)
	}

	token, err := o.repo.TokenByIDs(ctx, subject, tokenID)
	if err != nil {
		return errors.ThrowPermissionDenied(nil, "OIDC-Dsfb2", "token is not valid or has expired")
	}
	if token.ApplicationID != "" {
		if err = o.isOriginAllowed(ctx, token.ApplicationID, origin); err != nil {
			return err
		}
	}
	return o.setUserinfo(ctx, userInfo, token.UserID, token.ApplicationID, token.Scopes, nil)
}

func (o *OPStorage) SetUserinfoFromScopes(ctx context.Context, userInfo *oidc.UserInfo, userID, applicationID string, scopes []string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if applicationID != "" {
		app, err := o.query.AppByOIDCClientID(ctx, applicationID, false)
		if err != nil {
			return err
		}
		if app.OIDCConfig.AssertIDTokenRole {
			scopes, err = o.assertProjectRoleScopes(ctx, applicationID, scopes)
			if err != nil {
				return errors.ThrowPreconditionFailed(err, "OIDC-Dfe2s", "Errors.Internal")
			}
		}
	}
	return o.setUserinfo(ctx, userInfo, userID, applicationID, scopes, nil)
}

// SetUserinfoFromRequest extends the SetUserinfoFromScopes during the id_token generation.
// This is required for V2 tokens to be able to set the sessionID (`sid`) claim.
func (o *OPStorage) SetUserinfoFromRequest(ctx context.Context, userinfo *oidc.UserInfo, request op.IDTokenRequest, _ []string) error {
	switch t := request.(type) {
	case *AuthRequestV2:
		userinfo.AppendClaims("sid", t.SessionID)
	case *RefreshTokenRequestV2:
		userinfo.AppendClaims("sid", t.SessionID)
	}
	return nil
}

func (o *OPStorage) SetIntrospectionFromToken(ctx context.Context, introspection *oidc.IntrospectionResponse, tokenID, subject, clientID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if strings.HasPrefix(tokenID, command.IDPrefixV2) {
		token, err := o.query.ActiveAccessTokenByToken(ctx, tokenID)
		if err != nil {
			return err
		}
		projectID, err := o.query.ProjectIDFromClientID(ctx, clientID, false)
		if err != nil {
			return errors.ThrowPermissionDenied(nil, "OIDC-Adfg5", "client not found")
		}
		return o.introspect(ctx, introspection,
			tokenID, token.UserID, token.ClientID, projectID,
			token.Audience, token.Scope,
			token.AccessTokenCreation, token.AccessTokenExpiration)
	}

	token, err := o.repo.TokenByIDs(ctx, subject, tokenID)
	if err != nil {
		return errors.ThrowPermissionDenied(nil, "OIDC-Dsfb2", "token is not valid or has expired")
	}
	projectID, err := o.query.ProjectIDFromClientID(ctx, clientID, false)
	if err != nil {
		return errors.ThrowPermissionDenied(nil, "OIDC-Adfg5", "client not found")
	}
	if token.IsPAT {
		err = o.assertClientScopesForPAT(ctx, token, clientID, projectID)
		if err != nil {
			return errors.ThrowPreconditionFailed(err, "OIDC-AGefw", "Errors.Internal")
		}
	}
	return o.introspect(ctx, introspection,
		token.ID, token.UserID, token.ApplicationID, projectID,
		token.Audience, token.Scopes,
		token.CreationDate, token.Expiration)
}

func (o *OPStorage) ClientCredentialsTokenRequest(ctx context.Context, clientID string, scope []string) (op.TokenRequest, error) {
	loginname, err := query.NewUserLoginNamesSearchQuery(clientID)
	if err != nil {
		return nil, err
	}
	user, err := o.query.GetUser(ctx, false, false, loginname)
	if err != nil {
		return nil, err
	}
	scope, err = o.checkOrgScopes(ctx, user, scope)
	if err != nil {
		return nil, err
	}
	audience := domain.AddAudScopeToAudience(ctx, nil, scope)
	return &clientCredentialsRequest{
		sub:      user.ID,
		scopes:   scope,
		audience: audience,
	}, nil
}

func (o *OPStorage) ClientCredentials(ctx context.Context, clientID, clientSecret string) (op.Client, error) {
	loginname, err := query.NewUserLoginNamesSearchQuery(clientID)
	if err != nil {
		return nil, err
	}
	user, err := o.query.GetUser(ctx, false, false, loginname)
	if err != nil {
		return nil, err
	}
	if _, err := o.command.VerifyMachineSecret(ctx, user.ID, user.ResourceOwner, clientSecret); err != nil {
		return nil, err
	}
	return &clientCredentialsClient{
		id:        clientID,
		tokenType: accessTokenTypeToOIDC(user.Machine.AccessTokenType),
	}, nil
}

// isOriginAllowed checks whether a call by the client to the endpoint is allowed from the provided origin
// if no origin is provided, no error will be returned
func (o *OPStorage) isOriginAllowed(ctx context.Context, clientID, origin string) error {
	if origin == "" {
		return nil
	}
	app, err := o.query.AppByOIDCClientID(ctx, clientID, false)
	if err != nil {
		return err
	}
	if api_http.IsOriginAllowed(app.OIDCConfig.AllowedOrigins, origin) {
		return nil
	}
	return errors.ThrowPermissionDenied(nil, "OIDC-da1f3", "origin is not allowed")
}

func (o *OPStorage) introspect(
	ctx context.Context,
	introspection *oidc.IntrospectionResponse,
	tokenID, subject, clientID, projectID string,
	audience, scope []string,
	tokenCreation, tokenExpiration time.Time,
) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	for _, aud := range audience {
		if aud == clientID || aud == projectID {
			userInfo := new(oidc.UserInfo)
			err = o.setUserinfo(ctx, userInfo, subject, clientID, scope, []string{projectID}) // always
			if err != nil {
				return err
			}
			introspection.SetUserInfo(userInfo)
			introspection.Scope = scope
			introspection.ClientID = clientID
			introspection.TokenType = oidc.BearerToken
			introspection.Expiration = oidc.FromTime(tokenExpiration)
			introspection.IssuedAt = oidc.FromTime(tokenCreation)
			introspection.NotBefore = oidc.FromTime(tokenCreation)
			introspection.Audience = audience
			introspection.Issuer = op.IssuerFromContext(ctx)
			introspection.JWTID = tokenID
			return nil
		}
	}
	return errors.ThrowPermissionDenied(nil, "OIDC-sdg3G", "token is not valid for this client")
}

func (o *OPStorage) checkOrgScopes(ctx context.Context, user *query.User, scopes []string) ([]string, error) {
	for i := len(scopes) - 1; i >= 0; i-- {
		scope := scopes[i]
		if strings.HasPrefix(scope, domain.OrgDomainPrimaryScope) {
			var orgID string
			org, err := o.query.OrgByPrimaryDomain(ctx, strings.TrimPrefix(scope, domain.OrgDomainPrimaryScope))
			if err == nil {
				orgID = org.ID
			}
			if orgID != user.ResourceOwner {
				scopes[i] = scopes[len(scopes)-1]
				scopes[len(scopes)-1] = ""
				scopes = scopes[:len(scopes)-1]
			}
		}
		if strings.HasPrefix(scope, domain.OrgIDScope) {
			if strings.TrimPrefix(scope, domain.OrgIDScope) != user.ResourceOwner {
				scopes[i] = scopes[len(scopes)-1]
				scopes[len(scopes)-1] = ""
				scopes = scopes[:len(scopes)-1]
			}
		}
	}
	return scopes, nil
}

func (o *OPStorage) setUserinfo(ctx context.Context, userInfo *oidc.UserInfo, userID, applicationID string, scopes []string, roleAudience []string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	user, err := o.query.GetUserByID(ctx, true, userID, false)
	if err != nil {
		return err
	}
	var allRoles bool
	roles := make([]string, 0)
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.Subject = user.ID
		case oidc.ScopeEmail:
			setUserInfoEmail(userInfo, user)
		case oidc.ScopeProfile:
			o.setUserInfoProfile(ctx, userInfo, user)
		case oidc.ScopePhone:
			setUserInfoPhone(userInfo, user)
		case oidc.ScopeAddress:
			//TODO: handle address for human users as soon as implemented
		case ScopeUserMetaData:
			if err := o.setUserInfoMetadata(ctx, userInfo, userID); err != nil {
				return err
			}
		case ScopeResourceOwner:
			if err := o.setUserInfoResourceOwner(ctx, userInfo, userID); err != nil {
				return err
			}
		case ScopeProjectsRoles:
			allRoles = true
		default:
			if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
				roles = append(roles, strings.TrimPrefix(scope, ScopeProjectRolePrefix))
			}
			if strings.HasPrefix(scope, domain.OrgDomainPrimaryScope) {
				userInfo.AppendClaims(domain.OrgDomainPrimaryClaim, strings.TrimPrefix(scope, domain.OrgDomainPrimaryScope))
			}
			if strings.HasPrefix(scope, domain.OrgIDScope) {
				userInfo.AppendClaims(domain.OrgIDClaim, strings.TrimPrefix(scope, domain.OrgIDScope))
				if err := o.setUserInfoResourceOwner(ctx, userInfo, userID); err != nil {
					return err
				}
			}
		}
	}

	// if all roles are requested take the audience for those from the scopes
	if allRoles && len(roleAudience) == 0 {
		roleAudience = domain.AddAudScopeToAudience(ctx, roleAudience, scopes)
	}

	userGrants, projectRoles, err := o.assertRoles(ctx, userID, applicationID, roles, roleAudience)
	if err != nil {
		return err
	}
	o.setUserInfoRoleClaims(userInfo, projectRoles)

	return o.userinfoFlows(ctx, user.ResourceOwner, userGrants, userInfo)
}

func (o *OPStorage) setUserInfoProfile(ctx context.Context, userInfo *oidc.UserInfo, user *query.User) {
	userInfo.PreferredUsername = user.PreferredLoginName
	userInfo.UpdatedAt = oidc.FromTime(user.ChangeDate)
	if user.Machine != nil {
		userInfo.Name = user.Machine.Name
		return
	}
	userInfo.Name = user.Human.DisplayName
	userInfo.FamilyName = user.Human.LastName
	userInfo.GivenName = user.Human.FirstName
	userInfo.Nickname = user.Human.NickName
	userInfo.Gender = getGender(user.Human.Gender)
	userInfo.Locale = oidc.NewLocale(user.Human.PreferredLanguage)
	userInfo.Picture = domain.AvatarURL(o.assetAPIPrefix(ctx), user.ResourceOwner, user.Human.AvatarKey)
}

func setUserInfoEmail(userInfo *oidc.UserInfo, user *query.User) {
	if user.Human == nil {
		return
	}
	userInfo.UserInfoEmail = oidc.UserInfoEmail{
		Email:         string(user.Human.Email),
		EmailVerified: oidc.Bool(user.Human.IsEmailVerified)}
}

func setUserInfoPhone(userInfo *oidc.UserInfo, user *query.User) {
	if user.Human == nil {
		return
	}
	userInfo.UserInfoPhone = oidc.UserInfoPhone{
		PhoneNumber:         string(user.Human.Phone),
		PhoneNumberVerified: user.Human.IsPhoneVerified,
	}
}

func (o *OPStorage) setUserInfoMetadata(ctx context.Context, userInfo *oidc.UserInfo, userID string) error {
	userMetaData, err := o.assertUserMetaData(ctx, userID)
	if err != nil {
		return err
	}
	if len(userMetaData) > 0 {
		userInfo.AppendClaims(ClaimUserMetaData, userMetaData)
	}
	return nil
}

func (o *OPStorage) setUserInfoResourceOwner(ctx context.Context, userInfo *oidc.UserInfo, userID string) error {
	resourceOwnerClaims, err := o.assertUserResourceOwner(ctx, userID)
	if err != nil {
		return err
	}
	for claim, value := range resourceOwnerClaims {
		userInfo.AppendClaims(claim, value)
	}
	return nil
}

func (o *OPStorage) setUserInfoRoleClaims(userInfo *oidc.UserInfo, roles *projectsRoles) {
	if roles != nil && len(roles.projects) > 0 {
		if roles, ok := roles.projects[roles.requestProjectID]; ok {
			userInfo.AppendClaims(ClaimProjectRoles, roles)
		}
		for projectID, roles := range roles.projects {
			userInfo.AppendClaims(fmt.Sprintf(ClaimProjectRolesFormat, projectID), roles)
		}
	}
}

func (o *OPStorage) userinfoFlows(ctx context.Context, resourceOwner string, userGrants *query.UserGrants, userInfo *oidc.UserInfo) error {
	queriedActions, err := o.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeCustomiseToken, domain.TriggerTypePreUserinfoCreation, resourceOwner, false)
	if err != nil {
		return err
	}

	ctxFields := actions.SetContextFields(
		actions.SetFields("v1",
			actions.SetFields("claims", userinfoClaims(userInfo)),
			actions.SetFields("getUser", func(c *actions.FieldConfig) interface{} {
				return func(call goja.FunctionCall) goja.Value {
					user, err := o.query.GetUserByID(ctx, true, userInfo.Subject, false)
					if err != nil {
						panic(err)
					}
					return object.UserFromQuery(c, user)
				}
			}),
			actions.SetFields("user",
				actions.SetFields("getMetadata", func(c *actions.FieldConfig) interface{} {
					return func(goja.FunctionCall) goja.Value {
						resourceOwnerQuery, err := query.NewUserMetadataResourceOwnerSearchQuery(resourceOwner)
						if err != nil {
							logging.WithError(err).Debug("unable to create search query")
							panic(err)
						}
						metadata, err := o.query.SearchUserMetadata(
							ctx,
							true,
							userInfo.Subject,
							&query.UserMetadataSearchQueries{Queries: []query.SearchQuery{resourceOwnerQuery}},
							false,
						)
						if err != nil {
							logging.WithError(err).Info("unable to get md in action")
							panic(err)
						}
						return object.UserMetadataListFromQuery(c, metadata)
					}
				}),
				actions.SetFields("grants", func(c *actions.FieldConfig) interface{} {
					return object.UserGrantsFromQuery(c, userGrants)
				}),
			),
		),
	)

	for _, action := range queriedActions {
		actionCtx, cancel := context.WithTimeout(ctx, action.Timeout())
		claimLogs := []string{}

		apiFields := actions.WithAPIFields(
			actions.SetFields("v1",
				actions.SetFields("userinfo",
					actions.SetFields("setClaim", func(key string, value interface{}) {
						if userInfo.Claims[key] == nil {
							userInfo.AppendClaims(key, value)
							return
						}
						claimLogs = append(claimLogs, fmt.Sprintf("key %q already exists", key))
					}),
					actions.SetFields("appendLogIntoClaims", func(entry string) {
						claimLogs = append(claimLogs, entry)
					}),
				),
				actions.SetFields("claims",
					actions.SetFields("setClaim", func(key string, value interface{}) {
						if userInfo.Claims[key] == nil {
							userInfo.AppendClaims(key, value)
							return
						}
						claimLogs = append(claimLogs, fmt.Sprintf("key %q already exists", key))
					}),
					actions.SetFields("appendLogIntoClaims", func(entry string) {
						claimLogs = append(claimLogs, entry)
					}),
				),
				actions.SetFields("user",
					actions.SetFields("setMetadata", func(call goja.FunctionCall) goja.Value {
						if len(call.Arguments) != 2 {
							panic("exactly 2 (key, value) arguments expected")
						}
						key := call.Arguments[0].Export().(string)
						val := call.Arguments[1].Export()

						value, err := json.Marshal(val)
						if err != nil {
							logging.WithError(err).Debug("unable to marshal")
							panic(err)
						}

						metadata := &domain.Metadata{
							Key:   key,
							Value: value,
						}
						if _, err = o.command.SetUserMetadata(ctx, metadata, userInfo.Subject, resourceOwner); err != nil {
							logging.WithError(err).Info("unable to set md in action")
							panic(err)
						}
						return nil
					}),
				),
			),
		)

		err = actions.Run(
			actionCtx,
			ctxFields,
			apiFields,
			action.Script,
			action.Name,
			append(actions.ActionToOptions(action), actions.WithHTTP(actionCtx))...,
		)
		cancel()
		if err != nil {
			return err
		}
		if len(claimLogs) > 0 {
			userInfo.AppendClaims(fmt.Sprintf(ClaimActionLogFormat, action.Name), claimLogs)
		}
	}

	return nil
}

func (o *OPStorage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
	roles := make([]string, 0)
	var allRoles bool
	for _, scope := range scopes {
		switch scope {
		case ScopeUserMetaData:
			userMetaData, err := o.assertUserMetaData(ctx, userID)
			if err != nil {
				return nil, err
			}
			if len(userMetaData) > 0 {
				claims = appendClaim(claims, ClaimUserMetaData, userMetaData)
			}
		case ScopeResourceOwner:
			resourceOwnerClaims, err := o.assertUserResourceOwner(ctx, userID)
			if err != nil {
				return nil, err
			}
			for claim, value := range resourceOwnerClaims {
				claims = appendClaim(claims, claim, value)
			}
		case ScopeProjectsRoles:
			allRoles = true
		}
		if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
			roles = append(roles, strings.TrimPrefix(scope, ScopeProjectRolePrefix))
		}
		if strings.HasPrefix(scope, domain.OrgDomainPrimaryScope) {
			claims = appendClaim(claims, domain.OrgDomainPrimaryClaim, strings.TrimPrefix(scope, domain.OrgDomainPrimaryScope))
		}
		if strings.HasPrefix(scope, domain.OrgIDScope) {
			claims = appendClaim(claims, domain.OrgIDClaim, strings.TrimPrefix(scope, domain.OrgIDScope))
			resourceOwnerClaims, err := o.assertUserResourceOwner(ctx, userID)
			if err != nil {
				return nil, err
			}
			for claim, value := range resourceOwnerClaims {
				claims = appendClaim(claims, claim, value)
			}
		}
	}

	// If requested, use the audience as context for the roles,
	// otherwise the project itself will be used
	var roleAudience []string
	if allRoles {
		roleAudience = domain.AddAudScopeToAudience(ctx, roleAudience, scopes)
	}

	userGrants, projectRoles, err := o.assertRoles(ctx, userID, clientID, roles, roleAudience)
	if err != nil {
		return nil, err
	}

	if projectRoles != nil && len(projectRoles.projects) > 0 {
		if roles, ok := projectRoles.projects[projectRoles.requestProjectID]; ok {
			claims = appendClaim(claims, ClaimProjectRoles, roles)
		}
		for projectID, roles := range projectRoles.projects {
			claims = appendClaim(claims, fmt.Sprintf(ClaimProjectRolesFormat, projectID), roles)
		}
	}

	return o.privateClaimsFlows(ctx, userID, userGrants, claims)
}

func (o *OPStorage) privateClaimsFlows(ctx context.Context, userID string, userGrants *query.UserGrants, claims map[string]interface{}) (map[string]interface{}, error) {
	user, err := o.query.GetUserByID(ctx, true, userID, false)
	if err != nil {
		return nil, err
	}
	queriedActions, err := o.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeCustomiseToken, domain.TriggerTypePreAccessTokenCreation, user.ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	ctxFields := actions.SetContextFields(
		actions.SetFields("v1",
			actions.SetFields("claims", func(c *actions.FieldConfig) interface{} {
				return c.Runtime.ToValue(claims)
			}),
			actions.SetFields("getUser", func(c *actions.FieldConfig) interface{} {
				return func(call goja.FunctionCall) goja.Value {
					user, err := o.query.GetUserByID(ctx, true, userID, false)
					if err != nil {
						panic(err)
					}
					return object.UserFromQuery(c, user)
				}
			}),
			actions.SetFields("user",
				actions.SetFields("getMetadata", func(c *actions.FieldConfig) interface{} {
					return func(goja.FunctionCall) goja.Value {
						resourceOwnerQuery, err := query.NewUserMetadataResourceOwnerSearchQuery(user.ResourceOwner)
						if err != nil {
							logging.WithError(err).Debug("unable to create search query")
							panic(err)
						}
						metadata, err := o.query.SearchUserMetadata(
							ctx,
							true,
							userID,
							&query.UserMetadataSearchQueries{Queries: []query.SearchQuery{resourceOwnerQuery}},
							false,
						)
						if err != nil {
							logging.WithError(err).Info("unable to get md in action")
							panic(err)
						}
						return object.UserMetadataListFromQuery(c, metadata)
					}
				}),
				actions.SetFields("grants", func(c *actions.FieldConfig) interface{} {
					return object.UserGrantsFromQuery(c, userGrants)
				}),
			),
		),
	)

	for _, action := range queriedActions {
		claimLogs := []string{}
		actionCtx, cancel := context.WithTimeout(ctx, action.Timeout())

		apiFields := actions.WithAPIFields(
			actions.SetFields("v1",
				actions.SetFields("claims",
					actions.SetFields("setClaim", func(key string, value interface{}) {
						if _, ok := claims[key]; !ok {
							claims = appendClaim(claims, key, value)
							return
						}
						claimLogs = append(claimLogs, fmt.Sprintf("key %q already exists", key))
					}),
					actions.SetFields("appendLogIntoClaims", func(entry string) {
						claimLogs = append(claimLogs, entry)
					}),
				),
				actions.SetFields("user",
					actions.SetFields("setMetadata", func(call goja.FunctionCall) {
						if len(call.Arguments) != 2 {
							panic("exactly 2 (key, value) arguments expected")
						}
						key := call.Arguments[0].Export().(string)
						val := call.Arguments[1].Export()

						value, err := json.Marshal(val)
						if err != nil {
							logging.WithError(err).Debug("unable to marshal")
							panic(err)
						}

						metadata := &domain.Metadata{
							Key:   key,
							Value: value,
						}
						if _, err = o.command.SetUserMetadata(ctx, metadata, userID, user.ResourceOwner); err != nil {
							logging.WithError(err).Info("unable to set md in action")
							panic(err)
						}
					}),
				),
			),
		)

		err = actions.Run(
			actionCtx,
			ctxFields,
			apiFields,
			action.Script,
			action.Name,
			append(actions.ActionToOptions(action), actions.WithHTTP(actionCtx))...,
		)
		cancel()
		if err != nil {
			return nil, err
		}
		if len(claimLogs) > 0 {
			claims = appendClaim(claims, fmt.Sprintf(ClaimActionLogFormat, action.Name), claimLogs)
			claimLogs = nil
		}
	}

	return claims, nil
}

func (o *OPStorage) assertRoles(ctx context.Context, userID, applicationID string, requestedRoles, roleAudience []string) (*query.UserGrants, *projectsRoles, error) {
	if (applicationID == "" || len(requestedRoles) == 0) && len(roleAudience) == 0 {
		return nil, nil, nil
	}
	projectID, err := o.query.ProjectIDFromClientID(ctx, applicationID, false)
	// applicationID might contain a username (e.g. client credentials) -> ignore the not found
	if err != nil && !errors.IsNotFound(err) {
		return nil, nil, err
	}
	// ensure the projectID of the requesting is part of the roleAudience
	if projectID != "" {
		roleAudience = append(roleAudience, projectID)
	}
	queries := make([]query.SearchQuery, 0, 2)
	projectQuery, err := query.NewUserGrantProjectIDsSearchQuery(roleAudience)
	if err != nil {
		return nil, nil, err
	}
	queries = append(queries, projectQuery)
	userIDQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, nil, err
	}
	queries = append(queries, userIDQuery)
	grants, err := o.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: queries,
	}, true, false)
	if err != nil {
		return nil, nil, err
	}
	roles := new(projectsRoles)
	// if specific roles where requested, check if they are granted and append them in the roles list
	if len(requestedRoles) > 0 {
		for _, requestedRole := range requestedRoles {
			for _, grant := range grants.UserGrants {
				checkGrantedRoles(roles, grant, requestedRole, grant.ProjectID == projectID)
			}
		}
		return grants, roles, nil
	}
	// now specific roles were requested, so convert any grants into roles
	for _, grant := range grants.UserGrants {
		for _, role := range grant.Roles {
			roles.Add(grant.ProjectID, role, grant.ResourceOwner, grant.OrgPrimaryDomain, grant.ProjectID == projectID)
		}
	}
	return grants, roles, nil
}

func (o *OPStorage) assertUserMetaData(ctx context.Context, userID string) (map[string]string, error) {
	metaData, err := o.query.SearchUserMetadata(ctx, true, userID, &query.UserMetadataSearchQueries{}, false)
	if err != nil {
		return nil, err
	}

	userMetaData := make(map[string]string)
	for _, md := range metaData.Metadata {
		userMetaData[md.Key] = base64.RawURLEncoding.EncodeToString(md.Value)
	}
	return userMetaData, nil
}

func (o *OPStorage) assertUserResourceOwner(ctx context.Context, userID string) (map[string]string, error) {
	user, err := o.query.GetUserByID(ctx, true, userID, false)
	if err != nil {
		return nil, err
	}
	resourceOwner, err := o.query.OrgByID(ctx, true, user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		ClaimResourceOwner + "id":             resourceOwner.ID,
		ClaimResourceOwner + "name":           resourceOwner.Name,
		ClaimResourceOwner + "primary_domain": resourceOwner.Domain,
	}, nil
}

func checkGrantedRoles(roles *projectsRoles, grant *query.UserGrant, requestedRole string, isRequested bool) {
	for _, grantedRole := range grant.Roles {
		if requestedRole == grantedRole {
			roles.Add(grant.ProjectID, grantedRole, grant.ResourceOwner, grant.OrgPrimaryDomain, isRequested)
		}
	}
}

// projectsRoles contains all projects with all their roles for a user
type projectsRoles struct {
	// key is projectID
	projects map[string]projectRoles

	requestProjectID string
}

func (p *projectsRoles) Add(projectID, roleKey, orgID, domain string, isRequested bool) {
	if p.projects == nil {
		p.projects = make(map[string]projectRoles, 1)
	}
	if p.projects[projectID] == nil {
		p.projects[projectID] = make(projectRoles)
	}
	if isRequested {
		p.requestProjectID = projectID
	}
	p.projects[projectID].Add(roleKey, orgID, domain)
}

// projectRoles contains the roles of a project of multiple organisations
//
// key of the first map is the role key,
// key of the second map is the org id, value the org domain
type projectRoles map[string]map[string]string

func (p projectRoles) Add(roleKey, orgID, domain string) {
	if p[roleKey] == nil {
		p[roleKey] = make(map[string]string, 1)
	}
	p[roleKey][orgID] = domain
}

func getGender(gender domain.Gender) oidc.Gender {
	switch gender {
	case domain.GenderFemale:
		return "female"
	case domain.GenderMale:
		return "male"
	case domain.GenderDiverse:
		return "diverse"
	}
	return ""
}

func appendClaim(claims map[string]interface{}, claim string, value interface{}) map[string]interface{} {
	if claims == nil {
		claims = make(map[string]interface{})
	}
	claims[claim] = value
	return claims
}

func userinfoClaims(userInfo *oidc.UserInfo) func(c *actions.FieldConfig) interface{} {
	return func(c *actions.FieldConfig) interface{} {
		marshalled, err := json.Marshal(userInfo)
		if err != nil {
			panic(err)
		}

		claims := make(map[string]interface{}, 10)
		if err = json.Unmarshal(marshalled, &claims); err != nil {
			panic(err)
		}
		return c.Runtime.ToValue(claims)
	}
}
