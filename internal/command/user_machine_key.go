package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddMachineKey struct {
	Type           domain.AuthNKeyType
	ExpirationDate time.Time
}

type MachineKey struct {
	models.ObjectRoot

	KeyID          string
	Type           domain.AuthNKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
	PublicKey      []byte
}

func NewMachineKey(resourceOwner string, userID string, expirationDate time.Time, keyType domain.AuthNKeyType) *MachineKey {
	return &MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		ExpirationDate: expirationDate,
		Type:           keyType,
	}
}

func (key *MachineKey) SetPublicKey(publicKey []byte) {
	key.PublicKey = publicKey
}
func (key *MachineKey) SetPrivateKey(privateKey []byte) {
	key.PrivateKey = privateKey
}
func (key *MachineKey) GetExpirationDate() time.Time {
	return key.ExpirationDate
}
func (key *MachineKey) SetExpirationDate(t time.Time) {
	key.ExpirationDate = t
}

func (key *MachineKey) Detail() ([]byte, error) {
	if len(key.PrivateKey) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "KEY-sp2l2m", "Errors.Internal")
	}
	if key.Type == domain.AuthNKeyTypeJSON {
		return domain.MachineKeyMarshalJSON(key.KeyID, key.PrivateKey, key.ExpirationDate, key.AggregateID)
	}
	return nil, zerrors.ThrowPreconditionFailed(nil, "KEY-dsg52", "Errors.Internal")
}

func (key *MachineKey) content() error {
	if key.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-kqpoix", "Errors.ResourceOwnerMissing")
	}
	if key.AggregateID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-xuiwk2", "Errors.User.UserIDMissing")
	}
	if key.KeyID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-0p2m1h", "Errors.IDMissing")
	}
	return nil
}

func (key *MachineKey) valid() (err error) {
	if err := key.content(); err != nil {
		return err
	}
	// If a key is supplied, it should be a valid public key
	if len(key.PublicKey) > 0 {
		if _, err := crypto.BytesToPublicKey(key.PublicKey); err != nil {
			return zerrors.ThrowInvalidArgument(nil, "COMMAND-5F3h1", "Errors.User.Machine.Key.Invalid")
		}
	}
	key.ExpirationDate, err = domain.ValidateExpirationDate(key.ExpirationDate)
	return err
}

func (key *MachineKey) checkAggregate(ctx context.Context, filter preparation.FilterToQueryReducer) error {
	if exists, err := ExistsUser(ctx, filter, key.AggregateID, key.ResourceOwner); err != nil || !exists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-bnipwm1", "Errors.User.NotFound")
	}
	return nil
}

func (c *Commands) AddUserMachineKey(ctx context.Context, machineKey *MachineKey) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if machineKey.KeyID == "" {
		keyID, err := id_generator.Next()
		if err != nil {
			return nil, err
		}
		machineKey.KeyID = keyID
	}

	validation := prepareAddUserMachineKey(machineKey, c.machineKeySize)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareAddUserMachineKey(machineKey *MachineKey, keySize int) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if err := machineKey.valid(); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()

			if err := machineKey.checkAggregate(ctx, filter); err != nil {
				return nil, err
			}
			if len(machineKey.PublicKey) == 0 {
				if err = domain.SetNewAuthNKeyPair(machineKey, keySize); err != nil {
					return nil, err
				}
			}
			writeModel, err := getMachineKeyWriteModelByID(ctx, filter, machineKey.AggregateID, machineKey.KeyID, machineKey.ResourceOwner)
			if err != nil {
				return nil, err
			}
			if writeModel.Exists() {
				return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-091mops", "Errors.User.Machine.Key.AlreadyExists")
			}
			return []eventstore.Command{
				user.NewMachineKeyAddedEvent(
					ctx,
					UserAggregateFromWriteModel(&writeModel.WriteModel),
					machineKey.KeyID,
					machineKey.Type,
					machineKey.ExpirationDate,
					machineKey.PublicKey,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) RemoveUserMachineKey(ctx context.Context, machineKey *MachineKey) (*domain.ObjectDetails, error) {
	validation := prepareRemoveUserMachineKey(machineKey)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareRemoveUserMachineKey(machineKey *MachineKey) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if err := machineKey.content(); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineKeyWriteModelByID(ctx, filter, machineKey.AggregateID, machineKey.KeyID, machineKey.ResourceOwner)
			if err != nil {
				return nil, err
			}
			if !writeModel.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "COMMAND-4m77G", "Errors.User.Machine.Key.NotFound")
			}
			return []eventstore.Command{
				user.NewMachineKeyRemovedEvent(
					ctx,
					UserAggregateFromWriteModel(&writeModel.WriteModel),
					machineKey.KeyID,
				),
			}, nil
		}, nil
	}
}

func getMachineKeyWriteModelByID(ctx context.Context, filter preparation.FilterToQueryReducer, userID, keyID, resourceOwner string) (_ *MachineKeyWriteModel, err error) {
	writeModel := NewMachineKeyWriteModel(userID, keyID, resourceOwner)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
