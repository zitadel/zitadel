package oidc

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UserIDTokenType oidc.TokenType = "urn:zitadel:params:oauth:token-type:user_id"

	// TokenTypeNA is set when the returned Token Exchange access token value can't be used as an access token.
	// For example, when it is an ID Token.
	// See [RFC 8693, section 2.2.1, token_type](https://www.rfc-editor.org/rfc/rfc8693#section-2.2.1)
	TokenTypeNA = "N_A"
)

func init() {
	oidc.AllTokenTypes = append(oidc.AllTokenTypes, UserIDTokenType)
}

func (s *Server) TokenExchange(ctx context.Context, r *op.ClientRequest[oidc.TokenExchangeRequest]) (_ *op.Response, err error) {
	resp, err := s.tokenExchange(ctx, r)
	if err != nil {
		return nil, oidcError(err)
	}
	return resp, nil
}

func (s *Server) tokenExchange(ctx context.Context, r *op.ClientRequest[oidc.TokenExchangeRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !authz.GetFeatures(ctx).TokenExchange {
		return nil, zerrors.ThrowPreconditionFailed(nil, "OIDC-oan4I", "Errors.TokenExchange.FeatureDisabled")
	}
	if len(r.Data.Resource) > 0 {
		return nil, oidc.ErrInvalidTarget().WithDescription("resource parameter not supported")
	}

	client, ok := r.Client.(*Client)
	if !ok {
		// not supposed to happen, but just preventing a panic if it does.
		return nil, zerrors.ThrowInternal(nil, "OIDC-eShi5", "Error.Internal")
	}

	subjectToken, err := s.verifyExchangeToken(ctx, client, r.Data.SubjectToken, r.Data.SubjectTokenType, oidc.AllTokenTypes...)
	if err != nil {
		return nil, oidc.ErrInvalidRequest().WithParent(err).WithDescription("subject_token invalid")
	}

	actorToken := subjectToken // see [createExchangeTokens] comment.
	if subjectToken.tokenType == UserIDTokenType || subjectToken.tokenType == oidc.JWTTokenType || r.Data.ActorToken != "" {
		if !authz.GetInstance(ctx).EnableImpersonation() {
			return nil, zerrors.ThrowPermissionDenied(nil, "OIDC-Fae5w", "Errors.TokenExchange.Impersonation.PolicyDisabled")
		}
		actorToken, err = s.verifyExchangeToken(ctx, client, r.Data.ActorToken, r.Data.ActorTokenType, oidc.AccessTokenType, oidc.IDTokenType, oidc.RefreshTokenType)
		if err != nil {
			return nil, oidc.ErrInvalidRequest().WithParent(err).WithDescription("actor_token invalid")
		}
		ctx = authz.SetCtxData(ctx, authz.CtxData{
			UserID: actorToken.userID,
			OrgID:  actorToken.resourceOwner,
		})
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

// verifyExchangeToken verifies the passed token based on the token type. It is safe to pass both from the request as-is.
// A list of allowed token types must be passed to determine which types are trusted at a particular stage of the token exchange.
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
		if token.isPAT {
			if err = s.assertClientScopesForPAT(ctx, token, client.GetID(), client.client.ProjectID); err != nil {
				return nil, err
			}
		}
		return accessToExchangeToken(token, op.IssuerFromContext(ctx)), nil

	case oidc.IDTokenType:
		verifier := op.NewIDTokenHintVerifier(op.IssuerFromContext(ctx), s.idTokenHintKeySet)
		claims, err := op.VerifyIDTokenHint[*oidc.IDTokenClaims](ctx, token, verifier)
		if err != nil {
			return nil, zerrors.ThrowPermissionDenied(err, "OIDC-Rei0f", "Errors.TokenExchange.Token.Invalid")
		}
		resourceOwner, ok := claims.Claims[ClaimResourceOwnerID].(string)
		if !ok || resourceOwner == "" {
			user, err := s.query.GetUserByID(ctx, false, token)
			if err != nil {
				return nil, zerrors.ThrowPermissionDenied(err, "OIDC-aD0Oo", "Errors.TokenExchange.Token.Invalid")
			}
			resourceOwner = user.ResourceOwner
		}

		return idTokenClaimsToExchangeToken(claims, resourceOwner), nil

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

func validateTokenExchangeScopes(client *Client, requestedScopes, subjectScopes, actorScopes []string) ([]string, error) {
	// Scope always has 1 empty string is the space delimited array was an empty string.
	scopes := slices.DeleteFunc(requestedScopes, func(s string) bool {
		return s == ""
	})
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

	//nolint:gocritic
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
		userInfo   *oidc.UserInfo
		signingKey op.SigningKey
	)
	if slices.Contains(scopes, oidc.ScopeOpenID) || tokenType == oidc.JWTTokenType || tokenType == oidc.IDTokenType {
		projectID := client.client.ProjectID
		userInfo, err = s.userInfo(ctx, subjectToken.userID, scopes, projectID, client.client.ProjectRoleAssertion, false)
		if err != nil {
			return nil, err
		}
		signingKey, err = s.Provider().Storage().SigningKey(ctx)
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
	case oidc.AccessTokenType, "":
		resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.createExchangeAccessToken(ctx, client, subjectToken.resourceOwner, subjectToken.userID, audience, scopes, actorToken.authMethods, actorToken.authTime, reason, actor)
		resp.TokenType = oidc.BearerToken
		resp.IssuedTokenType = oidc.AccessTokenType

	case oidc.JWTTokenType:
		resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.createExchangeJWT(ctx, signingKey, client, subjectToken.resourceOwner, subjectToken.userID, audience, scopes, actorToken.authMethods, actorToken.authTime, reason, actor, userInfo.Claims)
		resp.TokenType = oidc.BearerToken
		resp.IssuedTokenType = oidc.JWTTokenType

	case oidc.IDTokenType:
		resp.AccessToken, resp.ExpiresIn, err = s.createExchangeIDToken(ctx, signingKey, client, subjectToken.userID, "", audience, userInfo, actorToken.authMethods, actorToken.authTime, reason, actor)
		resp.TokenType = TokenTypeNA
		resp.IssuedTokenType = oidc.IDTokenType
	case oidc.RefreshTokenType, UserIDTokenType:
		fallthrough
	default:
		err = zerrors.ThrowInvalidArgument(nil, "OIDC-wai5E", "Errors.TokenExchange.Token.TypeNotSupported")
	}
	if err != nil {
		return nil, err
	}

	if slices.Contains(scopes, oidc.ScopeOpenID) && tokenType != oidc.IDTokenType {
		resp.IDToken, _, err = s.createExchangeIDToken(ctx, signingKey, client, subjectToken.userID, resp.AccessToken, audience, userInfo, actorToken.authMethods, actorToken.authTime, reason, actor)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) createExchangeAccessToken(ctx context.Context, client *Client, resourceOwner, userID string, audience, scopes []string, authMethods []domain.UserAuthMethodType, authTime time.Time, reason domain.TokenReason, actor *domain.TokenActor) (accessToken string, refreshToken string, exp uint64, err error) {
	tokenInfo, refreshToken, err := s.createAccessTokenCommands(ctx, client, resourceOwner, userID, audience, scopes, authMethods, authTime, reason, actor)
	if err != nil {
		return "", "", 0, err
	}
	accessToken, err = op.CreateBearerToken(tokenInfo.TokenID, userID, s.Provider().Crypto())
	if err != nil {
		return "", "", 0, err
	}
	return accessToken, refreshToken, timeToOIDCExpiresIn(tokenInfo.Expiration), nil
}

func (s *Server) createExchangeJWT(ctx context.Context, signingKey op.SigningKey, client *Client, resourceOwner, userID string, audience, scopes []string, authMethods []domain.UserAuthMethodType, authTime time.Time, reason domain.TokenReason, actor *domain.TokenActor, privateClaims map[string]any) (accessToken string, refreshToken string, exp uint64, err error) {
	tokenInfo, refreshToken, err := s.createAccessTokenCommands(ctx, client, resourceOwner, userID, audience, scopes, authMethods, authTime, reason, actor)
	if err != nil {
		return "", "", 0, err
	}

	expTime := tokenInfo.Expiration.Add(client.ClockSkew())
	claims := oidc.NewAccessTokenClaims(op.IssuerFromContext(ctx), userID, tokenInfo.Audience, expTime, tokenInfo.TokenID, client.GetID(), client.ClockSkew())
	claims.Actor = actorDomainToClaims(tokenInfo.Actor)
	claims.Claims = privateClaims

	signer, err := op.SignerFromKey(signingKey)
	if err != nil {
		return "", "", 0, err
	}

	accessToken, err = crypto.Sign(claims, signer)
	if err != nil {
		return "", "", 0, nil
	}
	return accessToken, refreshToken, timeToOIDCExpiresIn(expTime), nil
}

func (s *Server) createExchangeIDToken(ctx context.Context, signingKey op.SigningKey, client *Client, userID, accessToken string, audience []string, userInfo *oidc.UserInfo, authMethods []domain.UserAuthMethodType, authTime time.Time, reason domain.TokenReason, actor *domain.TokenActor) (idToken string, exp uint64, err error) {
	expTime := time.Now().Add(client.IDTokenLifetime()).Add(client.ClockSkew())
	claims := oidc.NewIDTokenClaims(op.IssuerFromContext(ctx), userID, audience, expTime, authTime, "", "", AuthMethodTypesToAMR(authMethods), client.GetID(), client.ClockSkew())
	claims.Actor = actorDomainToClaims(actor)
	claims.SetUserInfo(userInfo)
	if accessToken != "" {
		claims.AccessTokenHash, err = oidc.ClaimHash(accessToken, signingKey.SignatureAlgorithm())
		if err != nil {
			return "", 0, err
		}
	}
	signer, err := op.SignerFromKey(signingKey)
	if err != nil {
		return "", 0, err
	}
	idToken, err = crypto.Sign(claims, signer)
	return idToken, timeToOIDCExpiresIn(expTime), err
}

func timeToOIDCExpiresIn(exp time.Time) uint64 {
	return uint64(time.Until(exp) / time.Second)
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
	return tokenInfo, "", err
}
