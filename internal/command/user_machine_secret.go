package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type GenerateMachineSecret struct {
	ClientSecret string
}

func (c *Commands) GenerateMachineSecret(ctx context.Context, userID string, resourceOwner string, generator crypto.Generator, set *GenerateMachineSecret) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareGenerateMachineSecret(agg, generator, set))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareGenerateMachineSecret(a *user.Aggregate, generator crypto.Generator, set *GenerateMachineSecret) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-x0992n", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-bzoqjs", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-x8910n", "Errors.User.NotExisting")
			}

			clientSecret, secretString, err := domain.NewMachineClientSecret(generator)
			if err != nil {
				return nil, err
			}
			set.ClientSecret = secretString

			return []eventstore.Command{
				user.NewMachineSecretSetEvent(ctx, &a.Aggregate, clientSecret),
			}, nil
		}, nil
	}
}

func (c *Commands) RemoveMachineSecret(ctx context.Context, userID string, resourceOwner string) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareRemoveMachineSecret(agg))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareRemoveMachineSecret(a *user.Aggregate) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-0qp2hus", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-bzosjs", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-x7s802", "Errors.User.NotExisting")
			}
			if writeModel.ClientSecret == nil {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-coi82n", "Errors.User.Machine.Secret.NotExisting")
			}
			return []eventstore.Command{
				user.NewMachineSecretRemovedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}

func (c *Commands) VerifyMachineSecret(ctx context.Context, userID string, resourceOwner string, secret string) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareVerifyMachineSecret(agg, secret, c.codeAlg))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	for _, cmd := range cmds {
		if cmd.Type() == user.MachineSecretCheckFailedType {
			logging.OnError(err).Error("could not push event MachineSecretCheckFailed")
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3kjh", "Errors.User.Machine.Secret.Invalid")
		}
	}
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareVerifyMachineSecret(a *user.Aggregate, secret string, algorithm crypto.HashAlgorithm) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-0qp2hus", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-bzosjs", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-569sh2o", "Errors.User.NotExisting")
			}
			if writeModel.ClientSecret == nil {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-x8910n", "Errors.User.Machine.Secret.NotExisting")
			}
			err = crypto.CompareHash(writeModel.ClientSecret, []byte(secret), algorithm)
			if err == nil {
				return []eventstore.Command{
					user.NewMachineSecretCheckSucceededEvent(ctx, &a.Aggregate),
				}, nil
			}
			return []eventstore.Command{
				user.NewMachineSecretCheckFailedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}
