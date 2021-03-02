package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddUserMachineKey(ctx context.Context, machineKey *domain.MachineKey, resourceOwner string) (*domain.MachineKey, error) {
	err := c.checkUserExists(ctx, machineKey.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	keyWriteModel := NewMachineKeyWriteModel(machineKey.AggregateID, keyID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, keyWriteModel)
	if err != nil {
		return nil, err
	}

	if err = domain.EnsureValidExpirationDate(machineKey); err != nil {
		return nil, err
	}

	if err = domain.SetNewAuthNKeyPair(machineKey, c.machineKeySize); err != nil {
		return nil, err
	}

	events, err := c.eventstore.PushEvents(ctx,
		user.NewMachineKeyAddedEvent(
			ctx,
			UserAggregateFromWriteModel(&keyWriteModel.WriteModel),
			keyID,
			machineKey.Type,
			machineKey.ExpirationDate,
			machineKey.PublicKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(keyWriteModel, events...)
	if err != nil {
		return nil, err
	}

	key := keyWriteModelToMachineKey(keyWriteModel)
	key.PrivateKey = machineKey.PrivateKey
	return key, nil
}

func (c *Commands) RemoveUserMachineKey(ctx context.Context, userID, keyID, resourceOwner string) (*domain.ObjectDetails, error) {
	keyWriteModel, err := c.machineKeyWriteModelByID(ctx, userID, keyID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !keyWriteModel.Exists() {
		return nil, errors.ThrowNotFound(nil, "COMMAND-4m77G", "Errors.User.Machine.Key.NotFound")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx,
		user.NewMachineKeyRemovedEvent(ctx, UserAggregateFromWriteModel(&keyWriteModel.WriteModel), keyID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(keyWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&keyWriteModel.WriteModel), nil
}

func (c *Commands) machineKeyWriteModelByID(ctx context.Context, userID, keyID, resourceOwner string) (writeModel *MachineKeyWriteModel, err error) {
	if userID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-4n8vs", "Errors.User.UserIDMissing")
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewMachineKeyWriteModel(userID, keyID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
