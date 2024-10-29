package oidc

import (
	"context"
	"encoding/base64"
	"slices"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	getUserInfo := s.getUserInfo(session.UserID, projectID, projectRoleAssertion, userInfoAssertion, session.Scope)
	getSigner := s.getSignerOnce()

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

// SignerFunc is a getter function that allows add-hoc retrieval of the instance's signer.
type SignerFunc func(ctx context.Context) (jose.Signer, jose.SignatureAlgorithm, error)

func (s *Server) getSignerOnce() SignerFunc {
	return GetSignerOnce(s.query.GetActiveSigningWebKey, s.Provider().Storage().SigningKey)
}

// GetSignerOnce returns a function which retrieves the instance's signer from the database once.
// Repeated calls of the returned function return the same results.
func GetSignerOnce(
	getActiveSigningWebKey func(ctx context.Context) (*jose.JSONWebKey, error),
	getSigningKey func(ctx context.Context) (op.SigningKey, error),
) SignerFunc {
	var (
		once    sync.Once
		signer  jose.Signer
		signAlg jose.SignatureAlgorithm
		err     error
	)
	return func(ctx context.Context) (jose.Signer, jose.SignatureAlgorithm, error) {
		once.Do(func() {
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()

			if authz.GetFeatures(ctx).WebKey {
				var webKey *jose.JSONWebKey
				webKey, err = getActiveSigningWebKey(ctx)
				if err != nil {
					return
				}
				signer, signAlg, err = signerFromWebKey(webKey)
				return
			}

			var signingKey op.SigningKey
			signingKey, err = getSigningKey(ctx)
			if err != nil {
				return
			}
			signAlg = signingKey.SignatureAlgorithm()
			signer, err = op.SignerFromKey(signingKey)
		})
		return signer, signAlg, err
	}
}

func signerFromWebKey(signingKey *jose.JSONWebKey) (jose.Signer, jose.SignatureAlgorithm, error) {
	signAlg := jose.SignatureAlgorithm(signingKey.Algorithm)
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: signAlg,
			Key:       signingKey,
		},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, "", zerrors.ThrowInternal(err, "OIDC-oaF0s", "Errors.Internal")
	}
	return signer, signAlg, nil
}

// userInfoFunc is a getter function that allows add-hoc retrieval of a user.
type userInfoFunc func(ctx context.Context, roleAssertion bool, triggerType domain.TriggerType) (*oidc.UserInfo, error)

// getUserInfo returns a function which retrieves userinfo from the database once.
// However, each time, role claims are asserted and also action flows will trigger.
func (s *Server) getUserInfo(userID, projectID string, projectRoleAssertion, userInfoAssertion bool, scope []string) userInfoFunc {
	userInfo := s.userInfo(userID, scope, projectID, projectRoleAssertion, userInfoAssertion, false)
	return func(ctx context.Context, roleAssertion bool, triggerType domain.TriggerType) (*oidc.UserInfo, error) {
		return userInfo(ctx, roleAssertion, triggerType)
	}
}

func (*Server) createIDToken(ctx context.Context, client op.Client, getUserInfo userInfoFunc, roleAssertion bool, getSigningKey SignerFunc, sessionID, accessToken string, audience []string, authMethods []domain.UserAuthMethodType, authTime time.Time, nonce string, actor *domain.TokenActor) (idToken string, exp uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userInfo, err := getUserInfo(ctx, roleAssertion, domain.TriggerTypePreUserinfoCreation)
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

func (s *Server) createJWT(ctx context.Context, client op.Client, session *command.OIDCSession, getUserInfo userInfoFunc, assertRoles bool, getSigner SignerFunc) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userInfo, err := getUserInfo(ctx, assertRoles, domain.TriggerTypePreAccessTokenCreation)
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

// decryptCode decrypts a code or refresh_token
func (s *Server) decryptCode(ctx context.Context, code string) (_ string, err error) {
	_, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	decoded, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return "", err
	}
	return s.encAlg.DecryptString(decoded, s.encAlg.EncryptionKeyID())
}
