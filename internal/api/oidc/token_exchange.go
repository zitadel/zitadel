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
		return nil, oidc.ErrInvalidTarget().WithDescription("resource parameter not supported")
	}

	client, ok := r.Client.(*Client)
	if !ok {
		// not supposed to happen, but just preventing a panic if it does.
		return nil, zerrors.ThrowInternal(nil, "OIDC-eShi5", "Error.Internal")
	}

	subjectToken, err := s.verifyExchangeToken(ctx, client, r.Data.SubjectToken, r.Data.SubjectTokenType)
	if err != nil {
		return nil, oidc.ErrInvalidRequest().WithParent(err).WithDescription("subject_token invalid")
	}

	actorToken := subjectToken // see [createExchangeTokens] comment.
	if subjectToken.tokenType == UserIDTokenType || r.Data.ActorToken != "" {
		if !authz.GetInstance(ctx).EnableImpersonation() {
			return nil, zerrors.ThrowPermissionDenied(nil, "OIDC-Fae5w", "Errors.TokenExchange.Impersonation.PolicyDisabled")
		}
		actorToken, err = s.verifyExchangeToken(ctx, client, r.Data.ActorToken, r.Data.ActorTokenType)
		if err != nil {
			return nil, oidc.ErrInvalidRequest().WithParent(err).WithDescription("actor_token invalid")
		}
	}

	audience, err := validateTokenExchangeAudience(r.Data.Audience, subjectToken.audience, actorToken.audience)
	if err != nil {
		return nil, err
	}
	scopes, err := validateTokenExchangeScopes(client, r.Data.Scopes, subjectToken.scopes, actorToken.scopes)
	if err != nil {
		return nil, err
	}

	resp, err := s.createExchangeTokens(ctx, r.Data.RequestedTokenType, client, subjectToken, actorToken, audience, scopes)
	if err != nil {
		return nil, err
	}

	return op.NewResponse(resp), nil
}

func validateTokenExchangeScopes(client *Client, requestedScopes, subjectScopes, actorScopes []string) ([]string, error) {
	scopes := requestedScopes
	if len(scopes) == 0 {
		scopes = subjectScopes
	}
	if len(scopes) == 0 {
		scopes = actorScopes
	}
	return op.ValidateAuthReqScopes(client, scopes)
}

func validateTokenExchangeAudience(requestedAudience, subjectAudience, actorAudience []string) ([]string, error) {
	if len(requestedAudience) == 0 {
		if len(subjectAudience) > 0 {
			return subjectAudience, nil
		}
		if len(actorAudience) > 0 {
			return actorAudience, nil
		}
	}
	if slices.Equal(requestedAudience, subjectAudience) || slices.Equal(requestedAudience, actorAudience) {
		return requestedAudience, nil
	}
	allowedAudience := append(subjectAudience, actorAudience...)
	for _, a := range requestedAudience {
		if !slices.Contains(allowedAudience, a) {
			return nil, oidc.ErrInvalidTarget().WithDescription("audience %q not found in subject or actor token", a)
		}
	}
	return requestedAudience, nil
}

// createExchangeTokens prepares the final tokens to be returned to the client.
// The subjectToken is used to set the new token's subject and resource owner.
// The actorToken is used to set the new token's auth time AMR and actor.
// Both tokens may point to the same object (subjectToken) in case of a regular Token Exchange.
// When the subject and actor Tokens point to different objects, the new tokens will be for impersonation / delegation.
func (s *Server) createExchangeTokens(ctx context.Context, tokenType oidc.TokenType, client *Client, subjectToken, actorToken *exchangeToken, audience, scopes []string) (_ *oidc.TokenExchangeResponse, err error) {
	var (
		userInfo *oidc.UserInfo
		signer   jose.Signer
	)
	if slices.Contains(scopes, oidc.ScopeOpenID) || tokenType == oidc.JWTTokenType || tokenType == oidc.IDTokenType {
		projectID := client.client.ProjectID
		userInfo, err = s.userInfo(ctx, subjectToken.userID, projectID, scopes, []string{projectID})
		if err != nil {
			return nil, err
		}
		signer, err = s.getSigner(ctx)
		if err != nil {
			return nil, err
		}
	}

	resp := &oidc.TokenExchangeResponse{
		Scopes: scopes,
	}

	reason := domain.TokenReasonExchange
	actor := actorToken.actor
	if subjectToken != actorToken {
		reason = domain.TokenReasonImpersonation
		actor = actorToken.nestedActor()
	}

	switch tokenType {
	case oidc.AccessTokenType:
		resp.AccessToken, resp.RefreshToken, err = s.createExchangeAccessToken(ctx, client, subjectToken.resourceOwner, subjectToken.userID, audience, scopes, actorToken.authMethods, actorToken.authTime, reason, actor)
		resp.TokenType = oidc.BearerToken
		resp.IssuedTokenType = oidc.AccessTokenType

	case oidc.JWTTokenType:
		resp.AccessToken, resp.RefreshToken, err = s.createExchangeJWT(ctx, signer, client, subjectToken.resourceOwner, subjectToken.userID, audience, scopes, actorToken.authMethods, actorToken.authTime, reason, actor, userInfo.Claims)
		resp.TokenType = oidc.BearerToken
		resp.IssuedTokenType = oidc.JWTTokenType
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
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
	// TODO: permission check
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
	return tokenInfo, "", err
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
	audience      []string
	scopes        []string
}

func (et *exchangeToken) nestedActor() *domain.TokenActor {
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
		audience:      token.audience,
		scopes:        token.scope,
	}
}

