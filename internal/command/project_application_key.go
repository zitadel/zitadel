package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddApplicationKeyWithID(ctx context.Context, key *domain.ApplicationKey, resourceOwner string) (_ *domain.ApplicationKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel, err := c.applicationKeyWriteModelByID(ctx, key.AggregateID, key.ApplicationID, key.KeyID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if writeModel.State != domain.AppStateUnspecified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-so20alo", "Errors.Project.App.Key.AlreadyExisting")
	}
	application, err := c.getApplicationWriteModel(ctx, key.AggregateID, key.ApplicationID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !application.State.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-sak24", "Errors.Project.App.NotFound")
	}
	return c.addApplicationKey(ctx, key, resourceOwner)
}

func (c *Commands) AddApplicationKey(ctx context.Context, key *domain.ApplicationKey, resourceOwner string) (_ *domain.ApplicationKey, err error) {
	if key.AggregateID == "" || key.ApplicationID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-55m9fs", "Errors.IDMissing")
	}
	application, err := c.getApplicationWriteModel(ctx, key.AggregateID, key.ApplicationID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !application.State.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-sak25", "Errors.Project.App.NotFound")
	}
	key.KeyID, err = id_generator.Next()
	if err != nil {
		return nil, err
	}

	return c.addApplicationKey(ctx, key, resourceOwner)
}

func (c *Commands) addApplicationKey(ctx context.Context, key *domain.ApplicationKey, resourceOwner string) (_ *domain.ApplicationKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	keyWriteModel := NewApplicationKeyWriteModel(key.AggregateID, key.ApplicationID, key.KeyID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, keyWriteModel)
	if err != nil {
		return nil, err
	}

	if !keyWriteModel.KeysAllowed {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Dff54", "Errors.Project.App.AuthMethodNoPrivateKeyJWT")
	}

	if err := domain.EnsureValidExpirationDate(key); err != nil {
		return nil, err
	}

	if len(key.PublicKey) == 0 {
		err = domain.SetNewAuthNKeyPair(key, c.applicationKeySize)
		if err != nil {
			return nil, err
		}
		key.ClientID = keyWriteModel.ClientID
	}

	pushedEvents, err := c.eventstore.Push(ctx,
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
	result := applicationKeyWriteModelToKey(keyWriteModel)
	if len(key.PrivateKey) > 0 {
		result.PrivateKey = key.PrivateKey
	}
	return result, nil
}

func (c *Commands) RemoveApplicationKey(ctx context.Context, projectID, applicationID, keyID, resourceOwner string) (*domain.ObjectDetails, error) {
	keyWriteModel := NewApplicationKeyWriteModel(projectID, applicationID, keyID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, keyWriteModel)
	if err != nil {
		return nil, err
	}
	if !keyWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4m77G", "Errors.Project.App.Key.NotFound")
	}

	pushedEvents, err := c.eventstore.Push(ctx, project.NewApplicationKeyRemovedEvent(ctx, ProjectAggregateFromWriteModel(&keyWriteModel.WriteModel), keyID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(keyWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&keyWriteModel.WriteModel), nil
}

func (c *Commands) applicationKeyWriteModelByID(ctx context.Context, projectID, appID, keyID, resourceOwner string) (writeModel *ApplicationKeyWriteModel, err error) {
	if appID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-029sn", "Errors.Project.App.NotFound")
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewApplicationKeyWriteModel(projectID, appID, keyID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
