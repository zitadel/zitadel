package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
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
	if !writeModel.Exists() || writeModel.Locked {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-G4LOrnjY7q", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUserState(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewLockedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) UnlockSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-krXtYscQZh", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForState(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() || !writeModel.Locked {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-gpBv46Lh9m", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUserState(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewUnlockedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) DeactivateSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-pjJhge86ZV", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForState(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserStateActive {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Ob6lR5iFTe", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUserState(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewDeactivatedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) ActivateSchemaUser(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-17XupGvxBJ", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserWMForState(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserStateInactive {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-rQjbBr4J3j", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUserState(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewActivatedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) checkPermissionUpdateUserState(ctx context.Context, resourceOwner, userID string) error {
	return c.checkPermission(ctx, domain.PermissionUserWrite, resourceOwner, userID)
}

func (c *Commands) getSchemaUserWMForState(ctx context.Context, resourceOwner, id string) (*UserV3WriteModel, error) {
	writeModel := NewExistsUserV3WriteModel(resourceOwner, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}
