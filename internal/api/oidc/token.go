package oidc

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/oidc/sign"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

/*
For each grant-type, tokens creation follows the same rough logical steps:

1. Information gathering: who is requesting the token, what do we put in the claims?
2. Decision making: is the request authorized? (valid exchange code, auth request completed, valid token etc...)
3. Build an OIDC session in storage: inform the eventstore we are creating tokens.
4. Use the OIDC session to encrypt and / or sign the requested tokens

In some cases step 1 till 3 are completely implemented in the command package,
for example the v2 code exchange and refresh token.
*/

func (s *Server) accessTokenResponseFromSession(ctx context.Context, client op.Client, session *command.OIDCSession, state, projectID string, projectRoleAssertion, accessTokenRoleAssertion, idTokenRoleAssertion, userInfoAssertion bool) (_ *oidc.AccessTokenResponse, err error) {
	getUserInfo := s.getUserInfo(session.UserID, projectID, client.GetID(), projectRoleAssertion, userInfoAssertion, session.Scope)
	return s.accessTokenResponse(ctx, client, session, state, getUserInfo, s.getSignerOnce(), accessTokenRoleAssertion, idTokenRoleAssertion)
}

// accessTokenResponse creates the token response for the session.
// The user info (including the action flows) and the signing key are taken from getUserInfo and getSigner,
// so callers can resolve them before the session is stored, see [Server.prepareRefreshTokenResponse].
func (s *Server) accessTokenResponse(ctx context.Context, client op.Client, session *command.OIDCSession, state string, getUserInfo userInfoFunc, getSigner sign.SignerFunc, accessTokenRoleAssertion, idTokenRoleAssertion bool) (_ *oidc.AccessTokenResponse, err error) {
	resp := &oidc.AccessTokenResponse{
		TokenType:    oidc.BearerToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    timeToOIDCExpiresIn(session.Expiration),
		State:        state,
	}

	// If the session does not have a token ID, it is an implicit ID-Token only response.
	if session.TokenID != "" {
		if client.AccessTokenType() == op.AccessTokenTypeJWT {
			resp.AccessToken, err = s.createJWT(ctx, client, session, getUserInfo, accessTokenRoleAssertion, getSigner)
		} else {
			resp.AccessToken, err = op.CreateBearerToken(session.TokenID, session.UserID, s.opCrypto)
		}
		if err != nil {
			return nil, err
		}
	}

	if slices.Contains(session.Scope, oidc.ScopeOpenID) {
		resp.IDToken, _, err = s.createIDToken(ctx, client, getUserInfo, idTokenRoleAssertion, getSigner, session.SessionID, resp.AccessToken, session.Audience, session.AuthMethods, session.AuthTime, session.Nonce, session.Actor)
	}
	return resp, err
}

func (s *Server) getSignerOnce() sign.SignerFunc {
	return sign.GetSignerOnce(s.query.GetActiveSigningWebKey)
}

// userInfoFunc is a getter function that allows add-hoc retrieval of a user.
type userInfoFunc func(ctx context.Context, roleAssertion bool, triggerType domain.TriggerType, actor *domain.TokenActor) (*oidc.UserInfo, error)

// getUserInfo returns a function which retrieves userinfo from the database once.
// However, each time, role claims are asserted and also action flows will trigger.
func (s *Server) getUserInfo(userID, projectID, clientID string, projectRoleAssertion, userInfoAssertion bool, scope []string) userInfoFunc {
	userInfo := s.userInfo(userID, scope, projectID, clientID, projectRoleAssertion, userInfoAssertion, false)
	return func(ctx context.Context, roleAssertion bool, triggerType domain.TriggerType, actor *domain.TokenActor) (*oidc.UserInfo, error) {
		return userInfo(ctx, roleAssertion, triggerType, actor)
	}
}

// cacheUserInfo returns a userInfoFunc which calls getUserInfo once per role assertion and trigger type,
// and returns the stored result, including an error, on later calls.
// This allows running the action flows ahead of building a token response without running them twice.
// The actor is not part of the key, because all calls for one token response pass the actor of the same session.
func cacheUserInfo(getUserInfo userInfoFunc) userInfoFunc {
	type key struct {
		roleAssertion bool
		triggerType   domain.TriggerType
	}
	type result struct {
		userInfo *oidc.UserInfo
		err      error
	}
	var mu sync.Mutex
	results := make(map[key]result)
	return func(ctx context.Context, roleAssertion bool, triggerType domain.TriggerType, actor *domain.TokenActor) (*oidc.UserInfo, error) {
		mu.Lock()
		defer mu.Unlock()
		k := key{roleAssertion: roleAssertion, triggerType: triggerType}
		if r, ok := results[k]; ok {
			return r.userInfo, r.err
		}
		userInfo, err := getUserInfo(ctx, roleAssertion, triggerType, actor)
		results[k] = result{userInfo: userInfo, err: err}
		return userInfo, err
	}
}

func (*Server) createIDToken(ctx context.Context, client op.Client, getUserInfo userInfoFunc, roleAssertion bool, getSigningKey sign.SignerFunc, sessionID, accessToken string, audience []string, authMethods []domain.UserAuthMethodType, authTime time.Time, nonce string, actor *domain.TokenActor) (idToken string, exp uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userInfo, err := getUserInfo(ctx, roleAssertion, domain.TriggerTypePreUserinfoCreation, actor)
	if err != nil {
		return "", 0, err
	}

	signer, signAlg, err := getSigningKey(ctx)
	if err != nil {
		return "", 0, err
	}

	expTime := time.Now().Add(client.IDTokenLifetime()).Add(client.ClockSkew())
	claims := oidc.NewIDTokenClaims(
		op.IssuerFromContext(ctx),
		"",
		audience,
		expTime,
		authTime,
		nonce,
		"",
		AuthMethodTypesToAMR(authMethods),
		client.GetID(),
		client.ClockSkew(),
	)
	claims.SessionID = sessionID
	claims.Actor = actorDomainToClaims(actor)
	claims.SetUserInfo(userInfo)
	if accessToken != "" {
		claims.AccessTokenHash, err = oidc.ClaimHash(accessToken, signAlg)
		if err != nil {
			return "", 0, err
		}
	}
	idToken, err = crypto.Sign(claims, signer)
	return idToken, timeToOIDCExpiresIn(expTime), err
}

func timeToOIDCExpiresIn(exp time.Time) uint64 {
	return uint64(time.Until(exp) / time.Second)
}

func (s *Server) createJWT(ctx context.Context, client op.Client, session *command.OIDCSession, getUserInfo userInfoFunc, assertRoles bool, getSigner sign.SignerFunc) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userInfo, err := getUserInfo(ctx, assertRoles, domain.TriggerTypePreAccessTokenCreation, session.Actor)
	if err != nil {
		return "", err
	}
	signer, _, err := getSigner(ctx)
	if err != nil {
		return "", err
	}

	expTime := session.Expiration.Add(client.ClockSkew())
	claims := oidc.NewAccessTokenClaims(
		op.IssuerFromContext(ctx),
		userInfo.Subject,
		session.Audience,
		expTime,
		session.TokenID,
		client.GetID(),
		client.ClockSkew(),
	)
	claims.Actor = actorDomainToClaims(session.Actor)
	claims.Claims = userInfo.Claims

	return crypto.Sign(claims, signer)
}
