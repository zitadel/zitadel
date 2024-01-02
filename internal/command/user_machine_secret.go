package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
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
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareGenerateMachineSecret(a *user.Aggregate, generator crypto.Generator, set *GenerateMachineSecret) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-x0992n", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-bzoqjs", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-x8910n", "Errors.User.NotExisting")
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
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareRemoveMachineSecret(a *user.Aggregate) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-0qp2hus", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-bzosjs", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-x7s802", "Errors.User.NotExisting")
			}
			if writeModel.ClientSecret == nil {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-coi82n", "Errors.User.Machine.Secret.NotExisting")
			}
			return []eventstore.Command{
				user.NewMachineSecretRemovedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}

func (c *Commands) MachineSecretCheckSucceeded(ctx context.Context, userID, resourceOwner string) {
	agg := user.NewAggregate(userID, resourceOwner)
	c.asyncPush(ctx, user.NewMachineSecretCheckSucceededEvent(ctx, &agg.Aggregate))
}

func (c *Commands) MachineSecretCheckFailed(ctx context.Context, userID, resourceOwner string) {
	agg := user.NewAggregate(userID, resourceOwner)
	c.asyncPush(ctx, user.NewMachineSecretCheckFailedEvent(ctx, &agg.Aggregate))
}
