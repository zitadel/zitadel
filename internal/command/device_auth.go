package command

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDeviceAuth(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes, audience []string) (*domain.ObjectDetails, error) {
	aggr := deviceauth.NewAggregate(deviceCode, authz.GetInstance(ctx).InstanceID())
	model := NewDeviceAuthWriteModel(deviceCode, aggr.ResourceOwner)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewAddedEvent(
		ctx,
		aggr,
		clientID,
		deviceCode,
		userCode,
		expires,
		scopes,
		audience,
	))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) ApproveDeviceAuth(ctx context.Context, deviceCode, userID, userOrgID string, authMethods []domain.UserAuthMethodType, authTime time.Time) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, deviceCode)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Hief9", "Errors.DeviceAuth.NotFound")
	}
	aggr := deviceauth.NewAggregate(model.AggregateID, model.InstanceID)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewApprovedEvent(ctx, aggr, userID, userOrgID, authMethods, authTime))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) CancelDeviceAuth(ctx context.Context, id string, reason domain.DeviceAuthCanceled) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, id)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-gee5A", "Errors.DeviceAuth.NotFound")
	}
	aggr := deviceauth.NewAggregate(model.AggregateID, model.InstanceID)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewCanceledEvent(ctx, aggr, reason))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) getDeviceAuthWriteModelByDeviceCode(ctx context.Context, deviceCode string) (*DeviceAuthWriteModel, error) {
	model := &DeviceAuthWriteModel{WriteModel: eventstore.WriteModel{AggregateID: deviceCode}}
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

type DeviceAuthStateError domain.DeviceAuthState

func (s DeviceAuthStateError) Error() string {
	return domain.DeviceAuthState(s).String()
}

func (c *Commands) CreateOIDCSessionFromDeviceAuth(ctx context.Context, deviceCode string) (_ *OIDCSession, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	deviceAuthModel, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, deviceCode)
	if err != nil {
		return nil, err
	}
	if deviceAuthModel.State != domain.DeviceAuthStateApproved {
		if deviceAuthModel.Expires.Before(time.Now()) {
			c.asyncPush(ctx, deviceauth.NewCanceledEvent(ctx, deviceAuthModel.aggregate, domain.DeviceAuthCanceledExpired))
			return nil, DeviceAuthStateError(domain.DeviceAuthStateExpired)
		}
		return nil, DeviceAuthStateError(deviceAuthModel.State)
	}

	cmd, err := c.newOIDCSessionAddEvents(ctx, deviceAuthModel.UserOrgID)
	if err != nil {
		return nil, err
	}

	cmd.AddSession(ctx,
		deviceAuthModel.UserID,
		"",
		deviceAuthModel.ClientID,
		deviceAuthModel.Audience,
		deviceAuthModel.Scopes,
		deviceAuthModel.UserAuthMethods,
		deviceAuthModel.AuthTime,
		nil, // TBD: should we use some kind of device fingerprint as useragent?
	)
	if err = cmd.AddAccessToken(ctx, deviceAuthModel.Scopes, domain.TokenReasonAuthRequest, nil); err != nil {
		return nil, err
	}

	if slices.Contains(deviceAuthModel.Scopes, "offline_access") {
		if err = cmd.AddRefreshToken(ctx, deviceAuthModel.UserID); err != nil {
			return nil, err
		}
	}
	cmd.DeviceAuthRequestDone(ctx, deviceAuthModel.aggregate)
	return cmd.PushEvents(ctx)
}

func (cmd *OIDCSessionEvents) DeviceAuthRequestDone(ctx context.Context, deviceAuthAggregate *eventstore.Aggregate) {
	cmd.events = append(cmd.events, deviceauth.NewDoneEvent(ctx, deviceAuthAggregate))
}
