package oidc

import (
	"context"
	"encoding/base64"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

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
for example the v2 code exchange.
*/

func (s *Server) CodeExchange(ctx context.Context, r *op.ClientRequest[oidc.AccessTokenRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	client, ok := r.Client.(*Client)
	if !ok {
		// not supposed to happen, but just preventing a panic if it does.
		return nil, zerrors.ThrowInternal(nil, "OIDC-eShi5", "Error.Internal")
	}

	plainCode, err := s.decryptCode(ctx, r.Data.Code)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "OIDC-ahLi2", "Errors.User.Code.Invalid")
	}

	var (
		session *command.OIDCSession
		state   string
	)
	if strings.HasPrefix(plainCode, command.IDPrefixV2) {
		session, state, err = s.command.CreateOIDCSessionFromCodeExchange(
			setContextUserSystem(ctx), plainCode, authRequestComplianceChecker(client, r.Data),
		)
	} else {
		session, state, err = s.codeExchangeV1(ctx, client, r.Data, plainCode)
	}
	if err != nil {
		return nil, err
	}

	accessToken, idToken, err := s.createTokensFromSession(ctx, client, session)
	if err != nil {
		return nil, err
	}
	return op.NewResponse(&oidc.AccessTokenResponse{
		AccessToken:  accessToken,
		TokenType:    oidc.BearerToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    timeToOIDCExpiresIn(session.Expiration),
		IDToken:      idToken,
		State:        state,
	}), nil
}

func (s *Server) createTokensFromSession(ctx context.Context, client *Client, session *command.OIDCSession) (accessToken, idToken string, err error) {
	getUserInfoAndSigner := s.getUserInfoAndSignerOnce(ctx, session.UserID, client.client.ProjectID, client.client.ProjectRoleAssertion, session.Scopes)

	if client.AccessTokenType() == op.AccessTokenTypeJWT {
		accessToken, err = s.createJWT(ctx, client, session, getUserInfoAndSigner)
	} else {
		accessToken, err = op.CreateBearerToken(session.TokenID, session.UserID, s.opCrypto)
	}
	if err != nil {
		return "", "", err
	}

	if slices.Contains(session.Scopes, oidc.ScopeOpenID) {
		idToken, _, err = s.createIDToken(ctx, client, getUserInfoAndSigner, accessToken, session.Audience, nil, time.Time{}, nil) // TODO: authmethods, authTime
	}

	return accessToken, idToken, err
}

// codeExchangeV1 creates a v2 token from a v1 auth request.
func (s *Server) codeExchangeV1(ctx context.Context, client *Client, req *oidc.AccessTokenRequest, plainCode string) (session *command.OIDCSession, state string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	authReq, err := s.getAuthRequestV1ByCode(ctx, plainCode)
	if err != nil {
		return nil, "", err
	}

	if challenge := authReq.GetCodeChallenge(); challenge != nil || client.AuthMethod() == oidc.AuthMethodNone {
		if err = op.AuthorizeCodeChallenge(req.CodeVerifier, challenge); err != nil {
			return nil, "", err
		}
	}
	if req.RedirectURI != authReq.GetRedirectURI() {
		return nil, "", oidc.ErrInvalidGrant().WithDescription("redirect_uri does not correspond")
	}
	userAgentID, _, userOrgID, authTime, authMethodsReferences, reason, actor := getInfoFromRequest(authReq)

	session, err = s.command.CreateOIDCSession(ctx,
		authReq.GetSubject(),
		userOrgID,
		client.client.ClientID,
		authReq.GetAudience(),
		authReq.GetScopes(),
		AMRToAuthMethodTypes(authMethodsReferences),
		authTime,
		&domain.UserAgent{
			FingerprintID: &userAgentID,
		},
		reason,
		actor,
	)
	return session, authReq.GetState(), err
}

func (s *Server) getAuthRequestV1ByCode(ctx context.Context, plainCode string) (op.AuthRequest, error) {
	authReq, err := s.repo.AuthRequestByCode(ctx, plainCode)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(authReq)
}

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

		userInfo, err = s.userInfo(ctx, userID, projectID, projectRoleAssertion, scope, []string{projectID})
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

func authRequestComplianceChecker(client *Client, req *oidc.AccessTokenRequest) command.AuthRequestComplianceChecker {
	return func(ctx context.Context, authReq *command.AuthRequestWriteModel) error {
		if authReq.CodeChallenge != nil || client.AuthMethod() == oidc.AuthMethodNone {
			err := op.AuthorizeCodeChallenge(req.CodeVerifier, CodeChallengeToOIDC(authReq.CodeChallenge))
			if err != nil {
				return err
			}
		}
		if req.RedirectURI != authReq.RedirectURI {
			return oidc.ErrInvalidGrant().WithDescription("redirect_uri does not correspond")
		}
		if err := authReq.CheckAuthenticated(); err != nil {
			return err
		}
		return nil
	}
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
