package oidc

import (
	"context"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (s *Server) JWTProfile(ctx context.Context, r *op.Request[oidc.JWTProfileGrantRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		span.EndWithError(err)
		err = oidcError(err)
	}()

	user, err := s.verifyJWTProfile(ctx, r.Data)
	if err != nil {
		return nil, err
	}

	client := &clientCredentialsClient{
		clientID:      user.Username,
		userID:        user.UserID,
		resourceOwner: user.ResourceOwner,
		tokenType:     user.TokenType,
	}
	scope, err := op.ValidateAuthReqScopes(client, r.Data.Scope)
	if err != nil {
		return nil, err
	}
	scope, err = s.checkOrgScopes(ctx, client.resourceOwner, scope)
	if err != nil {
		return nil, err
	}

	session, err := s.command.CreateOIDCSession(ctx,
		client.userID,
		client.resourceOwner,
		client.clientID,
		"", // backChannelLogoutURI not needed for service user session
		scope,
		domain.AddAudScopeToAudience(ctx, nil, r.Data.Scope),
		[]domain.UserAuthMethodType{domain.UserAuthMethodTypePrivateKey},
		time.Now(),
		"",
		nil,
		nil,
		domain.TokenReasonJWTProfile,
		nil,
		false,
		"",
		domain.OIDCResponseTypeUnspecified,
	)
	if err != nil {
		return nil, err
	}
	return response(s.accessTokenResponseFromSession(ctx, client, session, "", "", false, true, true, false))
}

func (s *Server) verifyJWTProfile(ctx context.Context, req *oidc.JWTProfileGrantRequest) (_ *query.AuthNKeyUser, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	storage := &jwtProfileKeyStorage{query: s.query}
	verifier := op.NewJWTProfileVerifier(
		storage, op.IssuerFromContext(ctx),
		time.Hour, time.Second,
	)
	_, err = op.VerifyJWTAssertion(ctx, req.Assertion, verifier)
	if err != nil {
		return nil, err
	}
	return storage.user, nil
}

type jwtProfileKeyStorage struct {
	query *query.Queries
	user  *query.AuthNKeyUser // only populated after GetKeyByIDAndClientID is called
}

func (s *jwtProfileKeyStorage) GetKeyByIDAndClientID(ctx context.Context, keyID, userID string) (_ *jose.JSONWebKey, err error) {
	s.user, err = s.query.GetAuthNKeyUser(ctx, keyID, userID)
	if err != nil {
		return nil, err
	}
	publicKey, err := crypto.BytesToPublicKey(s.user.PublicKey)
	if err != nil {
		return nil, err
	}
	return &jose.JSONWebKey{
		KeyID: keyID,
		Use:   "sig",
		Key:   publicKey,
	}, nil
}
