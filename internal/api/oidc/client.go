package oidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/actions/object"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ClaimPrefix                     = "urn:zitadel:iam"
	ScopeProjectRolePrefix          = "urn:zitadel:iam:org:project:role:"
	ScopeProjectsRoles              = "urn:zitadel:iam:org:projects:roles"
	ClaimProjectRoles               = "urn:zitadel:iam:org:project:roles"
	ClaimProjectRolesFormat         = "urn:zitadel:iam:org:project:%s:roles"
	ScopeUserMetaData               = "urn:zitadel:iam:user:metadata"
	ClaimUserMetaData               = ScopeUserMetaData
	ScopeResourceOwner              = "urn:zitadel:iam:user:resourceowner"
	ClaimResourceOwnerID            = ScopeResourceOwner + ":id"
	ClaimResourceOwnerName          = ScopeResourceOwner + ":name"
	ClaimResourceOwnerPrimaryDomain = ScopeResourceOwner + ":primary_domain"
	ClaimActionLogFormat            = "urn:zitadel:iam:action:%s:log"

	oidcCtx = "oidc"
)

func (o *OPStorage) GetClientByClientID(ctx context.Context, id string) (_ op.Client, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()
	client, err := o.query.ActiveOIDCClientByID(ctx, id, false)
	if err != nil {
		return nil, err
	}
	return ClientFromBusiness(client, o.defaultLoginURL, o.defaultLoginURLV2), nil
}

func (o *OPStorage) GetKeyByIDAndClientID(ctx context.Context, keyID, userID string) (_ *jose.JSONWebKey, err error) {
	return o.GetKeyByIDAndIssuer(ctx, keyID, userID)
}

func (o *OPStorage) GetKeyByIDAndIssuer(ctx context.Context, keyID, issuer string) (_ *jose.JSONWebKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()
	publicKeyData, err := o.query.GetAuthNKeyPublicKeyByIDAndIdentifier(ctx, keyID, issuer)
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

func (o *OPStorage) ValidateJWTProfileScopes(ctx context.Context, subject string, scopes []string) (_ []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()
	user, err := o.query.GetUserByID(ctx, true, subject)
	if err != nil {
		return nil, err
	}
	return o.checkOrgScopes(ctx, user, scopes)
}

func (o *OPStorage) AuthorizeClientIDSecret(ctx context.Context, id string, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()
	ctx = authz.SetCtxData(ctx, authz.CtxData{
		UserID: oidcCtx,
		OrgID:  oidcCtx,
	})
	app, err := o.query.AppByClientID(ctx, id)
	if err != nil {
		return err
	}
	if app.OIDCConfig != nil {
		return o.command.VerifyOIDCClientSecret(ctx, app.ProjectID, app.ID, secret)
	}
	return o.command.VerifyAPIClientSecret(ctx, app.ProjectID, app.ID, secret)
}

func (o *OPStorage) SetUserinfoFromToken(context.Context, *oidc.UserInfo, string, string, string) error {
	panic(o.panicErr("SetUserinfoFromToken"))
}

func (o *OPStorage) SetUserinfoFromScopes(context.Context, *oidc.UserInfo, string, string, []string) error {
	panic(o.panicErr("SetUserinfoFromScopes"))
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

// SetIntrospectionFromToken implements [op.OPStorage]
func (o *OPStorage) SetIntrospectionFromToken(ctx context.Context, introspection *oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
	panic(o.panicErr("SetIntrospectionFromToken"))
}

func (o *OPStorage) ClientCredentialsTokenRequest(ctx context.Context, clientID string, scope []string) (_ op.TokenRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()
	user, err := o.query.GetUserByLoginName(ctx, false, clientID)
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

// ClientCredentials method is kept to keep the storage interface implemented.
// However, it should never be called as the VerifyClient method on the Server is overridden.
func (o *OPStorage) ClientCredentials(context.Context, string, string) (op.Client, error) {
	return nil, zerrors.ThrowInternal(nil, "OIDC-Su8So", "Errors.Internal")
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

func (o *OPStorage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

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
	user, err := o.query.GetUserByID(ctx, true, userID)
	if err != nil {
		return nil, err
	}
	queriedActions, err := o.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeCustomiseToken, domain.TriggerTypePreAccessTokenCreation, user.ResourceOwner)
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
					return object.UserGrantsFromQuery(ctx, o.query, c, userGrants)
				}),
			),
			actions.SetFields("org",
				actions.SetFields("getMetadata", func(c *actions.FieldConfig) interface{} {
					return func(goja.FunctionCall) goja.Value {
						return object.GetOrganizationMetadata(ctx, o.query, c, user.ResourceOwner)
					}
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
						if strings.HasPrefix(key, ClaimPrefix) {
							return
						}
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
						if _, err = o.command.SetUserMetadata(ctx, metadata, userID, user.ResourceOwner); err != nil {
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
			append(actions.ActionToOptions(action), actions.WithHTTP(actionCtx), actions.WithUUID(actionCtx))...,
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
	projectID, err := o.query.ProjectIDFromClientID(ctx, applicationID)
	// applicationID might contain a username (e.g. client credentials) -> ignore the not found
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, nil, err
	}
	// ensure the projectID of the requesting is part of the roleAudience
	if projectID != "" {
		roleAudience = append(roleAudience, projectID)
	}
	projectQuery, err := query.NewUserGrantProjectIDsSearchQuery(roleAudience)
	if err != nil {
		return nil, nil, err
	}
	userIDQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, nil, err
	}
	activeQuery, err := query.NewUserGrantStateQuery(domain.UserGrantStateActive)
	if err != nil {
		return nil, nil, err
	}
	grants, err := o.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{
			projectQuery,
			userIDQuery,
			activeQuery,
		},
	}, true)
	if err != nil {
		return nil, nil, err
	}
	roles := new(projectsRoles)
	// if specific roles where requested, check if they are granted and append them in the roles list
	if len(requestedRoles) > 0 {
		for _, requestedRole := range requestedRoles {
			for _, grant := range grants.UserGrants {
				checkGrantedRoles(roles, *grant, requestedRole, grant.ProjectID == projectID)
			}
		}
		return grants, roles, nil
	}
	// no specific roles were requested, so convert any grants into roles
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
	user, err := o.query.GetUserByID(ctx, true, userID)
	if err != nil {
		return nil, err
	}
	resourceOwner, err := o.query.OrgByID(ctx, true, user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		ClaimResourceOwnerID:            resourceOwner.ID,
		ClaimResourceOwnerName:          resourceOwner.Name,
		ClaimResourceOwnerPrimaryDomain: resourceOwner.Domain,
	}, nil
}

