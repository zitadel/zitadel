package oidc

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"github.com/zitadel/zitadel/internal/command"
	errz "github.com/zitadel/zitadel/internal/errors"
)

func (s *Server) Introspect(ctx context.Context, r *op.Request[op.IntrospectionRequest]) (_ *op.Response, err error) {
	clientID, err := s.authenticateResourceClient(ctx, r.Data.ClientCredentials)
	if err != nil {
		return nil, err
	}
	response := new(oidc.IntrospectionResponse)
	tokenID, subject, err := s.getTokenIDAndSubject(ctx, r.Data.Token)
	if err != nil {
		// TODO: log error
		return op.NewResponse(response), nil
	}
	if strings.HasPrefix(tokenID, command.IDPrefixV2) {
		err = s.introspect(ctx, response, tokenID, subject, clientID)
		return op.NewResponse(response), nil
	}

	err = s.storage.SetIntrospectionFromToken(ctx, response, tokenID, subject, clientID)
	if err != nil {
		return op.NewResponse(response), nil
	}
	response.Active = true
	return op.NewResponse(response), nil
}

func (s *Server) authenticateResourceClient(ctx context.Context, cc *op.ClientCredentials) (clientID string, err error) {
	if cc.ClientAssertion != "" {
		verifier := op.NewJWTProfileVerifier(s.storage, op.IssuerFromContext(ctx), 1*time.Hour, time.Second)
		profile, err := op.VerifyJWTAssertion(ctx, cc.ClientAssertion, verifier)
		if err != nil {
			return "", err
		}
		return profile.Issuer, nil
	}

	if err = s.storage.AuthorizeClientIDSecret(ctx, cc.ClientID, cc.ClientSecret); err != nil {
		if err != nil {
			return "", err
		}
	}
	return cc.ClientID, nil
}

func (s *Server) getTokenIDAndSubject(ctx context.Context, accessToken string) (idToken, subject string, err error) {
	provider := s.Provider()
	tokenIDSubject, err := provider.Crypto().Decrypt(accessToken)
	if err == nil {
		splitToken := strings.Split(tokenIDSubject, ":")
		if len(splitToken) != 2 {
			return "", "", errors.New("invalid token format")
		}
		return splitToken[0], splitToken[1], nil
	}

	verifier := op.NewAccessTokenVerifier(op.IssuerFromContext(ctx), s.storage.keySet)
	accessTokenClaims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](ctx, accessToken, verifier)
	if err != nil {
		return "", "", err
	}
	return accessTokenClaims.JWTID, accessTokenClaims.Subject, nil
}

func (s *Server) introspect(ctx context.Context, introspection *oidc.IntrospectionResponse, tokenID, subject, clientID string) (err error) {
	// TODO: give clients their own aggregate, so we can skip this query
	projectID, err := s.storage.query.ProjectIDFromClientID(ctx, clientID, false)
	if err != nil {
		return errz.ThrowPermissionDenied(nil, "OIDC-Adfg5", "client not found")
	}

	token, err := s.storage.query.ActiveAccessTokenByToken(ctx, tokenID)
	if err != nil {
		return err
	}
	if !slices.ContainsFunc(token.Audience, func(aud string) bool {
		return aud == token.ClientID || aud == projectID
	}) {
		return errz.ThrowPermissionDenied(nil, "OIDC-sdg3G", "token is not valid for this client")
	}

	userInfo, err := s.storage.query.GetOIDCUserinfo(ctx, token.UserID, token.Scope, []string{projectID})
	if err != nil {
		return err
	}
	introspection.SetUserInfo(userinfoToOIDC(userInfo, token.Scope))
	introspection.Scope = token.Scope
	introspection.ClientID = token.ClientID
	introspection.TokenType = oidc.BearerToken
	introspection.Expiration = oidc.FromTime(token.AccessTokenExpiration)
	introspection.IssuedAt = oidc.FromTime(token.AccessTokenCreation)
	introspection.NotBefore = oidc.FromTime(token.AccessTokenCreation)
	introspection.Audience = token.Audience
	introspection.Issuer = op.IssuerFromContext(ctx)
	introspection.JWTID = tokenID

	return nil
}
