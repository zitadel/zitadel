package oidc

import (
	"context"
	"slices"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
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

	session, err := s.command.ExchangeOIDCSessionRefreshAndAccessToken(ctx, r.Data.RefreshToken, r.Data.Scopes, validateRefreshTokenScopes(r.Data.Scopes))
	if err != nil {
		return nil, err
	}
	return response(s.accessTokenResponseFromSession(ctx, client, session, "", client.client.ProjectID, client.client.ProjectRoleAssertion))
}

// validateRefreshTokenScopes validates that the requested scope is a subset of the original auth request scope.
func validateRefreshTokenScopes(requestedScope []string) command.RefreshTokenComplianceChecker {
	return func(_ context.Context, model *command.OIDCSessionWriteModel) error {
		if len(requestedScope) == 0 {
			return nil
		}
		for _, s := range requestedScope {
			if !slices.Contains(model.Scope, s) {
				return oidc.ErrInvalidScope()
			}
		}
		return nil
	}
}
