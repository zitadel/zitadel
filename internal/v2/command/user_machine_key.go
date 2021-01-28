package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
	"time"
)

const (
	yearLayout            = "2006-01-02"
	defaultExpirationDate = "9999-01-01"
)

func (r *CommandSide) AddUserMachineKey(ctx context.Context, machineKey *domain.MachineKey, resourceOwner string) (*domain.MachineKey, error) {
	err := r.checkUserExists(ctx, machineKey.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	keyID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	keyWriteModel := NewMachineKeyWriteModel(machineKey.AggregateID, keyID, resourceOwner)
	userAgg := UserAggregateFromWriteModel(&keyWriteModel.WriteModel)
	err = r.eventstore.FilterToQueryReducer(ctx, keyWriteModel)
	if err != nil {
		return nil, err
	}

	if machineKey.ExpirationDate.IsZero() {
		machineKey.ExpirationDate, err = time.Parse(yearLayout, defaultExpirationDate)
		if err != nil {
			logging.Log("COMMAND9-v8jMi").WithError(err).Warn("unable to set default date")
			return nil, errors.ThrowInternal(err, "COMMAND-38jfus", "Errors.Internal")
		}
	}
	if machineKey.ExpirationDate.Before(time.Now()) {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-38vns", "Errors.MachineKey.ExpireBeforeNow")
	}

	machineKey.GenerateNewMachineKeyPair(r.machineKeySize)

	userAgg.PushEvents(user.NewMachineKeyAddedEvent(ctx, keyID, machineKey.Type, machineKey.ExpirationDate, machineKey.PublicKey))
	err = r.eventstore.PushAggregate(ctx, keyWriteModel, userAgg)
	if err != nil {
		return nil, err
	}
	key := keyWriteModelToMachineKey(keyWriteModel)
	key.PrivateKey = machineKey.PrivateKey
	return key, nil
}

func (r *CommandSide) RemoveUserMachineKey(ctx context.Context, userID, keyID, resourceOwner string) error {
	keyWriteModel, err := r.machineKeyWriteModelByID(ctx, userID, keyID, resourceOwner)
	if err != nil {
		return err
	}
	if keyWriteModel.State == domain.MachineKeyStateUnspecified || keyWriteModel.State == domain.MachineKeyStateRemoved {
		return errors.ThrowNotFound(nil, "COMMAND-4m77G", "Errors.User.Machine.Key.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&keyWriteModel.WriteModel)
	userAgg.PushEvents(user.NewMachineKeyRemovedEvent(ctx, keyID))
	return r.eventstore.PushAggregate(ctx, keyWriteModel, userAgg)
}

func (r *CommandSide) machineKeyWriteModelByID(ctx context.Context, userID, keyID, resourceOwner string) (writeModel *MachineKeyWriteModel, err error) {
	if userID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-4n8vs", "Errors.User.UserIDMissing")
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewMachineKeyWriteModel(userID, keyID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
