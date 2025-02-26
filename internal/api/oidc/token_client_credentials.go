package oidc

import (
	"context"
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
		[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
		time.Now(),
		"",
		nil,
		nil,
		domain.TokenReasonClientCredentials,
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
