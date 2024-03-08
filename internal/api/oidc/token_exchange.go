package oidc

import (
	"context"
	"slices"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const UserIDTokenType oidc.TokenType = "urn:zitadel:params:oauth:token-type:user_id"

func (s *Server) TokenExchange(ctx context.Context, r *op.ClientRequest[oidc.TokenExchangeRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !authz.GetFeatures(ctx).TokenExchange {
		return nil, zerrors.ThrowPreconditionFailed(nil, "OIDC-oan4I", "Errors.Feature.Disabled.TokenExchange")
	}
	if len(r.Data.Resource) > 0 {
		return nil, zerrors.ThrowUnimplemented(nil, "OIDC-aeCh6", "Errors.TokenExchange.Unimplemented.Resource")
	}
	subjectToken, err := s.verifyExchangeToken(ctx, r.Data.SubjectToken, r.Data.SubjectTokenType)
	if err != nil {
		return nil, zerrors.ThrowUnauthenticated(err, "OIDC-phiX5", "Errors.TokenExchange.InvalidSubjectToken")
	}

	audience := r.Data.Audience // TODO
	scopes := r.Data.Scopes     // TODO

	if subjectToken.tokenType == UserIDTokenType || r.Data.ActorToken != "" {
		return s.tokenExchangeImpersonation(ctx, r, subjectToken, audience, scopes)
	}

	return s.LegacyServer.TokenExchange(ctx, r)
}

func (s *Server) tokenExchangeImpersonation(ctx context.Context, r *op.ClientRequest[oidc.TokenExchangeRequest], subjectToken *exchangeToken, audience, scope []string) (_ *op.Response, err error) {
	if !authz.GetInstance(ctx).EnableImpersonation() {
		return nil, zerrors.ThrowPermissionDenied(nil, "OIDC-Fae5w", "Errors.TokenExchange.Impersonation.PolicyDisabled")
	}
	token, err := s.verifyExchangeToken(ctx, r.Data.ActorToken, r.Data.ActorTokenType)
	if err != nil {
		return nil, oidc.ErrInvalidRequest().WithParent(err).WithDescription("actor_token actor_token_type=%s", r.Data.ActorTokenType)
	}
	// TODO: permission check

	resp, err := s.createExchangeTokens(ctx, r.Data.RequestedTokenType, r.Client, subjectToken.resourceOwner, subjectToken.userID, audience, scope, token.authMethods, token.authTime, domain.TokenReasonImpersonation, token.nestActor())
	if err != nil {
		return nil, err
	}

	return &op.Response{
		Data: resp,
	}, nil
}

func (s *Server) createExchangeTokens(ctx context.Context, tokenType oidc.TokenType, client op.Client, resourceOwner, userID string, audience, scopes []string, authMethods []domain.UserAuthMethodType, authTime time.Time, reason domain.TokenReason, actor *domain.TokenActor) (_ *oidc.TokenExchangeResponse, err error) {
	zClient, ok := client.(*Client)
	if !ok {
		// not supposed to happen, but just preventing a panic if it does.
		return nil, zerrors.ThrowInternal(nil, "OIDC-eShi5", "Error.Internal")
	}

	var (
		userInfo *oidc.UserInfo
		signer   jose.Signer
	)
	if slices.Contains(scopes, oidc.ScopeOpenID) || tokenType == oidc.JWTTokenType || tokenType == oidc.IDTokenType {
		projectID := zClient.client.ProjectID
		userInfo, err = s.userInfo(ctx, userID, projectID, scopes, []string{projectID})
		if err != nil {
			return nil, err
		}
		signer, err = s.getSigner(ctx)
		if err != nil {
			return nil, err
		}
	}

	var resp oidc.TokenExchangeResponse

	switch tokenType {
	case oidc.AccessTokenType:
		resp.AccessToken, resp.RefreshToken, err = s.createExchangeAccessToken(ctx, zClient, resourceOwner, userID, audience, scopes, authMethods, authTime, reason, actor)
	case oidc.JWTTokenType:
		resp.AccessToken, resp.RefreshToken, err = s.createExchangeJWT(ctx, signer, zClient, resourceOwner, userID, audience, scopes, authMethods, authTime, reason, actor, userInfo.Claims)
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Server) createExchangeAccessToken(ctx context.Context, client *Client, resourceOwner, userID string, audience, scopes []string, authMethods []domain.UserAuthMethodType, authTime time.Time, reason domain.TokenReason, actor *domain.TokenActor) (accessToken string, refreshToken string, err error) {
	tokenInfo, refreshToken, err := s.createAccessTokenCommands(ctx, client, resourceOwner, userID, audience, scopes, authMethods, authTime, reason, actor)
	if err != nil {
		return "", "", err
	}
	accessToken, err = op.CreateBearerToken(tokenInfo.TokenID, userID, s.Provider().Crypto())
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *Server) createExchangeJWT(ctx context.Context, signer jose.Signer, client *Client, resourceOwner, userID string, audience, scopes []string, authMethods []domain.UserAuthMethodType, authTime time.Time, reason domain.TokenReason, actor *domain.TokenActor, privateClaims map[string]any) (accessToken string, refreshToken string, err error) {
	tokenInfo, refreshToken, err := s.createAccessTokenCommands(ctx, client, resourceOwner, userID, audience, scopes, authMethods, authTime, reason, actor)
	if err != nil {
		return "", "", err
	}

	exp := tokenInfo.Expiration.Add(client.ClockSkew())
	claims := oidc.NewAccessTokenClaims(op.IssuerFromContext(ctx), userID, tokenInfo.Audience, exp, tokenInfo.TokenID, client.GetID(), client.ClockSkew())
	claims.Claims = privateClaims

	accessToken, err = crypto.Sign(claims, signer)
	if err != nil {
		return "", "", nil
	}
	return accessToken, refreshToken, nil
}

func (s *Server) createAccessTokenCommands(ctx context.Context, client *Client, resourceOwner, userID string, audience, scopes []string, authMethods []domain.UserAuthMethodType, authTime time.Time, reason domain.TokenReason, actor *domain.TokenActor) (tokenInfo *domain.Token, refreshToken string, err error) {
	settings := client.client.Settings
	if slices.Contains(scopes, oidc.ScopeOfflineAccess) {
		return s.command.AddAccessAndRefreshToken(
			ctx, resourceOwner, "", client.GetID(), userID, "", audience, scopes, AuthMethodTypesToAMR(authMethods),
			settings.AccessTokenLifetime, settings.RefreshTokenIdleExpiration, settings.RefreshTokenExpiration,
			authTime, reason, actor,
		)
	}
	tokenInfo, err = s.command.AddUserToken(
		ctx, resourceOwner, "", client.GetID(), userID, audience, scopes, AuthMethodTypesToAMR(authMethods),
		settings.AccessTokenLifetime,
		authTime, reason, actor,
	)
	return tokenInfo, "", nil
}

func (s *Server) getSigner(ctx context.Context) (jose.Signer, error) {
	key, err := s.Provider().Storage().SigningKey(ctx)
	if err != nil {
		return nil, err
	}
	return op.SignerFromKey(key)
}

type exchangeToken struct {
	tokenType     oidc.TokenType
	userID        string
	issuer        string
	resourceOwner string
	authTime      time.Time
	authMethods   []domain.UserAuthMethodType
	actor         *domain.TokenActor
}

func (et *exchangeToken) nestActor() *domain.TokenActor {
	return &domain.TokenActor{
		Actor:         et.actor,
		UserID:        et.userID,
		Issuer:        et.issuer,
		ResourceOwner: et.resourceOwner,
	}
}

func accessToExchangeToken(token *accessToken, issuer string) *exchangeToken {
	return &exchangeToken{
		tokenType:     oidc.AccessTokenType,
		userID:        token.userID,
		issuer:        issuer,
		resourceOwner: token.resourceOwner,
		authMethods:   token.authMethods,
		actor:         token.actor,
	}
}

func userToExchangeToken(user *query.User) *exchangeToken {
	return &exchangeToken{
		tokenType:     UserIDTokenType,
		userID:        user.ID,
		resourceOwner: user.ResourceOwner,
	}
}

func (s *Server) verifyExchangeToken(ctx context.Context, token string, tokenType oidc.TokenType, allowed ...oidc.TokenType) (*exchangeToken, error) {
	if token == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-lei0O", "Errors.TokenExchange.Token.Missing")
	}
	if tokenType == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-sei9V", "Errors.TokenExchange.Token.TypeMissing")
	}
	if !slices.Contains(allowed, tokenType) {
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-OZ1ie", "Errors.TokenExchange.Token.TypeNotAllowed")
	}

	switch tokenType {
	case oidc.AccessTokenType:
		token, err := s.verifyAccessToken(ctx, token)
		if err != nil {
			return nil, zerrors.ThrowPermissionDenied(err, "OIDC-Osh3t", "Errors.TokenExchange.Token.Invalid")
		}
		return accessToExchangeToken(token, op.IssuerFromContext(ctx)), nil

	case UserIDTokenType:
		user, err := s.query.GetUserByID(ctx, false, token)
		if err != nil {
			return nil, zerrors.ThrowPermissionDenied(err, "OIDC-Nee6r", "Errors.TokenExchange.Token.Invalid")
		}
		return userToExchangeToken(user), nil

	case oidc.IDTokenType, oidc.RefreshTokenType, oidc.JWTTokenType:
		fallthrough
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-oda4R", "Errors.TokenExchange.Token.TypeNotSupported")
	}
}
