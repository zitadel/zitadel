package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddUsername struct {
	ResourceOwner string
	UserID        string

	Username      string
	IsOrgSpecific bool
}

func (c *Commands) AddUsername(ctx context.Context, username *AddUsername) (*domain.ObjectDetails, error) {
	existing, err := existingSchemaUserWithPermission(ctx, c, username.ResourceOwner, username.UserID)
	if err != nil {
		return nil, err
	}
	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx,
		authenticator.NewUsernameCreatedEvent(ctx,
			&authenticator.NewAggregate(id, existing.ResourceOwner).Aggregate,
			existing.AggregateID,
			username.IsOrgSpecific,
			username.Username,
		),
	)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) DeleteUsername(ctx context.Context, resourceOwner, userID, id string) (_ *domain.ObjectDetails, err error) {
	existing, err := c.getSchemaUsernameExistsWithPermission(ctx, resourceOwner, userID, id)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx,
		authenticator.NewUsernameDeletedEvent(ctx,
			&authenticator.NewAggregate(id, existing.ResourceOwner).Aggregate,
			existing.IsOrgSpecific,
			existing.Username,
		),
	)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) getSchemaUsernameExistsWithPermission(ctx context.Context, resourceOwner, userID, id string) (*UsernameV3WriteModel, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-J6ybG5WZiy", "Errors.IDMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing")
	}
	writeModel := NewUsernameV3WriteModel(resourceOwner, userID, id)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if writeModel.Username == "" {
		return nil, zerrors.ThrowNotFound(nil, "TODO", "TODO")
	}

	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.UserID); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func existingSchemaUserWithPermission(ctx context.Context, c *Commands, resourceOwner, userID string) (*UserV3WriteModel, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing")
	}
	existingUser, err := c.getSchemaUserWMForState(ctx, resourceOwner, userID)
	if err != nil {
		return nil, err
	}
	if !existingUser.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-6T2xrOHxTx", "Errors.User.NotFound")
	}

	if err := c.checkPermissionUpdateUser(ctx, existingUser.ResourceOwner, existingUser.AggregateID); err != nil {
		return nil, err
	}

	existingSchema, err := c.getSchemaWriteModelByID(ctx, "", existingUser.SchemaID)
	if err != nil {
		return nil, err
	}
	if !existingSchema.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-6T2xrOHxTx", "TODO")
	}

	// TODO possible authenticators check
	return existingUser, nil
}
