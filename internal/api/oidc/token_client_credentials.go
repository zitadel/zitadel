package oidc

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) ClientCredentialsExchange(ctx context.Context, r *op.ClientRequest[oidc.ClientCredentialsRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		span.EndWithError(err)
		err = oidcError(err)
	}()
	client, ok := r.Client.(*clientCredentialsClient)
	if !ok {
		return nil, zerrors.ThrowInternal(nil, "OIDC-ga0EP", "Error.Internal")
	}
	scope, err := op.ValidateAuthReqScopes(client, r.Data.Scope)
	if err != nil {
		return nil, err
	}

	// TODO org scope sanitation?
	// https://github.com/zitadel/zitadel/blob/150d79af4767d62d7b9f90ab08af4eb26bde0f8b/internal/api/oidc/client.go#L302-L326

	session, err := s.command.CreateOIDCSession(ctx,
		client.user.ID,
		client.user.ResourceOwner,
		r.Data.ClientID,
		scope,
		domain.AddAudScopeToAudience(ctx, nil, r.Data.Scope),
		[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, // TBD: or nil?
		time.Now(),
		"",
		nil,
		domain.TokenReasonClientCredentials,
		nil,
		slices.Contains(scope, oidc.ScopeOfflineAccess),
	)

	return response(s.accessTokenResponseFromSession(ctx, client, session, "", "", false))
}
