package oidc

import (
	"context"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) clientCredentialsAuth(ctx context.Context, clientID, clientSecret string) (op.Client, error) {
	user, err := s.query.GetUserByLoginName(ctx, false, clientID)
	if zerrors.IsNotFound(err) {
		return nil, oidc.ErrInvalidClient().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError).WithDescription("client not found")
	}
	if err != nil {
		return nil, err // defaults to server error
	}
	if user.Machine == nil || user.Machine.EncodedSecret == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "OIDC-pieP8", "Errors.User.Machine.Secret.NotExisting")
	}
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := s.hasher.Verify(user.Machine.EncodedSecret, clientSecret)
	spanPasswordComparison.EndWithError(err)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "OIDC-VoXo6", "Errors.User.Machine.Secret.Invalid")
	}

	s.command.MachineSecretCheckSucceeded(ctx, user.ID, user.ResourceOwner, updated)
	return &clientCredentialsClient{
		clientID:      user.Username,
		userID:        user.ID,
		resourceOwner: user.ResourceOwner,
		tokenType:     user.Machine.AccessTokenType,
	}, nil
}

type clientCredentialsClient struct {
	clientID      string
	userID        string
	resourceOwner string
	tokenType     domain.OIDCTokenType
}

// AccessTokenType returns the AccessTokenType for the token to be created because of the client credentials request
// machine users currently only have opaque tokens ([op.AccessTokenTypeBearer])
func (c *clientCredentialsClient) AccessTokenType() op.AccessTokenType {
	return accessTokenTypeToOIDC(c.tokenType)
}

// GetID returns the client_id (username of the machine user) for the token to be created because of the client credentials request
func (c *clientCredentialsClient) GetID() string {
	return c.clientID
}

// RedirectURIs returns nil as there are no redirect uris
func (c *clientCredentialsClient) RedirectURIs() []string {
	return nil
}

// PostLogoutRedirectURIs returns nil as there are no logout redirect uris
func (c *clientCredentialsClient) PostLogoutRedirectURIs() []string {
	return nil
}

// ApplicationType returns [op.ApplicationTypeWeb] as the machine users is a confidential client
func (c *clientCredentialsClient) ApplicationType() op.ApplicationType {
	return op.ApplicationTypeWeb
}

// AuthMethod returns the allowed auth method type for machine user.
// It returns Basic Auth
func (c *clientCredentialsClient) AuthMethod() oidc.AuthMethod {
	return oidc.AuthMethodBasic
}

// ResponseTypes returns nil as the types are only required for an authorization request
func (c *clientCredentialsClient) ResponseTypes() []oidc.ResponseType {
	return nil
}

// GrantTypes returns the grant types supported by the machine users, which is currently only client credentials ([oidc.GrantTypeClientCredentials])
func (c *clientCredentialsClient) GrantTypes() []oidc.GrantType {
	return []oidc.GrantType{
		oidc.GrantTypeClientCredentials,
	}
}

// LoginURL returns an empty string as there is no login UI involved
func (c *clientCredentialsClient) LoginURL(_ string) string {
	return ""
}

// IDTokenLifetime returns 0 as there is no id_token issued
func (c *clientCredentialsClient) IDTokenLifetime() time.Duration {
	return 0
}

// DevMode returns false as there is no dev mode
func (c *clientCredentialsClient) DevMode() bool {
	return false
}

// RestrictAdditionalIdTokenScopes returns nil as no id_token is issued
func (c *clientCredentialsClient) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return nil
}

// RestrictAdditionalAccessTokenScopes returns the scope allowed for the token to be created because of the client credentials request
// currently it allows all scopes to be used in the access token
func (c *clientCredentialsClient) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}

func (c *clientCredentialsClient) IsScopeAllowed(scope string) bool {
	return isScopeAllowed(scope) || strings.HasPrefix(scope, ScopeProjectRolePrefix)
}

// IDTokenUserinfoClaimsAssertion returns null false as no id_token is issued
func (c *clientCredentialsClient) IDTokenUserinfoClaimsAssertion() bool {
	return false
}

// ClockSkew enable handling clock skew of the token validation. The duration (0-5s) will be added to exp claim and subtracted from iats,
// auth_time and nbf of the token to be created because of the client credentials request.
// It returns 0 as clock skew is not implemented on machine users.
func (c *clientCredentialsClient) ClockSkew() time.Duration {
	return 0
}
