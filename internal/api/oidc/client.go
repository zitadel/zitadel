package oidc

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/api/authz"
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

// GetClientByClientID implements the op.Storage interface to retrieve an OIDC client by its ID.
//
// TODO: Still used for Auth request creation for v1 login.
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

func (o *OPStorage) GetKeyByIDAndClientID(context.Context, string, string) (*jose.JSONWebKey, error) {
	panic(o.panicErr("GetKeyByIDAndClientID"))
}

func (o *OPStorage) ValidateJWTProfileScopes(context.Context, string, []string) ([]string, error) {
	panic(o.panicErr("ValidateJWTProfileScopes"))
}

func (o *OPStorage) AuthorizeClientIDSecret(context.Context, string, string) error {
	panic(o.panicErr("AuthorizeClientIDSecret"))
}

func (o *OPStorage) SetUserinfoFromToken(context.Context, *oidc.UserInfo, string, string, string) error {
	panic(o.panicErr("SetUserinfoFromToken"))
}

func (o *OPStorage) SetUserinfoFromScopes(context.Context, *oidc.UserInfo, string, string, []string) error {
	panic(o.panicErr("SetUserinfoFromScopes"))
}

func (o *OPStorage) SetUserinfoFromRequest(context.Context, *oidc.UserInfo, op.IDTokenRequest, []string) error {
	panic(o.panicErr("SetUserinfoFromRequest"))
}

func (o *OPStorage) SetIntrospectionFromToken(context.Context, *oidc.IntrospectionResponse, string, string, string) error {
	panic(o.panicErr("SetIntrospectionFromToken"))
}

func (o *OPStorage) ClientCredentialsTokenRequest(context.Context, string, []string) (op.TokenRequest, error) {
	panic(o.panicErr("ClientCredentialsTokenRequest"))
}

func (o *OPStorage) ClientCredentials(context.Context, string, string) (op.Client, error) {
	panic(o.panicErr("ClientCredentials"))
}

func (o *OPStorage) GetPrivateClaimsFromScopes(context.Context, string, string, []string) (map[string]interface{}, error) {
	panic(o.panicErr("GetPrivateClaimsFromScopes"))
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
		org, err := s.query.OrgByID(ctx, resourceOwner)
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