func checkGrantedRoles(roles *projectsRoles, grant query.UserGrant, requestedRole string, isRequested bool) {
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

func newProjectRoles(projectID string, grants []query.UserGrant, requestedRoles []string) *projectsRoles {
	roles := new(projectsRoles)
	// if specific roles where requested, check if they are granted and append them in the roles list
	if len(requestedRoles) > 0 {
		for _, requestedRole := range requestedRoles {
			for _, grant := range grants {
				checkGrantedRoles(roles, grant, requestedRole, grant.ProjectID == projectID)
			}
		}
		return roles
	}
	// no specific roles were requested, so convert any grants into roles
	for _, grant := range grants {
		for _, role := range grant.Roles {
			roles.Add(grant.ProjectID, role, grant.ResourceOwner, grant.OrgPrimaryDomain, grant.ProjectID == projectID)
		}
	}
	return roles
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

func (s *Server) VerifyClient(ctx context.Context, r *op.Request[op.ClientCredentials]) (_ op.Client, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	if oidc.GrantType(r.Form.Get("grant_type")) == oidc.GrantTypeClientCredentials {
		return s.clientCredentialsAuth(ctx, r.Data.ClientID, r.Data.ClientSecret)
	}

	clientID, assertion, err := clientIDFromCredentials(ctx, r.Data)
	if err != nil {
		return nil, err
	}
	client, err := s.query.ActiveOIDCClientByID(ctx, clientID, assertion)
	if zerrors.IsNotFound(err) {
		return nil, oidc.ErrInvalidClient().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError).WithDescription("no active client not found")
	}
	if err != nil {
		return nil, err // defaults to server error
	}
	if client.Settings == nil {
		client.Settings = &query.OIDCSettings{
			AccessTokenLifetime: s.defaultAccessTokenLifetime,
			IdTokenLifetime:     s.defaultIdTokenLifetime,
		}
	}

	switch client.AuthMethodType {
	case domain.OIDCAuthMethodTypeBasic, domain.OIDCAuthMethodTypePost:
		err = s.verifyClientSecret(ctx, client, r.Data.ClientSecret)
	case domain.OIDCAuthMethodTypePrivateKeyJWT:
		err = s.verifyClientAssertion(ctx, client, r.Data.ClientAssertion)
	case domain.OIDCAuthMethodTypeNone:
	}
	if err != nil {
		return nil, err
	}

	return ClientFromBusiness(client, s.defaultLoginURL, s.defaultLoginURLV2), nil
}

func (s *Server) verifyClientAssertion(ctx context.Context, client *query.OIDCClient, assertion string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if assertion == "" {
		return oidc.ErrInvalidClient().WithDescription("empty client assertion")
	}
	verifier := op.NewJWTProfileVerifierKeySet(keySetMap(client.PublicKeys), op.IssuerFromContext(ctx), time.Hour, client.ClockSkew)
	if _, err := op.VerifyJWTAssertion(ctx, assertion, verifier); err != nil {
		return oidc.ErrInvalidClient().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError).WithDescription("invalid assertion")
	}
	return nil
}

func (s *Server) verifyClientSecret(ctx context.Context, client *query.OIDCClient, secret string) (err error) {
	_, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if secret == "" {
		return oidc.ErrInvalidClient().WithDescription("empty client secret")
	}
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := s.hasher.Verify(client.HashedSecret, secret)
	spanPasswordComparison.EndWithError(err)
	if err != nil {
		return oidc.ErrInvalidClient().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError).WithDescription("invalid secret")
	}
	if updated != "" {
		s.command.OIDCUpdateSecret(ctx, client.AppID, client.ProjectID, client.Settings.ResourceOwner, updated)
	}
	return nil
}

func (s *Server) checkOrgScopes(ctx context.Context, resourceOwner string, scopes []string) ([]string, error) {
	if slices.ContainsFunc(scopes, func(scope string) bool {
		return strings.HasPrefix(scope, domain.OrgDomainPrimaryScope)
	}) {
		org, err := s.query.OrgByID(ctx, false, resourceOwner)
		if err != nil {
			return nil, err
		}
		scopes = slices.DeleteFunc(scopes, func(scope string) bool {
			if domain, ok := strings.CutPrefix(scope, domain.OrgDomainPrimaryScope); ok {
				return domain != org.Domain
			}
			return false
		})
	}
	return slices.DeleteFunc(scopes, func(scope string) bool {
		if orgID, ok := strings.CutPrefix(scope, domain.OrgIDScope); ok {
			return orgID != resourceOwner
		}
		return false
	}), nil
}
