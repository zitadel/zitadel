package oidc

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/op"

	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	DeviceAuthDefaultLifetime     = 5 * time.Minute
	DeviceAuthDefaultPollInterval = 5 * time.Second
)

type DeviceAuthorizationConfig struct {
	Lifetime     time.Duration
	PollInterval time.Duration
	UserCode     *UserCodeConfig
}

type UserCodeConfig struct {
	CharSet      string
	CharAmount   int
	DashInterval int
}

// toOPConfig converts DeviceAuthorizationConfig to a op.DeviceAuthorizationConfig,
// setting sane defaults for empty values.
// Safe to call when *DeviceAuthorizationConfig is nil.
func (c *DeviceAuthorizationConfig) toOPConfig() op.DeviceAuthorizationConfig {
	out := op.DeviceAuthorizationConfig{
		Lifetime:     DeviceAuthDefaultLifetime,
		PollInterval: DeviceAuthDefaultPollInterval,
		UserFormPath: login.EndpointDeviceAuth,
		UserCode:     op.UserCodeBase20,
	}
	if c == nil {
		return out
	}
	if c.Lifetime != 0 {
		out.Lifetime = c.Lifetime
	}
	if c.PollInterval != 0 {
		out.PollInterval = c.PollInterval
	}

	if c.UserCode == nil {
		return out
	}
	if c.UserCode.CharSet != "" {
		out.UserCode.CharSet = c.UserCode.CharSet
	}
	if c.UserCode.CharAmount != 0 {
		out.UserCode.CharAmount = c.UserCode.CharAmount
	}
	if c.UserCode.DashInterval != 0 {
		out.UserCode.DashInterval = c.UserCode.CharAmount
	}
	return out
}

func (o *OPStorage) StoreDeviceAuthorization(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes []string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	scopes, err = o.assertProjectRoleScopes(ctx, clientID, scopes)
	if err != nil {
		return errors.ThrowPreconditionFailed(err, "OIDC-She4t", "Errors.Internal")
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
		logging.WithError(err).WithFields(logrus.Fields{"client_id": clientID, "device_code": deviceCode}).Error()
		return nil, err
	}

	/*
		if deviceAuth.State != domain.DeviceAuthStateInitiated || deviceAuth.Expires.Before(time.Now()) {
			_, err = o.command.RemoveDeviceAuth(ctx, deviceAuth)
		}
	*/

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
