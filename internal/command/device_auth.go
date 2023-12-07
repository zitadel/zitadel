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

func (c *Commands) AddDeviceAuth(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes []string) (string, *domain.ObjectDetails, error) {
	aggrID, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	aggr := deviceauth.NewAggregate(aggrID, authz.GetInstance(ctx).InstanceID())
	model := NewDeviceAuthWriteModel(aggrID, aggr.ResourceOwner)

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
		return "", nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return "", nil, err
	}

	return model.AggregateID, writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) ApproveDeviceAuth(ctx context.Context, id, subject string) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Hief9", "Errors.DeviceAuth.NotFound")
	}
	aggr := deviceauth.NewAggregate(model.AggregateID, model.InstanceID)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewApprovedEvent(ctx, aggr, subject))
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
	model, err := c.getDeviceAuthWriteModelByID(ctx, id)
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

func (c *Commands) RemoveDeviceAuth(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	aggr := deviceauth.NewAggregate(model.AggregateID, model.InstanceID)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewRemovedEvent(ctx, aggr, model.ClientID, model.DeviceCode, model.UserCode))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) getDeviceAuthWriteModelByID(ctx context.Context, id string) (*DeviceAuthWriteModel, error) {
	model := &DeviceAuthWriteModel{WriteModel: eventstore.WriteModel{AggregateID: id}}
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
