package oidc

import (
	"context"
	"errors"
	"slices"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) RefreshToken(ctx context.Context, r *op.ClientRequest[oidc.RefreshTokenRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		span.EndWithError(err)
		err = oidcError(err)
	}()

	client, ok := r.Client.(*Client)
	if !ok {
		return nil, zerrors.ThrowInternal(nil, "OIDC-ga0EP", "Error.Internal")
	}

	session, err := s.command.ExchangeOIDCSessionRefreshAndAccessToken(ctx, r.Data.RefreshToken, r.Data.Scopes, refreshTokenComplianceChecker())
	if err == nil {
		return response(s.accessTokenResponseFromSession(ctx, client, session, "", client.client.ProjectID, client.client.ProjectRoleAssertion, client.client.AccessTokenRoleAssertion, client.client.IDTokenRoleAssertion, client.client.IDTokenUserinfoAssertion))
	} else if errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "OIDCS-JOI23", "Errors.OIDCSession.RefreshTokenInvalid")) {
		// We try again for v1 tokens when we encountered specific parsing error
		return s.refreshTokenV1(ctx, client, r)
	}
	return nil, err
}

// refreshTokenV1 verifies a v1 refresh token.
// When valid a v2 OIDC session is created and v2 tokens are returned.
// This "upgrades" existing v1 sessions to v2 session without requiring users to re-login.
//
// This function can be removed when we retire the v1 token repo.
func (s *Server) refreshTokenV1(ctx context.Context, client *Client, r *op.ClientRequest[oidc.RefreshTokenRequest]) (_ *op.Response, err error) {
	refreshToken, err := s.repo.RefreshTokenByToken(ctx, r.Data.RefreshToken)
	if err != nil {
		return nil, err
	}
	scope, err := validateRefreshTokenScopes(refreshToken.Scopes, r.Data.Scopes)
	if err != nil {
		return nil, err
	}
	session, err := s.command.CreateOIDCSession(ctx,
		refreshToken.UserID,
		refreshToken.ResourceOwner,
		refreshToken.ClientID,
		"", // backChannelLogoutURI is not in refresh token view
		scope,
		refreshToken.Audience,
		AMRToAuthMethodTypes(refreshToken.AuthMethodsReferences),
		refreshToken.AuthTime,
		"",
		nil, // Preferred language not in refresh token view
		&domain.UserAgent{
			FingerprintID: &refreshToken.UserAgentID,
			Description:   &refreshToken.UserAgentID,
		},
		domain.TokenReasonRefresh,
		refreshToken.Actor,
		true,
		"",
		domain.OIDCResponseTypeUnspecified,
	)
	if err != nil {
		return nil, err
	}

	// make sure the v1 refresh token can't be reused.
	_, err = s.command.RevokeRefreshToken(ctx, refreshToken.UserID, refreshToken.ResourceOwner, refreshToken.ID)
	if err != nil {
		return nil, err
	}

	return response(s.accessTokenResponseFromSession(ctx, client, session, "", client.client.ProjectID, client.client.ProjectRoleAssertion, client.client.AccessTokenRoleAssertion, client.client.IDTokenRoleAssertion, client.client.IDTokenUserinfoAssertion))
}

// refreshTokenComplianceChecker validates that the requested scope is a subset of the original auth request scope.
func refreshTokenComplianceChecker() command.RefreshTokenComplianceChecker {
	return func(_ context.Context, model *command.OIDCSessionWriteModel, requestedScope []string) ([]string, error) {
		return validateRefreshTokenScopes(model.Scope, requestedScope)
	}
}

func validateRefreshTokenScopes(currentScope, requestedScope []string) ([]string, error) {
	if len(requestedScope) == 0 {
		return currentScope, nil
	}
	for _, s := range requestedScope {
		if !slices.Contains(currentScope, s) {
			return nil, oidc.ErrInvalidScope()
		}
	}
	return requestedScope, nil
}