func idTokenClaimsToExchangeToken(claims *oidc.IDTokenClaims) *exchangeToken {
	resourceOwner, _ := claims.Claims[ClaimResourceOwnerID].(string)
	return &exchangeToken{
		tokenType:     oidc.IDTokenType,
		userID:        claims.Subject,
		issuer:        claims.Issuer,
		resourceOwner: resourceOwner,
		authTime:      claims.GetAuthTime(),
		authMethods:   AMRToAuthMethodTypes(claims.AuthenticationMethodsReferences),
		actor:         actorClaimsToDomain(claims.Actor),
		audience:      claims.Audience,
	}
}

func actorClaimsToDomain(actor *oidc.ActorClaims) *domain.TokenActor {
	if actor == nil {
		return nil
	}
	resourceOwner, _ := actor.Claims[ClaimResourceOwnerID].(string)
	return &domain.TokenActor{
		Actor:         actorClaimsToDomain(actor.Actor),
		UserID:        actor.Subject,
		Issuer:        actor.Issuer,
		ResourceOwner: resourceOwner,
	}
}

func jwtToExchangeToken(jwt *oidc.JWTTokenRequest, resourceOwner string) *exchangeToken {
	return &exchangeToken{
		tokenType:     oidc.JWTTokenType,
		userID:        jwt.Subject,
		issuer:        jwt.Issuer,
		resourceOwner: resourceOwner,
		scopes:        jwt.Scopes,
		authTime:      jwt.IssuedAt.AsTime(),
		// audience omitted as we don't thrust audiences not signed by us
	}
}

func userToExchangeToken(user *query.User) *exchangeToken {
	return &exchangeToken{
		tokenType:     UserIDTokenType,
		userID:        user.ID,
		resourceOwner: user.ResourceOwner,
	}
}

func (s *Server) verifyExchangeToken(ctx context.Context, client *Client, token string, tokenType oidc.TokenType, allowed ...oidc.TokenType) (*exchangeToken, error) {
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

	case oidc.IDTokenType:
		verifier := op.NewIDTokenHintVerifier(op.IssuerFromContext(ctx), s.idTokenHintKeySet)
		claims, err := op.VerifyIDTokenHint[*oidc.IDTokenClaims](ctx, token, verifier)
		if err != nil {
			return nil, zerrors.ThrowPermissionDenied(err, "OIDC-Rei0f", "Errors.TokenExchange.Token.Invalid")
		}
		return idTokenClaimsToExchangeToken(claims), nil

	case oidc.JWTTokenType:
		resourceOwner := new(string)
		verifier := op.NewJWTProfileVerifierKeySet(keySetMap(client.client.PublicKeys), op.IssuerFromContext(ctx), time.Hour, client.client.ClockSkew, s.jwtProfileUserCheck(ctx, resourceOwner))
		jwt, err := op.VerifyJWTAssertion(ctx, token, verifier)
		if err != nil {
			return nil, zerrors.ThrowPermissionDenied(err, "OIDC-eiS6o", "Errors.TokenExchange.Token.Invalid")
		}
		return jwtToExchangeToken(jwt, *resourceOwner), nil

	case UserIDTokenType:
		user, err := s.query.GetUserByID(ctx, false, token)
		if err != nil {
			return nil, zerrors.ThrowPermissionDenied(err, "OIDC-Nee6r", "Errors.TokenExchange.Token.Invalid")
		}
		return userToExchangeToken(user), nil

	case oidc.RefreshTokenType:
		fallthrough
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-oda4R", "Errors.TokenExchange.Token.TypeNotSupported")
	}
}

func (s *Server) jwtProfileUserCheck(ctx context.Context, resourceOwner *string) op.JWTProfileVerifierOption {
	return op.SubjectCheck(func(request *oidc.JWTTokenRequest) error {
		user, err := s.query.GetUserByID(ctx, false, request.Subject)
		if err != nil {
			return zerrors.ThrowPermissionDenied(err, "OIDC-Nee6r", "Errors.TokenExchange.Token.Invalid")
		}
		*resourceOwner = user.ResourceOwner
		return nil
	})
}
