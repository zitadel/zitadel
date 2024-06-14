package oidc

import (
	"context"
	"errors"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) DeviceToken(ctx context.Context, r *op.ClientRequest[oidc.DeviceAccessTokenRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		span.EndWithError(err)
		err = oidcError(err)
	}()

	client, ok := r.Client.(*Client)
	if !ok {
		return nil, zerrors.ThrowInternal(nil, "OIDC-Ae2ph", "Error.Internal")
	}
	session, err := s.command.CreateOIDCSessionFromDeviceAuth(ctx, r.Data.DeviceCode)
	if err == nil {
		return response(s.accessTokenResponseFromSession(ctx, client, session, "", client.client.ProjectID, client.client.ProjectRoleAssertion, client.client.AccessTokenRoleAssertion, client.client.IDTokenRoleAssertion, client.client.IDTokenUserinfoAssertion))
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, oidc.ErrSlowDown().WithParent(err)
	}

	var target command.DeviceAuthStateError
	if errors.As(err, &target) {
		state := domain.DeviceAuthState(target)
		if state == domain.DeviceAuthStateInitiated {
			return nil, oidc.ErrAuthorizationPending()
		}
		if state == domain.DeviceAuthStateExpired {
			return nil, oidc.ErrExpiredDeviceCode()
		}
	}
	return nil, oidc.ErrAccessDenied().WithParent(err)
}
