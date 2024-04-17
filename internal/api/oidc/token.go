package oidc

import (
	"context"
	"encoding/base64"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

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
for example the v2 code exchange.
*/

// userInfoAndSignerFunc is a getter function that allows add-hoc retrieval of userinfo and signer,
// when a JWT token needs to be created in the token endpoint.
type userInfoAndSignerFunc func(ctx context.Context) (*oidc.UserInfo, jose.Signer, jose.SignatureAlgorithm, error)

// getUserInfoAndSignerOnce returns a function which retrieves the userinfo and signer from the database once.
// Repeated calls of the returned function return the same results.
func (s *Server) getUserInfoAndSignerOnce(ctx context.Context, userID, projectID string, projectRoleAssertion bool, scope []string) userInfoAndSignerFunc {
	var (
		userInfo *oidc.UserInfo
		signer   jose.Signer
		signAlg  jose.SignatureAlgorithm
		err      error
	)
	call := sync.OnceFunc(func() {
		ctx, span := tracing.NewSpan(ctx)
		defer func() { span.EndWithError(err) }()

		userInfo, err = s.userInfo(ctx, userID, scope, projectID, projectRoleAssertion, false)
		if err != nil {
			return
		}
		var signingKey op.SigningKey
		signingKey, err = s.Provider().Storage().SigningKey(ctx)
		if err != nil {
			return
		}
		signAlg = signingKey.SignatureAlgorithm()

		signer, err = op.SignerFromKey(signingKey)
		if err != nil {
			return
		}
	})
	return func(ctx context.Context) (*oidc.UserInfo, jose.Signer, jose.SignatureAlgorithm, error) {
		call()
		return userInfo, signer, signAlg, err
	}
}

func (*Server) createIDToken(ctx context.Context, client *Client, getUserInfoAndSigningKey userInfoAndSignerFunc, accessToken string, audience []string, authMethods []domain.UserAuthMethodType, authTime time.Time, actor *domain.TokenActor) (idToken string, exp uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userInfo, signer, signAlg, err := getUserInfoAndSigningKey(ctx)
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
		"",
		"",
		AuthMethodTypesToAMR(authMethods),
		client.GetID(),
		client.ClockSkew(),
	)
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

func (*Server) createJWT(ctx context.Context, client *Client, session *command.OIDCSession, getUserInfoAndSigningKey userInfoAndSignerFunc) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userInfo, signer, _, err := getUserInfoAndSigningKey(ctx)
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
