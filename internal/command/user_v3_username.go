package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddUsername struct {
	Details *domain.ObjectDetails

	ResourceOwner string
	UserID        string

	Username      string
	IsOrgSpecific bool
}

func (c *Commands) AddUsername(ctx context.Context, username *AddUsername) (err error) {
	existing, err := existingSchemaUserWithPermission(ctx, c, username.ResourceOwner, username.UserID)
	if err != nil {
		return err
	}
	id, err := c.idGenerator.Next()
	if err != nil {
		return err
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
		return err
	}
	username.Details = pushedEventsToObjectDetails(events)
	return nil
}

func (c *Commands) DeleteUsername(ctx context.Context, resourceOwner, id string) (_ *domain.ObjectDetails, err error) {
	existing, err := c.getSchemaUsernameExists(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if existing.Username == "" {
		return nil, zerrors.ThrowNotFound(nil, "TODO", "TODO")
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

func (c *Commands) getSchemaUsernameExists(ctx context.Context, resourceOwner, id string) (*UsernameV3WriteModel, error) {
	writeModel := NewUsernameV3WriteModel(resourceOwner, id)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func existingSchemaUserWithPermission(ctx context.Context, c *Commands, resourceOwner, userID string) (*UserV3WriteModel, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing")
	}
	existingUser, err := c.getSchemaUserExists(ctx, resourceOwner, userID)
	if err != nil {
		return nil, err
	}
	if !existingUser.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-6T2xrOHxTx", "Errors.User.NotFound")
	}

	_, err = c.getSchemaWriteModelByID(ctx, "", existingUser.SchemaID)
	if err != nil {
		return nil, err
	}
	if !existingUser.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-6T2xrOHxTx", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUser(ctx, existingUser.ResourceOwner, existingUser.AggregateID); err != nil {
		return nil, err
	}
	return existingUser, nil
}
