package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddUsername struct {
	ResourceOwner string
	UserID        string

	Username      string
	IsOrgSpecific bool
}

func (c *Commands) AddUsername(ctx context.Context, username *AddUsername) (*domain.ObjectDetails, error) {
	if username.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing")
	}
	schemauser, err := existingSchemaUser(ctx, c, username.ResourceOwner, username.UserID)
	if err != nil {
		return nil, err
	}

	_, err = existingSchema(ctx, c, "", schemauser.SchemaID)
	if err != nil {
		return nil, err
	}
	// TODO check for possible authenticators

	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	writeModel, err := c.getSchemaUsernameWM(ctx, schemauser.ResourceOwner, schemauser.AggregateID, id)
	if err != nil {
		return nil, err
	}
	events, err := writeModel.NewCreate(ctx, username.IsOrgSpecific, username.Username)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) DeleteUsername(ctx context.Context, resourceOwner, userID, id string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-J6ybG5WZiy", "Errors.IDMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing")
	}

	writeModel, err := c.getSchemaUsernameWM(ctx, resourceOwner, userID, id)
	if err != nil {
		return nil, err
	}
	events, err := writeModel.NewDelete(ctx)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) getSchemaUsernameWM(ctx context.Context, resourceOwner, userID, id string) (*UsernameV3WriteModel, error) {
	writeModel := NewUsernameV3WriteModel(resourceOwner, userID, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}
