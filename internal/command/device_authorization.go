package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
)

func (c *Commands) AddDeviceAuth(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes []string) (string, *domain.ObjectDetails, error) {
	aggrID, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	model := &DeviceAuthWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			AggregateID:   aggrID,
		},
	}
	aggr := eventstore.AggregateFromWriteModel(&model.WriteModel, deviceauth.AggregateType, deviceauth.AggregateVersion)

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
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Hief9", "Errors.DeviceAuth.NotFound")
	}
	aggr := eventstore.AggregateFromWriteModel(&model.WriteModel, deviceauth.AggregateType, deviceauth.AggregateVersion)

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

func (c *Commands) DenyDeviceAuth(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-gee5A", "Errors.DeviceAuth.NotFound")
	}
	aggr := eventstore.AggregateFromWriteModel(&model.WriteModel, deviceauth.AggregateType, deviceauth.AggregateVersion)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewDeniedEvent(ctx, aggr))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) RemoveDeviceAuth(ctx context.Context, deviceAuth *domain.DeviceAuth) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByID(ctx, deviceAuth.AggregateID)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-yo9Ie", "Errors.DeviceAuth.NotFound")
	}
	aggr := eventstore.AggregateFromWriteModel(&model.WriteModel, deviceauth.AggregateType, deviceauth.AggregateVersion)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewRemovedEvent(ctx, aggr, deviceAuth.ClientID, deviceAuth.DeviceCode, deviceAuth.UserCode))
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
