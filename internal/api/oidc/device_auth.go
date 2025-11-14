package oidc

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/ui/login"
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
		out.UserCode.DashInterval = c.UserCode.DashInterval
	}
	return out
}

// StoreDeviceAuthorization creates a new Device Authorization request.
// Implements the op.DeviceAuthorizationStorage interface.
func (o *OPStorage) StoreDeviceAuthorization(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scope []string) (err error) {
	const logMsg = "store device authorization"
	logger := logging.WithFields("client_id", clientID, "device_code", deviceCode, "user_code", userCode, "expires", expires, "scope", scope)

	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		logger.OnError(err).Error(logMsg)
		span.EndWithError(err)
	}()
	scope, audience, err := o.createAuthRequestScopeAndAudience(ctx, clientID, scope)
	if err != nil {
		return err
	}
	details, err := o.command.AddDeviceAuth(ctx, clientID, deviceCode, userCode, expires, scope, audience, slices.Contains(scope, oidc.ScopeOfflineAccess))
	if err == nil {
		logger.SetFields("details", details).Debug(logMsg)
	}

	return err
}

func (o *OPStorage) GetDeviceAuthorizatonState(context.Context, string, string) (*op.DeviceAuthorizationState, error) {
	panic(o.panicErr("GetDeviceAuthorizatonState"))
}
