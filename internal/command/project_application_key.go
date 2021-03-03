package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/project"
)

func (c *Commands) AddApplicationKey(ctx context.Context, key *domain.ApplicationKey, resourceOwner string) (_ *domain.ApplicationKey, err error) {
	application, err := c.getApplicationWriteModel(ctx, key.AggregateID, key.ApplicationID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !application.State.Exists() {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-sak25", "Errors.Application.NotFound")
	}
	key.KeyID, err = c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	keyWriteModel := NewApplicationKeyWriteModel(key.AggregateID, key.ApplicationID, key.KeyID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, keyWriteModel)
	if err != nil {
		return nil, err
	}

	if !keyWriteModel.KeysAllowed {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-Dff54", "Errors.Project.App.AuthMethodNoPrivateKeyJWT")
	}

	if err := domain.EnsureValidExpirationDate(key); err != nil {
		return nil, err
	}

	err = domain.SetNewAuthNKeyPair(key, c.applicationKeySize)
	if err != nil {
		return nil, err
	}
	key.ClientID = keyWriteModel.ClientID

	pushedEvents, err := c.eventstore.PushEvents(ctx,
		project.NewApplicationKeyAddedEvent(
			ctx,
			ProjectAggregateFromWriteModel(&keyWriteModel.WriteModel),
			key.ApplicationID,
			key.ClientID,
			key.KeyID,
			key.Type,
			key.ExpirationDate,
			key.PublicKey),
	)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(keyWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	result := applicationKeyWriteModelToKey(keyWriteModel, key.PrivateKey)
	return result, nil
}

func (c *Commands) RemoveApplicationKey(ctx context.Context, projectID, applicationID, keyID, resourceOwner string) (*domain.ObjectDetails, error) {
	keyWriteModel := NewApplicationKeyWriteModel(projectID, applicationID, keyID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, keyWriteModel)
	if err != nil {
		return nil, err
	}
	if !keyWriteModel.State.Exists() {
		return nil, errors.ThrowNotFound(nil, "COMMAND-4m77G", "Errors.Application.Key.NotFound")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewApplicationKeyRemovedEvent(ctx, ProjectAggregateFromWriteModel(&keyWriteModel.WriteModel), keyID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(keyWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&keyWriteModel.WriteModel), nil
}
