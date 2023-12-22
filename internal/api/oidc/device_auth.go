package oidc

import (
	"context"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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

// toOPConfig converts DeviceAuthorizationConfig to a [op.DeviceAuthorizationConfig],
// setting sane defaults for empty values.
// Safe to call when c is nil.
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

// StoreDeviceAuthorization creates a new Device Authorization request.
// Implements the op.DeviceAuthorizationStorage interface.
func (o *OPStorage) StoreDeviceAuthorization(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes []string) (err error) {
	const logMsg = "store device authorization"
	logger := logging.WithFields("client_id", clientID, "device_code", deviceCode, "user_code", userCode, "expires", expires, "scopes", scopes)

	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		logger.OnError(err).Error(logMsg)
		span.EndWithError(err)
	}()

	scopes, err = o.assertProjectRoleScopes(ctx, clientID, scopes)
	if err != nil {
		return zerrors.ThrowPreconditionFailed(err, "OIDC-She4t", "Errors.Internal")
	}
	details, err := o.command.AddDeviceAuth(ctx, clientID, deviceCode, userCode, expires, scopes)
	if err == nil {
		logger.SetFields("details", details).Debug(logMsg)
	}

	return err
}

func newDeviceAuthorizationState(d *query.DeviceAuth) *op.DeviceAuthorizationState {
	return &op.DeviceAuthorizationState{
		ClientID: d.ClientID,
		Scopes:   d.Scopes,
		Expires:  d.Expires,
		Done:     d.State.Done(),
		Denied:   d.State.Denied(),
		Subject:  d.Subject,
		AMR:      AuthMethodTypesToAMR(d.UserAuthMethods),
		AuthTime: d.AuthTime,
	}
}

// GetDeviceAuthorizatonState retrieves the current state of the Device Authorization process.
// It implements the [op.DeviceAuthorizationStorage] interface and is used by devices that
// are polling until they successfully receive a token or we indicate a denied or expired state.
// As generated user codes are of low entropy, this implementation also takes care or
// device authorization request cleanup, when it has been Approved, Denied or Expired.
func (o *OPStorage) GetDeviceAuthorizatonState(ctx context.Context, clientID, deviceCode string) (state *op.DeviceAuthorizationState, err error) {
	const logMsg = "get device authorization state"
	logger := logging.WithFields("device_code", deviceCode)

	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		if err != nil {
			logger.WithError(err).Error(logMsg)
		}
		span.EndWithError(err)
	}()

	deviceAuth, err := o.query.DeviceAuthByDeviceCode(ctx, deviceCode)
	if err != nil {
		return nil, err
	}
	logger.SetFields(
		"expires", deviceAuth.Expires, "scopes", deviceAuth.Scopes,
		"subject", deviceAuth.Subject, "state", deviceAuth.State,
	).Debug("device authorization state")

	// Cancel the request if it is expired, only if it wasn't Done meanwhile
	if !deviceAuth.State.Done() && deviceAuth.Expires.Before(time.Now()) {
		_, err = o.command.CancelDeviceAuth(ctx, deviceAuth.DeviceCode, domain.DeviceAuthCanceledExpired)
		if err != nil {
			return nil, err
		}
		deviceAuth.State = domain.DeviceAuthStateExpired
	}

	return newDeviceAuthorizationState(deviceAuth), nil
}
