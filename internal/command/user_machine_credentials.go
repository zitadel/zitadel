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

type SetMachineCredentials struct {
	ClientID     string
	ClientSecret string
}

func (c *Commands) SetMachineCredentials(ctx context.Context, userID string, resourceOwner string, generator crypto.Generator, set *SetMachineCredentials) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareSetMachineCredentials(agg, generator, set))
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

func prepareSetMachineCredentials(a *user.Aggregate, generator crypto.Generator, set *SetMachineCredentials) preparation.Validation {
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
			set.ClientID = writeModel.UserName

			clientSecret, secretString, err := domain.NewMachineClientSecret(generator)
			if err != nil {
				return nil, err
			}
			set.ClientSecret = secretString

			return []eventstore.Command{
				user.NewMachineCredentialsSetEvent(ctx, &a.Aggregate, clientSecret),
			}, nil
		}, nil
	}
}

func (c *Commands) RemoveMachineCredentials(ctx context.Context, userID string, resourceOwner string) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareRemoveMachineCredentials(agg))
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

func prepareRemoveMachineCredentials(a *user.Aggregate) preparation.Validation {
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
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-coi82n", "Errors.User.Credentials.NotFound")
			}
			return []eventstore.Command{
				user.NewMachineCredentialsRemovedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}

func (c *Commands) VerifyMachineCredentials(ctx context.Context, userID string, resourceOwner string, secret string) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareVerifyMachineCredentials(agg, secret, c.userPasswordAlg))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	for _, cmd := range cmds {
		if cmd.Type() == user.MachineCredentialsCheckFailedType {
			logging.OnError(err).Error("could not push event MachineCredentialsCheckFailed")
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3kjh", "Errors.User.Credentials.ClientSecretInvalid")
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

func prepareVerifyMachineCredentials(a *user.Aggregate, secret string, algorithm crypto.HashAlgorithm) preparation.Validation {
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
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-x8910n", "Errors.User.Credentials.NotFound")
			}
			err = crypto.CompareHash(writeModel.ClientSecret, []byte(secret), algorithm)
			if err == nil {
				return []eventstore.Command{
					user.NewMachineCredentialsCheckSucceededEvent(ctx, &a.Aggregate),
				}, nil
			}
			return []eventstore.Command{
				user.NewMachineCredentialsCheckFailedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}
