package oidc

import (
	"context"
	"time"

	"github.com/zitadel/oidc/v2/pkg/op"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (o *OPStorage) StoreDeviceAuthorization(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes []string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	scopes, err = o.assertProjectRoleScopes(ctx, clientID, scopes)
	if err != nil {
		return errors.ThrowPreconditionFailed(err, "OIDC-< TODO: code>", "Errors.Internal")
	}
	_, _, err = o.command.AddDeviceAuth(ctx, clientID, deviceCode, userCode, expires, scopes)
	return err
}

func createDeviceAuthorizationState(d *domain.DeviceAuth) *op.DeviceAuthorizationState {
	return &op.DeviceAuthorizationState{
		ClientID: d.ClientID,
		Scopes:   d.Scopes,
		Expires:  d.Expires,
		Done:     d.State.Done(),
		Subject:  d.Subject,
		Denied:   d.State.Denied(),
	}
}

func (o *OPStorage) GetDeviceAuthorizatonState(ctx context.Context, clientID, deviceCode string) (state *op.DeviceAuthorizationState, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	deviceAuth, err := o.query.DeviceAuthByDeviceCode(ctx, clientID, deviceCode)
	if err != nil {
		return nil, err
	}
	if deviceAuth.State != domain.DeviceAuthStateInitiated || deviceAuth.Expires.Before(time.Now()) {
		_, err = o.command.RemoveDeviceAuth(ctx, deviceAuth)
	}
	return createDeviceAuthorizationState(deviceAuth), nil
}

// This is actually not used, as the current implementation operates on the storage directly from the handlers.
func (o *OPStorage) GetDeviceAuthorizationByUserCode(ctx context.Context, userCode string) (state *op.DeviceAuthorizationState, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	deviceAuth, err := o.query.DeviceAuthByUserCode(ctx, userCode)
	if err != nil {
		return nil, err
	}

	return createDeviceAuthorizationState(deviceAuth), err
}

// This is actually not used, as the current implementation operates on the storage directly from the handlers.
func (o *OPStorage) CompleteDeviceAuthorization(ctx context.Context, userCode, subject string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, err = o.command.ApproveDeviceAuth(ctx, userCode, subject)
	return err
}

// This is actually not used, as the current implementation operates on the storage directly from the handlers.
func (o *OPStorage) DenyDeviceAuthorization(ctx context.Context, userCode string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, err = o.command.DenyDeviceAuth(ctx, userCode)
	return err
}
