package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddUsername struct {
	ResourceOwner string
	UserID        string

	Username *Username
}

type Username struct {
	Username      string
	IsOrgSpecific bool
}

func (c *Commands) AddUsername(ctx context.Context, add *AddUsername) (*domain.ObjectDetails, error) {
	if add.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing")
	}
	schemauser, err := existingSchemaUser(ctx, c, add.ResourceOwner, add.UserID)
	if err != nil {
		return nil, err
	}
	add.ResourceOwner = schemauser.ResourceOwner

	_, err = existingSchema(ctx, c, "", schemauser.SchemaID)
	if err != nil {
		return nil, err
	}
	// TODO check for possible authenticators

	writeModel, events, err := c.addUsername(ctx, add.ResourceOwner, add.UserID, add.Username)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) addUsername(ctx context.Context, resourceOwner, userID string, add *Username) (*UsernameV3WriteModel, []eventstore.Command, error) {
	if resourceOwner == "" || userID == "" || add == nil {
		return nil, nil, nil
	}
	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	writeModel, err := c.getSchemaUsernameWM(ctx, resourceOwner, userID, id)
	if err != nil {
		return nil, nil, err
	}
	events, err := writeModel.NewCreate(ctx, add.IsOrgSpecific, add.Username)
	if err != nil {
		return nil, nil, err
	}
	return writeModel, events, nil
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
