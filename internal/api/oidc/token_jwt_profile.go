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

	user, jwtReq, err := s.verifyJWTProfile(ctx, r.Data)
	if err != nil {
		return nil, err
	}

	client := &clientCredentialsClient{
		id:   jwtReq.Subject,
		user: user,
	}
	scope, err := op.ValidateAuthReqScopes(client, r.Data.Scope)
	if err != nil {
		return nil, err
	}
	scope, err = s.checkOrgScopes(ctx, client.user, scope)
	if err != nil {
		return nil, err
	}

	session, err := s.command.CreateOIDCSession(ctx,
		user.ID,
		user.ResourceOwner,
		"",
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
	)
	if err != nil {
		return nil, err
	}
	return response(s.accessTokenResponseFromSession(ctx, client, session, "", "", false, true, false, false))
}

func (s *Server) verifyJWTProfile(ctx context.Context, req *oidc.JWTProfileGrantRequest) (user *query.User, tokenRequest *oidc.JWTTokenRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	checkSubject := func(jwt *oidc.JWTTokenRequest) (err error) {
		user, err = s.query.GetUserByID(ctx, true, jwt.Subject)
		return err
	}
	verifier := op.NewJWTProfileVerifier(
		&jwtProfileKeyStorage{query: s.query},
		op.IssuerFromContext(ctx),
		time.Hour, time.Second,
		op.SubjectCheck(checkSubject),
	)
	tokenRequest, err = op.VerifyJWTAssertion(ctx, req.Assertion, verifier)
	if err != nil {
		return nil, nil, err
	}
	return user, tokenRequest, nil
}

type jwtProfileKeyStorage struct {
	query *query.Queries
}

func (s *jwtProfileKeyStorage) GetKeyByIDAndClientID(ctx context.Context, keyID, userID string) (*jose.JSONWebKey, error) {
	publicKeyData, err := s.query.GetAuthNKeyPublicKeyByIDAndIdentifier(ctx, keyID, userID)
	if err != nil {
		return nil, err
	}
	publicKey, err := crypto.BytesToPublicKey(publicKeyData)
	if err != nil {
		return nil, err
	}
	return &jose.JSONWebKey{
		KeyID: keyID,
		Use:   "sig",
		Key:   publicKey,
	}, nil
}
