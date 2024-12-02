package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) LockSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Eu8I2VAfjF", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForState(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	events, err := writeModel.NewLock(ctx)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) UnlockSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-krXtYscQZh", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForState(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	events, err := writeModel.NewUnlock(ctx)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) DeactivateSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-pjJhge86ZV", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForState(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	events, err := writeModel.NewDeactivate(ctx)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) ActivateSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-17XupGvxBJ", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForState(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	events, err := writeModel.NewActivate(ctx)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) getSchemaUserWMForState(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewExistsUserV3WriteModel(resourceOwner, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func existingSchemaUser(ctx context.Context, c *Commands, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel, err := c.getSchemaUserWMForState(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	return writeModel, writeModel.Exists()
}
