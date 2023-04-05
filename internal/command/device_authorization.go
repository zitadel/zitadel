package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
)

func (c *Commands) AddDeviceAuth(ctx context.Context, ada *domain.DeviceAuthorization) (_ string, _ *domain.ObjectDetails, err error) {
	// is this required? user input check is already part of the framework
	/*
		if !ada.IsValid() {
			return "", nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-eg2gf", "Errors.Action.Invalid")
		}
	*/

	ada.AggregateID, err = c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	model := &DeviceAuthWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: ada.AggregateID,
		},
	}
	aggr := eventstore.AggregateFromWriteModel(&model.WriteModel, deviceauth.AggregateType, deviceauth.AggregateVersion)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewAddedEvent(
		ctx,
		aggr,
		ada.ClientID,
		ada.DeviceCode,
		ada.UserCode,
		ada.Expires,
		ada.Scopes,
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

func (c *Commands) ApproveDeviceAuth(ctx context.Context, userCode, subject string) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByUsercode(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Sfg2t", "Errors.Action.NotFound")
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

func (c *Commands) DenyDeviceAuth(ctx context.Context, userCode string) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByUsercode(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Sfg2t", "Errors.Action.NotFound")
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

func (c *Commands) RemoveAction(ctx context.Context, clientID, deviceCode, userCode string) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, clientID, deviceCode)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Sfg2t", "Errors.Action.NotFound")
	}
	aggr := eventstore.AggregateFromWriteModel(&model.WriteModel, deviceauth.AggregateType, deviceauth.AggregateVersion)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewRemovedEvent(ctx, aggr, clientID, deviceCode, userCode))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) getDeviceAuthWriteModelByUsercode(ctx context.Context, userCode string) (*DeviceAuthWriteModel, error) {
	model := &DeviceAuthWriteModel{UserCode: userCode}
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (c *Commands) getDeviceAuthWriteModelByDeviceCode(ctx context.Context, clietID, deviceCode string) (*DeviceAuthWriteModel, error) {
	model := &DeviceAuthWriteModel{
		ClientID:   clietID,
		DeviceCode: deviceCode,
	}
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
