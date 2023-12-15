package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDeviceAuth(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes []string) (*domain.ObjectDetails, error) {
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

func (c *Commands) ApproveDeviceAuth(ctx context.Context, deviceCode, subject string, authMethods []domain.UserAuthMethodType, authTime time.Time) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, deviceCode)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Hief9", "Errors.DeviceAuth.NotFound")
	}
	aggr := deviceauth.NewAggregate(model.AggregateID, model.InstanceID)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewApprovedEvent(ctx, aggr, subject, authMethods, authTime))
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
